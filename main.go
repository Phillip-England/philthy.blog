package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/phillip-england/vbf"
)

const KeyTemplates = "KEYTEMPLATES"

func main() {

	mux, gCtx := vbf.VeryBestFramework()

	strEquals := func(input string, value string) bool {
		return input == value
	}

	funcMap := template.FuncMap{
		"strEquals": strEquals,
	}

	templates, err := vbf.ParseTemplates("./templates", funcMap)
	if err != nil {
		panic(err)
	}

	vbf.SetGlobalContext(gCtx, KeyTemplates, templates)
	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		templates, _ := vbf.GetContext(KeyTemplates, r).(*template.Template)
		mdContent, err := vbf.LoadMarkdown("./posts/index.md", "dracula")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		if r.URL.Path == "/" {
			vbf.ExecuteTemplate(w, templates, "root.html", map[string]interface{}{
				"Title":       "Welcome - Philthy Blog",
				"Content":     template.HTML(mdContent),
				"ReqPath":     r.URL.Path,
				"ArticleName": "philthy.blog",
				"SubText":     "Saying the things I'm afraid to say",
				"ImageSrc":    "./static/img/profile-sm.png",
			})
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		vbf.WriteHTML(w, `/posts`)
	}, vbf.MwLogger)

	vbf.AddRoute("GET /post/{postNumber}", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		postNumber := r.PathValue("postNumber")
		vbf.WriteHTML(w, `/post/`+postNumber)

	}, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}
