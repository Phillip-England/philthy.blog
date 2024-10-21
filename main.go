package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Phillip-England/vbf"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/exp/rand"
)

func main() {

	executionCmd := GetArg(1)
	if executionCmd == "" {
		executionCmd = "--serve"
	}

	if executionCmd == "--serve" {
		serveApp()
		return
	}

	if executionCmd == "--build" {
		var endpoints []PhilthyEndpoint
		outputURL := "https://philthy.blog"
		serverURL := "http://localhost:8080"
		outputDir := "./output"
		postStartIndex := 1
		postEndIndex := 5
		endpoints = append(endpoints, PhilthyEndpoint{
			ReqPath:           serverURL + "/",
			HrefSearchPattern: "/",
			OutputHref:        outputURL + "/index.html",
			SaveTo:            "./output/index.html",
			HTMLContent:       "",
		})
		endpoints = append(endpoints, PhilthyEndpoint{
			ReqPath:           serverURL + "/posts",
			HrefSearchPattern: "/posts",
			OutputHref:        outputURL + "/posts.html",
			SaveTo:            "./output/posts.html",
			HTMLContent:       "",
		})
		for i := postStartIndex; i < postEndIndex; i++ {
			endpoints = append(endpoints, PhilthyEndpoint{
				ReqPath:           fmt.Sprintf("%s/post/%d", serverURL, i),
				HrefSearchPattern: fmt.Sprintf("/post/%d", i),
				OutputHref:        fmt.Sprintf("%s/%d.html", outputURL, i),
				SaveTo:            fmt.Sprintf("./output/%d.html", i),
				HTMLContent:       "",
			})
		}
		err := os.RemoveAll(outputDir)
		if err != nil {
			fmt.Printf("Error removing directory: %v\n", err)
		}
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = copyDir("./static", "./output/static")
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(endpoints); i++ {
			endpoint := endpoints[i]

			// Make the GET request
			resp, err := http.Get(endpoint.ReqPath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
			htmlContent := string(body)
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
			if err != nil {
				panic(err)
			}
			doc.Find("a").Each(func(index int, item *goquery.Selection) {
				href, _ := item.Attr("href")
				for i2 := 0; i2 < len(endpoints); i2++ {
					endpoint2 := endpoints[i2]
					if href == endpoint2.HrefSearchPattern {
						item.SetAttr("href", endpoint2.OutputHref)
					}
				}
			})
			doc.Find("link").Each(func(index int, item *goquery.Selection) {
				href, _ := item.Attr("href")
				if href == "" {
					return
				}
				if strings.Contains(href, "https://") {
					return
				}
				item.SetAttr("href", outputURL+href)
			})
			doc.Find("img").Each(func(index int, item *goquery.Selection) {
				src, _ := item.Attr("src")
				if src == "" {
					return
				}
				if strings.Contains(src, "https://") {
					return
				}
				item.SetAttr("src", outputURL+src)
			})
			doc.Find("script").Each(func(index int, item *goquery.Selection) {
				src, _ := item.Attr("src")
				if src == "" {
					return
				}
				if strings.Contains(src, "https://") {
					return
				}
				item.SetAttr("src", outputURL+src)
			})
			outputHTML, err := doc.Html()
			if err != nil {
				panic(err)
			}
			file, err := os.OpenFile(endpoint.SaveTo, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			_, err = file.WriteString(outputHTML)
			if err != nil {
				panic(err)
			}
		}
		return
	}

	fmt.Println("command: " + executionCmd + " is not a valid command")

}

func serveApp() {

	mux, gCtx := vbf.VeryBestFramework()
	vbf.HandleStaticFiles(mux)
	vbf.HandleFavicon(mux)

	err := GeneratePostsPage()
	if err != nil {
		panic(err)
	}

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			question := getRandomChristianQuestion()
			vbf.WriteHTML(w, Layout("Home", "philthy.blog", "Saying the things I'm afraid to say", r.URL.Path, `
				<section class="flex-grow flex flex-col md:pl-[300px] relative top-[65px] bg-white dark:bg-black text-black dark:text-white" style="min-height: calc(100vh - 160px)">
					<h2 class='text-2xl md:text-3xl md:w-[550px] p-4 md:p-12'>`+question+`</h2>
				</section>
			`))
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		article, err := MarkdownArticle("./static/docs/posts.md")
		if err != nil {
			vbf.WriteString(w, "failed to load .md content")
		}
		vbf.WriteHTML(w, Layout("Posts", "philthy.blog", "Saying the things I'm afraid to say", r.URL.Path, article))
	}, vbf.MwLogger)

	vbf.AddRoute("GET /post/{postNumber}", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		postNumber := r.PathValue("postNumber")
		article, err := MarkdownArticle(fmt.Sprintf("./static/docs/%s.md", postNumber))
		if err != nil {
			vbf.WriteString(w, "failed to load markdown content")
			return
		}
		mdContentTitle, err := ExtractH1Text(article)
		if err != nil {
			vbf.WriteString(w, "failed to locate title in markdown content")
		}
		vbf.WriteHTML(w, Layout(mdContentTitle, "philthy.blog", "Saying the things I'm afraid to say", r.URL.Path, article))
	}, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}

// represnts a markdown file which will serve as a blog post for the site
type PhilthyMarkdownFile struct {
	Path       string
	FileName   string
	PostNumber string
	Title      string
}

// represents a static endpoint which will be used to generate static files
type PhilthyEndpoint struct {
	ReqPath           string // where a http request will be made to (assumes the server is running)
	HrefSearchPattern string // will search the DOM for this pattern in <a> tags, and will replace them with OutputHref
	OutputHref        string // see HrefSearchPattern
	SaveTo            string // where the final html file will save on disk
	HTMLContent       string // the HTML to be saved on disk
}

// Copy a directory from src to dst
func copyDir(src string, dst string) error {
	// Walk through the source directory
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the new path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(dst, relPath)

		// Handle directories
		if info.IsDir() {
			// Create the destination directory
			if err := os.MkdirAll(newPath, info.Mode()); err != nil {
				return err
			}
		} else {
			// Handle files
			if err := copyFile(path, newPath); err != nil {
				return err
			}
		}

		return nil
	})
}

// Copy a single file from src to dst
func copyFile(src string, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the file contents
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy the file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// will search an html string and extract the text of the first <h1> in the string
func ExtractH1Text(mdContent string) (string, error) {
	indexOfH1 := strings.Index(mdContent, "<h1") + 1
	if indexOfH1 == -1 {
		return "", fmt.Errorf("markdown content did not contain a valid title")
	}
	mdContentTitle := ""
	collectText := false
	for i := indexOfH1; i < len(mdContent); i++ {
		char := string(mdContent[i])
		if char == ">" {
			collectText = true
			continue
		}
		if char == "<" {
			break
		}
		if collectText {
			mdContentTitle = mdContentTitle + string(char)
		}
	}
	return mdContentTitle, nil
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
	mdContent = strings.ReplaceAll(mdContent, "<blockquote", "<blockquote class='italic pl-4'")
	mdContent = strings.ReplaceAll(mdContent, "<a", "<a class='underline text-blue-500 visited:text-purple-500'")
	return mdContent
}

func GetPhilthyMarkdownFiles() ([]PhilthyMarkdownFile, error) {
	var philthyMarkdownFiles []PhilthyMarkdownFile
	err := filepath.Walk("./static/docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		pathParts := strings.Split(path, "/")
		mdFileName := pathParts[len(pathParts)-1]
		mdFileNameParts := strings.Split(mdFileName, ".")
		firstMdFileNamePart := mdFileNameParts[0]
		if _, err := strconv.Atoi(firstMdFileNamePart); err != nil {
			// Skip files that don't start with a number
		} else {
			philthyMdFile := PhilthyMarkdownFile{
				Path:       path,
				FileName:   mdFileName,
				PostNumber: firstMdFileNamePart,
			}
			philthyMarkdownFiles = append(philthyMarkdownFiles, philthyMdFile)
		}
		return nil
	})
	if err != nil {
		return philthyMarkdownFiles, err
	}
	return philthyMarkdownFiles, nil
}

func GeneratePostsPage() error {

	// Delete the output file to ensure a clean rewrite
	err := os.Remove("./static/docs/posts.md")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Open the output file for writing (it will be recreated)
	outputMdFile, err := os.OpenFile("./static/docs/posts.md", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outputMdFile.Close()

	// writing the header of the output file
	outputMdFile.Write([]byte("# All Posts\n"))

	// Collecting all valid .md files for the blog
	var philthyMarkdownFiles []PhilthyMarkdownFile
	filepath.Walk("./static/docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		pathParts := strings.Split(path, "/")
		mdFileName := pathParts[len(pathParts)-1]
		mdFileNameParts := strings.Split(mdFileName, ".")
		firstMdFileNamePart := mdFileNameParts[0]
		if _, err := strconv.Atoi(firstMdFileNamePart); err != nil {
			// Skip files that don't start with a number
		} else {
			philthyMdFile := PhilthyMarkdownFile{
				Path:       path,
				FileName:   mdFileName,
				PostNumber: firstMdFileNamePart,
			}
			philthyMarkdownFiles = append(philthyMarkdownFiles, philthyMdFile)
		}
		return nil
	})

	// extracting the titles out of each markdown file
	for i := 0; i < len(philthyMarkdownFiles); i++ {
		currentMdFile := philthyMarkdownFiles[i]
		mdFile, err := os.Open(currentMdFile.Path)
		if err != nil {
			return err
		}
		defer mdFile.Close()
		scanner := bufio.NewScanner(mdFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Count(line, `#`) == 1 {
				lineParts := strings.Split(line, " ")
				mdFileTitleSlice := lineParts[1:]
				mdFileTitle := strings.Join(mdFileTitleSlice, " ")
				currentMdFile.Title = mdFileTitle
				outputMdFile.Write([]byte(fmt.Sprintf(`%s. [%s](/post/%s)`, currentMdFile.PostNumber, currentMdFile.Title, currentMdFile.PostNumber) + "\n"))
			}
		}
	}

	// writing out the "Recent Posts" section
	// for i := 0; i < len(philthyMarkdownFiles); i++ {
	// 	currentMdFile := philthyMarkdownFiles[i]
	// 	fmt.Println(currentMdFile)
	// }

	return nil
}

// getRandomChristianQuestion returns a random thought-provoking question that might challenge beliefs about Christianity.
func getRandomChristianQuestion() string {
	// List of thought-provoking questions that could provoke doubt or critical thinking about Christianity
	questions := []string{
		"Why does an all-powerful, loving God allow so much suffering and evil in the world?",
		"If God is omniscient and knows everything, why create people He knows will suffer eternally?",
		"Why are there so many conflicting interpretations of the Bible if it's divinely inspired?",
		"How do we reconcile the Bible's teachings on morality with the actions of Christians throughout history?",
		"Why do miracles seem to be rare or absent in the modern world compared to biblical times?",
		"How can a loving God condemn people to eternal punishment for finite actions in their lives?",
		"Why do prayers often seem unanswered, even for things that appear to be good or just?",
		"If God is unchanging, why does His behavior appear to differ between the Old and New Testaments?",
		"How do we explain the existence of other religions with similar moral teachings and claims of divine truth?",
		"What if our faith has made promises it cannot keep?",
	}

	// Seed the random number generator to ensure different results on each run
	rand.Seed(uint64(time.Now().UnixNano()))

	// Select a random index from the questions list
	randomIndex := rand.Intn(len(questions))
	return questions[randomIndex]
}

// GetArg returns the command-line argument at the specified position.
// If the argument does not exist, it returns an empty string.
func GetArg(pos int) string {
	if pos < len(os.Args) {
		return os.Args[pos]
	}
	return ""
}

func Layout(title string, headerText string, subText string, currentPath string, contentComponents ...string) string {
	return fmt.Sprintf(`
		<html lang="en">
	    <head>
	        <meta charset="UTF-8" />
	        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	        <link rel="stylesheet" href="/static/css/output.css" />
	        <script src="/static/js/index.js"></script>
	        <script src="https://unpkg.com/hyperscript.org@0.9.12"></script>
			%s
	        <title>%s | Philthy</title>
	    </head>
	    <body class="h-screen bg-white dark:bg-black">
	        <div id="root" class="min-h-screen">
	            <div class="h-[95px]"></div>
				%s%s%s%s%s
	        </div>
	    </body>
	</html>
	`, ThemeScript(), title, Header(headerText, subText), SocialMediaBenner(), NavMenu(currentPath), strings.Join(contentComponents, ""), Overlay())
}

func ThemeScript() string {
	return "<script>" +
		"(function () {" +
		"  const getCookie = (name) => {" +
		"    const value = '; ' + document.cookie;" +
		"    const parts = value.split('; ' + name + '=');" +
		"    if (parts.length === 2) return parts.pop().split(';').shift();" +
		"  };" +
		"" +
		"  const theme = getCookie('theme');" +
		"  if (theme === 'light') {" +
		"    document.documentElement.classList.remove('dark');" +
		"  } else if (theme === 'dark') {" +
		"    document.documentElement.classList.add('dark');" +
		"  } else if (!theme && window.matchMedia('(prefers-color-scheme: dark)').matches) {" +
		"    document.documentElement.classList.add('dark');" +
		"  }" +
		"})();" +
		"</script>"
}

func MarkdownArticle(mdFilePath string) (string, error) {
	mdContent, err := vbf.LoadMarkdown(mdFilePath)
	if err != nil {
		return "", err
	}
	mdContent = StyleMarkdownContent(mdContent)
	return Article(mdContent), nil
}

func Article(children ...string) string {
	childrenHTML := strings.Join(children, "")
	return fmt.Sprintf(`
		<article class="flex-grow p-4 flex flex-col md:pl-[300px] lg:pr-[300px] relative top-[65px] bg-white dark:bg-black text-black dark:text-white" style="min-height: calc(100vh - 160px)">
		    <div class="flex flex-col md:pl-12 md:p-8">
		        %s
		    </div>
		</article>
	`, childrenHTML)
}

func Overlay() string {
	return `
		<div id="overlay" _="on click toggle .hidden on me then toggle .hidden on #icon-bars then toggle .hidden on #icon-x then toggle .hidden on #navmenu" class="h-full w-full bg-black opacity-50 fixed top-0 hidden z-30 md:hidden"></div>
	`
}

func Header(headerText string, subText string) string {
	return fmt.Sprintf(`
		<header class="flex flex-row justify-between p-4 dark:border-gray-800 border-b select-none z-40 bg-white fixed w-full top-0 h-[95px] dark:bg-black dark:text-white">
			<a href='/' class='flex dark:hidden'>
				<img src='/static/img/logo-black.svg' width="50" />
			</a>
			<a href='/' class='dark:flex hidden'>
				<img src='/static/img/logo-white.svg' width="50" />
			</a>
		    <div class="flex flex-col gap-2 md:items-end">
		        <h1 class="font-bold text-2xl">%s</h1>
		        <p class="text-sm">%s</p>
		    </div>
		    <div class="flex items-center md:hidden">
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
		<nav id="navmenu" class="flex top-[95px] fixed bg-white w-2/3 h-full border-r dark:border-gray-800 z-40 hidden md:flex md:w-[300px] bg-white dark:bg-black text-black dark:text-white">
		    <ul class="flex flex-col gap-2 p-2 w-full text-sm">
				%s%s
		    </ul>
		</nav>

	`, NavItem(currentPath, "/", "Home"), NavItem(currentPath, "/posts", "Posts"))
}

func NavItem(currentPath string, href string, text string) string {
	if currentPath == href {
		return fmt.Sprintf(`
		  <li class="rounded border dark:border-gray-800 w-full flex bg-gray-200 text-black">
		      <a class="w-full p-4" href="%s">%s</a>
		  </li>
		`, href, text)
	} else {
		return fmt.Sprintf(`
		  <li class="rounded border dark:border-gray-800 w-full flex bg-white dark:bg-black text-black dark:text-white">
		      <a class="w-full p-4" href="%s">%s</a>
		  </li>
		`, href, text)
	}
}

func SocialMediaBenner() string {
	return `
		<div class="text-black bg-white dark:bg-black dark:text-white border-b dark:border-gray-800 md:pl-[300px] flex flex-row fixed top-[95px] h-[65px] w-full z-30 justify-between">
			<div class='flex flex-row gap-4'>
				<a href='https://www.youtube.com/@phillip-england' target="_blank">
					<div class='p-4'>
						<svg class="w-8 h-8" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
						  <path fill-rule="evenodd" d="M21.7 8.037a4.26 4.26 0 0 0-.789-1.964 2.84 2.84 0 0 0-1.984-.839c-2.767-.2-6.926-.2-6.926-.2s-4.157 0-6.928.2a2.836 2.836 0 0 0-1.983.839 4.225 4.225 0 0 0-.79 1.965 30.146 30.146 0 0 0-.2 3.206v1.5a30.12 30.12 0 0 0 .2 3.206c.094.712.364 1.39.784 1.972.604.536 1.38.837 2.187.848 1.583.151 6.731.2 6.731.2s4.161 0 6.928-.2a2.844 2.844 0 0 0 1.985-.84 4.27 4.27 0 0 0 .787-1.965 30.12 30.12 0 0 0 .2-3.206v-1.516a30.672 30.672 0 0 0-.202-3.206Zm-11.692 6.554v-5.62l5.4 2.819-5.4 2.801Z" clip-rule="evenodd"/>
						</svg>
					</div>
				</a>
			</div>
			<div>
				<div class='p-4 flex flex-row gap-4'>
					<svg _='on click toggle .dark on <html/> then set cookies["theme"] to {value: "dark", path: "/"}' id='moon-icon' class="w-8 h-8 cursor-pointer flex dark:hidden"  aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
						<path fill-rule="evenodd" d="M11.675 2.015a.998.998 0 0 0-.403.011C6.09 2.4 2 6.722 2 12c0 5.523 4.477 10 10 10 4.356 0 8.058-2.784 9.43-6.667a1 1 0 0 0-1.02-1.33c-.08.006-.105.005-.127.005h-.001l-.028-.002A5.227 5.227 0 0 0 20 14a8 8 0 0 1-8-8c0-.952.121-1.752.404-2.558a.996.996 0 0 0 .096-.428V3a1 1 0 0 0-.825-.985Z" clip-rule="evenodd"/>
					</svg>
					<svg _='on click toggle .dark on <html/> then set cookies["theme"] to {value: "light", path: "/"}' id='sun-icon' class="w-8 h-8 hidden cursor-pointer dark:flex" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
						<path fill-rule="evenodd" d="M13 3a1 1 0 1 0-2 0v2a1 1 0 1 0 2 0V3ZM6.343 4.929A1 1 0 0 0 4.93 6.343l1.414 1.414a1 1 0 0 0 1.414-1.414L6.343 4.929Zm12.728 1.414a1 1 0 0 0-1.414-1.414l-1.414 1.414a1 1 0 0 0 1.414 1.414l1.414-1.414ZM12 7a5 5 0 1 0 0 10 5 5 0 0 0 0-10Zm-9 4a1 1 0 1 0 0 2h2a1 1 0 1 0 0-2H3Zm16 0a1 1 0 1 0 0 2h2a1 1 0 1 0 0-2h-2ZM7.757 17.657a1 1 0 1 0-1.414-1.414l-1.414 1.414a1 1 0 1 0 1.414 1.414l1.414-1.414Zm9.9-1.414a1 1 0 0 0-1.414 1.414l1.414 1.414a1 1 0 0 0 1.414-1.414l-1.414-1.414ZM13 19a1 1 0 1 0-2 0v2a1 1 0 1 0 2 0v-2Z" clip-rule="evenodd"/>
					</svg>
				</div>
			</div>
		</div>
	`
}
