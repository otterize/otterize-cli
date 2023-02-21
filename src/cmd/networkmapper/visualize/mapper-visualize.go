package visualize

import (
	"context"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/nfnt/resize"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
	GraphFormatKey      = "format"
)

type Encoder interface {
}

type Visualizer struct {
	*graphviz.Graphviz
	graph           *cgraph.Graph
	nodeCache       map[string]*cgraph.Node
	singleNamespace bool
}

func NewVisualizer(singleNamespace bool) *Visualizer {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	return &Visualizer{
		Graphviz:        g,
		graph:           graph,
		nodeCache:       make(map[string]*cgraph.Node, 0),
		singleNamespace: singleNamespace,
	}
}

func (v *Visualizer) addToCache(nodeName string) error {
	if _, ok := v.nodeCache[nodeName]; !ok {
		node, err := v.graph.CreateNode(nodeName)
		if err != nil {
			return err
		}
		node.SetMargin(0.2)
		v.nodeCache[nodeName] = node
	}
	return nil
}

func (v *Visualizer) populateNodeCache(serviceIntents []mapperclient.ServiceIntentsUpToMapperV017ServiceIntents) error {
	for _, service := range serviceIntents {
		clientName := lo.Ternary(v.singleNamespace, service.Client.Name, fmt.Sprintf("%s.%s", service.Client.Name, service.Client.Namespace))
		if err := v.addToCache(clientName); err != nil {
			return err
		}
		for _, intent := range service.Intents {
			targetServiceName := v.formatTargetServiceName(service.Client.Namespace, intent)
			if err := v.addToCache(targetServiceName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visualizer) buildEdges(serviceIntents []mapperclient.ServiceIntentsUpToMapperV017ServiceIntents) error {
	for _, service := range serviceIntents {
		clientName := lo.Ternary(v.singleNamespace, service.Client.Name, fmt.Sprintf("%s.%s", service.Client.Name, service.Client.Namespace))
		for _, intent := range service.Intents {
			targetNameWithNS := v.formatTargetServiceName(service.Client.Namespace, intent)
			_, err := v.graph.CreateEdge(
				fmt.Sprintf("%s to %s", clientName, targetNameWithNS),
				v.nodeCache[clientName],
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
			graphFormat, err := getGraphvizFormat(format)
			if err != nil {
				return err
			}

			outFile := viper.GetString(config.OutputPathKey)
			servicesIntents, err := c.ServiceIntents(context.Background(), namespacesFilter)
			if err != nil {
				return err
			}
			visualizer := NewVisualizer(len(namespacesFilter) == 1)
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
			if err := visualizer.RenderFilename(visualizer.graph, graphFormat, outFile); err != nil {
				return err
			}
			if err := visualizer.addWatermark(outFile, graphFormat); err != nil {
				return err
			}

			output.PrintStdout("Exported graph as %s format to path %s", format, outFile)
			return nil
		})
	},
}

func getGraphvizFormat(format string) (graphviz.Format, error) {
	switch strings.ToLower(format) {
	case "png":
		return graphviz.PNG, nil
	case "jpg", "jpeg":
		return graphviz.JPG, nil
	case "svg":
		return graphviz.SVG, nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func init() {
	VisualizeCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	VisualizeCmd.Flags().String(GraphFormatKey, "jpg", "Graph output format (png/svg/jpg)")
	VisualizeCmd.Flags().StringP(config.OutputPathKey, config.OutputPathShorthand, "", "exported graph output file path")
	cobra.CheckErr(VisualizeCmd.MarkFlagRequired(config.OutputPathKey))
}

func (v *Visualizer) formatTargetServiceName(clientNS string, target mapperclient.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) string {
	ns := lo.Ternary(len(target.Namespace) != 0, target.Namespace, clientNS)
	return fmt.Sprintf("%s.%s", target.Name, ns)
}

func (v *Visualizer) addWatermark(graphPath string, format graphviz.Format) error {
	graphFile, err := os.Open(graphPath)
	if err != nil {
		return err
	}
	defer graphFile.Close()
	graphImg, err := jpeg.Decode(graphFile)
	if err != nil {
		return err
	}

	watermarkFile, err := os.Open("/home/evya/watermark.png")
	if err != nil {
		return err
	}
	defer watermarkFile.Close()
	watermarkImg, err := png.Decode(watermarkFile)
	if err != nil {
		return err
	}

	graphBounds := graphImg.Bounds()
	height := graphBounds.Max.Y / 8
	wmWidsh := height * (watermarkImg.Bounds().Max.Y / watermarkImg.Bounds().Max.X)
	resizedWatermark := resize.Resize(uint(wmWidsh), uint(height), watermarkImg, resize.Lanczos3)

	offset := image.Pt(graphBounds.Dx()-resizedWatermark.Bounds().Dx(), graphBounds.Dy()-resizedWatermark.Bounds().Dy())
	graphImgBounds := graphImg.Bounds()
	graphImgBounds.Max.X = graphImgBounds.Max.X + wmWidsh
	graphImgBounds.Max.Y = graphImgBounds.Max.Y + height

	graphImgWithWatermark := image.NewRGBA(graphImgBounds)
	draw.Draw(graphImgWithWatermark, graphImgBounds, graphImg, image.Point{}, draw.Src)
	draw.Draw(graphImgWithWatermark, resizedWatermark.Bounds().Add(offset), resizedWatermark, image.Point{}, draw.Over)

	result, err := os.Create(graphPath)
	defer result.Close()
	if err != nil {
		return err
	}

	if err := jpeg.Encode(result, graphImgWithWatermark, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
		return err
	}
	return nil
}
