package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/Phillip-England/vbf"
)

const KEY_TEMPLATES = "TEMPLATES"

func main() {

	//==================================
	// INIT
	//==================================

	// Define the global FuncMap with the custom "eq" function
	funcMap := template.FuncMap{
		"eq": func(a, b string) bool { return a == b },
	}

	// Parse templates and apply the FuncMap to all templates
	templates, err := ParseTemplatesWithFuncs("./templates", funcMap)
	if err != nil {
		panic(err)
	}

	mux, gCtx := vbf.VeryBestFramework()

	// Store the parsed templates in the global context
	vbf.SetGlobalContext(gCtx, KEY_TEMPLATES, templates)

	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	//==================================
	// ROUTES
	//==================================

	// Define the route with middleware and template execution
	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			t, _ := vbf.GetContext(KEY_TEMPLATES, r).(*template.Template)
			mdContent, err := vbf.LoadMarkdown("./static/docs/content.md")
			if err != nil {
				vbf.WriteString(w, "failed to load markdown content")
				return
			}
			styledMdContent := StyleMarkdownContent(mdContent)
			err = vbf.ExecuteTemplate(w, t, "layout.html", map[string]interface{}{
				"Title":       "Why am I Writing?",
				"HeaderText":  "Philthy Blog",
				"SubText":     "Saying the things I'm afraid to say",
				"Content":     template.HTML(styledMdContent),
				"CurrentPath": r.URL.Path,
			})
			if err != nil {
				vbf.WriteString(w, err.Error())
				return
			}
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		t, _ := vbf.GetContext(KEY_TEMPLATES, r).(*template.Template)
		err = vbf.ExecuteTemplate(w, t, "layout.html", map[string]interface{}{
			"Title":       "Posts",
			"HeaderText":  "Philthy Blog",
			"SubText":     "Saying the things I'm afraid to say",
			"Content":     "",
			"CurrentPath": r.URL.Path,
		})
		if err != nil {
			vbf.WriteString(w, err.Error())
			return
		}
	}, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}

// ParseTemplatesWithFuncs parses templates and applies a FuncMap to them
func ParseTemplatesWithFuncs(templateDir string, funcMap template.FuncMap) (*template.Template, error) {
	// Load and parse templates from the directory
	tmpl := template.New("").Funcs(funcMap) // Apply the FuncMap globally
	parsedTmpl, err := tmpl.ParseGlob(templateDir + "/*.html")
	if err != nil {
		return nil, err
	}
	return parsedTmpl, nil
}

func StyleMarkdownContent(mdContent string) string {
	// Apply styles to the Markdown content
	mdContent = strings.ReplaceAll(mdContent, "<p><em>", "<p class='text-xs italic'><em>")
	mdContent = strings.ReplaceAll(mdContent, "<h1>", "<h1 class='font-bold text-2xl pt-4'>")
	mdContent = strings.Replace(mdContent, "<h1 class='font-bold text-2xl pt-4'>", "<h1 class='font-bold text-2xl'>", 1)
	mdContent = strings.ReplaceAll(mdContent, "<h2>", "<h2 class='font-bold text-xl pt-4'>")
	mdContent = strings.Replace(mdContent, "<h2 class='font-bold text-xl pt-4'>", "<h2 class='font-bold text-xl'>", 1)
	mdContent = strings.ReplaceAll(mdContent, "<ol>", "<ol class='list-decimal list-inside pl-4 flex flex-col gap-2'>")

	return mdContent
}
