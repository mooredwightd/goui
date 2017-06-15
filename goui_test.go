package goui

import (
	"testing"
	"github.com/mooredwightd/gotestutil"
	"reflect"
	"strings"
)

var (
	testElement *UIObject
)

func TestNewElement(t *testing.T) {
	testElement = NewElement("test", "test1", "strong", "text")
	gotestutil.AssertNotNil(t, testElement, "Expected valid UI Element.")

	gotestutil.AssertNotEmptyString(t, testElement.contentType, "Expected non-empty contentType.")
	gotestutil.AssertStringsEqual(t, testElement.contentType, "test",
		"contentType does not match. Actual: %s", testElement.contentType)

	gotestutil.AssertNotEmptyString(t, testElement.id, "Expected non-empty id.")
	gotestutil.AssertStringsEqual(t, testElement.id, "test1", "id does not match. Actual: %s", testElement.id)

	gotestutil.AssertNotEmptyString(t, testElement.class, "Expected non-empty class field.")
	gotestutil.AssertStringsEqual(t, testElement.class, "strong",
		"class field does not match. Actual: %s", testElement.class)

	gotestutil.AssertNotEmptyString(t, testElement.text, "Expected non-empty text field.")
	gotestutil.AssertStringsEqual(t, testElement.text, "text",
		"text field does not match. Actual: %s", testElement.text)
}

func TestNewElementFromJSON(t *testing.T) {
	t.Skipped()
}

type testInterfaceObj struct {
	id string
}

func (ti *testInterfaceObj) SetId(id string) {
	ti.id = id
}
func (ti *testInterfaceObj) Id() string {
	return ti.id
}
func TestUIObject_GetHTMLElementWriter(t *testing.T) {

	t.Run("A1", func(t *testing.T) {
		dType := reflect.TypeOf(testElement.GetHTMLElementWriter())
		valid := dType.Implements(reflect.TypeOf(new(HTMLElementWriter)).Elem())
		gotestutil.AssertTrue(t, valid, "Expected valid HTMLElementWriter interface.")
	})

	t.Run("B1", func(t *testing.T) {
		var tio testInterfaceObj
		dType := reflect.TypeOf(tio)
		valid := dType.Implements(reflect.TypeOf(new(HTMLElementWriter)).Elem())
		gotestutil.AssertFalse(t, valid, "Expected invalid HTMLElementWriter interface.")
	})
}

func TestUIObject_SetText(t *testing.T) {
	testElement.SetText("set_text")
	gotestutil.AssertStringsEqual(t, testElement.text, "set_text", "Text does not match. Actual: %s", testElement.text)
}

func TestUIObject_Text(t *testing.T) {
	gotestutil.AssertStringsEqual(t, testElement.Text(), "set_text",
		"Text does not match. Actual: %s", testElement.text)
}

func TestUIObject_SetId(t *testing.T) {
	testElement.SetId("this_id")
	gotestutil.AssertStringsEqual(t, testElement.id, "this_id",
		"Id field does not match. Actual: %s", testElement.id)
}

func TestUIObject_Id(t *testing.T) {
	gotestutil.AssertStringsEqual(t, testElement.Id(), "this_id",
		"Id() does not match. Actual: %s", testElement.id)
}

func TestUIObject_Class(t *testing.T) {
	gotestutil.AssertStringsEqual(t, testElement.class, "strong",
		"class field does not match. Actual: \"%s\".", testElement.class)
}

func TestUIObject_AddCssClass(t *testing.T) {
	testElement.AddCssClass("weak")
	v := strings.Contains(testElement.class, "strong")
	gotestutil.AssertTrue(t, v, "Expected class \"strong\" in result.")
	v = strings.Contains(testElement.class, "weak")
	gotestutil.AssertTrue(t, v, "Expected class \"weak\" in result.")
}

func TestUIObject_RemoveCssClass(t *testing.T) {
	testElement.RemoveCssClass("weak")
	v := strings.Contains(testElement.class, "weak")
	gotestutil.AssertFalse(t, v, "Found class \"weak\". Expected only \"strong\". Actual: %s", testElement.class)

	testElement.RemoveCssClass("other_class")
	v = strings.Contains(testElement.class, "other_class")
	gotestutil.AssertFalse(t, v, "Found class \"other_class\". Expected only \"strong\". Actual: %s", testElement.class)

	v = strings.Contains(testElement.class, "strong")
	gotestutil.AssertTrue(t, v, "Expected class \"strong\" in result.")
}

func TestUIObject_AddAttribute(t *testing.T) {
	testElement.AddAttribute("href", "value1")
	value, found := testElement.attrs["href"]
	gotestutil.AssertTrue(t, found, "Expected attribute \"href\".")
	gotestutil.AssertStringsEqual(t, value, "value1", "Expected \"href\" attribute to have value \"value1\".")

	testElement.AddAttribute("testAttr1", "value1")
	_, found = testElement.attrs["testAttr1"]
	gotestutil.AssertTrue(t, found, "Expected attribute \"testAttr1\".")
}

func TestUIObject_GetAttribute(t *testing.T) {
	t.Run("A1", func(t *testing.T) {
		value := testElement.GetAttribute("href")
		gotestutil.AssertStringsEqual(t, value, "value1", "Expected \"href\" attribute to have value \"value1\".")
	})

	t.Run("B1", func(t *testing.T) {
		value := testElement.GetAttribute("bad_attribute_name")
		gotestutil.AssertEmptyString(t, value, "Expected empty string for attribute value.")
	})
}

func TestUIObject_RemoveAttribute(t *testing.T) {

	t.Run("A1", func(t *testing.T) {
		value, found := testElement.attrs["testAttr1"]
		gotestutil.AssertTrue(t, found, "Expected to find attribute \"testAttr1\".")
		gotestutil.AssertStringsEqual(t, value, "value1", "Expected \"testAttr1\" attribute with value \"value1\".")
	})

	t.Run("B1", func(t *testing.T) {
		testElement.RemoveAttribute("testAttr1")
		value, found := testElement.attrs["testAttr1"]
		gotestutil.AssertFalse(t, found, "Expected to find attribute \"testAttr1\".")
		gotestutil.AssertEmptyString(t, value, "Expected empty string for \"testAttr1\" attribute value.")
	})
}

func TestUIObject_Attributes(t *testing.T) {
	t.Run("A1", func(t *testing.T) {
		x := testElement.GetAttribute("href")
		gotestutil.AssertNotEmptyString(t, x, "Expected valid value for attribute \"href\".")
		gotestutil.AssertStringsEqual(t, x, "value1", "Expected attribute \"href\" with value \"value1\".")
	})

	t.Run("B1", func(t *testing.T) {
		x := testElement.GetAttribute("testAttr1")
		gotestutil.AssertEmptyString(t, x, "Expected valid value for attribute \"href\".")
	})
}

func TestUIObject_SetContentType(t *testing.T) {
	t.Run("A1", func(t *testing.T) {
		testElement.SetContentType("testtype")
		gotestutil.AssertNotEmptyString(t, testElement.contentType,
			"Expected valid content type for attribute. Actual: \"%s\".", testElement.contentType)
		gotestutil.AssertStringsEqual(t, testElement.contentType, "testtype",
			"Expected valid content  type \"testtype\". Actual: \"%s\".", testElement.contentType)
	})
}

func TestUIObject_ContentType(t *testing.T) {
	ct := testElement.ContentType()
	gotestutil.AssertNotEmptyString(t, ct, "Expected valid value for content type. Actual: \"%s\".", ct)
	gotestutil.AssertStringsEqual(t, testElement.contentType, ct,
		"Expected content type match. Actual: \"%s\".", ct)
}

func TestUIObject_AddChild(t *testing.T) {
	gotestutil.AssertEqual(t, len(testElement.children), 0, "Expected child count of zero (0). Actual: %d",
		len(testElement.children))
	ch := NewElement("testch", "ch1", "", "Child 1")
	testElement.AddChild(ch)
	gotestutil.AssertEqual(t, len(testElement.children), 1, "Expected child count of one (1). Actual: %d",
		len(testElement.children))
}

func TestUIObject_AddChildren(t *testing.T) {
	gotestutil.AssertEqual(t, len(testElement.children), 1, "Expected child count of one (1). Actual: %d",
		len(testElement.children))
	ch2 := NewElement("testch", "ch2", "", "Child 2")
	ch3 := NewElement("testch", "ch3", "", "Child 3")
	testElement.AddChildren([]HTMLElementWriter{ch2, ch3})
	gotestutil.AssertEqual(t, len(testElement.children), 3, "Expected child count of three (3). Actual: %d",
		len(testElement.children))
}

func TestUIObject_ChildCount(t *testing.T) {
	gotestutil.AssertEqual(t, testElement.ChildCount(), 3,
		"Expected three children. Actual: %d.", testElement.ChildCount())
}

func TestUIObject_GetChildById(t *testing.T) {
	t.Run("A1", func(t *testing.T) {
		ch := testElement.GetChildById("ch2")
		gotestutil.AssertNotNil(t, ch, "Expected valid HTMLElementWriter for id ch2.")
		gotestutil.AssertStringsEqual(t, ch.Id(), "ch2", "Expected id to be \"ch2\".")
		gotestutil.AssertStringsEqual(t, ch.ContentType(), "testch", "Expected conent type \"%s\" for id ch2.")
	})
	t.Run("B1", func(t *testing.T) {
		ch := testElement.GetChildById("dummy")
		gotestutil.AssertNil(t, ch, "Expected empty HTMLWlementWriter for id \"dummy\".")
	})
}

func TestUIObject_SearchChildrenById(t *testing.T) {
	t.Run("A1", func(t *testing.T) {
		ch3 := testElement.GetChildById("ch3")
		ch3.AddChild(NewElement("testch", "subch1", "", "Sub-Child 2"))
		subch := testElement.SearchChildrenById("subch1")
		gotestutil.AssertNotNil(t, subch, "Expected valid element for id \"subch1\".")
		gotestutil.AssertStringsEqual(t, subch.Id(), "subch1", "Expected id of \"subch1\". Actual: %s.", subch.Id())
	})
	t.Run("B1", func(t *testing.T) {
		subch := testElement.SearchChildrenById("dummy")
		gotestutil.AssertNil(t, subch, "Expected nil element for id \"dummy\".")
	})
}

func TestUIObject_ChildrenByOrder(t *testing.T) {
	chList := testElement.ChildrenByOrder()
	gotestutil.AssertEqual(t, len(chList), 3, "Expected slice of three (3) children. Actual: %d.",
		len(chList))
	gotestutil.AssertStringsEqual(t, chList[0].Id(), "ch1", "Expected id \"ch1\". Actual: %s.", chList[0].Id())
	gotestutil.AssertStringsEqual(t, chList[1].Id(), "ch2", "Expected id \"ch2\". Actual: %s.", chList[1].Id())
	gotestutil.AssertStringsEqual(t, chList[2].Id(), "ch3", "Expected id \"ch3\". Actual: %s.", chList[2].Id())
}