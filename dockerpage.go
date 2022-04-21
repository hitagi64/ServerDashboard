package main

import "net/http"

type dockerData struct {
}

func (tp PageTemplates) dockerPageHandler(w http.ResponseWriter, r *http.Request){
	data := dockerData{

	}
	tp.runBasePage(w, "Docker", tp.docker, data)
}