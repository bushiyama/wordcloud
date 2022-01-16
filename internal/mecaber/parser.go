package mecaber

import (
	"fmt"
	"strings"

	"github.com/bluele/mecab-golang"
	"github.com/bushiyama/wordcloud/internal/image"
)

const BOSEOS = "BOS/EOS"

func ParseToNode(m *mecab.MeCab, text string) ([]string, error) {
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
	image.GenCloud()
	return ret, nil
}
