package cli

import "github.com/spf13/cobra"

const (
	ClustersKey           = "clusters"
	ClustersShortHand     = "c"
	EnvironmentsKey       = "envs"
	EnvironmentsShorthand = "e"
	NamespacesKey         = "namespaces"
	NamespacesShorthand   = "n"
	ServicesKey           = "services"
	ServicesShorthand     = "s"
)

func RegisterStandardFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(ClustersKey, ClustersShortHand, nil, "filter by clusters")
	cmd.Flags().StringSliceP(EnvironmentsKey, EnvironmentsShorthand, nil, "filter by environments")
	cmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter by namespaces")
	cmd.Flags().StringSliceP(ServicesKey, ServicesShorthand, nil, "filter by services")
}
