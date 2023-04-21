package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Error(err)
	}
	ping(rr, r)

	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(rs.Body)

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) != "OK" {
		t.Errorf("want %q; got %q", "OK", string(b))
	}
}
