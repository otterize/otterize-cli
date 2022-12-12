package get

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/integrations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func getIntegrationText(i integrations.IntegrationWithStatus) string {
	integration := i.IntegrationFields
	s := fmt.Sprintf("id: %s\n", integration.Id)
	s += fmt.Sprintf("type: %s\n", integration.IntegrationType)
	s += fmt.Sprintf("name: %s\n", integration.Name)

	envNames := lo.Map(integration.Environments, func(env integrations.IntegrationFieldsEnvironmentsEnvironment, i int) string {
		return fmt.Sprintf("%s (%s)", env.Name, env.Id)
	})

	s += fmt.Sprintf("environments: %s\n", strings.Join(envNames, ", "))

	if i.Status.Id != "" {
		s += fmt.Sprintf("controller last seen: %s\n", i.Status.LastSeen)
		s += fmt.Sprintf("intents last applied: %s\n", i.Status.IntentsStatus.AppliedAt)
		applyError := i.Status.IntentsStatus.ApplyError
		if applyError != "" {
			s += fmt.Sprintf("error applying intents: %s\n", applyError)
		}
	}

	return s
}

var GetIntegrationCmd = &cobra.Command{
	Use:          "get",
	Short:        `Gets details for an integration.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := integrations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := viper.GetString(IdKey)
		name := viper.GetString(NameKey)
		var integration integrations.IntegrationWithStatus
		var err error
		if id != "" {
			integration, err = c.GetIntegration(ctxTimeout, id)
		} else {
			integration, err = c.GetIntegrationByName(ctxTimeout, name)
		}
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(integration, getIntegrationText)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	config.RegisterStringArg(GetIntegrationCmd, IdKey, "integration ID", false)
	config.RegisterStringArg(GetIntegrationCmd, NameKey, "integration name", false)
}
