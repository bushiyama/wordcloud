package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bushiyama/wordcloud/internal/cloudimg"
	"github.com/bushiyama/wordcloud/internal/mecaber"
	cli "github.com/urfave/cli/v2"
)

var (
	word string
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "w",
			Usage:       "required word",
			Destination: &word,
			Required:    true,
		},
	}

	app.Action = func(c *cli.Context) error {
		words, err := mecaber.ParseToNode(word)
		if err != nil {
			return fmt.Errorf("mecaber: %w", err)
		}

		cloudimg.GenCloud(analyzeNodeToMap(words, "名詞"))

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
		w := strings.Split(v, " ")
		g := strings.Split(w[1], ",")
		if g[0] == tgtGrammar {
			retMap[w[0]] += 1
		}
	}
	return retMap
}
