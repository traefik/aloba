package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/containous/aloba/options"
)

type server struct {
	options *options.Report
}

func (s *server) ListenAndServe() error {
	return http.ListenAndServe(":"+strconv.Itoa(s.options.ServerPort), s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid http method: %s", r.Method)
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := launch(s.options)
	if err != nil {
		log.Printf("Report error: %v", err)
		http.Error(w, "Report error.", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Scheluded.")
}
