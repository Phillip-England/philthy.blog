package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	mdFiles, err := NewMarkdownFilesFromDir("./posts")
	if err != nil {
		panic(err)
	}

	err = CreatePostsMdFile(mdFiles)
	if err != nil {
		panic(err)
	}

	vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		templates, _ := vbf.GetContext(KeyTemplates, r).(*template.Template)
		mdContent, err := vbf.LoadMarkdown("./posts/index.md", "dracula")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		var recentPosts []*MarkdownFile
		if len(mdFiles) <= 5 {
			recentPosts = mdFiles
		} else {
			recentPosts = mdFiles[len(mdFiles)-5:]
		}
		if r.URL.Path == "/" {
			vbf.ExecuteTemplate(w, templates, "root.html", map[string]interface{}{
				"Title":       "Welcome - Philthy Blog",
				"Content":     template.HTML(mdContent),
				"ReqPath":     r.URL.Path,
				"ArticleName": "philthy.blog",
				"SubText":     "Saying the things I'm afraid to say",
				"ImageSrc":    "./static/img/profile-sm.png",
				"RecentPosts": recentPosts,
			})
		} else {
			vbf.WriteString(w, "404 not found")
		}
	}, vbf.MwLogger)

	vbf.AddRoute("GET /posts", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		templates, _ := vbf.GetContext(KeyTemplates, r).(*template.Template)
		mdContent, err := vbf.LoadMarkdown("./posts/posts.md", "dracula")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		var recentPosts []*MarkdownFile
		if len(mdFiles) <= 5 {
			recentPosts = mdFiles
		} else {
			recentPosts = mdFiles[len(mdFiles)-5:]
		}
		vbf.ExecuteTemplate(w, templates, "root.html", map[string]interface{}{
			"Title":       "Welcome - Philthy Blog",
			"Content":     template.HTML(mdContent),
			"ReqPath":     r.URL.Path,
			"ArticleName": "philthy.blog",
			"SubText":     "Saying the things I'm afraid to say",
			"ImageSrc":    "./static/img/profile-sm.png",
			"RecentPosts": recentPosts,
		})
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

type MarkdownFile struct {
	Path       string
	ImagePath  string
	Title      string
	PostNumber string
	Content    string
	Href       string
}

func NewMarkdownFilesFromDir(dir string) ([]*MarkdownFile, error) {
	var files []*MarkdownFile
	var potErr error
	potErr = nil
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if path == "posts/index.md" || path == "posts/posts.md" {
			return nil
		}
		file, err := NewMarkdownFileFromPath(path)
		if err != nil {
			potErr = err
			return nil
		}
		files = append(files, file)
		return nil
	})
	if potErr != nil {
		return files, potErr
	}
	return files, nil
}

func NewMarkdownFileFromPath(path string) (*MarkdownFile, error) {
	var file *MarkdownFile
	mdContent, err := vbf.LoadMarkdown("/"+path, "dracula")
	if err != nil {
		return file, err
	}
	fileName := strings.Split(path, "/")[1]
	fileNumber := strings.Split(fileName, "_")[0]
	fileTitle := strings.Split(strings.Split(fileName, "_")[1], ".")[0]
	imagePath := strings.Replace(path, "posts", "static/post_img", 1)
	imagePath = strings.Replace(imagePath, ".md", ".jpg", 1)
	file = &MarkdownFile{
		Path:       path,
		PostNumber: fileNumber,
		Title:      fileTitle,
		Content:    mdContent,
		ImagePath:  "./" + imagePath,
		Href:       "/post/" + fileNumber,
	}
	return file, nil
}

func GetMarkdownFileByPostNumber(mdFiles []*MarkdownFile, number string) (*MarkdownFile, error) {
	for _, file := range mdFiles {
		if file.PostNumber == number {
			return file, nil
		}
	}
	return nil, fmt.Errorf(`post number %s does not exist`, number)
}

func CreatePostsMdFile(mdFiles []*MarkdownFile) error {
	filename := "./posts/posts.md"
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			log.Fatalf("Failed to delete existing file: %v", err)
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()
	file.WriteString("## All Posts\n")
	for _, mdFile := range mdFiles {
		file.WriteString(fmt.Sprintf(`%s. [%s](%s)`, mdFile.PostNumber, mdFile.Title, mdFile.Href))
	}
	return nil
}
