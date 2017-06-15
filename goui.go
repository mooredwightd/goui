package goui

import (
	"strings"
	"fmt"
	"math/rand"
	"html/template"
	"container/list"
	"encoding/json"
	"log"
)

// Used for template errors that bubble up.
type Error struct {
	Msg string
	Err error
}

func errorf(msg string, cause error) *Error {
	return &Error{Msg: msg, Err: cause}
}

func (e *Error) Error() string {
	return fmt.Sprintf("goui: %s. %s", e.Msg, e.Err.Error())
}

// General framework
//

const (
	ContentTypeLink string = "link"
	ContentTypeMenu string = "menu"
	ContentTypeImage string = "image"
	ContentTypeSeparator string = "separator"

	// Input types
	ContentInputButton string = "button_input"
	ContentInputCheckbox string = "checkbox_input"
	ContentInputColor string = "color_input"
	ContentInputDate string = "date_input"
	ContentInputDateTimeLoc string = "datetimeloc_input"
	ContentInputEmail string = "email_input"
	ContentInputFile string = "file_input"
	ContentInputHidden string = "hidden_input"
	ContentInputImage string = "image_input"
	ContentInputMonth string = "month_input"
	ContentInputNumber string = "number_input"
	ContentInputPassword string = "password_input"
	ContentInputRadio string = "radio_input"
	ContentInputRange string = "range_input"
	ContentInputReset string = "reset_input"
	ContentInputSearch string = "search_input"
	ContentInputSubmit string = "submit_input"
	ContentInputTel string = "tel_input"
	ContentInputText string = "text_input"
	ContentInputTime string = "time_input"
	ContentInputUrl string = "url_input"
	ContentInputWeek string = "week_input"
)

type AttributeInterface interface {
	Attributes() template.HTML
	GetAttribute(k string) string
	AddAttribute(attName string, attValue string) (HTMLElementWriter)
	RemoveAttribute(attName string) (HTMLElementWriter)
}

type ClassInterface interface {
	Class() string
	AddCssClass(className string) (HTMLElementWriter)
	RemoveCssClass(className string) (HTMLElementWriter)
}

type IdInterface interface {
	SetId(i string) HTMLElementWriter
	Id() string
}

type TextInterface interface {
	SetText(t string) HTMLElementWriter
	Text() string
}

type ChildrenInterface interface {
	Children() []HTMLElementWriter
	ChildrenByOrder() []HTMLElementWriter
	ChildCount() int
	AddChild(ui HTMLElementWriter) HTMLElementWriter
	SetChildOrder(thisId, beforeId string) HTMLElementWriter
	GetChildById(id string) HTMLElementWriter
	SearchChildrenById(id string) HTMLElementWriter
}

type HTMLElementWriter interface {
	GetHTMLElementWriter() HTMLElementWriter
	SetContentType(s string) HTMLElementWriter
	ContentType() string
	GetContentByType(t string) []HTMLElementWriter
	ClassInterface
	AttributeInterface
	IdInterface
	TextInterface
	ChildrenInterface
}

func init() {
	uic := GetUIConfig()
	uic.LoadTemplates("")
}

// Attributes allows a template to include ad-hoc element attributes for template.Execute()
// An attribute can be any name=value pairs.
type AttributeMap map[string]string

// Render the attributes into a string suitable for use inside an HTML element tag in a template.
// The template call call the String() function, e.e. {{Attr.String}}
func (attr AttributeMap) String() template.HTML {
	x := ""
	for k, v := range attr {
		x = x + " " + k + "=\"" + v + "\""
	}
	return template.HTML(x)
}

// Structure for UIObjects via JSON
// Example: `{"text":"A", "id":"button1", "type":"button", "class":"btn-primary", "attributes":{"href":"/"}}`*/
type elementStruct struct {
	// Text is typically for labels, or can be used in templates for various means.
	Text       string `json:"text"`
	// Id is the element id on the page.
	Id         string `json:"id"`
	// The elemen type. This can be one of the ContentType* or ContentInput* values or arbitrary. Used in templates.
	Etype      string `json:"type"`
	// A class name that can be added to the item in the rendered template, and progammatically changed.
	ClassName  string `json:"class"`
	// Additional HTML tag attributes. Used in templates.
	Attributes AttributeMap `json:"attributes"`
	// Child elements of this element. E.g. items in a menu. content elements in a composite panel element.
	Children   []elementStruct `json:"children"`
}

// To be used within a template
type UIObject struct {
	// Text, ID, Class are separate from Attrs to make templates easier, and are the most
	// common standard elements
	text        string
	id          string
	class       string
	attrs       AttributeMap
	// Settable element type name. This can be used in methods or templates
	contentType string
	// For a composite object, these are the children objects. The index is the "id"
	children    map[string]*UIObject
	// Ordered list of children. Order set by the when AddChild() is called, or SetOrder()
	childOrder  *list.List
}

func NewElement(cType string, id string, className string, text string) *UIObject {
	h := &UIObject{contentType: cType,
		id: id,
		class: className,
		text: text,
		children: make(map[string]*UIObject, 1),
		attrs: make(AttributeMap, 1),
		childOrder: list.New(),
	}
	return h
}

// JSON format for creating an element.
// Format is a JSON object with fields: "text", "id", "class", "type", and "attributes".
// Attributes is an AttributeMap (map[string]string) of additional HTML attributes.
//
// The "type" attribute can be one of the ContentType* or ContentInput* values, or another arbitrary values. The type
// is used by templates in rendering.
func NewElementFromJSON(s string) (*UIObject, error) {
	var eBuf elementStruct
	err := json.Unmarshal([]byte(s), &eBuf)
	if err != nil {
		log.Printf("goui.NewElementFromJSON: %s", err)
		return nil, err
	}

	uio := NewElement(eBuf.Etype, eBuf.Id, eBuf.ClassName, eBuf.Text).AddAttributeMap(eBuf.Attributes)

	for _, v := range eBuf.Children {
		cuio := NewElement(v.Etype, v.Id, v.ClassName, v.Text).AddAttributeMap(v.Attributes)
		uio.AddChild(cuio)
	}
	return uio.(*UIObject), nil
}

// Returns the HTMLElementWriter interface from a *UIObject.
func (he *UIObject) GetHTMLElementWriter() HTMLElementWriter {
	return HTMLElementWriter(he)
}

// Set the Text field of the object.
// Implements Text interface
func (he *UIObject) SetText(t string) HTMLElementWriter {
	he.text = t
	return HTMLElementWriter(he)
}

// Retrieves the text field of the object. Templates access this data via {{.Text}} pipeline
// Implements Text interface
func (he UIObject) Text() string {
	return he.text
}

// Sets the HTML "id" attribute. This is used in templates.
// Implements Id interface
func (he *UIObject) SetId(i string) HTMLElementWriter {
	he.id = i
	return he
}


// Retrievs the Id field of the object. Templates access this field via the {{.Id}} pipeline.
// Implements Id interface
func (he UIObject) Id() string {
	return he.id
}

// Returns the current CSS classes for the element. Templates access this field via the {{.Class}} pipeline.
// Implements the Class interface
func (he UIObject) Class() string {
	return he.class
}

// Add a CSS class to an object.
// Implements the Class interface
func (he *UIObject) AddCssClass(className string) (HTMLElementWriter) {
	if !strings.Contains(he.class, className) {
		he.class = strings.TrimSpace(he.class) + " " + strings.TrimSpace(className)
	}
	return he
}

// Removes a CSS class from an object. This removes all instances of the className string.
// Implements the Class interface
func (he *UIObject) RemoveCssClass(className string) (HTMLElementWriter) {
	he.class = strings.TrimSpace(strings.Replace(he.class, strings.TrimSpace(className), "", -1))
	return he
}

// Returns a string of the attribute name/value pairs.
// Templates access this value using the {{.Attributes}} pipeline.
//     ["dir":"ltr", "data-toggle":"f1", "draggable":"true"] returns
//     `dir="ltr" data-toggle="f1" draggable="true"`
// Implements the Class interface
func (he *UIObject) Attributes() template.HTML {
	return he.attrs.String()
}

// Retrieve a HTML attribute from an object.
// Template access the values using the {{.GetAttribute "attrName"}} pipeline.
// Example:: <a href='{{.GetAttribute "href"}}'>text</a>
// Implements Attribute interface
func (he *UIObject) GetAttribute(attrName string) string {
	if x, ok := he.attrs[attrName]; ok {
		return x
	}
	return ""
}

// Add a HTML attribute to an object. If the attribute already exists, it is replaced.
// Templates can retrieve an attribute using .GetAttribute pipeline (See GetAttribute)
// Implements Attribute interface
func (he *UIObject) AddAttribute(attName string, attValue string) (HTMLElementWriter) {
	he.attrs[attName] = attValue
	return he
}

// Add attributes from an AttributeMap.
func (he *UIObject) AddAttributeMap(a AttributeMap) (HTMLElementWriter) {
	for k, v := range a {
		he.attrs[k] = v
	}
	return he
}

// Remove a CSS class style from the object
// Implements Attribute interface
func (he *UIObject) RemoveAttribute(attName string) HTMLElementWriter {
	delete(he.attrs, attName)
	return he
}

// Set the content type. This is any arbitrary text that can be used to establish a content type
// Implements the HTMLElementWriter interface
func (he *UIObject) SetContentType(s string) HTMLElementWriter {
	he.contentType = s
	return he
}

// Retrieve the content type. Templates access this feild via the {{.ContentType}} pipeline.
// Implements the HTMLElementWriter interface
func (he *UIObject) ContentType() string {
	return he.contentType
}

// Return all items that have the specific content type t. Items are returned in the order they are added.
// Implements the HTMLElementWriter interface
func (he *UIObject) GetContentByType(t string) []HTMLElementWriter {
	var x []HTMLElementWriter
	for l := he.childOrder.Front(); l != nil; l = l.Next() {
		if l.Value.(*UIObject).contentType == t {
			x = append(x, l.Value.(*UIObject))
		}
	}
	return x
}

// Return an unordered list of children objects. Templates access the slice of children via
// .Children pipeline.
// Example: {{range .Children}}{{printf "Object id=%q" .Id}}{{end}}
// Example: {{range $index, $element := .Children}}{{printf "Object #%q id=%q" .index .element.Id}}{{end}}
// Implements the Child interface.
func (he *UIObject) Children() []HTMLElementWriter {
	var x []HTMLElementWriter
	for _, v := range he.children {
		x = append(x, v)
	}
	return x
}

// Return and ordered list of children objects.
// Initial order is determined by the sequence they are added. The order of a specific object
// can be set with SetOrder()
// Example: {{range .ChildrenByOrder}}{{printf "Object id=%q" .Id}}{{end}}
// Implements the Child interface.
func (he *UIObject) ChildrenByOrder() []HTMLElementWriter {
	var x []HTMLElementWriter
	for l := he.childOrder.Front(); l != nil; l = l.Next() {
		x = append(x, l.Value.(*UIObject))
	}
	return x
}

// Returns the number of children data objects. Useful in templates where you need to know the
// number of items.
// Example: {{printf "Number of items: %q" .ChildCount}}
// Implements the Child Interface.
func (he *UIObject) ChildCount() int {
	return len(he.children)
}

// Lookup a child object by its Id.
// Implements the Child Interface.
func (uio *UIObject) GetChildById(id string) HTMLElementWriter {
	for _, k := range uio.children {
		if k.Id() == id {
			return k
		}
	}
	return nil
}

// Traverse the object and all children to find an descendant by id.
// Implements the Child Interface.
func (uio *UIObject) SearchChildrenById(id string) HTMLElementWriter {
	var x = uio.GetHTMLElementWriter()
	if x.Id() == id {
		return x
	}
	// Seach each child
	for _, c := range uio.children {
		// No match, try a deep search
		if x = c.SearchChildrenById(id); x != nil {
			return x
		}
	}

	return nil
}

// Add a child object to a parent object. This creates a hierarchy of data objects.
// If a UI control "panel" is a content block of header/title, text, and date,
// The panel object has three children, one child for each.
// Implements the Child Interface.
func (he *UIObject) AddChild(ui HTMLElementWriter) HTMLElementWriter {
	if len(ui.Id()) == 0 {
		ui.SetId(generateId(ui.ContentType(), 100))
	}
	he.children[ui.Id()] = ui.(*UIObject)
	he.childOrder.PushBack(ui)
	return he
}

// Add multiple children
// Implements the Child Interface.
func (he *UIObject) AddChildren(ui []HTMLElementWriter) HTMLElementWriter {
	for _, v := range ui {
		he.AddChild(v)
	}
	return he
}

// Set the sequence (order) of an individual object by Id before a specified object.
// For example, SetOrder("itemX", "item2") will set ItemX before Item2
// Implements the Child Interface.
func (he *UIObject) SetChildOrder(thisId, beforeId string) HTMLElementWriter {
	var obj, mark *list.Element
	x := he.childOrder.Front()
	for ; x != nil || (obj != nil && mark != nil); x = x.Next() {
		if x.Value.(HTMLElementWriter).Id() == thisId {
			obj = x
		}
		if x.Value.(HTMLElementWriter).Id() == beforeId {
			mark = x
		}
	}
	if obj != nil && mark != nil {
		he.childOrder.Remove(x)
		he.childOrder.InsertBefore(x.Value, mark)
	}
	return he
}

// Generate a unique Id based on a prefix, affixing a random number.
// The max parameter determines the bounded range [0, max)
func generateId(prefix string, max int) string {
	return fmt.Sprintf("%s#%d", prefix, rand.Intn(max))
}

func StripWhitespace(s string) string {
	return strings.Trim(s, "\n\t")
}