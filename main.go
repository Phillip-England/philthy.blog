package main

import (
	"fmt"
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
			mdContent, err := vbf.LoadMarkdown("./static/docs/me.md")
			if err != nil {
				vbf.WriteString(w, "failed to load markdown content")
				return
			}
			mdContent = StyleMarkdownContent(mdContent)
			vbf.WriteHTML(w, Layout("Philthy", "Philthy", "Saying the things I'm afraid to say", r.URL.Path, mdContent))
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
			"HeaderText":  "Philthy",
			"SubText":     "Saying the things I'm afraid to say",
			"Content":     "",
			"CurrentPath": r.URL.Path,
		})
		if err != nil {
			vbf.WriteString(w, err.Error())
			return
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts/{postNumber}", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		postNumber := r.PathValue("postNumber")
		mdContent, err := vbf.LoadMarkdown(fmt.Sprintf("./static/docs/%s.md", postNumber))
		mdContent = StyleMarkdownContent(mdContent)
		if err != nil {
			vbf.WriteString(w, "failed to load markdown content")
			return
		}
		t, _ := vbf.GetContext(KEY_TEMPLATES, r).(*template.Template)
		err = vbf.ExecuteTemplate(w, t, "layout.html", map[string]interface{}{
			"Title":       "Posts",
			"HeaderText":  "Philthy",
			"SubText":     "Saying the things I'm afraid to say",
			"Content":     template.HTML(mdContent),
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

// takes plain markdown HTML and adds tailwind classes for styling
func StyleMarkdownContent(mdContent string) string {
	mdContent = strings.ReplaceAll(mdContent, "<p><em>", "<p class='text-xs italic'><em>")
	mdContent = strings.ReplaceAll(mdContent, "<h1>", "<h1 class='font-bold text-2xl pt-4'>")
	mdContent = strings.Replace(mdContent, "<h1 class='font-bold text-2xl pt-4'>", "<h1 class='font-bold text-2xl'>", 1)
	mdContent = strings.ReplaceAll(mdContent, "<h2>", "<h2 class='font-bold text-xl pt-4'>")
	mdContent = strings.Replace(mdContent, "<h2 class='font-bold text-xl pt-4'>", "<h2 class='font-bold text-xl'>", 1)
	mdContent = strings.ReplaceAll(mdContent, "<ol>", "<ol class='list-decimal list-inside pl-4 flex flex-col gap-2'>")
	mdContent = strings.ReplaceAll(mdContent, "<img", "<img class='pt-4 w-fit md:max-h-sm md:max-w-sm'")
	return mdContent
}

//==================================
// COMPONENTS
//==================================

func Layout(title string, headerText string, subText string, currentPath string, mdContent string) string {
	return fmt.Sprintf(`
		<html lang="en">
	    <head>
	        <meta charset="UTF-8" />
	        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	        <link rel="stylesheet" href="/static/css/output.css" />
	        <script src="/static/js/index.js"></script>
	        <script src="https://unpkg.com/hyperscript.org@0.9.12"></script>
	        <title>%s | Philthy</title>
	    </head>
	    <body class="h-screen">
	        <div id="root" class="min-h-screen">
	            <div class="h-[95px]"></div>
					%s%s
	            <main
	                class="flex-grow p-4 flex flex-col gap-4 md:pl-[300px]"
	                style="min-height: calc(100vh - 190px)"
	            >
	                <div class="flex flex-col md:pl-12 md:p-8 gap-4">
	                    %s
	                </div>
	            </main>
	            %s
	            <div class="h-[100px]"></div>
	        </div>
	    </body>
	</html>
	`, title, Header(headerText, subText), NavMenu(currentPath), mdContent, Overlay())
}

func Overlay() string {
	return `
		<div id="overlay" _="on click toggle .hidden on me then toggle .hidden on #icon-bars then toggle .hidden on #icon-x then toggle .hidden on #navmenu" class="h-full w-full bg-black opacity-50 fixed top-0 hidden z-30"></div>
	`
}

func Header(headerText string, subText string) string {
	return fmt.Sprintf(`
		<header class="flex flex-row justify-between p-4 border-b select-none z-40 bg-white fixed w-full top-0 h-[95px]">
    		<img src="/static/img/logo.svg" />
		    <div class="flex flex-col gap-2">
		        <h1 class="font-bold text-2xl">%s</h1>
		        <p class="text-sm">%s</p>
		    </div>
		    <div class="flex items-center">
		        <div id="icon-bars" _="on click toggle .hidden on me then toggle .hidden on #icon-x then toggle .hidden on #overlay then toggle .hidden on #navmenu" class="md:hidden">
		            <svg class="w-8 h-8" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
		                <path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M5 7h14M5 12h14M5 17h14" />
		            </svg>
		        </div>
		        <div id="icon-x" class="hidden md:hidden" _="on click toggle .hidden on me then toggle .hidden on #icon-bars then toggle .hidden on #overlay then toggle .hidden on #navmenu">
		            <svg class="w-8 h-8" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
		                <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18 17.94 6M18 18 6.06 6" />
		            </svg>
		        </div>
		    </div>
		</header>
	`, headerText, subText)
}

func NavMenu(currentPath string) string {
	return fmt.Sprintf(`
		<nav id="navmenu" class="flex top-[95px] fixed bg-white w-1/3 h-full border-r z-40 hidden md:flex md:w-[300px]">
		    <ul class="flex flex-col gap-2 p-2 w-full text-sm">
				%s%s
		    </ul>
		</nav>

	`, NavItem(currentPath, "/", "Home"), NavItem(currentPath, "/posts", "Posts"))
}

func NavItem(currentPath string, href string, text string) string {
	if currentPath == href {
		return fmt.Sprintf(`
		  <li class="rounded border w-full flex bg-gray-200">
		      <a class="w-full p-4" href="%s">%s</a>
		  </li>
		`, href, text)
	} else {
		return fmt.Sprintf(`
		  <li class="rounded border w-full flex bg-white">
		      <a class="w-full p-4" href="%s">%s</a>
		  </li>
		`, href, text)
	}
}
