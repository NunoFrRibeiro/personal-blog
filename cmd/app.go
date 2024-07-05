package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func Start() {
	sidebarData, err := handleSideBarData("./posts/")
	if err != nil {
		log.Fatal(err)
	}

	posts, err := HandleMarkDownPosts("./posts/")
	if err != nil {
		log.Fatal(err)
	}

	funcMap := template.FuncMap{
		"dict": dict,
	}

	tmpl := template.New("").Funcs(funcMap)
	tmpl, errs := tmpl.ParseGlob("./templates/*")
	if errs != nil {
		log.Fatal(errs)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath := "./posts/index.md"
		indexContent, err := os.ReadFile(indexPath)
		if err != nil {
			log.Printf("Error occurred during operation: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		post, err := parseMarkDownFile(indexContent)
		if err != nil {
			log.Printf("Error occurred during operation: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		sidebarLinks := handleSideBarLinks(post.Headers)

		data := map[string]interface{}{
			"Title":                   post.Title,
			"Content":                 post.Content,
			"SidebarData":             sidebarData,
			"Headers":                 post.Headers,
			"SidebarLinks":            sidebarLinks,
			"CurrentSlug":             post.Slug,
			"MetaDescription":         post.MetaDescription,
			"MetaPropertyTitle":       post.MetaPropertyTitle,
			"MetaPropertyDescription": post.MetaPropertyDescription,
			"MetaOgURL":               post.MetaOgURL,
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			log.Printf("Error executing template: %v\n", err)
		}
	})

	for _, post := range posts {
		localPost := post
		if localPost.Slug != "" {
			http.HandleFunc("/"+localPost.Slug, func(w http.ResponseWriter, r *http.Request) {
				sidebarLinks := handleSideBarLinks(localPost.Headers)
				data := map[string]interface{}{
					"Title":                   localPost.Title,
					"Content":                 localPost.Content,
					"SidebarData":             sidebarData,
					"Headers":                 localPost.Headers,
					"Description":             localPost.Description,
					"SidebarLinks":            sidebarLinks,
					"CurrentSlug":             localPost.Slug,
					"MetaDescription":         localPost.MetaDescription,
					"MetaPropertyTitle":       localPost.MetaPropertyTitle,
					"MetaPropertyDescription": localPost.MetaPropertyDescription,
					"MetaOgURL":               localPost.MetaOgURL,
				}

				if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
					log.Printf("Error executing template: %v\n", err)
				}
			})
		} else {
			log.Printf("Warning: Post titled '%s' has an empty slug and will not be accessible via a unique URL.\n", localPost.Title)
		}
	}

	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Page Not Found",
		}

		if err := tmpl.ExecuteTemplate(w, "404.html", data); err != nil {
			log.Printf("Error executing template: %v\n", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid number of arguments for dict function")
	}
	d := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		d[key] = values[i+1]
	}
	return d, nil
}
