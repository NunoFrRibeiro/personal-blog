package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	logger "github.com/NunoFrRibeiro/personal-blog/pkg/logger"
)

func Start() {
	router := http.NewServeMux()

	router.HandleFunc("/", HomeHandler)

	pwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("filePath abs read error: %v", err)
	}
	indexPath := "static"
	path := filepath.Join(pwd, indexPath)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path))))

	posts, errf := HandleMarkDownPosts("posts")
	if errf != nil {
		logger.Errorf("error parsing markdown posts: %v", err)
	}

	for _, post := range posts {
		indexPath := "templates/index.html"
		path := filepath.Join(pwd, indexPath)
		indexContent, err := os.ReadFile(path)
		if err != nil {
			logger.Errorf("error reading index.html: %v", err)
		}
		if post.Slug != "" {
			router.HandleFunc("/"+post.Slug, func(w http.ResponseWriter, r *http.Request) {
				RenderTemplate(w, string(indexContent), post)
			})
		} else {
			logger.Warnf("post titled: %s has empty slug and will not be accessible", post.Title)
		}
	}

	log.Fatal(http.ListenAndServe(":8081", nil))
}
