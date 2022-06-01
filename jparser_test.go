package jparser

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
)

type MockRes struct {
	Data string `json:"key,omitempty"`
	Err  string `json:"error,omitempty"`
}

func (m *MockRes) GetStatusCode() int {
	return 200
}

type MockReqBody struct {
	Key string `json:"key"`
}

func TestSend(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Send(w, r, &MockRes{Data: "test"})
	}))
	defer svr.Close()

	res, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("unable to complete Get request %v", err)
	}
	if res.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("invalid json header")
	}
	if res.StatusCode != 200 {
		t.Errorf("invalid status code recieved")
	}
}

func TestSendWithStatusCode(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendWithStatusCode(w, r, &MockRes{Data: "test"}, 201)
	}))
	defer svr.Close()

	res, err := http.Get(svr.URL)
	if err != nil {
		t.Errorf("unable to complete Get request %v", err)
	}
	if res.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("invalid json header")
	}
	if res.StatusCode != 201 {
		t.Errorf("invalid status code recieved")
	}
}

func TestGet(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body MockReqBody
		if err := Get(r, &body); err != nil {
			t.Errorf("unable to get request body %v", err)
		}

		if body.Key != "value" {
			t.Errorf("unable to get request key")
		}
	}))
	defer svr.Close()
	postBody, _ := json.Marshal(map[string]string{
		"key": "value",
	})
	responseBody := bytes.NewBuffer(postBody)
	_, err := http.Post(svr.URL, "application/json", responseBody)
	if err != nil {
		t.Errorf("unable to complete Get request %v", err)
	}
}
