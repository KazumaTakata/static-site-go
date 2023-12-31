package render

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
)

type fontMatter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
}

type ViewData struct {
	Title    string
	Tags     []string
	Content  template.HTML
	PostList []Post
}

type Post struct {
	Title string
	Url   string
}

func renderListPageHTML(path string) {

	destinationPath := "dist"
	pathWithoutRoot := strings.Join(strings.Split(path, "/")[1:], "/")

	fmt.Println("info.IsDir()")
	fmt.Println(path)

	_ = os.MkdirAll(filepath.Join(destinationPath, pathWithoutRoot), os.ModePerm)

	tmpl, err := template.ParseFiles("template/template.html", "template/header.html", "template/list.html")

	if err != nil {
		panic("Parse failed!!")
	}

	fileOrDirs, err := os.ReadDir(path)

	var postList []Post

	for _, fileOrDir := range fileOrDirs {
		if !fileOrDir.IsDir() {
			fileName := fileOrDir.Name()
			filepathUnderDir := filepath.Join(path, fileName)
			dat, _ := os.ReadFile(filepathUnderDir)
			var matter fontMatter
			_, err := frontmatter.Parse(strings.NewReader(string(dat)), &matter)
			if err != nil {
			}
			fmt.Printf("%+v\n", matter)

			fileExtention := filepath.Ext(fileName)

			newFilePath := fileName[0 : len(fileName)-len(fileExtention)]
			post := Post{Title: matter.Title, Url: newFilePath}
			postList = append(postList, post)
		}
	}

	viewData := ViewData{PostList: postList}

	var htmlOutputBuffer bytes.Buffer

	tmpl.Execute(&htmlOutputBuffer, viewData)

}

func RenderHTML() {

	destinationPath := "dist"
	sourcePath := "source"

	_ = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			pathWithoutRoot := strings.Join(strings.Split(path, "/")[1:], "/")

			if info.IsDir() {
				renderListPageHTML(path)
				return nil
			}

			dat, _ := os.ReadFile(path)
			var matter fontMatter
			rest, err := frontmatter.Parse(strings.NewReader(string(dat)), &matter)
			if err != nil {
			}
			fmt.Printf("%+v\n", matter)

			var buf bytes.Buffer
			if err := goldmark.Convert(rest, &buf); err != nil {
				panic(err)
			}

			viewData := ViewData{Content: template.HTML(buf.String()), Title: matter.Title, Tags: matter.Tags}

			var htmlOutputBuffer bytes.Buffer

			tmpl, err := template.ParseFiles("template/template.html", "template/header.html", "template/post.html")

			if err != nil {
				panic("Parse failed!!")
			}

			tmpl.Execute(&htmlOutputBuffer, viewData)

			replacedHTML, _ := replaceCodeParts(htmlOutputBuffer.Bytes(), styles.Monokai)

			destinationFilePath := filepath.Join(destinationPath, pathWithoutRoot)

			fmt.Println(filepath.Ext(destinationFilePath))
			fileExtention := filepath.Ext(destinationFilePath)

			newFilePath := destinationFilePath[0:len(destinationFilePath)-len(fileExtention)] + ".html"

			os.WriteFile(newFilePath, []byte(replacedHTML), 0644)

		}

		return nil
	})
}
