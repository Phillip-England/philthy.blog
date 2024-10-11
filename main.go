package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Phillip-England/vbf"
)

func main() {

	//==================================
	// INIT
	//==================================

	mux, gCtx := vbf.VeryBestFramework()
	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	//==================================
	// ROUTES
	//==================================

	// Define the route with middleware and template execution
	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			mdContent, err := vbf.LoadMarkdown("./static/docs/home.md")
			if err != nil {
				vbf.WriteString(w, "failed to load markdown content")
				return
			}
			mdContent = StyleMarkdownContent(mdContent)
			vbf.WriteHTML(w, Layout("Philthy", "Philthy", "Saying the things I'm afraid to say", r.URL.Path, mdContent))
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		mdContent, err := vbf.LoadMarkdown("./static/docs/home.md")
		if err != nil {
			vbf.WriteString(w, "failed to load markdown content")
			return
		}
		mdContent = StyleMarkdownContent(mdContent)
		vbf.WriteHTML(w, Layout("Philthy", "Philthy", "Saying the things I'm afraid to say", r.URL.Path, mdContent))
	}, vbf.MwLogger)

	// vbf.AddRoute("GET /posts/{postNumber}", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
	// 	postNumber := r.PathValue("postNumber")
	// 	mdContent, err := vbf.LoadMarkdown(fmt.Sprintf("./static/docs/%s.md", postNumber))
	// 	mdContent = StyleMarkdownContent(mdContent)
	// 	if err != nil {
	// 		vbf.WriteString(w, "failed to load markdown content")
	// 		return
	// 	}
	// 	err = vbf.ExecuteTemplate(w, t, "layout.html", map[string]interface{}{
	// 		"Title":       "Posts",
	// 		"HeaderText":  "Philthy",
	// 		"SubText":     "Saying the things I'm afraid to say",
	// 		"Content":     template.HTML(mdContent),
	// 		"CurrentPath": r.URL.Path,
	// 	})
	// 	if err != nil {
	// 		vbf.WriteString(w, err.Error())
	// 		return
	// 	}
	// }, vbf.MwLogger)

	err := vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}
}

// takes plain markdown HTML and adds tailwind classes for styling
func StyleMarkdownContent(mdContent string) string {
	mdContent = strings.Replace(mdContent, "<p><em>", "<p class='text-xs italic pb-8'><em>", 1) // the first <em> is for the date the article was written
	mdContent = strings.ReplaceAll(mdContent, "<p><em>", "<p class='italic'><em>")              // all the other <em>'s in the markdown file
	mdContent = strings.ReplaceAll(mdContent, "<h1>", "<h1 class='font-bold text-3xl pt-4 pb-1'>")
	mdContent = strings.ReplaceAll(mdContent, "<h2>", "<h2 class='font-bold text-2xl pt-8 pb-8'>")
	mdContent = strings.ReplaceAll(mdContent, "<ol>", "<ol class='list-decimal list-inside pl-4 flex flex-col gap-4 pb-8 pt-4'>")
	mdContent = strings.ReplaceAll(mdContent, "<p><img", "<p class='flex mb-8 items-center'><img class='max-w-[200px]'")
	mdContent = strings.ReplaceAll(mdContent, "<p>", "<p class='mb-4'>")
	mdContent = strings.ReplaceAll(mdContent, "<blockquote", "<blockquote class='italic'")

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
					%s
     			<div class="h-[65px]"></div>
					%s%s
	            <main
	                class="flex-grow p-4 flex flex-col md:pl-[300px] lg:pr-[300px] mt-4"
	                style="min-height: calc(100vh - 190px)"
	            >
	                <div class="flex flex-col md:pl-12 md:p-8">
	                    %s
	                </div>
	            </main>
	            %s
	            <div class="h-[100px]"></div>
	        </div>
	    </body>
	</html>
	`, title, Header(headerText, subText), SocialMediaBenner(), NavMenu(currentPath), mdContent, Overlay())
}

func Overlay() string {
	return `
		<div id="overlay" _="on click toggle .hidden on me then toggle .hidden on #icon-bars then toggle .hidden on #icon-x then toggle .hidden on #navmenu" class="h-full w-full bg-black opacity-50 fixed top-0 hidden z-30"></div>
	`
}

func Header(headerText string, subText string) string {
	return fmt.Sprintf(`
		<header class="flex flex-row justify-between p-4 border-b select-none z-40 bg-white fixed w-full top-0 h-[95px]">
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
		<nav id="navmenu" class="flex top-[95px] fixed bg-white w-2/3 h-full border-r z-40 hidden md:flex md:w-[300px]">
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

func SocialMediaBenner() string {
	return `
		<div class="bg-black text-white md:pl-[300px] flex flex-row fixed top-[95px] h-[65px] w-full">
			<a href='https://www.youtube.com/@phillip-england' target="_blank">
				<div class='p-4'>
					<svg class="w-8 h-8" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
					  <path fill-rule="evenodd" d="M21.7 8.037a4.26 4.26 0 0 0-.789-1.964 2.84 2.84 0 0 0-1.984-.839c-2.767-.2-6.926-.2-6.926-.2s-4.157 0-6.928.2a2.836 2.836 0 0 0-1.983.839 4.225 4.225 0 0 0-.79 1.965 30.146 30.146 0 0 0-.2 3.206v1.5a30.12 30.12 0 0 0 .2 3.206c.094.712.364 1.39.784 1.972.604.536 1.38.837 2.187.848 1.583.151 6.731.2 6.731.2s4.161 0 6.928-.2a2.844 2.844 0 0 0 1.985-.84 4.27 4.27 0 0 0 .787-1.965 30.12 30.12 0 0 0 .2-3.206v-1.516a30.672 30.672 0 0 0-.202-3.206Zm-11.692 6.554v-5.62l5.4 2.819-5.4 2.801Z" clip-rule="evenodd"/>
					</svg>
				</div>
			</a>

		</div>
	`
}
