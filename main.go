package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Phillip-England/vbf"
	"github.com/PuerkitoBio/goquery"
)

const KeyTemplates = "KEYTEMPLATES"
const TitleCatchPhrase = "Philthy Blog: Christianity, Doubt, and Faith"

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
		if r.URL.Path == "/" {
			vbf.ExecuteTemplate(w, templates, "root.html", BaseTemplate{
				Title:            "Welcome - " + TitleCatchPhrase,
				Content:          template.HTML(mdContent),
				ReqPath:          r.URL.Path,
				ArticleName:      "philthy.blog",
				SubText:          "Saying the things I'm afraid to say",
				ImageSrc:         "./static/img/profile-sm.png",
				DateWritten:      "1/20/2025",
				MetaDescription:  "Philthy Blog: A heartfelt exploration of Christianity, doubt, and the personal struggle with faith and God. Join an honest journey seeking truth, understanding, and peace.",
				MetaKeywords:     "Christianity, God, Faith",
				EmbeddedVideoSrc: "https://www.youtube.com/embed/k8S1DCdhneQ?si=1L3ZlcHVjXJ24RvK",
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
		vbf.ExecuteTemplate(w, templates, "root.html", BaseTemplate{
			Title:           "Posts - " + TitleCatchPhrase,
			Content:         template.HTML(mdContent),
			ReqPath:         r.URL.Path,
			ArticleName:     "philthy.blog",
			SubText:         "Things I've written",
			ImageSrc:        "./static/img/posts-sm.webp",
			DateWritten:     "1/20/2025",
			MetaDescription: "Philthy Blog: A heartfelt exploration of Christianity, doubt, and the personal struggle with faith and God. Here, see all the posts I've written.",
			MetaKeywords:    "Christianity, God, Faith",
		})
	}, vbf.MwLogger)

	vbf.AddRoute("GET /post/{postNumber}", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		postStr := r.PathValue("postNumber")
		postNumber, err := strconv.Atoi(postStr)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		mdFile := mdFiles[postNumber]
		vbf.ExecuteTemplate(w, templates, "root.html", BaseTemplate{
			Title:            mdFile.Title + " - " + TitleCatchPhrase,
			Content:          mdFile.Content,
			ReqPath:          r.URL.Path,
			ArticleName:      mdFile.Title,
			SubText:          mdFile.SubText,
			ImageSrc:         mdFile.ImagePath,
			DateWritten:      mdFile.DateWritten,
			MetaDescription:  "Philthy Blog: " + mdFile.MetaDescription,
			MetaKeywords:     mdFile.MetaKeywords,
			EmbeddedVideoSrc: mdFile.EmbeddedVideoSrc,
		})

	}, vbf.MwLogger)

	vbf.AddRoute("GET /screenplay", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
		vbf.ExecuteTemplate(w, templates, "screenplay-root.html", ScreenplayTemplate{
			Title: "Philthy Blog: Screenplay",
		})

	}, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}

type ScreenplayTemplate struct {
	Title string
}

type BaseTemplate struct {
	Title            string
	Content          template.HTML
	ReqPath          string
	ArticleName      string
	SubText          string
	ImageSrc         string
	DateWritten      string
	MetaDescription  string
	MetaKeywords     string
	EmbeddedVideoSrc string
}

type MarkdownFile struct {
	Path             string
	ImagePath        string
	Title            string
	PostNumber       string
	Content          template.HTML
	Href             string
	SubText          string
	DateWritten      string
	MetaDescription  string
	MetaKeywords     string
	EmbeddedVideoSrc string
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
	fileNameWithoutExtension := strings.Split(fileName, ".")[0]
	fileTitleWithUnderscoresSlice := strings.Split(fileNameWithoutExtension, "_")[1:]
	fileTitleWithUnderscores := strings.Join(fileTitleWithUnderscoresSlice, " ")
	fileTitle := strings.Replace(fileTitleWithUnderscores, "_", " ", 1)
	fileTitle = strings.ReplaceAll(fileTitle, "_", " ")
	imagePath := strings.Replace(path, "posts", "/static/post_img", 1)
	imagePath = strings.Replace(imagePath, ".md", ".webp", 1)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(mdContent))
	if err != nil {
		return file, err
	}
	metaDataElm := doc.Find("#meta")
	var subText string
	var dob string
	var metaDescription string
	var metaKeywords string
	var videoSrc string
	metaDataElm.Find("*").Each(func(i int, sel *goquery.Selection) {
		key, _ := sel.Attr("key")
		if key == "video-src" {
			val, _ := sel.Attr("value")
			videoSrc = val
		}
		if key == "subtext" {
			val, _ := sel.Attr("value")
			subText = val
		}
		if key == "dob" {
			val, _ := sel.Attr("value")
			dob = val
		}
		if key == "description" {
			val, _ := sel.Attr("value")
			metaDescription = val
		}
		if key == "keywords" {
			val, _ := sel.Attr("value")
			metaKeywords = val
		}
	})
	file = &MarkdownFile{
		Path:             path,
		PostNumber:       fileNumber,
		Title:            fileTitle,
		Content:          template.HTML(mdContent),
		ImagePath:        imagePath,
		Href:             "/post/" + fileNumber,
		SubText:          subText,
		DateWritten:      dob,
		MetaDescription:  metaDescription,
		MetaKeywords:     metaKeywords,
		EmbeddedVideoSrc: videoSrc,
	}
	return file, nil
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
		file.WriteString(fmt.Sprintf(`%s. [%s](%s) *%s*`, mdFile.PostNumber, mdFile.Title, mdFile.Href, mdFile.DateWritten) + "\n")
	}
	return nil
}
