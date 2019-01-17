package cli

import "github.com/fatih/color"

type Theme struct {
	Title   *color.Color
	Error   *color.Color
	Keyword *color.Color
	Value   *color.Color
	Option  *color.Color
}

func NewDefaultTheme() *Theme {
	return &Theme{
		Title:   color.New(color.Bold, color.Underline),
		Error:   color.New(color.FgRed, color.Bold, color.Underline),
		Keyword: color.New(color.FgWhite, color.Underline),
		Value:   color.New(color.FgYellow),
		Option:  color.New(color.FgGreen),
	}
}
