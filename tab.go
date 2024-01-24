package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Tab struct {
	GetEndpoint string
	Content     string
}

func tabbar(tabs []Tab, selected int) (string, error) {
	if selected >= len(tabs) {
		log.Printf("Tried to select %d tab out of %d tabs", selected, len(tabs))
		return "", errors.New("Tab index out of range")
	}
	html := `<div class="tab-list" role="tablist">`
	for i, t := range tabs {
		html += fmt.Sprintf(`<button hx-get="%s" role="tab" aria-controls="tab-content"`, t.GetEndpoint)
		if i == selected {
			html += `class="selected" aria-selected="true"`
		}
		html += fmt.Sprintf(">%s</button>\n", t.Content)
	}

	html += "</div>"
	return html, nil
}

var mainTabBar []Tab = []Tab{
	{
		GetEndpoint: "/portal/tab1",
		Content:     "Tab 1",
	},
	{
		GetEndpoint: "/portal/tab2",
		Content:     "Tab 2",
	},
	{
		GetEndpoint: "/portal/tab3",
		Content:     "Tab 3",
	},
}

func tab1(w http.ResponseWriter, r *http.Request) {
	html, err := tabbar(mainTabBar, 0)
	if err != nil {
		return
	}

	html += `
<div id="tab-content" role="tabpanel" class="tab-content">
	TAB 1
</div>
`
	fmt.Fprint(w, html)
}
func tab2(w http.ResponseWriter, r *http.Request) {
	html, err := tabbar(mainTabBar, 1)
	if err != nil {
		return
	}

	html += `
<div id="tab-content" role="tabpanel" class="tab-content">
	TAB 2
</div>
`
	fmt.Fprint(w, html)
}
func tab3(w http.ResponseWriter, r *http.Request) {
	html, err := tabbar(mainTabBar, 2)
	if err != nil {
		return
	}

	html += `
<div id="tab-content" role="tabpanel" class="tab-content">
	TAB 3
</div>
`
	fmt.Fprint(w, html)
}
