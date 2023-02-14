package visualize

import (
	"context"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
	GraphFormatKey      = "graph-format"
)

var AllowedFormats = []string{"png", "dot", "svg", "jpg"}

type Visualizer struct {
	*graphviz.Graphviz
	graph     *cgraph.Graph
	nodeCache map[string]*cgraph.Node
}

func NewVisualizer() *Visualizer {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	return &Visualizer{
		Graphviz:  g,
		graph:     graph,
		nodeCache: make(map[string]*cgraph.Node, 0),
	}
}

func (v *Visualizer) addToCache(nodeName string) error {
	if _, ok := v.nodeCache[nodeName]; !ok {
		node, err := v.graph.CreateNode(nodeName)
		if err != nil {
			return err
		}
		v.nodeCache[nodeName] = node
	}
	return nil
}

func (v *Visualizer) populateNodeCache(serviceIntents []mapperclient.ServiceIntentsUpToMapperV017ServiceIntents) error {
	for _, service := range serviceIntents {
		clientNameWithNS := fmt.Sprintf("%s.%s", service.Client.Name, service.Client.Namespace)
		if err := v.addToCache(clientNameWithNS); err != nil {
			return err
		}
		for _, intent := range service.Intents {
			targetNameWithNS := formatTargetServiceName(service.Client.Namespace, intent)
			if err := v.addToCache(targetNameWithNS); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visualizer) buildEdges(serviceIntents []mapperclient.ServiceIntentsUpToMapperV017ServiceIntents) error {
	for _, service := range serviceIntents {
		clientNameWithNS := fmt.Sprintf("%s.%s", service.Client.Name, service.Client.Namespace)
		for _, intent := range service.Intents {
			targetNameWithNS := formatTargetServiceName(service.Client.Namespace, intent)
			_, err := v.graph.CreateEdge(
				fmt.Sprintf("%s to %s", clientNameWithNS, targetNameWithNS),
				v.nodeCache[clientNameWithNS],
				v.nodeCache[targetNameWithNS])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var VisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize an access graph for network mapper intents using go-graphviz",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			namespacesFilter := viper.GetStringSlice(NamespacesKey)
			format := viper.GetString(GraphFormatKey)
			if !lo.Contains(AllowedFormats, strings.ToLower(format)) {
				return fmt.Errorf("unsupported format: %s", format)
			}

			outFile := viper.GetString(config.OutputPathKey)
			servicesIntents, err := c.ServiceIntents(context.Background(), namespacesFilter)
			if err != nil {
				return err
			}
			visualizer := NewVisualizer()
			defer func() {
				if err := visualizer.graph.Close(); err != nil {
					panic(err)
				}
				if err := visualizer.Close(); err != nil {
					panic(err)
				}
			}()
			if err := visualizer.populateNodeCache(servicesIntents); err != nil {
				return err
			}

			if err := visualizer.buildEdges(servicesIntents); err != nil {
				return err
			}

			if err := visualizer.RenderFilename(visualizer.graph, graphviz.Format(format), outFile); err != nil {
				return err
			}

			output.PrintStdout("Exported graph as %s format to path %s", format, outFile)
			return nil
		})
	},
}

func init() {
	VisualizeCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	VisualizeCmd.Flags().String(GraphFormatKey, "svg", "Graph output format (png/svg/jpg/dot)")
	VisualizeCmd.Flags().StringP(config.OutputPathKey, config.OutputPathShorthand, "", "exported graph output file path")
	cobra.CheckErr(VisualizeCmd.MarkFlagRequired(config.OutputPathKey))
}

func formatTargetServiceName(clientNS string, target mapperclient.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) string {
	ns := lo.Ternary(len(target.Namespace) != 0, target.Namespace, clientNS)
	return fmt.Sprintf("%s.%s", target.Name, ns)
}
