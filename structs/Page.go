package structs

import "shinkyuShotokan/models"

type Page struct {
	PageName string
	Tabs     []Tab
	Classes  []models.Class
}
