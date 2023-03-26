package istiomapper

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/istiomapper/list"
	"github.com/spf13/cobra"
)

var IstioMapperCmd = &cobra.Command{
	Use:     "istio-mapper",
	GroupID: groups.OSSGroup.ID,
	Aliases: []string{"istio"},
	Short:   "Map Istio traffic with ",
	Hidden:  true,
}

func init() {
	IstioMapperCmd.AddCommand(list.ListCmd)
}
