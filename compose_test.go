package mixer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/mixer"
)

func TestClone(t *testing.T) {
	f := mixer.Classic()

	var result string
	f.Before(func(c mixer.Context) {
		result += "f"
	})

	s := f.Clone()
	s.Before(func(c mixer.Context) {
		result += "s"
	})

	h1 := func(c mixer.Context) {

	}
	r2 := httptest.NewRecorder()
	s.Handler(h1).ServeHTTP(r2, (*http.Request)(nil))
	expect(t, result, "fs")

	r1 := httptest.NewRecorder()
	f.Handler(h1).ServeHTTP(r1, (*http.Request)(nil))
	expect(t, result, "fsf")
}
