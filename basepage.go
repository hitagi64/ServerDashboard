package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

type basePageData struct {
	PageTitle string
	PageContent template.HTML
}

func (tp PageTemplates) runBasePage(w http.ResponseWriter, title string, t *template.Template, data interface{}) {
	var contenthtml = bytes.NewBufferString("")
	t.Execute(contenthtml, data)

	var pagedata basePageData
	pagedata.PageTitle = title
	pagedata.PageContent = template.HTML(contenthtml.String())
	err := tp.basehtml.Execute(w, pagedata); 
	if err != nil {
		fmt.Println("2:" + err.Error())
		http.Error(w, "500 Internal server error", 500)
		return;
	}
}