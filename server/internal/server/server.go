package server

import (
	"fmt"
	"net/http"
)

func StartServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
