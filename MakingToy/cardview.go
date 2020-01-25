package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// Card ...
type Card struct {
	vecty.Core
	Display bool   `vecty:"prop"`
	State   string `vecty:"prop"`
	Content string
}

// OnAnimationEnd ...
func (c *Card) OnAnimationEnd(event *vecty.Event) {
	switch c.State {
	case "forwardOut":
		c.State = "prev"
		c.Display = false
	case "forwardIn", "reverseIn":
		c.State = "active"
	case "reverseOut":
		c.State = "next"
		c.Display = false
	}
	vecty.Rerender(c)
}

// Render ...
func (c *Card) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			event.AnimationEnd(c.OnAnimationEnd),
			vecty.Class("card", c.State),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-body"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("content"),
					vecty.UnsafeHTML(c.Content),
				),
			),
		),
	)
}
