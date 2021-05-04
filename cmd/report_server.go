package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/traefik/aloba/options"
)

type server struct {
	options *options.Report
}

func (s *server) ListenAndServe() error {
	return http.ListenAndServe(":"+strconv.Itoa(s.options.ServerPort), s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Error().Msgf("Invalid http method: %s", r.Method)
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
		log.Error().Msg("Slack channel is mandatory.")
		http.Error(w, "Slack channel is mandatory.", http.StatusBadRequest)
		return
	}

	err := launch(s.options)
	if ok {
		s.options.Slack.ChannelID = ChannelIDOld
	}
	if err != nil {
		log.Error().Err(err).Msg("Report error.")
		http.Error(w, "Report error.", http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprint(w, "Scheluded.")
}
