package main

import (
	"log"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gowebapi/webapi"
)

// Controller ...
type Controller struct {
	vecty.Core
	parent      *TopView
	CardCount   int
	Prev        string `vecty:"prop"`
	Next        string `vecty:"prop"`
	HasFragment bool   `vecty:"prop"`
}

func (c *Controller) prev() string {
	n, err := strconv.ParseInt(GetURL().Path, 10, 64)
	if err != nil {
		return ""
	}
	r := int(n - 1)
	if r < 1 {
		return ""
	}
	return strconv.Itoa(r)
}

func (c *Controller) next() string {
	n, err := strconv.ParseInt(GetURL().Path, 10, 64)
	if err != nil {
		return ""
	}
	r := int(n + 1)
	if r > c.CardCount {
		return ""
	}
	return strconv.Itoa(r)
}

// Update ...
func (c *Controller) Update() {
	c.Prev = "#" + c.prev()
	c.Next = "#" + c.next()
	c.HasFragment = webapi.GetDocument().QuerySelector(".active .fragment") != nil
	c.HasFragment = c.HasFragment || webapi.GetDocument().QuerySelector(".forwardIn .fragment") != nil
	log.Printf("%#v", *c)
}

// Render ...
func (c *Controller) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("controller"),
		),
		elem.Button(
			vecty.Markup(
				vecty.ClassMap{
					"btn":      true,
					"btn-link": true,
					"btn-lg":   true,
				},
				event.Click(c.parent.Prev).PreventDefault(),
			),
			vecty.Text("<"),
		),
		elem.Button(
			vecty.Markup(
				vecty.ClassMap{
					"btn":      true,
					"btn-link": true,
					"btn-lg":   true,
				},
				event.Click(c.parent.Next).PreventDefault(),
			),
			vecty.Text(">"),
		),
	)
}
