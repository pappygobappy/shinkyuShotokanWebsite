package structs

type Tab struct {
	Name    string
	GetUrl  string
	SubTabs []Tab
}
