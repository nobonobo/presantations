package main

import (
	"log"
	"syscall/js"

	"github.com/gowebapi/webapi"
	"github.com/gowebapi/webapi/fetch"
)

func main() {
	webapi.GetWindow().Fetch(
		webapi.UnionFromJS(js.ValueOf("/wasm_exec.js")), nil,
	).Then(
		fetch.PromiseResponseOnFulfilledToJS(
			func(r *fetch.Response) {
				log.Println(r.StatusText(), r.Body())
			},
		),
		fetch.PromiseResponseOnRejectedToJS(
			func(err js.Value) {
				log.Println("error:", err)
			},
		),
	)
	select {}
}
