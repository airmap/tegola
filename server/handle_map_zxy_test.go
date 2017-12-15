package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimfeld/httptreemux"
	"github.com/airmap/tegola/server"
)

func TestHandleMapZXY(t *testing.T) {
	//	setup a new provider
	testcases := []struct {
		handler      http.Handler
		uri          string
		uriPattern   string
		reqMethod    string
		expected     []byte
		expectedCode int
	}{
		{
			handler:      server.HandleMapZXY{},
			uri:          "/maps/test-map/1/2/3.pbf",
			uriPattern:   "/maps/:map_name/:z/:x/:y",
			reqMethod:    "GET",
			expectedCode: http.StatusOK,
		},
		{ // issue-163
			handler:      server.HandleMapZXY{},
			uri:          "/maps/test-map/-1/0/0.pbf",
			uriPattern:   "/maps/:map_name/:z/:x/:y",
			reqMethod:    "GET",
			expectedCode: http.StatusBadRequest,
			expected:     []byte("negative zoom levels are not allowed"),
		},
	}

	for i, test := range testcases {
		var err error

		//	setup a new router. this handles parsing our URL wildcards (i.e. :map_name, :z, :x, :y)
		router := httptreemux.New()
		//	setup a new router group
		group := router.NewGroup("/")
		group.UsingContext().Handler(test.reqMethod, test.uriPattern, server.HandleMapZXY{})

		r, err := http.NewRequest(test.reqMethod, test.uri, nil)
		if err != nil {
			t.Fatal(err)
			continue
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		if w.Code != test.expectedCode {
			t.Errorf("failed test %v. handler returned wrong status code: got (%v) expected (%v)", i, w.Code, test.expectedCode)
		}
		// Only try to decode as string for errors.
		if len(test.expected) > 0 && test.expectedCode >= 400 {
			wbody := strings.TrimSpace(w.Body.String())

			if string(test.expected) != wbody {
				t.Errorf("failed test %v. handler returned wrong body: got (%v) expected (%v)", i, wbody, test.expected)
			}
		}
	}
}
