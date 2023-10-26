package requests

import (
	"context"
	"fmt"

	"github.com/frankli0324/go-requests/request"
)

func ExampleCtxDo() {
	done, resp, err := CtxDo(context.Background(), request.Get("https://www.google.com", nil))
	if err != nil {
		panic(err)
	}
	defer done()
	fmt.Println(resp.Header.Get("Date"))
}
