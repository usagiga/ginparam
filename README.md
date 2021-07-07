# ginparam

Query parameter binder for gin-gonic/gin


## Installation

```sh
$ go get github.com/usagiga/ginparam
```


## Example

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/usagiga/ginparam"
	"log"
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

func main() {
	ctx := gin.New()

	ctx.GET("/foo", func(ctx *gin.Context) {
		req := &Request{}
		err := ginparam.Read(ctx, req)
		if err != nil {
			log.Fatalf("Can't read params: %+v", err)
		}

		log.Println("StringVal: ", req.StringVal)
	})

	log.Println("Server started on :8080")
	err := ctx.Run(":8080")
	if err != nil {
		log.Fatalf("Something wrong in server: %+v", err)
    }
}
```

If there's no `query` struct tag, no value in specified environment keys or field applied `query:"-"`, ginparam will ignore it.

See also [example](./example) .


## Features

- Compatible with `xerrors`
- Auto type detection

### Supported types

- `int`
- `string`
- `bool`
- `slice` of them
- `struct` contained them

## Dependencies

- Go (1.16 or higher)
- [golang.org/x/xerrors](https://pkg.go.dev/golang.org/x/xerrors)


## License

MIT
