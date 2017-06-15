package goui

import (
	"html/template"
	"net/http"
	"fmt"
)

const (
	PageData = "Data"
	PageTitle = "Title"
	PageNav = "Nav"
)

//
//
type UIPage struct {
	defaultTmpl string
	t           *template.Template
	PageData    map[string]interface{}
}

// Create a new page object with the provided title, and default template.
func NewPage(uic *UIContext, title string, defTmpl string) *UIPage{
	t, _ := uic.Templates().Clone()
	p := &UIPage{
		defaultTmpl: defTmpl,
		t: t,
		PageData: make(map[string]interface{}, 1),
	}
	p.AddPageData(map[string]interface{}{
		PageTitle: title,
	})

	return p
}

// Execute the page rendering with data. If a template name is provided, then that template is rendered.
// If no template name is provided, then it first checks to see if the page's default template is non empty.
// If there is no default page template, the default template for the configuration is used.
// The PageData field is used to render the template.
func (uip *UIPage) ExecuteTemplate(wr http.ResponseWriter, tmplName string) error {
	tmpl := tmplName
	if len(tmpl) == 0 {
		if len(uip.defaultTmpl) == 0 {
			tmpl = GetUIConfig().p.GetString(CfgHomepage)
		} else {
			tmpl = uip.defaultTmpl
		}
	}
	// Execute the page with the data
	if err := uip.t.ExecuteTemplate(wr, tmpl, uip.PageData); err != nil {
		return errorf("Error on Execute.", err)
	}
	return nil
}

// Addes a template to the tree. If the template already exists, it is replaced.
func (uip *UIPage) AddTemplates(tmpl...string) (error) {
	for _, t := range tmpl {
		_, err := uip.t.Parse(t)
		if err != nil {
			return errorf(fmt.Sprintf("goui: Error adding template to page %s.", uip.PageData["Title"]), err)
		}
	}
	return nil
}

func (uip *UIPage) SetPageTitle(s string) *UIPage {
	uip.PageData[PageTitle] = s
	return uip
}

func (uip *UIPage) SetPageData(v interface{}) *UIPage {
	uip.PageData[PageData] = v
	return uip
}

func (uip *UIPage) AddPageData(m map[string]interface{}) *UIPage {
	for k, v := range m {
		uip.PageData[k] = v
	}
	return uip
}

// Add a navigation object to the page.
func (uip *UIPage) AddNavigation(uio HTMLElementWriter) *UIPage {
	uip.PageData[PageNav] = uio
	return uip
}

// Retrive the navigation object. Templates use this to render navigation.
func (uip *UIPage) Navigation() HTMLElementWriter {
	return uip.PageData[PageNav].(HTMLElementWriter)
}
