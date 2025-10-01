package aweFunc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type respDTO struct {
	OK    bool `json:"ok"`
	Value int  `json:"value"`
}

func TestRequestPostJSON_Success(t *testing.T) {
	t.Parallel()

	wantBody := map[string]string{"a": "b"}
	saw := struct {
		Accept      string
		ContentType string
		Body        map[string]string
	}{}

	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s", r.Method)
				}
				saw.Accept = r.Header.Get("Accept")
				saw.ContentType = r.Header.Get("Content-Type")
				var m map[string]string
				if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
					t.Fatalf("request body decode: %v", err)
				}
				saw.Body = m

				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, `{"ok":true,"value":123}`)
				if err != nil {
					t.Fatalf("WriteString: %v", err)
				}
			},
		),
	)
	defer srv.Close()

	var out respDTO
	if err := RequestPostJSON(context.Background(), srv.URL, wantBody, &out); err != nil {
		t.Fatalf("RequestPostJSON: %v", err)
	}

	if !out.OK || out.Value != 123 {
		t.Fatalf("out mismatch: %+v", out)
	}
	if saw.Accept != "application/json" {
		t.Fatalf("missing/invalid Accept: %q", saw.Accept)
	}
	if saw.ContentType != "application/json" {
		t.Fatalf("missing/invalid Content-Type: %q", saw.ContentType)
	}
	if got := saw.Body["a"]; got != "b" {
		t.Fatalf("request body mismatch: %v", saw.Body)
	}
}

func TestRequestPostJSON_MarshalError(t *testing.T) {
	t.Parallel()

	badBody := map[string]float64{"x": math.NaN()}
	var out respDTO

	err := RequestPostJSON(context.Background(), "http://unused.local", badBody, &out)
	if err == nil || !strings.Contains(err.Error(), "RequestPostJSON: marshal body") {
		t.Fatalf("want marshal error, got: %v", err)
	}
}

func TestRequestPostJSON_MakeRequestError(t *testing.T) {
	t.Parallel()

	badURL := "http://%zz"
	var out respDTO

	err := RequestPostJSON(context.Background(), badURL, map[string]int{"x": 1}, &out)
	if err == nil || !strings.Contains(err.Error(), "RequestPostJSON: make request") {
		t.Fatalf("want to make request error, got: %v", err)
	}
}

func TestRequestPostJSON_DecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, "not-json")
				if err != nil {
					t.Fatalf("WriteString: %v", err)
				}
			},
		),
	)
	defer srv.Close()

	var out respDTO
	err := RequestPostJSON(context.Background(), srv.URL, map[string]int{"x": 1}, &out)
	if err == nil || !strings.Contains(err.Error(), "decode") {
		t.Fatalf("want decode error, got: %v", err)
	}
}

func TestRequestPostJSON_NoContent(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
		),
	)
	defer srv.Close()

	var out respDTO
	if err := RequestPostJSON(context.Background(), srv.URL, nil, &out); err != nil {
		t.Fatalf("RequestPostJSON: %v", err)
	}
}

func TestRequestPostJSON_OutNil(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, err := io.WriteString(w, `{"ok":true,"value":7}`)
				if err != nil {
					t.Fatalf("WriteString: %v", err)
				}
			},
		),
	)
	defer srv.Close()

	if err := RequestPostJSON[respDTO](context.Background(), srv.URL, nil, nil); err != nil {
		t.Fatalf("RequestPostJSON: %v", err)
	}
}

func TestRequestPostJSON_Non2xx_ShortBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "bad req\n", http.StatusBadRequest)
			},
		),
	)
	defer srv.Close()

	var out respDTO
	err := RequestPostJSON(context.Background(), srv.URL, nil, &out)
	if err == nil || !strings.Contains(err.Error(), "status 400") || !strings.Contains(err.Error(), "bad req") {
		t.Fatalf("want status 400 with body, got: %v", err)
	}
}

func TestRequestPostJSON_Non2xx_LongBodyTrimmed(t *testing.T) {
	t.Parallel()

	long := bytes.Repeat([]byte("a"), 6000)
	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write(long)
				if err != nil {
					t.Fatalf("WriteHeader: %v", err)
				}
			},
		),
	)
	defer srv.Close()

	var out respDTO
	err := RequestPostJSON(context.Background(), srv.URL, nil, &out)
	if err == nil {
		t.Fatalf("expected error")
	}

	msg := err.Error()
	i := strings.Index(msg, ": ")
	if i < 0 {
		t.Fatalf("unexpected error format: %q", msg)
	}
	bodyPart := msg[i+2:]

	if len(bodyPart) > 4096+32 {
		t.Fatalf("body not trimmed enough, len=%d", len(bodyPart))
	}
}

func TestRequestPostJSON_NoBody_RequestHeaders(t *testing.T) {
	t.Parallel()

	var sawCT, sawAccept string
	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				sawCT = r.Header.Get("Content-Type")
				sawAccept = r.Header.Get("Accept")
				_, err := io.WriteString(w, `{"ok":true,"value":1}`)
				if err != nil {
					t.Fatalf("WriteHeader: %v", err)
				}
			},
		),
	)
	defer srv.Close()

	var out respDTO
	if err := RequestPostJSON(context.Background(), srv.URL, nil, &out); err != nil {
		t.Fatalf("RequestPostJSON: %v", err)
	}
	if sawCT != "" {
		t.Fatalf("Content-Type should be empty when body=nil, got %q", sawCT)
	}
	if sawAccept != "application/json" {
		t.Fatalf("Accept header missing, got %q", sawAccept)
	}
}

func TestRequestPostJSON_RespectsCallerTimeout(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(50 * time.Millisecond)
				_, err := io.WriteString(w, `{"ok":true,"value":1}`)
				if err != nil {
					t.Fatalf("WriteHeader: %v", err)
				}
			},
		),
	)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	var out respDTO
	err := RequestPostJSON(ctx, srv.URL, map[string]int{"x": 1}, &out)
	if err == nil {
		t.Fatalf("expected timeout error")
	}
}
