package main

import (
	"encoding/json"
	"github.com/nsf/termbox-go"
)

type Style struct {
	parent *Style
	bg     *termbox.Attribute
	fg     *termbox.Attribute
	childs map[string]*Style
}

func (s *Style) Bg() termbox.Attribute {
	for s != nil {
		if s.bg != nil {
			return *s.bg
		}
		s = s.parent
	}
	return termbox.ColorDefault
}

func (s *Style) Fg() termbox.Attribute {
	for s != nil {
		if s.fg != nil {
			return *s.fg
		}
		s = s.parent
	}
	return termbox.ColorDefault
}

func (s *Style) Sub(name ...string) *Style {
	if len(name) > 0 {
		c, ok := s.childs[name[0]]
		if ok {
			return c.Sub(name[1:]...)
		}
	}
	return s
}

func (s *Style) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil
	}
	if fg, ok := raw["FG"]; ok {
		delete(raw, "FG")
		f := new(termbox.Attribute)
		if err = json.Unmarshal(fg, f); err != nil {
			return err
		}
		s.fg = f
	}
	if bg, ok := raw["BG"]; ok {
		delete(raw, "BG")
		f := new(termbox.Attribute)
		if err = json.Unmarshal(bg, f); err != nil {
			return err
		}
		s.bg = f
	}

	for k, v := range raw {
		if s.childs == nil {
			s.childs = make(map[string]*Style)
		}
		sub := new(Style)
		if err := json.Unmarshal(v, sub); err != nil {
			return err
		}
		s.childs[k] = sub
		sub.parent = s
	}
	return nil
}

func StyleFromJson(js string) (*Style, error) {
	style := new(Style)
	err := json.Unmarshal([]byte(js), style)
	return style, err
}
