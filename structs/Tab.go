package structs

import "html/template"

type Tab struct {
	Name    template.HTML
	GetUrl  string
	SubTabs []Tab
}
