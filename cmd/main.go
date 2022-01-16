package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bluele/mecab-golang"
	"github.com/bushiyama/wordcloud/internal/mecaber"
)

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

		ret, err := mecaber.ParseToNode(m, p)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, strings.Join(ret, ""))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("http.ListenAndServe:", err)
	}
}
