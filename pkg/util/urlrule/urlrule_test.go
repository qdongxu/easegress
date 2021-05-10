package urlrule

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/megaease/easegateway/pkg/context"
	"github.com/megaease/easegateway/pkg/tracing"
)

func TestURLRULEMatch(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"GET",
			"POST",
		},
		URL: StringMatch{
			Prefix: "/",
		},
	}

	u.Init()
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/user/api", bytes.NewReader(buff))
	if err != nil {
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if !u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req is %#v, urlrule is %#v", ctx.Request(), u)
	}
}

func TestURLRegxMatch(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"GET",
			"POST",
		},
		URL: StringMatch{
			RegEx: "^\\/app\\/.+$",
		},
	}

	u.Init()
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/app/api", bytes.NewReader(buff))
	if err != nil {
		t.Error("create http failed, err: ", err)
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if !u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req path is %s, urlrule is regEx: %s", ctx.Request().Path(), u.URL.RegEx)
	}
}

func TestURLExactMatch(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"GET",
			"POST",
		},
		URL: StringMatch{
			Exact: "/app/v2/user",
		},
	}

	u.Init()
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/app/v2/user", bytes.NewReader(buff))
	if err != nil {
		t.Error("create http failed, err: ", err)
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if !u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req path is %s, urlrule is exact : %s", ctx.Request().Path(), u.URL.Exact)
	}
}

func TestURLExactNotMatch(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"GET",
			"POST",
		},
		URL: StringMatch{
			Exact: "/app/v2/user",
		},
	}

	u.Init()
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/app/v3/user", bytes.NewReader(buff))
	if err != nil {
		t.Error("create http failed, err: ", err)
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req path is %s, urlrule is exact : %s", ctx.Request().Path(), u.URL.Exact)
	}
}

func TestURLPrefixNotMatch(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"GET",
			"POST",
		},
		URL: StringMatch{
			Prefix: "/app/v3",
		},
	}

	u.Init()
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/app/v2/user", bytes.NewReader(buff))
	if err != nil {
		t.Error("create http failed, err: ", err)
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req path is %s, urlrule is exact : %s", ctx.Request().Path(), u.URL.Exact)
	}
}

func TestURLRULENoMatchMethod(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"DELETE",
		},
		URL: StringMatch{
			Prefix: "/",
		},
	}
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/user/api", bytes.NewReader(buff))
	if err != nil {
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req mehtod is %#v, urlrule method is required: %v", ctx.Request().Method(), u.Methods)
	}
}

func TestURLRULENoMatchURL(t *testing.T) {
	u := &URLRule{
		Methods: []string{
			"POST",
		},
		URL: StringMatch{
			Exact: "/user",
		},
	}
	buff := []byte("ok")
	w := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "http://localhost/user/api", bytes.NewReader(buff))
	if err != nil {
		return
	}

	ctx := context.New(w, request, tracing.NoopTracing, "")
	if u.Match(ctx.Request()) {
		t.Errorf("match HTTP URL and Method failed, req path is %s, urlrule path is required exact : %s", ctx.Request().Path(), u.URL.Exact)
	}
}