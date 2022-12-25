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
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	UserToken    string `json:"user_token,omitempty"`
}

type OrganizationConfig struct {
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
const ApiOrganizationFilename = "organization"

func LoadApiCredentialsFile() {
	if viper.InConfig(ApiClientIdKey) || viper.InConfig(ApiClientSecretKey) {
		logrus.Warningf("API client ID and client secret should not be set in config. Instead, store the credentials in ~/%s/%s", OtterizeConfigDirName, ApiCredentialsFilename)
	}

	if viper.IsSet(ApiClientIdKey) && viper.IsSet(ApiClientSecretKey) {
		// auth was provided as a flag or env var
		return
	}

	// try to read the auth from file
	configDir, err := OtterizeConfigDirPath()
	if errors.Is(err, os.ErrNotExist) {
		// no home dir probably
		return
	} else if err != nil {
		logrus.Warningf("Failed to find otterize config dir: %s", err)
		return
	}

	credentialsFilePath := filepath.Join(configDir, ApiCredentialsFilename)
	file, err := os.Open(credentialsFilePath)
	if errors.Is(err, os.ErrNotExist) {
		// no auth file
		return
	}

	if err != nil {
		must.Must(fmt.Errorf("failed reading otterize api credentials file: %w", err))
		return
	}

	var secretConfig SecretConfig

	err = json.NewDecoder(file).Decode(&secretConfig)
	if err != nil {
		must.Must(fmt.Errorf("failed to decode secret config: %w", err))
		return
	}

	if secretConfig.UserToken != "" {
		viper.Set(ApiUserTokenKey, secretConfig.UserToken)
		return
	}

	if secretConfig.ClientId == "" || secretConfig.ClientSecret == "" {
		return
	}

	viper.Set(ApiClientIdKey, secretConfig.ClientId)
	viper.Set(ApiClientSecretKey, secretConfig.ClientSecret)
}

func SaveSecretConfig(conf SecretConfig) error {
	dirPath, err := OtterizeConfigDirPath()
	if err != nil {
		return fmt.Errorf("unable to get config dir path: %w", err)
	}

	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to create config path: %w", err)
	}

	tokenPath := path.Join(dirPath, ApiCredentialsFilename)
	file, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w", tokenPath, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(conf)
	if err != nil {
		return fmt.Errorf("unable to write auth to path %s: %w", tokenPath, err)
	}

	return nil
}

func SaveSelectedOrganization(orgConf OrganizationConfig) error {
	dirPath, err := OtterizeConfigDirPath()
	if err != nil {
		return fmt.Errorf("unable to get config dir path: %w", err)
	}

	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to create config path: %w", err)
	}

	tokenPath := path.Join(dirPath, ApiOrganizationFilename)
	file, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w", tokenPath, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&orgConf)
	if err != nil {
		return fmt.Errorf("unable to write auth to path %s: %w", tokenPath, err)
	}

	viper.Set(ApiSelectedOrganizationId, orgConf.OrganizationId)

	return nil
}

func LoadSelectedOrganizationFile() {
	if viper.IsSet(ApiSelectedOrganizationId) {
		// org ID was provided as a flag or env var
		return
	}

	// try to read the auth from file
	configDir, err := OtterizeConfigDirPath()
	if errors.Is(err, os.ErrNotExist) {
		// no home dir probably
		return
	} else if err != nil {
		logrus.Warningf("Failed to find otterize config dir: %s", err)
		return
	}

	orgFilePath := filepath.Join(configDir, ApiSelectedOrganizationId)
	file, err := os.Open(orgFilePath)
	if errors.Is(err, os.ErrNotExist) {
		// no auth file
		return
	}

	if err != nil {
		must.Must(fmt.Errorf("failed reading otterize selected org ID file: %w", err))
		return
	}

	var orgConfig OrganizationConfig

	err = json.NewDecoder(file).Decode(&orgConfig)
	if err != nil {
		must.Must(fmt.Errorf("failed to decode org config: %w", err))
		return
	}

	if orgConfig.OrganizationId != "" {
		viper.Set(ApiSelectedOrganizationId, orgConfig.OrganizationId)
		return
	}
}
