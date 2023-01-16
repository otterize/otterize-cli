package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"path"
	"path/filepath"
)

type SecretConfig struct {
	ClientId       string `json:"client_id,omitempty"`
	ClientSecret   string `json:"client_secret,omitempty"`
	UserToken      string `json:"user_token,omitempty"`
	OrganizationId string `json:"organization_id"`
}

func GetAPITokenSource(ctx context.Context) oauth2.TokenSource {
	if viper.IsSet(ApiUserTokenKey) {
		return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString(ApiUserTokenKey)})
	}
	cfg := clientcredentials.Config{
		ClientID:     viper.GetString(ApiClientIdKey),
		ClientSecret: viper.GetString(ApiClientSecretKey),
		TokenURL:     fmt.Sprintf("%s/auth/tokens/token", viper.GetString(OtterizeAPIAddressKey)),
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	return cfg.TokenSource(ctx)
}

func GetAPIToken(ctx context.Context) string {
	tokenSrc := GetAPITokenSource(ctx)
	token, err := tokenSrc.Token()
	must.Must(err)
	return token.AccessToken
}

const ApiCredentialsFilename = "credentials"

func LoadConfigFile(output any, filename string) (bool, error) {
	configDir, err := OtterizeConfigDirPath()
	if errors.Is(err, os.ErrNotExist) {
		// no home dir probably
		return false, nil
	} else if err != nil {
		logrus.Warningf("Failed to find otterize config dir: %s", err)
		return false, nil
	}

	configFilePath := filepath.Join(configDir, filename)
	file, err := os.Open(configFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed reading otterize api credentials file: %w", err)
	}

	err = json.NewDecoder(file).Decode(output)
	if err != nil {
		return false, fmt.Errorf("failed to decode secret config: %w", err)
	}

	return true, nil
}

func LoadApiCredentialsFile() {
	if viper.InConfig(ApiClientIdKey) || viper.InConfig(ApiClientSecretKey) {
		logrus.Warningf("API client ID and client secret should not be set in config. Instead, store the credentials in ~/%s/%s", OtterizeConfigDirName, ApiCredentialsFilename)
	}

	if viper.IsSet(ApiClientIdKey) && viper.IsSet(ApiClientSecretKey) {
		// auth was provided as a flag or env var
		return
	}

	var secretConfig SecretConfig
	loaded, err := LoadConfigFile(&secretConfig, ApiCredentialsFilename)
	must.Must(err)

	if !loaded {
		return
	}

	if secretConfig.UserToken != "" {
		viper.Set(ApiUserTokenKey, secretConfig.UserToken)
		return
	}

	if secretConfig.ClientId == "" || secretConfig.ClientSecret == "" {
		return
	}

	if secretConfig.OrganizationId != "" {
		viper.Set(ApiSelectedOrganizationId, secretConfig.OrganizationId)
		return
	}

	viper.Set(ApiClientIdKey, secretConfig.ClientId)
	viper.Set(ApiClientSecretKey, secretConfig.ClientSecret)
}

func SaveSecretConfig(conf SecretConfig) error {
	return SaveJSONConfig(conf, ApiCredentialsFilename)
}

func SaveJSONConfig(config any, filename string) error {
	dirPath, err := OtterizeConfigDirPath()
	if err != nil {
		return fmt.Errorf("unable to get config dir path: %w", err)
	}

	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to create config path: %w", err)
	}

	tokenPath := path.Join(dirPath, filename)
	file, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w", tokenPath, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("unable to write auth to path %s: %w", tokenPath, err)
	}

	return nil
}
