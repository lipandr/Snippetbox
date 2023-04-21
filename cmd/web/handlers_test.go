package main

import (
	"bytes"
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

func TestPingEndToEnd(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want %q; got %q", "OK", string(body))
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		url      string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("This is a snippet")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.url)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if !bytes.Contains(body, []byte(tt.wantBody)) {
				t.Errorf("want %q; got %q", tt.wantBody, string(body))
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Log("csrfToken:", csrfToken)

	// TODO investigate why csrfToken is not being passed to signup handler and handler returns 400
	//tests := []struct {
	//	name         string
	//	userName     string
	//	userEmail    string
	//	userPassword string
	//	csrfToken    string
	//	wantCode     int
	//	wantBody     []byte
	//}{
	//	{"Valid submission", "bob", "bob@example.com", "validPa$$word",
	//		csrfToken, http.StatusSeeOther, nil},
	//	{"Empty name", "", "bob@example.com", "validPa$$word",
	//		csrfToken, http.StatusOK, []byte("This field cannot be blank")},
	//	{"Empty email", "bob", "", "validPa$$word",
	//		csrfToken, http.StatusOK, []byte("This field cannot be blank")},
	//	{"Empty password", "bob", "bob@example.com", "",
	//		csrfToken, http.StatusOK, []byte("This field cannot be blank")},
	//	{"Invalid email (incomplete domain)", "bob", "bob@example", "validPa$$word",
	//		csrfToken, http.StatusOK, []byte("This field is invalid")},
	//	{"Invalid email (missing @)", "bob", "bobexample.com", "validPa$$word",
	//		csrfToken, http.StatusOK, []byte("This field is invalid")},
	//	{"Invalid email (missing local part)", "bob", "@example.com", "validPa$$word",
	//		csrfToken, http.StatusOK, []byte("This field is invalid")},
	//	{"Short password", "bob", "bob@example.con", "pass",
	//		csrfToken, http.StatusOK, []byte("This field too short (minimum is 10 characters)")},
	//	{"Duplicate email", "bob", "bob@example", "validPa$$word",
	//		csrfToken, http.StatusOK, []byte("Address is already in use")},
	//	{"Invalid CSRF token", "", "", "",
	//		"invalidToken", http.StatusBadRequest, nil},
	//}
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		form := url.Values{}
	//		form.Add("name", tt.userName)
	//		form.Add("email", tt.userEmail)
	//		form.Add("password", tt.userPassword)
	//		form.Add("csrf_token", tt.csrfToken)
	//
	//		code, _, body := ts.postForm(t, "/user/signup", form)
	//
	//		if code != tt.wantCode {
	//			t.Errorf("want %d; got %d", tt.wantCode, code)
	//		}
	//		if !bytes.Contains(body, tt.wantBody) {
	//			t.Errorf("want %q; got %s", tt.wantBody, body)
	//		}
	//	})
	//}
}
