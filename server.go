package quayd

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Handler
}

func NewServer(s *StatusesService) *Server {
	if s == nil {
		s = DefaultStatusesService
	}

	m := mux.NewRouter()

	m.Handle("/quay/{status}", &Webhook{s}).Methods("POST")

	return &Server{m}
}

type Webhook struct {
	*StatusesService
}

type WebhookForm struct {
	Repository  string `json:"repository"`
	TriggerKind string `json:"trigger_kind"`
	BuildName   string `json:"build_name"`
}

func (wh *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var form WebhookForm

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := wh.StatusesService.Create(form.Repository, form.BuildName, vars["status"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
