package mecaber

import (
	"fmt"
	"strings"

	"github.com/bluele/mecab-golang"
)

const BOSEOS = "BOS/EOS"

// ex) text = "朝も夜も君に会いたい"
//   [朝 名詞,副詞可能,*,*,*,*,朝,アサ,アサ
//    も 助詞,係助詞,*,*,*,*,も,モ,モ
//    夜 名詞,副詞可能,*,*,*,*,夜,ヨル,ヨル
//    も 助詞,係助詞,*,*,*,*,も,モ,モ
//    君 名詞,代名詞,一般,*,*,*,君,キミ,キミ
//    に 助詞,格助詞,一般,*,*,*,に,ニ,ニ
//    会い 動詞,自立,*,*,五段・ワ行促音便,連用形,会う,アイ,アイ
//    たい 助動詞,*,*,*,特殊・タイ,基本形,たい,タイ,タイ
//   ]%
func ParseToNode(text string) ([]string, error) {
	m, err := mecab.New("-Owakati")
	if err != nil {
		panic(err)
	}
	defer m.Destroy()

	ret := []string{}
	tg, err := m.NewTagger()
	if err != nil {
		return ret, err
	}
	defer tg.Destroy()

	lt, err := m.NewLattice(text)
	if err != nil {
		return ret, err
	}
	defer lt.Destroy()

	node := tg.ParseToNode(lt)
	for {
		features := strings.Split(node.Feature(), ",")
		if features[0] != BOSEOS {
			ret = append(ret, fmt.Sprintf("%s %s\n", node.Surface(), node.Feature()))

		}
		if node.Next() != nil {
			break
		}
	}
	return ret, nil
}
