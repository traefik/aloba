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

	ChannelIDOld := s.options.Slack.ChannelID
	values := r.URL.Query()
	channel, ok := values["channel"]
	if ok {
		s.options.Slack.ChannelID = channel[0]
	}
	if len(s.options.Slack.ChannelID) == 0 {
		log.Printf("Slack channel is mandatory.")
		http.Error(w, "Slack channel is mandatory.", http.StatusBadRequest)
		return
	}

	err := launch(s.options)
	if ok {
		s.options.Slack.ChannelID = ChannelIDOld
	}
	if err != nil {
		log.Printf("Report error: %v", err)
		http.Error(w, "Report error.", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Scheluded.")
}
