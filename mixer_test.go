package mixer_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/codegangsta/mixer"
)

func ExampleClassic() {
	m := mixer.Classic()

	http.Handle("/", m.Handler(func(c mixer.Context) {
		fmt.Fprint(c.ResponseWriter(), "Hello, world")
	}))

	http.ListenAndServe(":3000", nil)
}

type Context struct {
	mixer.Context

	Logger *log.Logger
}

func ExampleNew() {
	// Inject any dependencies this context needs before it is passed to the
	// handler
	m := mixer.New(func(c mixer.Context) *Context {
		return &Context{
			Context: c,
			Logger:  log.New(os.Stdout, "[my logger]", 0),
		}
	})

	http.Handle("/", m.Handler(func(c *Context) {
		c.Logger.Println("Hello world")
		fmt.Fprintln(c.ResponseWriter(), "Hello, world")
	}))

	http.ListenAndServe(":3000", nil)
}

func TestClassic(t *testing.T) {
	m := mixer.Classic()

	h := m.Handler(func(c mixer.Context) {
		fmt.Fprint(c.ResponseWriter(), "Hello, world")
	})

	response := httptest.NewRecorder()
	h.ServeHTTP(response, (*http.Request)(nil))

	expect(t, response.Code, 200)
	expect(t, response.Body.String(), "Hello, world")
}

type CustomContext struct {
	mixer.Context
}

func TestNew(t *testing.T) {
	m := mixer.New(func(c mixer.Context) *CustomContext {
		return &CustomContext{
			Context: c,
		}
	})

	h := m.Handler(func(c *CustomContext) {
		fmt.Fprint(c.ResponseWriter(), "Hello, world")
	})

	response := httptest.NewRecorder()
	h.ServeHTTP(response, (*http.Request)(nil))

	expect(t, response.Code, 200)
	expect(t, response.Body.String(), "Hello, world")
}

func TestBeforeAfter(t *testing.T) {
	var result string

	m := mixer.Classic()
	m.Before(func(c mixer.Context) {
		result += "foo"
	})
	m.Before(func(c mixer.Context) {
		result += "bar"
	})
	m.After(func(c mixer.Context) {
		result += "baz"
	})

	h := func(c mixer.Context) {
		result += "bat"
	}

	response := httptest.NewRecorder()
	m.Handler(h).ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbaz")
}

func TestEarlyReturn(t *testing.T) {
	var result string

	m := mixer.Classic()
	m.Before(func(c mixer.Context) {
		result += "foo"
	})
	m.Before(func(c mixer.Context) {
		result += "bar"
		c.ResponseWriter().Write([]byte("Hello world"))
	})
	m.After(func(c mixer.Context) {
		result += "baz"
	})

	h := func(c mixer.Context) {
		result += "bat"
	}

	response := httptest.NewRecorder()
	m.Handler(h).ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbaz")
}

/* Test Helpers */
func expect(t *testing.T, a any, b any) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a any, b any) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
