package main

import (
	"github.com/gin-gonic/gin"
	"github.com/usagiga/ginparam"
	"log"
	"net/http"
)

type Request struct {
	// These fields will be loaded
	BoolVal   bool   `query:"bool_val"`
	StringVal string `query:"string_val"`
	IntVal    int    `query:"bool_val"`

	// These fields will be ignored
	IgnoredVal1 string `query:"-"`
	IgnoredVal2 string
}

// Usage:
// Run main() and type in your terminal as below
//
// ```sh
// $ curl -X GET 'http://localhost:8080/foo?string_val=abc'
// ```
func main() {
	ctx := gin.New()

	ctx.GET("/foo", func(ctx *gin.Context) {
		// Read query parameters
		req := &Request{}
		err := ginparam.Read(ctx, req)
		if err != nil {
			log.Fatalf("Can't read params: %+v", err)
		}

		// Dump query parameters
		ctx.JSON(http.StatusOK, req)
	})

	log.Println("Server started on :8080")
	err := ctx.Run(":8080")
	if err != nil {
		log.Fatalf("Something wrong in server: %+v", err)
	}
}
