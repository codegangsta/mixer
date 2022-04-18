package mixer

import "net/http"

type Context interface {
	Request() *http.Request
	SetRequest(*http.Request)

	ResponseWriter() http.ResponseWriter
	SetResponseWriter(http.ResponseWriter)
}

type context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *context) Request() *http.Request {
	return c.r
}

func (c *context) SetRequest(r *http.Request) {
	c.r = r
}

func (c *context) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *context) SetResponseWriter(w http.ResponseWriter) {
	c.w = w
}
