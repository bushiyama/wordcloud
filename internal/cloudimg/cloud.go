package cloudimg

import (
	"flag"
	"image/color"
	"image/png"
	"os"

	"github.com/psykhi/wordclouds"
)

const (
	defaultFontPath = "./fonts/marumonica/x12y16pxMaruMonica.ttf"
)

type MaskConf struct {
	File  string     `json:"file"`
	Color color.RGBA `json:"color"`
}
type Conf struct {
	FontMaxSize     int          `json:"font_max_size"`
	FontMinSize     int          `json:"font_min_size"`
	RandomPlacement bool         `json:"random_placement"`
	FontFile        string       `json:"font_file"`
	Colors          []color.RGBA `json:"colors"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	Mask            MaskConf     `json:"mask"`
}

var DefaultColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x59, 0x3a, 0xee, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}

func GenCloud(wordCounts map[string]int) {
	conf := Conf{
		FontMaxSize:     64 * 10,
		FontMinSize:     64,
		RandomPlacement: false,
		FontFile:        defaultFontPath,
		Colors:          DefaultColors,
		Width:           1024 * 3,
		Height:          1024 * 3,
		Mask: MaskConf{"", color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}},
	}
	var boxes []*wordclouds.Box
	colors := make([]color.Color, 0)
	for _, c := range conf.Colors {
		colors = append(colors, c)
	}

	w := wordclouds.NewWordcloud(wordCounts,
		wordclouds.FontFile(defaultFontPath),
		wordclouds.FontMaxSize(conf.FontMaxSize),
		wordclouds.FontMinSize(conf.FontMinSize),
		wordclouds.Colors(colors),
		wordclouds.MaskBoxes(boxes),
		wordclouds.Height(conf.Height),
		wordclouds.Width(conf.Width),
		wordclouds.RandomPlacement(conf.RandomPlacement),
	)

	img := w.Draw()
	var output = flag.String("output", "output.png", "path to output image")
	outputFile, _ := os.Create(*output)
	png.Encode(outputFile, img)
	outputFile.Close()
}
