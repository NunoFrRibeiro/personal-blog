package app

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	parsers "github.com/gomarkdown/markdown/parser"

	errorf "github.com/NunoFrRibeiro/personal-blog/pkg/errors"
	"github.com/NunoFrRibeiro/personal-blog/pkg/logger"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		errf        *errorf.AppError
		sideBarData SideBarData
	)
	pwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("filePath abs read error: %v", err)
	}
	indexPath := "posts/index.md"
	path := filepath.Join(pwd, indexPath)
	indexContent, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	post, err := parseMarkDownFile(indexContent)
	if err != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	// Sidebar Data
	sideBarData, errf = handleSideBarData("personal-blog/posts/")
	if errf != nil {
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	sideBarLinks := handleSideBarLinks(post.Headers)
	data := map[string]interface{}{
		"Title":                   post.Title,
		"Content":                 post.Content,
		"SidebarData":             sideBarData,
		"Headers":                 post.Headers,
		"SidebarLinks":            sideBarLinks,
		"CurrentSlug":             post.Slug,
		"MetaDescription":         post.MetaDescription,
		"MetaPropertyTitle":       post.MetaPropertyTitle,
		"MetaPropertyDescription": post.MetaPropertyDescription,
		"MetaOgURL":               post.MetaOgURL,
	}

	RenderTemplate(w, indexPath, data)
}

func HandleMarkDownPosts(dir string) ([]BlogPost, *errorf.AppError) {
	var posts []BlogPost
	pwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("filePath abs read error: %v", err)
		return nil, errorf.ValidationError("filePath abs read error")
	}
	path := filepath.Join(pwd, dir)
	files, err := os.ReadDir(path)
	if err != nil {
		logger.Errorf("os read dir error: %v", err)
		return nil, errorf.ValidationError("os read dir error")
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			path := dir + "/" + file.Name()
			content, err := os.ReadFile(path)
			if err != nil {
				logger.Errorf("os read file error: %v", err.Error())
				return nil, errorf.ValidationError("os read file error")
			}

			post, err := parseMarkDownFile(content)
			if err != nil {
				logger.Errorf("error parsing markdown file: %v", err.Error())
				return nil, errorf.ValidationError("os read file error")
			}

			posts = append(posts, post)
		}
	}
	return posts, nil
}

func parseMarkDownFile(content []byte) (BlogPost, error) {
	sections := strings.SplitN(string(content), "---", 2)
	if len(sections) < 2 {
		return BlogPost{}, errors.New("invalid markdown format")
	}

	metadata := sections[0]
	mdData := sections[1]

	metadata = strings.ReplaceAll(metadata, "\r", "")
	mdData = strings.ReplaceAll(mdData, "\r", "")

	title,
		slug,
		parent,
		description,
		order,
		metaDescription,
		metaPropertyTitle,
		MetaPropertyDescription,
		metaOgUrl := parseMetadata(metadata)

	htmlContent := mdToHTML([]byte(mdData))
	headers := handleHeaders([]byte(mdData))

	return BlogPost{
		Title:                   title,
		Slug:                    slug,
		Parent:                  parent,
		Description:             description,
		Content:                 template.HTML(htmlContent),
		Headers:                 headers,
		Order:                   order,
		MetaDescription:         metaDescription,
		MetaPropertyTitle:       metaPropertyTitle,
		MetaPropertyDescription: MetaPropertyDescription,
		MetaOgURL:               metaOgUrl,
	}, nil
}

func handleSideBarData(dir string) (SideBarData, *errorf.AppError) {
	var sidebarData SideBarData
	categoriesMap := make(map[string]*Category)

	posts, err := HandleMarkDownPosts(dir)
	if err != nil {
		logger.Errorf("error handling markdown post: %v", err)
		return sidebarData, errorf.ValidationError("error handling markdown posts")
	}

	for _, post := range posts {
		if post.Parent != "" {
			if _, exists := categoriesMap[post.Parent]; !exists {
				categoriesMap[post.Parent] = &Category{
					Name:  post.Parent,
					Pages: []BlogPost{post},
					Order: post.Order,
				}
			} else {
				categoriesMap[post.Parent].Pages = append(categoriesMap[post.Parent].Pages, post)
			}
		}
	}

	for _, cat := range categoriesMap {
		sidebarData.Categories = append(sidebarData.Categories, *cat)
	}

	sort.Slice(sidebarData.Categories, func(i, j int) bool {
		return sidebarData.Categories[i].Order < sidebarData.Categories[j].Order
	})

	return sidebarData, nil
}

func handleSideBarLinks(headers []string) template.HTML {
	var linksHTML string
	for _, header := range headers {
		sanitizedHeader := sanitizeHeaderForId(header)
		link := fmt.Sprintf(`<li><a href="#%s">%s</a></li>`, sanitizedHeader, header)
		linksHTML += link
	}
	return template.HTML(linksHTML)
}

func handleHeaders(content []byte) []string {
	var headers []string

	re := regexp.MustCompile(`(?m)^##\s+(.*)`)
	matches := re.FindAllSubmatch(content, -1)

	for _, match := range matches {
		headers = append(headers, string(match[1]))
	}

	return headers
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mdToHTML(md []byte) []byte {
	extensions := parsers.CommonExtensions | parsers.AutoHeadingIDs
	parser := parsers.NewWithExtensions(extensions)

	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}

	renderer := html.NewRenderer(opts)

	doc := parser.Parse(md)

	output := markdown.Render(doc, renderer)

	return output
}

func parseMetadata(metadata string) (
	title string,
	slug string,
	parent string,
	description string,
	order int,
	metadataDescription string,
	metadataPropertyTitle string,
	metadataPropertyDescription string,
	metaOgURL string,
) {
	re := regexp.MustCompile(`(?m)^(\w+):\s*(.+)`)
	matches := re.FindAllStringSubmatch(metadata, -1)

	metadataMap := make(map[string]string)
	for _, match := range matches {
		if len(match) == 3 {
			metadataMap[match[1]] = match[2]
		}
	}

	title = metadataMap["Title"]
	slug = metadataMap["Slug"]
	parent = metadataMap["Parent"]
	description = metadataMap["Description"]
	orderStr := metadataMap["Order"]
	metaDescriptionStr := metadataMap["MetaDescription"]
	metaPropertyTitleStr := metadataMap["MetaPropertyTitle"]
	metaPropertyDescriptionStr := metadataMap["MetaPropertyDescription"]
	metaOgURLStr := metadataMap["MetaOgURL"]

	orderStr = strings.TrimSpace(orderStr)
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		logger.Errorf("error converting order from string: %s", err)
		order = 9999
	}

	return title, slug, parent, description, order, metaDescriptionStr, metaPropertyTitleStr, metaPropertyDescriptionStr, metaOgURLStr
}

func sanitizeHeaderForId(header string) string {
	// characters to lowercase
	header = strings.ToLower(header)

	// replace spaces with hyphens
	header = strings.ReplaceAll(header, " ", "-")

	header = regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(header, "")

	return header
}
