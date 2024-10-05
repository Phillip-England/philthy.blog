package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/Phillip-England/vbf"
)

const KEY_TEMPLATES = "TEMPLATES"

func main() {

	mux, gCtx := vbf.VeryBestFramework()

	templates, err := vbf.ParseTemplates("./templates")
	if err != nil {
		panic(err)
	}

	vbf.SetGlobalContext(gCtx, KEY_TEMPLATES, templates)
	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			templates, _ := vbf.GetContext(KEY_TEMPLATES, r).(*template.Template)
			mdContent, err := vbf.LoadMarkdown("./static/docs/content.md")
			if err != nil {
				vbf.WriteString(w, "failed to load markdown content")
				return
			}
			styledMdContent := StyleMarkdownContent(mdContent)
			vbf.ExecuteTemplate(w, templates, "layout.html", map[string]interface{}{
				"Title":      "Why am I Writing?",
				"HeaderText": "Philthy Blog",
				"SubText":    "Saying the things I'm afraid to say",
				"Content":    template.HTML(styledMdContent),
			})
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}

func StyleMarkdownContent(mdContent string) string {

	// styling all italics
	mdContent = strings.ReplaceAll(mdContent, "<p><em>", "<p class='text-xs italic'><em>")

	// styling all <h1> elements
	mdContent = strings.ReplaceAll(mdContent, "<h1>", "<h1 class='font-bold text-2xl pt-4'>")

	// styling all <h2> elements, but the first <h2> does not get top padding
	mdContent = strings.ReplaceAll(mdContent, "<h2>", "<h2 class='font-bold text-xl pt-4'>")
	mdContent = strings.Replace(mdContent, "<h2 class='font-bold text-xl pt-4'>", "<h2 class='font-bold text-xl'>", 1)

	// styling all <ol>
	mdContent = strings.ReplaceAll(mdContent, "<ol>", "<ol class='list-decimal list-inside pl-4 flex flex-col gap-2'>")

	return mdContent
}
