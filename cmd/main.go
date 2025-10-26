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

		wordMap := analyzeNodeToMap(words, "名詞")

		cloudimg.GenCloud(wordMap)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

/*
AnalyzeNodeToMap

		gen input param for wordcloud.

		ex)
		node = [
			朝 名詞,副詞可能,*,*,*,*,朝,アサ,アサ
			も 助詞,係助詞,*,*,*,*,も,モ,モ
			夜 名詞,副詞可能,*,*,*,*,夜,ヨル,ヨル
			も 助詞,係助詞,*,*,*,*,も,モ,モ
			君 名詞,代名詞,一般,*,*,*,君,キミ,キミ
			に 助詞,格助詞,一般,*,*,*,に,ニ,ニ
			会い 動詞,自立,*,*,五段・ワ行促音便,連用形,会う,アイ,アイ
			たい 助動詞,*,*,*,特殊・タイ,基本形,たい,タイ,タイ
			と 助詞,格助詞,引用,*,*,*,と,ト,ト
			朝 名詞,副詞可能,*,*,*,*,朝,アサ,アサ
			思う 動詞,自立,*,*,五段・ワ行促音便,基本形,思う,オモウ,オモウ
	  	]
		tgtGrammar = "名詞"

		=> return
		map[
	    	朝:2
			夜:1
			君:1
		]
*/
func analyzeNodeToMap(node []string, tgtGrammar string) map[string]int {
	retMap := make(map[string]int)
	for _, v := range node {
		// フォーマット: "単語 [品詞 詳細1 詳細2...]"
		// 例: "メロス [名詞 一般 * * * * *]"
		parts := strings.SplitN(v, " ", 2)
		if len(parts) < 2 {
			continue
		}

		word := parts[0]
		// "[名詞 一般 ...]" から "[" と "]" を削除して "名詞 一般 ..." にする
		features := strings.Trim(parts[1], "[]")
		grammarParts := strings.Split(features, " ")

		if len(grammarParts) > 0 && grammarParts[0] == tgtGrammar {
			retMap[word] += 1
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
