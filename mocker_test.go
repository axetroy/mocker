package mocker_test

import (
	"github.com/axetroy/mocker"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAllMethod(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		router.Handle(method, "/foo", func(context *gin.Context) {
			context.Header("foo", "bar")
			context.String(http.StatusOK, "bar")
		})
	}

	m := mocker.New(router)

	for _, method := range methods {
		res := m.Request(method, "/foo", nil, nil)
		assert.Equal(t, "bar", res.Body.String())
		assert.Equal(t, "bar", res.Header().Get("foo"))
	}
}

func TestBasic(t *testing.T) {
	type route struct {
		Method   string
		Path     string
		Body     []byte
		Header   *mocker.Header
		Response string
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	routes := []route{
		{
			Method:   http.MethodGet,
			Path:     "/foo",
			Body:     nil,
			Header:   nil,
			Response: "bar",
		},
		{
			Method: http.MethodPost,
			Path:   "/foo",
			Body:   nil,
			Header: &mocker.Header{
				"foo": "bar",
			},
			Response: "bar",
		},
		{
			Method:   http.MethodPut,
			Path:     "/foo",
			Body:     []byte("foo"),
			Response: "bar",
		},
	}

	for _, r := range routes {
		func(r route) {
			router.Handle(r.Method, r.Path, func(context *gin.Context) {
				context.Header("X-Request-Path", r.Path)

				if r.Header != nil {
					for key, val := range *r.Header {
						assert.Equal(t, val, context.GetHeader(key))
					}
				}

				if len(r.Body) > 0 {
					data, _ := context.GetRawData()
					assert.Equal(t, r.Body, data)
				}

				context.String(http.StatusOK, r.Response)
			})
		}(r)
	}

	m := mocker.New(router)

	for _, r := range routes {
		res := m.Request(r.Method, r.Path, r.Body, r.Header)

		// response body validate
		assert.Equal(t, r.Response, res.Body.String())

		// response header validate
		assert.Equal(t, r.Path, res.Header().Get("X-Request-Path"))
	}

}
