package mocker

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMocker_Get(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	type Input struct {
		Foo string
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello")
	})

	router.GET("/with_header", func(context *gin.Context) {
		context.String(http.StatusOK, context.GetHeader("foo"))
	})

	router.GET("/with_body", func(context *gin.Context) {

		var input = Input{}

		if err := context.BindJSON(&input); err != nil {
			context.String(http.StatusOK, err.Error())
		}

		context.String(http.StatusOK, input.Foo)
	})

	router.GET("/get", func(context *gin.Context) {
		context.String(http.StatusOK, context.Request.Method)
	})

	router.POST("/method", func(context *gin.Context) {
		context.String(http.StatusOK, context.Request.Method)
	})

	router.PUT("/put", func(context *gin.Context) {
		context.String(http.StatusOK, context.Request.Method)
	})

	router.DELETE("/delete", func(context *gin.Context) {
		context.String(http.StatusOK, context.Request.Method)
	})

	router.PATCH("/patch", func(context *gin.Context) {
		context.String(http.StatusOK, context.Request.Method)
	})

	m := New(router)

	// request with nothing
	r1 := m.Get("/", []byte(""), nil)

	assert.Equal(t, http.StatusOK, r1.Code)
	assert.Equal(t, "hello", r1.Body.String())

	// request with header
	r2 := m.Get("/with_header", []byte(""), &Header{
		"foo": "bar",
	})

	assert.Equal(t, "bar", r2.Body.String())

	// request with body
	input := Input{
		Foo: "bar",
	}

	if b, err := json.Marshal(&input); err != nil {
		assert.Error(t, err, err.Error())
	} else {
		r3 := m.Get("/with_body", b, nil)

		assert.Equal(t, "bar", r3.Body.String())
	}

	// request with another method
	r4 := m.Post("/method", []byte(""), nil)

	assert.Equal(t, http.MethodPost, r4.Body.String())
}
