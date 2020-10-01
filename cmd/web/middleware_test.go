package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	want := "deny"
	got := rs.Header.Get("X-Frame-Options")
	if got != want {
		t.Errorf("want %q; got %q", want, got)
	}

	want = "1; mode=block"
	got = rs.Header.Get("X-XSS-Protection")
	if got != want {
		t.Errorf("want %q; got %q", want, got)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %q; got %q", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)

	if string(body) != "OK" {
		t.Errorf("want OK; got %q", string(body))
	}
}
