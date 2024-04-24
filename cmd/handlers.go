package app

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
 fmt.Println("<h3>Hello Goland!</h3>") 
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
  
}
