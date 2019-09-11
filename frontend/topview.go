package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gowebapi/webapi"
	"github.com/gowebapi/webapi/dom/domcore"
	"github.com/gowebapi/webapi/html/htmlevent"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// TopView ...
type TopView struct {
	vecty.Core
	Cards       vecty.List   `vecty:"prop"`
	Index       int          `vecty:"prop"`
	Last        int          `vecty:"prop"`
	Controller  *Controller  `vecty:"prop"`
	SearchLight *SearchLight `vecty:"prop"`
	Cursor      string       `vecty:"prop"`
}

// AddPage ...
func (c *TopView) AddPage(content string) {
	c.Cards = append(c.Cards, &Card{State: "next", Content: content})
}

// OnHashChange ...
func (c *TopView) OnHashChange(event *domcore.Event) {
	ev := htmlevent.HashChangeEventFromJS(event.JSValue())
	oldID, err := strconv.Atoi(ParseHash(ev.OldURL()).Path)
	if err != nil {
		oldID = c.Index
	}
	newID, err := strconv.Atoi(ParseHash(ev.NewURL()).Path)
	if err != nil {
		newID = c.Index
	}
	log.Println(oldID, "->", newID)
	c.Index = newID
	c.Last = oldID
	vecty.Rerender(c)
	vecty.Rerender(c.Controller)
}

// Prev ...
func (c *TopView) Prev(event *vecty.Event) {
	if c.Controller.Prev == "#" {
		return
	}
	webapi.GetWindow().Location().SetHash(c.Controller.Prev)
}

// Next ...
func (c *TopView) Next(event *vecty.Event) {
	if elem := webapi.GetDocument().QuerySelector(".active .fragment"); elem != nil {
		cs := elem.ClassList()
		cs.Remove("fragment")
		cs.Add("appeared")
		if v, ok := c.Cards[c.Index-1].(vecty.Component); ok {
			vecty.Rerender(v)
		}
		return
	}
	if c.Controller.Next == "#" {
		return
	}
	webapi.GetWindow().Location().SetHash(c.Controller.Next)
}

// OnKeyDown ...
func (c *TopView) OnKeyDown(event *domcore.Event) {
	switch v := event.JSValue().Get("code").String(); v {
	case "ArrowLeft":
		log.Println("prev:", v)
		c.Prev(nil)
	case "ArrowRight":
		log.Println("next:", v)
		c.Next(nil)
	case "KeyF":
		webapi.GetDocument().Body().RequestFullscreen(nil)
	case "Space":
		c.SearchLight.Enabled = !c.SearchLight.Enabled
		vecty.Rerender(c.SearchLight)
	default:
		log.Println(v)
	}
}

// OnMouseMove ...
func (c *TopView) OnMouseMove(event *vecty.Event) {
	c.SearchLight.OnMouseMove(event)
	cursor := "inherit"
	if c.SearchLight.Active {
		cursor = "none"
	}
	if cursor != c.Cursor {
		c.Cursor = cursor
		c.Last = c.Index
		vecty.Rerender(c)
	}
}

// Render ...
func (c *TopView) Render() vecty.ComponentOrHTML {
	c.Index, _ = strconv.Atoi(GetURL().Path)
	if c.Last == 0 {
		c.Last = c.Index
	}
	vecty.SetTitle(fmt.Sprintf("Page(%d/%d)", c.Index, len(c.Cards)))
	for i, component := range c.Cards {
		i++
		card := component.(*Card)
		switch {
		case i < c.Index-1:
			card.State = "prev"
		case i == c.Index-1:
			switch {
			case i == c.Last:
				card.State = "forwardOut"
			default:
				card.State = "prev"
			}
		case i == c.Index:
			switch {
			case i-1 == c.Last:
				card.State = "forwardIn"
			case i+1 == c.Last:
				card.State = "reverseIn"
			default:
				card.State = "active"
			}
		case i == c.Index+1:
			switch {
			case i == c.Last:
				card.State = "reverseOut"
			default:
				card.State = "next"
			}
		case i > c.Index+1:
			card.State = "next"
		}
	}
	c.Controller.Update()
	return elem.Body(
		vecty.Markup(
			event.MouseDown(c.SearchLight.OnMouseDown).PreventDefault().StopPropagation(),
			event.MouseOut(c.SearchLight.OnMouseOut),
			event.MouseMove(c.OnMouseMove),
			event.Wheel(c.SearchLight.OnWheel),
			vecty.Style("cursor", c.Cursor),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("container"),
			),
			c.Cards,
		),
		c.Controller,
		c.SearchLight,
	)
}
