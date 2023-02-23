package visualize

import (
	"bytes"
	"context"
	_ "embed"
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
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

const (
	NamespacesKey                    = "namespaces"
	NamespacesShorthand              = "n"
	GraphFormatKey                   = "format"
	WatermarkHeightPercentageOfGraph = 20
)

//go:embed watermark.png
var watermarkFile []byte

type Visualizer struct {
	*graphviz.Graphviz
	graph       *cgraph.Graph
	nodeCache   map[string]*cgraph.Node
	graphFormat graphviz.Format
}

func NewVisualizer(format graphviz.Format) *Visualizer {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	graph.SetRankDir(cgraph.LRRank)
	return &Visualizer{
		Graphviz:    g,
		graph:       graph,
		nodeCache:   make(map[string]*cgraph.Node, 0),
		graphFormat: format,
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
		clientName := fmt.Sprintf("%s.%s", service.Client.Name, service.Client.Namespace)
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
		clientName := fmt.Sprintf("%s.%s", service.Client.Name, service.Client.Namespace)
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
			visualizer := NewVisualizer(graphFormat)
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
			if err := visualizer.addWatermark(outFile); err != nil {
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

func (v *Visualizer) formatTargetServiceName(clientNS string, target mapperclient.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) string {
	ns := lo.Ternary(len(target.Namespace) != 0, target.Namespace, clientNS)
	return fmt.Sprintf("%s.%s", target.Name, ns)
}

func (v *Visualizer) addWatermark(graphPath string) error {
	graphImg, err := v.decodeImage(graphPath)
	if err != nil {
		return err
	}

	watermarkImg, err := png.Decode(bytes.NewReader(watermarkFile))
	if err != nil {
		return err
	}

	graphBounds := graphImg.Bounds()
	watermarkHeight := graphBounds.Max.Y / WatermarkHeightPercentageOfGraph
	watermarkWidth := watermarkHeight * (watermarkImg.Bounds().Max.Y / watermarkImg.Bounds().Max.X)

	resizedWatermark := resize.Resize(uint(watermarkWidth), uint(watermarkHeight), watermarkImg, resize.Lanczos3)

	graphImgBounds := graphImg.Bounds()
	graphImgBounds.Max.X = graphImgBounds.Max.X + watermarkWidth
	graphImgBounds.Max.Y = graphImgBounds.Max.Y + watermarkHeight

	watermarkOffset := image.Pt(graphImgBounds.Dx()-resizedWatermark.Bounds().Dx(), graphImgBounds.Dy()-resizedWatermark.Bounds().Dy())
	whiteOffset := image.Pt(0, watermarkHeight)

	graphImgWithWatermark := image.NewRGBA(graphImgBounds)
	draw.Draw(graphImgWithWatermark, graphImgBounds, graphImg, image.Point{}, draw.Src)
	// Add a white offset matching watermark size, so we can add the watermark image under the graph
	draw.Draw(graphImgWithWatermark, graphImgBounds.Bounds().Add(whiteOffset),
		&image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255}}, image.Point{}, draw.Over)
	draw.Draw(graphImgWithWatermark, resizedWatermark.Bounds().Add(watermarkOffset), resizedWatermark, image.Point{}, draw.Over)

	return v.encodeImage(graphPath, graphImgWithWatermark)
}

func (v *Visualizer) encodeImage(path string, img image.Image) error {
	outFile, err := os.Create(path)
	defer outFile.Close()
	if err != nil {
		return err
	}

	switch v.graphFormat {
	case graphviz.JPG:
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case graphviz.PNG:
		return png.Encode(outFile, img)
	default:
		return fmt.Errorf("unsupported format: %s", v.graphFormat)
	}
}

func (v *Visualizer) decodeImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	switch v.graphFormat {
	case graphviz.JPG:
		return jpeg.Decode(file)
	case graphviz.PNG:
		return png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported format: %s", v.graphFormat)
	}
}

func init() {
	VisualizeCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	VisualizeCmd.Flags().String(GraphFormatKey, "jpg", "Graph output format (png/jpg)")
	VisualizeCmd.Flags().StringP(config.OutputPathKey, config.OutputPathShorthand, "", "exported graph output file path")
	cobra.CheckErr(VisualizeCmd.MarkFlagRequired(config.OutputPathKey))
}
