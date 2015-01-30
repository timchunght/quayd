package quayd

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var validStatuses = []string{"pending", "success", "error", "failure"}

type Server struct {
	http.Handler
}

func NewServer(q *Quayd) *Server {
	if q == nil {
		q = Default
	}

	m := mux.NewRouter()

	m.Handle("/quay/{status}", &Webhook{q}).Methods("POST")

	return &Server{m}
}

type Webhook struct {
	*Quayd
}

type WebhookForm struct {
	Repository  string   `json:"repository"`
	TriggerKind string   `json:"trigger_kind"`
	IsManual    bool     `json:"is_manual"`
	DockerTags  []string `json:"docker_tags"`
	BuildName   string   `json:"build_name"`
	BuildURL    string   `json:"homepage"`
}

func (wh *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]
	if !validStatus(status) {
		http.Error(w, "Invalid status: "+status, 400)
		return
	}

	var form WebhookForm

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// We don't want to process manually triggered builds.
	if !(!form.IsManual && form.TriggerKind == "github") {
		w.WriteHeader(204)
		return
	}

	if status == "success" {
		if err := wh.Quayd.LoadImageTags(form.DockerTags[0], form.Repository, form.BuildName); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err := wh.Quayd.Handle(form.Repository, form.BuildName, form.BuildURL, status); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func validStatus(a string) bool {
	for _, b := range validStatuses {
		if b == a {
			return true
		}
	}
	return false
}
