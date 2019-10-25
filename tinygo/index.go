package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

var (
	document = js.Global.Get("document")
)

func appendCSS(url string) {
	link := document.Call("createElement", "link")
	link.Set("rel", "stylesheet")
	link.Set("href", url)
	document.Get("head").Call("appendChild", link)
}

func appendJS(url string) {
	script := document.Call("createElement", "script")
	script.Set("src", url)
	document.Get("head").Call("appendChild", script)
}

func main() {
	// この時点でheadタグまで読み込み完了
	appendCSS("./assets/lib/css/reveal.css")
	appendCSS("./assets/lib/css/theme/sky.css")
	appendCSS("./assets/lib/css/zenburn.css")
	appendCSS("./assets/site.css")
	appendJS("./assets/lib/js/head.min.js")
	appendJS("./assets/lib/js/reveal.js")
	document.Call("addEventListener", "DOMContentLoaded", func(*js.Object) {
		// この時点でbodyとhtmlタグまで読み込み完了
		document.Set("title", "TinyGo")
		holder := document.Call("createElement", "div")
		holder.Set("className", "reveal")
		document.Get("body").Call("appendChild", holder)
		slides := document.Call("createElement", "div")
		slides.Set("className", "slides")
		holder.Call("appendChild", slides)
		section := document.Call("createElement", "section")
		section.Call("setAttribute", "data-markdown", "./contents.md?tm="+time.Now().String())
		section.Call("setAttribute", "data-separator", "^====\n")
		section.Call("setAttribute", "data-separator-vertical", "^- - -\n")
		section.Call("setAttribute", "data-separator-notes", "^Note:")
		section.Call("setAttribute", "data-charset", "utf-8")
		script := document.Call("createElement", "script")
		script.Call("setAttribute", "type", "text/template")
		section.Call("appendChild", script)
		//<pre style="display:none"><code></code></pre>
		pre := document.Call("createElement", "pre")
		pre.Call("setAttribute", "style", "display:none")
		pre.Call("appendChild", document.Call("createElement", "code"))
		section.Call("appendChild", pre)
		slides.Call("appendChild", section)

		js.Global.Set("onload", func(ev *js.Object) {
			// この時点でappendXXXしたcssとjsが読み込み完了
			log.Println("initializing")
			Reveal := js.Global.Get("Reveal")
			Reveal.Call("initialize", js.M{
				"width":           1336,
				"height":          768,
				"progress":        true,
				"center":          false,
				"history":         true,
				"mouseWheel":      false,
				"transition":      "fade",
				"transitionSpeed": "fast",
				"dependencies": js.S{
					js.M{"src": "./assets/plugin/markdown/marked.js"},
					js.M{"src": "./assets/plugin/markdown/markdown.js"},
					js.M{"src": "./assets/plugin/notes/notes.js"},
					js.M{"src": "./assets/plugin/highlight/highlight.js"},
				},
			})
			Reveal.Call("addEventListener", "ready", func(*js.Object) {
				fmt.Println("reveal ready")
				element := document.Call("getElementById", "searchlight")
				var spotlight bool
				toggel := func() {
					if spotlight {
						element.Get("style").Set("opacity", "0")
						document.Get("body").Get("style").Set("cursor", "auto")
					} else {
						element.Get("style").Set("opacity", "0.5")
						document.Get("body").Get("style").Set("cursor", "none")
					}
					spotlight = !spotlight

				}
				document.Get("body").Call("addEventListener", "mousemove", func(ev *js.Object) {
					element.Get("classList").Call("add", "on")
					element.Get("style").Set("margin-left", ev.Get("pageX").Int()-150)
					element.Get("style").Set("margin-top", ev.Get("pageY").Int()-150)
				})
				document.Get("body").Call("addEventListener", "mouseout", func(ev *js.Object) {
					element.Get("classList").Call("remove", "on")
				})
				document.Get("body").Call("addEventListener", "wheel", func(ev *js.Object) {
					if spotlight {
						// spotlight-onの時はwheelイベントを親に伝搬させない
						ev.Call("stopPropagation")
					}
				})
				document.Get("body").Call("addEventListener", "click", func(ev *js.Object) {
					if ev.Get("which").Int() == 2 {
						toggel()
					}
					ev.Call("preventDefault")
				}, true)
			})
		})
	})
}
