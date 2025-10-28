package cloudimg

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
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

// calculateDynamicParams は単語数に応じて動的にパラメータを計算する
func calculateDynamicParams(wordCounts map[string]int) (width, height, fontMax, fontMin int) {
	// 単語の種類数
	uniqueWords := len(wordCounts)

	// 総単語数（出現回数の合計）
	totalWords := 0
	maxFreq := 0
	for _, count := range wordCounts {
		totalWords += count
		if count > maxFreq {
			maxFreq = count
		}
	}

	fmt.Printf("=== 動的パラメータ計算 ===\n")
	fmt.Printf("ユニーク単語数: %d\n", uniqueWords)
	fmt.Printf("総単語数: %d\n", totalWords)
	fmt.Printf("最大出現回数: %d\n", maxFreq)

	// 基本サイズの計算（単語数に応じてスケール）
	// 単語数が少ない：小さいキャンバス
	// 単語数が多い：大きいキャンバス
	baseSize := 2048 // 基本サイズ

	if uniqueWords <= 20 {
		// 少ない単語数：小さめのキャンバス
		width = 2048
		height = 2048
		fontMax = 300
		fontMin = 40
	} else if uniqueWords <= 50 {
		// 中程度の単語数
		width = 2560
		height = 2560
		fontMax = 280
		fontMin = 35
	} else if uniqueWords <= 100 {
		// やや多い単語数
		width = 3072
		height = 3072
		fontMax = 250
		fontMin = 30
	} else if uniqueWords <= 200 {
		// 多い単語数
		width = 3584
		height = 3584
		fontMax = 200
		fontMin = 25
	} else if uniqueWords <= 400 {
		// 非常に多い単語数
		width = 4096
		height = 4096
		fontMax = 160
		fontMin = 20
	} else {
		// 極端に多い単語数：さらに大きく、フォントは小さく
		// 対数スケールで増加
		scale := math.Log10(float64(uniqueWords) / 200.0)
		width = int(float64(baseSize) * (2.0 + scale*0.4))
		height = width
		fontMax = int(160.0 / (1.0 + scale*0.3))
		fontMin = int(20.0 / (1.0 + scale*0.15))

		// 上限を設定
		if width > 8192 {
			width = 8192
			height = 8192
		}
		if fontMax < 80 {
			fontMax = 80
		}
		if fontMin < 12 {
			fontMin = 12
		}
	}

	fmt.Printf("計算されたパラメータ:\n")
	fmt.Printf("  キャンバス: %dx%d\n", width, height)
	fmt.Printf("  フォントサイズ: %d - %d\n", fontMin, fontMax)
	fmt.Printf("========================\n")

	return width, height, fontMax, fontMin
}

func GenCloud(wordCounts map[string]int) {
	if len(wordCounts) == 0 {
		fmt.Println("ERROR: wordCounts is empty!")
		return
	}

	// 動的にパラメータを計算
	width, height, fontMax, fontMin := calculateDynamicParams(wordCounts)

	conf := Conf{
		FontMaxSize:     fontMax,
		FontMinSize:     fontMin,
		RandomPlacement: false,
		FontFile:        defaultFontPath,
		Colors:          DefaultColors,
		Width:           width,
		Height:          height,
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

	fmt.Println("Drawing wordcloud...")
	img := w.Draw()
	fmt.Printf("Image bounds: %v\n", img.Bounds())
	outputFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	fmt.Println("Encoding image to output.png...")
	if err := png.Encode(outputFile, img); err != nil {
		panic(err)
	}
	fmt.Println("✓ Successfully created output.png")
}
