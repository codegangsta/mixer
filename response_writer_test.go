package mixer_test

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codegangsta/mixer"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func TestResponseWriterBeforeWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	expect(t, rw.Status(), 0)
	expect(t, rw.Written(), false)
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello world")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 11)
	expect(t, rw.Written(), true)
}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))
	rw.Write([]byte("foo bar bat baz"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello worldfoo bar bat baz")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 26)
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
}

func TestResponseWriterWritingHeaderTwice(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)
	rw.WriteHeader(http.StatusInternalServerError)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
}

func TestResponseWriterHijack(t *testing.T) {
	hijackable := newHijackableResponse()
	rw := mixer.NewResponseWriter(hijackable)
	hijacker, ok := rw.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()
	if err != nil {
		t.Error(err)
	}
	expect(t, hijackable.Hijacked, true)
}

func TestResponseWriteHijackNotOK(t *testing.T) {
	hijackable := new(http.ResponseWriter)
	rw := mixer.NewResponseWriter(*hijackable)
	hijacker, ok := rw.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()

	refute(t, err, nil)
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	rw := mixer.NewResponseWriter(rec)
	closed := false
	notifier := rw.(http.CloseNotifier).CloseNotify()
	rec.close()
	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}
	expect(t, closed, true)
}

func TestResponseWriterNonCloseNotify(t *testing.T) {
	rw := mixer.NewResponseWriter(httptest.NewRecorder())
	_, ok := rw.(http.CloseNotifier)
	expect(t, ok, false)
}

func TestResponseWriterFlusher(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	_, ok := rw.(http.Flusher)
	expect(t, ok, true)
}

func TestResponseWriter_Flush_marksWritten(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := mixer.NewResponseWriter(rec)

	rw.Flush()
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Written(), true)
}
