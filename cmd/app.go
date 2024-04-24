package app

import (
	"fmt"
	"net/http"

	// 	errors "github.com/NunoFrRibeiro/personal-go-blog/pkg/errors"
	logger "github.com/NunoFrRibeiro/personal-blog/pkg/logger"
)

func Start() {

	router := http.NewServeMux()
	router.HandleFunc("/", index)
	router.HandleFunc("/posts", getAllPosts)
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Errorf(fmt.Sprintf("listen and serve error: %s", err))
	}
}
