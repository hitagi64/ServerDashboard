package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func (tp PageTemplates) webhandler(w http.ResponseWriter, r *http.Request){
	// Strip ?xx=xxx&aa=etc from path
	urlstrparts := strings.Split(r.URL.String(), "?")
	var urlstr string
	if (len(urlstrparts) < 1) {
		urlstr = r.URL.String()
	}else{
		urlstr = urlstrparts[0]
	}

	if (urlstr == "/") {
		tp.dashboardPageHandler(w, r)
	}else if (urlstr == "/dashboardRefreshData") {
		tp.dashboardRefreshPageHandler(w, r)
	}else if (urlstr == "/disks") {
		tp.diskPageHandler(w, r)
	}else{
		fmt.Fprintf(w, "404: %s", urlstr)
	}
}

func main() {
	// Logs cpu usage over a span of 31 seconds
	cl1.maxlen = 31
	cl1.waittime = time.Second * 1
	go cl1.cpuLoggingTask()

	tp := loadPageTemplates()

	// Register web request handlers
	http.HandleFunc("/", tp.webhandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files/"))))

	// Start server
	fmt.Printf("Starting dashboard at port %d\n", 5000)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}