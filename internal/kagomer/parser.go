package kagomer

import (
	"fmt"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

const BOSEOS = "BOS/EOS"

// ex) text = "朝も夜も君に会いたい"
//
//	[朝 名詞,副詞可能,*,*,*,*,朝,アサ,アサ
//	 も 助詞,係助詞,*,*,*,*,も,モ,モ
//	 夜 名詞,副詞可能,*,*,*,*,夜,ヨル,ヨル
//	 も 助詞,係助詞,*,*,*,*,も,モ,モ
//	 君 名詞,代名詞,一般,*,*,*,君,キミ,キミ
//	 に 助詞,格助詞,一般,*,*,*,に,ニ,ニ
//	 会い 動詞,自立,*,*,五段・ワ行促音便,連用形,会う,アイ,アイ
//	 たい 助動詞,*,*,*,特殊・タイ,基本形,たい,タイ,タイ
//	]%
func ParseToNode(text string) ([]string, error) {
	// デフォルトのシステム辞書（IPA辞書）で初期化
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}

	tokens := t.Tokenize(text)

	rets := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if token.Surface == BOSEOS {
			continue
		}
		rets = append(rets, fmt.Sprintf("%s %s\n", token.Surface, token.Features()))
	}
	return rets, nil
}
