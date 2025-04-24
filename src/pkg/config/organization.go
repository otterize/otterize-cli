package config

import (
	"github.com/otterize/intents-operator/src/shared/errors"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/spf13/viper"
)

var ErrNoOrganization = errors.New("no organization exists in config or as parameter")

func ResolveOrgID() (string, error) {
	if viper.GetString(ApiSelectedOrganizationId) != "" {
		return viper.GetString(ApiSelectedOrganizationId), nil
	}

	var c Config
	loaded, err := LoadConfigFile(&c, ApiCredentialsFilename)
	must.Must(err)

	if !loaded {
		return "", ErrNoOrganization
	}

	if c.OrganizationId != "" {
		return c.OrganizationId, nil
	}

	return "", ErrNoOrganization
}
