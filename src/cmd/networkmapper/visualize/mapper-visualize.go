package visualize

import (
	"aqwari.net/xml/xmltree"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/nfnt/resize"
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
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
	GraphFormatKey      = "format"
	OutputPathKey       = "output-path"
	OutputPathShorthand = "o"
)

//go:embed watermark.png
var watermarkFile []byte

//go:embed watermark.svg
var watermarkSVG []byte

type Visualizer struct {
	graphviz    *graphviz.Graphviz
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
		graphviz:    g,
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

			outputFilepath := viper.GetString(OutputPathKey)
			servicesIntents, err := c.ServiceIntents(context.Background(), namespacesFilter)
			if err != nil {
				return err
			}
			visualizer := NewVisualizer(graphFormat)
			defer func() {
				if err := visualizer.graph.Close(); err != nil {
					panic(err)
				}
				if err := visualizer.graphviz.Close(); err != nil {
					panic(err)
				}
			}()

			if err := visualizer.populateNodeCache(servicesIntents); err != nil {
				return err
			}
			if err := visualizer.buildEdges(servicesIntents); err != nil {
				return err
			}
			outputData, err := visualizer.renderOutputByFormat(graphFormat)
			if err != nil {
				return err
			}

			outputFile, err := os.Create(outputFilepath)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = outputFile.Write(outputData)
			if err != nil {
				return err
			}

			output.PrintStderr("Exported graph as %s format to path %s", format, outputFilepath)
			return nil
		})
	},
}

func (v *Visualizer) renderOutputByFormat(format graphviz.Format) ([]byte, error) {
	if format == graphviz.SVG {
		return v.renderSVG()
	}

	return v.renderImage()
}

func (v *Visualizer) renderSVG() ([]byte, error) {
	outputBuf := &bytes.Buffer{}
	err := v.graphviz.Render(v.graph, graphviz.SVG, outputBuf)
	if err != nil {
		return nil, err
	}

	doc, err := xmltree.Parse(outputBuf.Bytes())
	if err != nil {
		return nil, err
	}

	watermark, err := xmltree.Parse(watermarkSVG)
	if err != nil {
		return nil, err
	}

	doc.Children = append(doc.Children, *watermark)

	outputWithWatermark := xmltree.Marshal(doc)
	return outputWithWatermark, err
}

func (v *Visualizer) renderImage() ([]byte, error) {
	outputImg, err := v.graphviz.RenderImage(v.graph)
	if err != nil {
		return nil, err
	}
	watermarkedImg, err := v.addWatermarkToImage(outputImg)
	if err != nil {
		return nil, err
	}

	outputImgBytes, err := v.encodeImage(watermarkedImg)
	if err != nil {
		return nil, err
	}

	return outputImgBytes, nil
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

func (v *Visualizer) addWatermarkToImage(graphImg image.Image) (image.Image, error) {
	watermarkImg, err := png.Decode(bytes.NewReader(watermarkFile))
	if err != nil {
		return nil, err
	}

	graphBounds := graphImg.Bounds()
	watermarkWidth := int(float64(graphBounds.Max.X) * (20.0 / 100.0))
	watermarkHeight := int(float64(watermarkWidth) * (float64(watermarkImg.Bounds().Dy()) / float64(watermarkImg.Bounds().Dx())))

	resizedWatermark := resize.Resize(uint(watermarkWidth), uint(watermarkHeight), watermarkImg, resize.Lanczos3)

	graphImgWithWatermarkBounds := graphImg.Bounds()
	graphImgWithWatermarkBounds.Max.Y = graphImgWithWatermarkBounds.Max.Y + watermarkHeight

	graphImgWithWatermark := image.NewRGBA(graphImgWithWatermarkBounds)
	whiteBounds := graphImgWithWatermark.Bounds()
	whiteBounds.Min.Y = graphImgWithWatermark.Bounds().Dy() - watermarkHeight - 1
	watermarkBounds := graphImgWithWatermark.Bounds()
	watermarkBounds.Min.X = watermarkBounds.Min.X + (watermarkBounds.Max.X - resizedWatermark.Bounds().Dx())
	watermarkBounds.Min.Y = watermarkBounds.Min.Y + (watermarkBounds.Max.Y - resizedWatermark.Bounds().Dy())
	draw.Draw(graphImgWithWatermark, graphImgWithWatermarkBounds, graphImg, image.Point{}, draw.Src)
	// Add a white offset matching watermark size, so we can add the watermark image under the graph
	draw.Draw(graphImgWithWatermark, whiteBounds,
		&image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)
	fiftyPercentOpacityMask := image.NewUniform(color.Alpha{A: 128})
	draw.DrawMask(graphImgWithWatermark, watermarkBounds, resizedWatermark, image.Point{}, fiftyPercentOpacityMask, image.Point{}, draw.Over)

	return graphImgWithWatermark, nil
}

func (v *Visualizer) encodeImage(img image.Image) ([]byte, error) {
	out := make([]byte, 0)
	writer := bytes.NewBuffer(out)

	switch v.graphFormat {
	case graphviz.JPG:
		err := jpeg.Encode(writer, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			return nil, err
		}
	case graphviz.PNG:
		err := png.Encode(writer, img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", v.graphFormat)
	}

	return writer.Bytes(), nil
}

func init() {
	VisualizeCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	VisualizeCmd.Flags().String(GraphFormatKey, "png", "Graph output format (png/jpg)")
	VisualizeCmd.Flags().StringP(OutputPathKey, OutputPathShorthand, "", "exported graph output file path")
	cobra.CheckErr(VisualizeCmd.MarkFlagRequired(OutputPathKey))
}
