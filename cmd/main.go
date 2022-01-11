package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bluele/mecab-golang"
)

const BOSEOS = "BOS/EOS"

func parseToNode(m *mecab.MeCab, text string) ([]string, error) {
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
			fmt.Printf("%s %s\n", node.Surface(), node.Feature())
			ret = append(ret, fmt.Sprintf("%s %s\n", node.Surface(), node.Feature()))
		}
		if node.Next() != nil {
			break
		}
	}
	return ret, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		var p string
		ary, ok := r.Form["hoge"]
		if !ok {
			panic(fmt.Errorf("not hoge"))
		}
		p = ary[0]

		m, err := mecab.New("-Owakati")
		if err != nil {
			panic(err)
		}
		defer m.Destroy()

		ret, err := parseToNode(m, p)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, strings.Join(ret, ""))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("http.ListenAndServe:", err)
	}
}
