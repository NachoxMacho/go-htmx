package main

import (
	"fmt"
	"github.com/NachoxMacho/go-htmx/api"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func index(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	html := `
<html>
<head>
<script src="https://unpkg.com/htmx.org@1.9.10"></script>
<style>
.progress {
	height: 20px;
	margin-bottom: 20px;
	overflow: hidden;
	background-color: #f5f5f5;
	border-radius: 4px;
	box-shadow: inset 0 1px 2px rgba(0,0,0,.1);
}
.progress-bar {
	float: left;
	width: 0%;
	height: 100%;
	font-size: 12px;
	line-height: 20px;
	color: #fff;
	text-align: center;
	background-color: #337ab7;
	box-shadow 0 -1px 0 rgba(0,0,0,.15);
	transition: width .5s ease;
}

#tabs > .tab-list button {
  border: none;
  display: inline-block;
  padding: 5px 10px;
  cursor: pointer;
  background-color: transparent;
}

#tabs > .tab-list button.selected {
  background-color: #eee;
}
</style>
</head>
<body>
<button hx-get="/portal/boldHello" hx-swap="innerHTML" hx-target="#message-box">Get the time</button>
<div id="message-box">Press the Button</div>
<div id="naked-time" hx-get="/portal/nakedTime" hx-swap="outerHTML" hx-trigger="load">Shooting your shot</div>
<div id="tabs" hx-get="/portal/tab1" hx-trigger="load delay:100ms" hx-target="#tabs" hx-swap="innerHTML"></div>
<body>
</html>
	`
	w.Write([]byte(html))
}

func boldHello(w http.ResponseWriter, _ *http.Request) {
	t := time.Now()
	fmt.Fprintf(w, "<b>HELLO ITS %s</b>", t)
}

var timeOfNextChance time.Time

func randomTimer(w http.ResponseWriter, _ *http.Request) {
	random := rand.Int()
	if random%100 > 96 {
		fmt.Fprint(w, `<div id="naked-time">It's time!</div>`)
		return
	}
	futureTime := time.Now().Add(time.Second * 30)
	timeOfNextChance = futureTime
	fmt.Fprintf(w, `<div id="naked-time" hx-get="/portal/nakedTime" hx-swap="outerHTML" hx-trigger="load delay:30s">
	Failure, trying again at %d:%d:%d
	<div id="test" hx-get="/portal/randomTimer/progress" hx-trigger="every 500ms" hx-swap="innerHTML">
	<div class="progress"><div id="check-progress" class=progress-bar style="width:0%%"></div></div>
	</div>
	</div>
	</div>`, futureTime.Hour(), futureTime.Minute(), futureTime.Second())
}

func randomTimeProgress(w http.ResponseWriter, _ *http.Request) {
	now := time.Now()
	if now.After(timeOfNextChance) {
		fmt.Fprintf(w, `<div class="progress"><div id="check-progress" class=progress-bar style="width:100%%"></div></div>`)
	} else {
		percent := (30 - timeOfNextChance.Sub(now).Seconds()) / 30 * 100
		// log.Printf("Time Difference: %f", timeOfNextChance.Sub(now).Seconds())
		fmt.Fprintf(w, `<div class="progress"><div id="check-progress" class=progress-bar style="width:%d%%"></div></div>`, int(percent))
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Welcome to the Portal")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", root)
	mux.HandleFunc("/portal/", index)
	mux.HandleFunc("/portal/boldHello", boldHello)
	mux.HandleFunc("/portal/randomTimer", randomTimer)
	mux.HandleFunc("/portal/randomTimer/progress", randomTimeProgress)
	mux.HandleFunc("/portal/tab1", tab1)
	mux.HandleFunc("/portal/tab2", tab2)
	mux.HandleFunc("/portal/tab3", tab3)
	api.Mux("/api/", mux)

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
