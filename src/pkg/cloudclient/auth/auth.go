package auth

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"path/filepath"
	"time"
)

func GetAPITokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	if viper.IsSet(config.ApiClientIdKey) && viper.IsSet(config.ApiClientSecretKey) && viper.IsSet(config.OtterizeAPIAddressKey) {
		cfg := clientcredentials.Config{
			ClientID:     viper.GetString(config.ApiClientIdKey),
			ClientSecret: viper.GetString(config.ApiClientSecretKey),
			TokenURL:     fmt.Sprintf("%s/auth/tokens/token", viper.GetString(config.OtterizeAPIAddressKey)),
			AuthStyle:    oauth2.AuthStyleInParams,
		}
		return cfg.TokenSource(ctx), nil
	}
	if viper.IsSet(config.ApiUserTokenKey) {
		if viper.IsSet(config.ApiUserTokenExpiryKey) && viper.GetTime(config.ApiUserTokenExpiryKey).Before(time.Now()) {
			return nil, fmt.Errorf("your Otterize session token is expired")
		}

		return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString(config.ApiUserTokenKey)}), nil
	}

	return nil, fmt.Errorf("no API credentials were provided")
}

func GetAPIToken(ctx context.Context) string {
	tokenSrc, err := GetAPITokenSource(ctx)
	if err != nil {
		Fail(err)
	}

	token, err := tokenSrc.Token()
	must.Must(err)
	return token.AccessToken
}

func Fail(err error) {
	logrus.Errorf("Authentication failed: %s. To refresh your credentials, run '%s login'.",
		err.Error(), filepath.Base(os.Args[0]))
	logrus.Exit(1)
}
