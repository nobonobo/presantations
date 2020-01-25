package main

import (
	"bytes"
	"log"
	"net/url"
	"time"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/gopherjs/vecty"
	"github.com/gowebapi/webapi"
	"github.com/gowebapi/webapi/dom"
	"github.com/gowebapi/webapi/dom/domcore"
	"github.com/vincent-petithory/dataurl"
	bf "gopkg.in/russross/blackfriday.v2"
	fetch "marwan.io/wasm-fetch"
)

var (
	exts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
		bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
		bf.DefinitionLists | bf.Footnotes

	flags = bf.UseXHTML | bf.HrefTargetBlank | bf.Smartypants |
		bf.SmartypantsFractions | bf.SmartypantsDashes | bf.SmartypantsLatexDashes
)

func makeContent(b []byte, r *bfchroma.Renderer) string {
	unsafeHTML := bf.Run(
		b,
		bf.WithRenderer(r),
		bf.WithExtensions(exts),
	)
	return string(unsafeHTML)
}

// ParseHash ...
func ParseHash(s string) *url.URL {
	orig, _ := url.Parse(s)
	u, _ := url.Parse(orig.Fragment)
	if u.Path == "" {
		u.Path = "1"
	}
	return u
}

// GetURL ...
func GetURL() *url.URL {
	return ParseHash(webapi.GetWindow().Location().Hash())
}

func main() {
	log.Println("start")
	meta := webapi.GetDocument().CreateElement("meta", nil)
	meta.SetAttribute("name", "viewport")
	meta.SetAttribute("content", "width=device-width,initial-scale=1")
	webapi.GetDocument().Head().Append(dom.UnionFromJS(meta.JSValue()))
	top := &TopView{}
	webapi.GetWindow().AddEventListener(
		"hashchange", domcore.NewEventListenerFunc(top.OnHashChange), nil,
	)
	webapi.GetWindow().AddEventListener(
		"keydown", domcore.NewEventListenerFunc(top.OnKeyDown), nil,
	)
	vecty.AddStylesheet("css/spectre.min.css")
	vecty.AddStylesheet("css/spectre-exp.min.css")
	vecty.AddStylesheet("css/spectre-icons.min.css")
	vecty.AddStylesheet("app.css")
	renderer := bfchroma.NewRenderer(
		bfchroma.WithoutAutodetect(),
		bfchroma.Style("monokai"),
		bfchroma.ChromaOptions(html.WithClasses()),
		bfchroma.Extend(
			bf.NewHTMLRenderer(bf.HTMLRendererParameters{
				Flags: flags,
			}),
		),
	)
	w := new(bytes.Buffer)
	if err := renderer.Formatter.WriteCSS(w, renderer.Style); err != nil {
		log.Println(err)
	}
	data := dataurl.New(w.Bytes(), "text/css")
	vecty.AddStylesheet(data.String())
	uri := "contents.md?" + time.Now().Format(time.RFC3339Nano)
	resp, err := fetch.Fetch(uri, &fetch.Opts{})
	if err != nil {
		log.Println(err)
	}
	for _, chunk := range bytes.Split(resp.Body, []byte("\n====\n")) {
		chunk = bytes.Trim(chunk, "\r\n\t ")
		top.AddPage(makeContent(chunk, renderer))
	}
	top.Controller = &Controller{parent: top, CardCount: len(top.Cards)}
	top.SearchLight = &SearchLight{}
	top.Controller.Update()
	vecty.RenderBody(top)
	select{}
}
