package structs

type Tab struct {
	Name     string
	GetUrl   string
	ScrollTo string
	SubTabs  []Tab
}
