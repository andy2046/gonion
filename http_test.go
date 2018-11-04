package gonion

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func funcsEqual(f1, f2 interface{}) bool {
	val1 := reflect.ValueOf(f1)
	val2 := reflect.ValueOf(f2)
	return val1.Pointer() == val2.Pointer()
}

var testApp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test\n"))
})

func tagMiddleware(tag string) HTTPHandler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

func TestHTTP_Then_WithNoMiddleware(t *testing.T) {
	if !funcsEqual(NewHTTP().Then(testApp), testApp) {
		t.Error("Then does not work with no middleware")
	}
}

func TestHTTP_Then_NilAsDefaultServeMux(t *testing.T) {
	if NewHTTP().Then(nil) != http.DefaultServeMux {
		t.Error("Then does not replace nil with DefaultServeMux")
	}
}

func TestHTTP_Then_Handlers(t *testing.T) {
	t1 := tagMiddleware("t1\n")
	t2 := tagMiddleware("t2\n")
	t3 := tagMiddleware("t3\n")

	chained := NewHTTP(t1, t2, t3).Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)

	if w.Body.String() != "t1\nt2\nt3\ntest\n" {
		t.Error("Then does not order handlers correctly")
	}
}

func TestHTTP_Append_Handlers(t *testing.T) {
	chain := NewHTTP(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	newChain := chain.Append(tagMiddleware("t3\n"), tagMiddleware("t4\n"))

	if len(chain.handlers) != 2 {
		t.Error("chain should have 2 handlers")
	}
	if len(newChain.handlers) != 4 {
		t.Error("newChain should have 4 handlers")
	}

	chained := newChain.Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)

	if w.Body.String() != "t1\nt2\nt3\nt4\ntest\n" {
		t.Error("Append does not add handlers correctly")
	}
}

func TestHTTP_Extend_Handlers(t *testing.T) {
	chain1 := NewHTTP(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	chain2 := NewHTTP(tagMiddleware("t3\n"), tagMiddleware("t4\n"))
	newChain := chain1.Extend(chain2)

	if len(chain1.handlers) != 2 {
		t.Error("chain1 should contain 2 handlers")
	}
	if len(chain2.handlers) != 2 {
		t.Error("chain2 should contain 2 handlers")
	}
	if len(newChain.handlers) != 4 {
		t.Error("newChain should contain 4 handlers")
	}

	chained := newChain.Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)

	if w.Body.String() != "t1\nt2\nt3\nt4\ntest\n" {
		t.Error("Extend does not add handlers correctly")
	}
}
