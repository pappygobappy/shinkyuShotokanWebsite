package models

type Class struct {
	Name        string
	Description string
	Annotations []string
	Location    Location
	GetUrl      string
	StartAge    int
	EndAge      int
	Schedule    string
}
