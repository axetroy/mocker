package mocker_test

import (
	"github.com/axetroy/mocker"
	"net/http"
	"reflect"
	"testing"
)

type Route struct {
	Method   string
	Path     string
	Body     []byte
	Header   *mocker.Header
	Response string
}

type Handler struct {
	router []Route
}

func (h Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range h.router {
		if route.Path == req.URL.Path && route.Method == req.Method {
			if route.Header != nil {
				for key, value := range *route.Header {
					res.Header().Set(key, value)
				}
			}

			_, _ = res.Write(route.Body)
		}
	}
}

func TestAllMethod(t *testing.T) {
	m := mocker.New(Handler{})

	m.Head("/foo", nil, nil)
	m.Options("/foo", nil, nil)
	m.Get("/foo", nil, nil)
	m.Put("/foo", nil, nil)
	m.Post("/foo", nil, nil)
	m.Delete("/foo", nil, nil)
	m.Patch("/foo", nil, nil)
	m.Trace("/foo", nil, nil)
}

func TestMocker_Request(t *testing.T) {
	routes := []Route{
		{
			Method:   http.MethodGet,
			Path:     "/foo",
			Body:     []byte("bar"),
			Header:   nil,
			Response: "bar",
		},
		{
			Method: http.MethodPost,
			Path:   "/hello",
			Body:   []byte("world"),
			Header: &mocker.Header{
				"foo": "bar",
			},
			Response: "world",
		},
		{
			Method:   http.MethodPut,
			Path:     "/123",
			Body:     []byte("123"),
			Response: "123",
		},
	}

	handler := Handler{router: routes}

	type fields struct {
		router http.Handler
	}
	type args struct {
		method string
		path   string
		body   []byte
		header *mocker.Header
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "/foo",
			fields: fields{
				router: handler,
			},
			args: args{
				method: http.MethodGet,
				path:   "/foo",
				body:   []byte("bar"),
				header: nil,
			},
			want: []byte("bar"),
		},
		{
			name: "/hello",
			fields: fields{
				router: handler,
			},
			args: args{
				method: http.MethodPost,
				path:   "/hello",
				body:   []byte("world"),
				header: nil,
			},
			want: []byte("world"),
		},
		{
			name: "/123",
			fields: fields{
				router: handler,
			},
			args: args{
				method: http.MethodPut,
				path:   "/123",
				body:   []byte("123"),
				header: nil,
			},
			want: []byte("123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &mocker.Mocker{
				Router: tt.fields.router,
			}
			if got := c.Request(tt.args.method, tt.args.path, tt.args.body, tt.args.header); !reflect.DeepEqual(got.Body.Bytes(), tt.want) {
				t.Errorf("Request() = %v, want %v", got.Body.Bytes(), tt.want)
			}
		})
	}
}
