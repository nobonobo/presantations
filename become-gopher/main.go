package main

import (
	"log"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

var (
	Window   = js.Global()
	Document = js.Global().Get("document")
)

// TopView ...
type TopView struct {
	vecty.Core
}

// Render ...
func (c *TopView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Div(
			vecty.Markup(
				prop.ID("main"),
				vecty.Style("display", ""),
			),
			elem.Video(
				vecty.Markup(
					prop.ID("video"),
					vecty.Attribute("playsinline", ""),
					vecty.Style("display", "none"),
				),
			),
			elem.Canvas(
				vecty.Markup(
					prop.ID("output"),
				),
			),
		),
		elem.Script(
			vecty.Markup(
				prop.Src("js/camera.283d5d54.js?t=1"),
			),
		),
		elem.Script(
			vecty.Markup(
				prop.Src("js/CapsuleGeometry.js"),
			),
		),
		elem.Script(
			vecty.Markup(
				prop.Src("js/three.min.js"),
				event.Load(func(_ *vecty.Event) {
					SetupGopher()
				}),
			),
		),
	)
}

func main() {
	log.Println("start")
	meta := Document.Call("createElement", "meta")
	meta.Call("setAttribute", "name", "viewport")
	meta.Call("setAttribute", "content", "width=device-width,initial-scale=1")
	Document.Get("head").Call("append", meta)
	top := &TopView{}
	vecty.AddStylesheet("css/spectre.min.css")
	vecty.AddStylesheet("css/spectre-exp.min.css")
	vecty.AddStylesheet("css/spectre-icons.min.css")
	vecty.AddStylesheet("app.css")
	vecty.RenderBody(top)
	select {}
}
