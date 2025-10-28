package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/bushiyama/wordcloud/internal/cloudimg"
	"github.com/bushiyama/wordcloud/internal/kagomer"
)

var (
	word     string
	filePath string
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

		wordMap := analyzeNodeToMapAdvanced(words)
		fmt.Fprintf(os.Stdout, "=== 抽出された単語（上位20件）===\n")
		count := 0
		for word, freq := range wordMap {
			if count >= 20 {
				break
			}
			fmt.Fprintf(os.Stdout, "%s: %d回\n", word, freq)
			count++
		}
		fmt.Fprintf(os.Stdout, "合計: %d種類の単語\n", len(wordMap))
		fmt.Fprintf(os.Stdout, "================================\n")

		cloudimg.GenCloud(wordMap)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

// ストップワード（除外する一般的な単語）
var stopWords = map[string]bool{
	"こと": true, "もの": true, "ため": true, "よう": true,
	"そう": true, "これ": true, "それ": true, "あれ": true,
	"どれ": true, "ここ": true, "そこ": true, "あそこ": true,
	"どこ": true, "とき": true, "時": true, "中": true,
	"人": true, "私": true, "僕": true, "俺": true, "あなた": true,
	"彼": true, "彼女": true, "方": true, "者": true, "上": true,
	"下": true, "前": true, "後": true, "間": true, "所": true,
	"場合": true, "今": true, "今日": true, "明日": true, "昨日": true,
	"ところ": true, "何": true, "の": true, "お": true, "ご": true,
	"さん": true, "くん": true, "ちゃん": true, "様": true,
}

// analyzeNodeToMapAdvanced は改善版の単語抽出関数
// 名詞の細分類でフィルタリングし、ストップワードを除外し、意味のある単語のみを抽出
func analyzeNodeToMapAdvanced(node []string) map[string]int {
	retMap := make(map[string]int)

	for _, v := range node {
		// フォーマット: "単語 [品詞 詳細1 詳細2...]"
		// 例: "メロス [名詞 一般 * * * * *]"
		parts := strings.SplitN(v, " ", 2)
		if len(parts) < 2 {
			continue
		}

		word := parts[0]

		// 1文字の単語を除外（助詞や記号が混ざるのを防ぐ）
		if len([]rune(word)) < 2 {
			continue
		}

		// ストップワードを除外
		if stopWords[word] {
			continue
		}

		// "[名詞 一般 ...]" から "[" と "]" を削除
		features := strings.Trim(parts[1], "[]")
		grammarParts := strings.Split(features, " ")

		if len(grammarParts) < 2 {
			continue
		}

		pos := grammarParts[0]       // 品詞（名詞、動詞など）
		posDetail := grammarParts[1] // 品詞細分類

		// 名詞の場合：有意義な名詞のみを抽出
		if pos == "名詞" {
			// 除外する名詞の細分類
			excludedTypes := map[string]bool{
				"非自立":  true, // 「こと」「もの」など
				"代名詞":  true, // 「これ」「それ」など
				"接尾":   true, // 接尾辞
				"数":    true, // 数詞
				"副詞可能": true, // 「今」「昨日」など副詞的な名詞
			}

			// 含める名詞の細分類
			if !excludedTypes[posDetail] && posDetail != "*" {
				retMap[word] += 1
			}
		}

		// 動詞の場合：自立語のみ（補助動詞を除外）
		if pos == "動詞" && posDetail == "自立" {
			// 動詞の基本形を取得（7番目の要素）
			if len(grammarParts) >= 7 && grammarParts[6] != "*" {
				baseForm := grammarParts[6]
				if len([]rune(baseForm)) >= 2 && !stopWords[baseForm] {
					retMap[baseForm] += 1
				}
			}
		}

		// 形容詞も追加（オプション）
		if pos == "形容詞" && posDetail == "自立" {
			if len(grammarParts) >= 7 && grammarParts[6] != "*" {
				baseForm := grammarParts[6]
				if len([]rune(baseForm)) >= 2 && !stopWords[baseForm] {
					retMap[baseForm] += 1
				}
			}
		}
	}

	return retMap
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
