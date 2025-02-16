package intentsvisualizer

import (
	"aqwari.net/xml/xmltree"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/nfnt/resize"
	"github.com/otterize/intents-operator/src/operator/api/v2beta1"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
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
	GraphFormatKey      = "format"
	OutputPathKey       = "output-path"
	OutputPathShorthand = "o"
)

//go:embed watermark.png
var watermarkFile []byte

//go:embed watermark.svg
var watermarkSVG []byte

type Visualizer struct {
	graphviz       *graphviz.Graphviz
	graph          *cgraph.Graph
	nodeCache      map[string]*cgraph.Node
	graphFormat    graphviz.Format
	outputFilepath string
}

func NewVisualizer() (*Visualizer, error) {
	format := viper.GetString(GraphFormatKey)
	graphFormat, err := getGraphvizFormat(format)
	if err != nil {
		return nil, err
	}

	outputFilepath := viper.GetString(OutputPathKey)

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}
	graph.SetRankDir(cgraph.LRRank)

	v := Visualizer{
		graphviz:       g,
		graph:          graph,
		nodeCache:      make(map[string]*cgraph.Node, 0),
		graphFormat:    graphFormat,
		outputFilepath: outputFilepath,
	}
	return &v, nil
}

func (v *Visualizer) Close() {
	if err := v.graph.Close(); err != nil {
		panic(err)
	}
	if err := v.graphviz.Close(); err != nil {
		panic(err)
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

func (v *Visualizer) populateNodeCache(intents []v2beta1.ClientIntents) error {
	for _, intent := range intents {
		clientNS := intent.Namespace
		clientName := getServiceNameWithNamespace(clientNS, intent.GetWorkloadName())
		if err := v.addToCache(clientName); err != nil {
			return err
		}
		for _, call := range intent.GetTargetList() {
			targetServiceName := getServiceNameWithNamespace(clientNS, call.GetTargetServerNameAsWritten())
			if err := v.addToCache(targetServiceName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visualizer) buildEdges(intents []v2beta1.ClientIntents) error {
	for _, intent := range intents {
		clientNS := intent.Namespace
		clientName := getServiceNameWithNamespace(clientNS, intent.GetWorkloadName())
		for _, call := range intent.GetTargetList() {
			targetServiceName := getServiceNameWithNamespace(clientNS, call.GetTargetServerNameAsWritten())
			_, err := v.graph.CreateEdge(
				fmt.Sprintf("%s to %s", clientName, targetServiceName),
				v.nodeCache[clientName],
				v.nodeCache[targetServiceName])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Visualizer) renderOutputByFormat() ([]byte, error) {
	if v.graphFormat == graphviz.SVG {
		return v.renderSVG()
	}

	return v.renderImage()
}

func (v *Visualizer) RenderOutputToFile() error {
	outputData, err := v.renderOutputByFormat()
	if err != nil {
		return err
	}

	outputFile, err := os.Create(v.outputFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(outputData)
	if err != nil {
		return err
	}

	prints.PrintCliStderr("Exported graph as %s format to path %s", v.graphFormat, v.outputFilepath)
	return nil
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

func (v *Visualizer) Build(intents []v2beta1.ClientIntents) error {
	if err := v.populateNodeCache(intents); err != nil {
		return err
	}
	if err := v.buildEdges(intents); err != nil {
		return err
	}

	return nil
}

func getServiceNameWithNamespace(clientNS, name string) string {
	if len(strings.Split(name, ".")) > 1 {
		return name
	}
	return fmt.Sprintf("%s.%s", name, clientNS)
}

func InitVisualizeOutputFlags(cmd *cobra.Command) {
	cmd.Flags().String(GraphFormatKey, "png", "Graph output format (png/jpg/svg)")
	cmd.Flags().StringP(OutputPathKey, OutputPathShorthand, "", "exported graph output file path")
	cobra.CheckErr(cmd.MarkFlagRequired(OutputPathKey))
}
