package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/opentracing/opentracing-go"

	addendpoint "github.com/peterbourgon/go-microservices/addsvc/pkg/endpoint"
	addhttp "github.com/peterbourgon/go-microservices/addsvc/pkg/http"
	addservice "github.com/peterbourgon/go-microservices/addsvc/pkg/service"
)

func TestWiring(t *testing.T) {
	svc := addservice.New(log.NewNopLogger(), discard.NewCounter(), discard.NewCounter())
	eps := addendpoint.New(svc, log.NewNopLogger(), discard.NewHistogram(), opentracing.GlobalTracer())
	mux := addhttp.NewHandler(context.Background(), eps, log.NewNopLogger(), opentracing.GlobalTracer())
	srv := httptest.NewServer(mux)
	defer srv.Close()

	for _, testcase := range []struct {
		method, url, body, want string
	}{
		{"GET", srv.URL + "/concat", `{"a":"1","b":"2"}`, `{"v":"12"}`},
		{"GET", srv.URL + "/sum", `{"a":1,"b":2}`, `{"v":3}`},
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
			t.Errorf("%s %s %s: want %q, have %q", testcase.method, testcase.url, testcase.body, want, have)
		}
	}
}
