package main

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "html/template"
	"log"
	"net/http"
	"runtime"
	"text/template"
)

// see issue https://github.com/golang/go/issues/70681 for further explanation
// if you want to see the difference, just replace the import of text/template with html/template,
// no code switch or tags needed

//go:embed *.html *.css
var assets embed.FS

// type Map = map[string]any    // <- this is how my application defines a map

type MapOfStr = map[string]string // <- this map works as expected with text/templates
type MapOfAny = map[string]any    // <- this map behaves different between HTML and Text templates

// Deliver writes a template result on the http.ResponseWriter, if no arg is given it handles the content as is
// This function is the same as in my application
func Deliver(w http.ResponseWriter, path string, arg ...any) (err error) {
	var content []byte
	if content, err = assets.ReadFile(path); err != nil {
		return
	}

	if len(arg) > 0 {
		// with arg we deliver a template result
		var tpl *template.Template
		tpl, err = template.New(path).Option("missingkey=zero").Parse(string(content))
		if err != nil {
			return
		}

		if err = tpl.Execute(w, arg[0]); err != nil {
			return
		}
	} else {
		// without arg it's just an ordinary file, so we deliver it as is
		_, err = w.Write(content)
	}

	return
}

// TestTemplateWith returns a http.HandlerFunc see main.main below
func TestTemplateWith[T any](mapOf T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		if err := Deliver(w, "test.html", mapOf); err != nil {
			log.Fatal(err)
		}
	}
}

// now it's time for or test
func main() {
	router := chi.NewMux()
	router.Get("/", Index())
	router.Get("/map-of-str", TestTemplateWith(MapOfStr{}))
	router.Get("/map-of-any", TestTemplateWith(MapOfAny{}))
	router.Get("/common.css", CSS("common.css"))

	fmt.Println("listening on http://0.0.0.0:8080 ....")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

// ====================================================================================================
// DON'T CARE MUCH ABOUT STUFF BELOW THIS COMMENT. NOT NECESSARY TO UNDERSTAND THE ISS<UE

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := Deliver(w, "index.html", map[string]string{
			"Version": runtime.Version(),
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func CSS(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		if err := Deliver(w, path); err != nil {
			log.Fatal(err)
		}
	}
}
