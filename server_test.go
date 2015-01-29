package quayd

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func loadFixture(fixture string, t testing.TB) io.Reader {
	body, err := ioutil.ReadFile("test-fixtures/quay.io/" + fixture + ".json")
	if err != nil {
		t.Fatalf("Unable to load fixture %s: %s", fixture, err)
	}

	return bytes.NewReader(body)
}

func TestWebhook(t *testing.T) {
	r := DefaultStatusesRepository
	s := NewServer(nil)
	defer r.Reset()

	tests := []struct {
		status   string
		fixture  string
		expected Status
	}{
		{"pending", "pending_build", Status{Repo: "ejholmes/docker-statsd", Ref: "long-f1fb3b0", State: "pending", Context: "Docker Image", TargetURL: "https://quay.io/repository/ejholmes/docker-statsd/build?current=077f3664-35d3-48e6-9da7-889f9be73070", Description: "Quay: Image Building"}},
		{"success", "pending_build", Status{Repo: "ejholmes/docker-statsd", Ref: "long-f1fb3b0", State: "success", Context: "Docker Image", TargetURL: "https://quay.io/repository/ejholmes/docker-statsd/build?current=077f3664-35d3-48e6-9da7-889f9be73070", Description: "Quay: Image Built"}},
	}

	for _, tt := range tests {
		r.Reset()

		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/quay/"+tt.status, loadFixture(tt.fixture, t))

		s.ServeHTTP(resp, req)

		if len(r.statuses) != 1 {
			t.Fatal("Expected 1 commit status")
		}

		if got, want := r.statuses[0], &tt.expected; !reflect.DeepEqual(got, want) {
			t.Fatalf("Status => %q; want %q", got, want)
		}
	}
}

func TestWebhook_InvalidStatus(t *testing.T) {
	r := DefaultStatusesRepository
	s := NewServer(nil)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quay/foo", loadFixture("pending_build", t))

	s.ServeHTTP(resp, req)

	if len(r.statuses) != 0 {
		t.Fatal("Expected 0 commit statuses")
	}
}

func TestWebhook_ManualTrigger(t *testing.T) {
	r := DefaultStatusesRepository
	s := NewServer(nil)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quay/pending", loadFixture("pending_build.manual", t))

	s.ServeHTTP(resp, req)

	if len(r.statuses) != 0 {
		t.Fatal("Expected 0 commit statuses")
	}
}

func TestWebhook_TagsImageID(t *testing.T) {
	s := NewServer(nil)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/quay/success", loadFixture("pending_build", t))

	s.ServeHTTP(resp, req)
}
