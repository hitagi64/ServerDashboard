package main

import "html/template"

type PageTemplates struct {
    basehtml *template.Template
	dashboard *template.Template
	dashboardRefreshData *template.Template
	disks *template.Template
}

func loadPageTemplates() PageTemplates {
	var pt PageTemplates
	pt.basehtml = template.Must(template.ParseFiles("html/base.gohtml"))
	pt.dashboard = template.Must(template.ParseFiles("html/dashboard.gohtml"))
	pt.dashboardRefreshData = template.Must(template.ParseFiles("html/dashboardRefreshData.gohtml"))
	pt.disks = template.Must(template.ParseFiles("html/disks.gohtml"))

	return pt
}