package quayd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var body = `{"build_id": "077f3664-35d3-48e6-9da7-889f9be73070", "trigger_kind": "github", "name": "docker-statsd", "repository": "ejholmes/docker-statsd", "namespace": "ejholmes", "docker_url": "quay.io/ejholmes/docker-statsd", "visibility": "public", "docker_tags": ["test"], "build_name": "f1fb3b0", "trigger_id": "ffcbfaef-c7fe-4721-b69e-2e78fb6d29d5", "homepage": "https://quay.io/repository/ejholmes/docker-statsd/build?current=077f3664-35d3-48e6-9da7-889f9be73070"}`

func TestWebhook(t *testing.T) {
	r := DefaultStatusesRepository
	s := NewServer(nil)
	defer r.Reset()

	tests := []struct {
		status   string
		body     string
		expected Status
	}{
		{"pending", body, Status{Repo: "ejholmes/docker-statsd", Ref: "long-f1fb3b0", State: "pending", Context: "docker"}},
		{"success", body, Status{Repo: "ejholmes/docker-statsd", Ref: "long-f1fb3b0", State: "success", Context: "docker"}},
	}

	for _, tt := range tests {
		r.Reset()

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/quay/"+tt.status, bytes.NewBufferString(tt.body))

		s.ServeHTTP(resp, req)

		if len(r.Statuses) != 1 {
			t.Fatal("Expected 1 commit status")
		}

		if got, want := r.Statuses[0], &tt.expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("Status => %q; want %q", got, want)
		}
	}
}

func TestWebhook_InvalidStatus(t *testing.T) {
	r := DefaultStatusesRepository
	s := NewServer(nil)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quay/foo", bytes.NewBufferString(body))

	s.ServeHTTP(resp, req)

	if len(r.Statuses) != 0 {
		t.Fatal("Expected 0 commit statuses")
	}
}
