package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	cli "github.com/urfave/cli/v2"

	"github.com/bushiyama/wordcloud/internal/analyzer"
	"github.com/bushiyama/wordcloud/internal/cloudimg"
	"github.com/bushiyama/wordcloud/internal/converter"
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
			text, err = converter.ConvertToUTF8(data)
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
