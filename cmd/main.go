package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/bushiyama/wordcloud/internal/analyzer"
	"github.com/bushiyama/wordcloud/internal/cloudimg"
	"github.com/bushiyama/wordcloud/internal/kagomer"
)

var (
	word       string
	filePath   string
	outputFile string
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "w",
			Usage:       "input text directly",
			Destination: &word,
		},
		&cli.StringFlag{
			Name:        "f",
			Aliases:     []string{"file"},
			Usage:       "input text file path",
			Destination: &filePath,
		},
		&cli.StringFlag{
			Name:        "o",
			Aliases:     []string{"output"},
			Usage:       "output image file path (default: output.png)",
			Destination: &outputFile,
			Value:       "output.png",
		},
	}

	app.Action = func(c *cli.Context) error {
		// どちらか一方だけが指定されているかをチェック
		if word == "" && filePath == "" {
			return fmt.Errorf("either -w or -f option is required")
		}
		if word != "" && filePath != "" {
			return fmt.Errorf("cannot use both -w and -f options at the same time")
		}

		var text string
		if filePath != "" {
			// ファイルから読み込み
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			// エンコーディング変換
			text, err = convertToUTF8(data)
			if err != nil {
				return fmt.Errorf("failed to convert encoding: %w", err)
			}
		} else {
			text = word
		}

		words, err := kagomer.ParseToNode(text)
		if err != nil {
			return fmt.Errorf("mecaber: %w", err)
		}

		wordMap := analyzer.AnalyzeNodeToMapAdvanced(words)

		// 頻度順にソート
		type wordFreq struct {
			word string
			freq int
		}
		wordList := make([]wordFreq, 0, len(wordMap))
		for w, f := range wordMap {
			wordList = append(wordList, wordFreq{w, f})
		}
		sort.Slice(wordList, func(i, j int) bool {
			return wordList[i].freq > wordList[j].freq
		})

		// 上位20件を表示
		displayCount := 20
		if len(wordList) < displayCount {
			displayCount = len(wordList)
		}

		fmt.Fprintf(os.Stdout, "=== 抽出された単語（頻度順・上位%d件）===\n", displayCount)
		for i := 0; i < displayCount; i++ {
			fmt.Fprintf(os.Stdout, "%d. %s: %d回\n", i+1, wordList[i].word, wordList[i].freq)
		}
		fmt.Fprintf(os.Stdout, "合計: %d種類の単語\n", len(wordMap))
		fmt.Fprintf(os.Stdout, "================================\n")

		cloudimg.GenCloud(wordMap, outputFile)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

// convertToUTF8 は様々なエンコーディング（JIS、Shift_JIS、EUC-JP、UTF-8）からUTF-8に変換します
func convertToUTF8(data []byte) (string, error) {
	// まずUTF-8として試す
	if isValidUTF8(data) {
		return string(data), nil
	}

	// ISO-2022-JP (JIS) として試す
	decoder := japanese.ISO2022JP.NewDecoder()
	decoded, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err == nil && isValidUTF8(decoded) {
		return string(decoded), nil
	}

	// Shift_JIS として試す
	decoder = japanese.ShiftJIS.NewDecoder()
	decoded, err = io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err == nil && isValidUTF8(decoded) {
		return string(decoded), nil
	}

	// EUC-JP として試す
	decoder = japanese.EUCJP.NewDecoder()
	decoded, err = io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err == nil && isValidUTF8(decoded) {
		return string(decoded), nil
	}

	// どのエンコーディングでも失敗した場合、元のデータをそのまま返す
	return string(data), nil
}

// isValidUTF8 はバイト列が有効なUTF-8かどうかをチェックします
func isValidUTF8(data []byte) bool {
	// 文字列に変換して、元のバイト列と比較
	s := string(data)
	return len(s) > 0 && string([]byte(s)) == string(data)
}
