package quayd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/negroni"
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

	n := negroni.Classic()
	n.UseHandler(m)

	return &Server{n}
}

type Webhook struct {
	*Quayd
}

type WebhookForm struct {
	Repository      string                 `json:"repository"`
	TriggerKind     string                 `json:"trigger_kind"`
	IsManual        bool                   `json:"is_manual"`
	DockerTags      []string               `json:"docker_tags"`
	BuildName       string                 `json:"build_name"`
	BuildURL        string                 `json:"homepage"`
	TriggerMetadata map[string]interface{} `json:"trigger_metadata"`
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
		errorResponse(w, err)
		return
	}

	fmt.Println(form)
	// We don't want to process manually triggered builds.
	if !(!form.IsManual && form.TriggerKind == "github") {
		w.WriteHeader(204)
		return
	}

	if status == "success" {
		commitID, ok := form.TriggerMetadata["commit"].(string)
		if !ok {
			errorResponse(w, errors.New("Missing commit"))
			return
		}
		fmt.Println("Commit: ", commitID)
		if err := wh.Quayd.LoadImageTags(commitID, form.DockerTags[0], form.Repository, form.BuildName); err != nil {
			errorResponse(w, err)
			return
		}
	}

	// if err := wh.Quayd.Handle(form.Repository, form.BuildName, form.BuildURL, status); err != nil {
	// 	errorResponse(w, err)
	// 	return
	// }

}

func validStatus(a string) bool {
	for _, b := range validStatuses {
		if b == a {
			return true
		}
	}
	return false
}

func errorResponse(w http.ResponseWriter, err error) {
	fmt.Println(err)
	http.Error(w, err.Error(), 500)
}
