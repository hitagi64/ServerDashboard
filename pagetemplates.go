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
	pt.basehtml = template.Must(template.ParseFiles("html/base.html.tpl"))
	pt.dashboard = template.Must(template.ParseFiles("html/dashboard.html.tpl"))
	pt.dashboardRefreshData = template.Must(template.ParseFiles("html/dashboardRefreshData.json.tpl"))
	pt.disks = template.Must(template.ParseFiles("html/disks.html.tpl"))

	return pt
}