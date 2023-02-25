module github.com/nobonobo/presantations/frontend

go 1.12

require (
	github.com/Depado/bfchroma v1.1.2
	github.com/alecthomas/chroma v0.6.4
	github.com/gopherjs/vecty v0.0.0-20190701174234-2b6fc20f8913
	github.com/gowebapi/webapi v0.0.0-20190711074835-2928b664209b
	github.com/mattn/go-isatty v0.0.8 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20160330182126-9a301d65acbb
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
	marwan.io/wasm-fetch v0.0.0-20190701140803-40ec60ab1603
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
