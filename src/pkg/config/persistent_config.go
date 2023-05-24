package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"time"
)

const ApiCredentialsFilename = "credentials"
const OtterizeContextIdFileName = "contextId"

type Config struct {
	ClientId       string    `json:"client_id,omitempty"`
	ClientSecret   string    `json:"client_secret,omitempty"`
	UserToken      string    `json:"user_token,omitempty"`
	Expiry         time.Time `json:"expiry,omitempty"`
	OrganizationId string    `json:"organization_id"`
}

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
		logrus.Warningf("API client id and client secret should not be set in config. Instead, store the credentials in ~/%s/%s", OtterizeConfigDirName, ApiCredentialsFilename)
	}

	if viper.IsSet(ApiClientIdKey) && viper.IsSet(ApiClientSecretKey) {
		// auth was provided as a flag or env var
		return
	}

	var Config Config
	loaded, err := LoadConfigFile(&Config, ApiCredentialsFilename)
	must.Must(err)

	if !loaded {
		return
	}

	if Config.UserToken != "" {
		viper.Set(ApiUserTokenKey, Config.UserToken)
	}
	if !Config.Expiry.IsZero() {
		viper.Set(ApiUserTokenExpiryKey, Config.Expiry)
	}
	if Config.OrganizationId != "" {
		viper.Set(ApiSelectedOrganizationId, Config.OrganizationId)
	}
	if Config.ClientId != "" || Config.ClientSecret != "" {
		viper.Set(ApiClientIdKey, Config.ClientId)
		viper.Set(ApiClientSecretKey, Config.ClientSecret)
	}
}

func SaveConfig(conf Config) error {
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

func getContextId() string {
	return uuid.NewString()
}

func InitContextId() {
	contextId := getContextId()
	viper.Set(ContextIdKey, contextId)
}
