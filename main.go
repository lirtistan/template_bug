package main

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "html/template"
	"net/http"
	"os"
	"text/template"
)

//go:embed *.html
var assets embed.FS

type Map = map[string]string

func Deliver(w http.ResponseWriter, path string, data ...any) {
	content, err := assets.ReadFile(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	} else {
		tpl, err2 := template.New(path).Option("missingkey=zero").Parse(string(content))
		if err2 != nil {
			fmt.Fprint(os.Stderr, err2)
		}
		var arg any
		if len(data) > 0 {
			arg = data[0]
		}
		if ok := tpl.Execute(w, arg); ok != nil {
			fmt.Fprint(os.Stderr, ok)
		}
	}
}

func main() {
	router := chi.NewMux()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Deliver(w, "index.html", Map{})
	})
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
