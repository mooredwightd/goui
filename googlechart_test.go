package goui

import (
	"testing"
	"github.com/mooredwightd/gotestutil"
	"log"
)

var (
	gdt *GoogleDataTable
	tData = [][]GoogleDataItem{
		{{Value:"a"}, {Value:10}}, {{Value:"b"}, {Value:15}}, {{Value:"c"}, {Value:30}},
	}
)

func TestNewGoogleDataTable(t *testing.T) {
	gdt = NewGoogleDataTable()
	gotestutil.AssertNotNil(t, gdt, "Expected valid GoogleDataTable.")
}

func TestGoogleDataTable_AddColumns(t *testing.T) {
	c1 := GoogleColumn{Id:"x", Type:ColTypeString, Label: "XCol"}
	c2 := GoogleColumn{Id:"y", Type:ColTypeNumber, Label: "YCol"}

	gdt.AddColumns(c1, c2)
	gotestutil.AssertEqual(t, len(gdt.Cols), 2, "Expected 2 defined columns. Actual: %d.", len(gdt.Cols))
	gotestutil.AssertStringsEqual(t, gdt.Cols[0].Type, ColTypeString,
		"Expected string type for \"x\". Actual: %s.", gdt.Cols[0].Type)
	gotestutil.AssertStringsEqual(t, gdt.Cols[0].Id, "x",
		"Expected id type for \"y\". Actual: %s.", gdt.Cols[0].Id)
	gotestutil.AssertStringsEqual(t, gdt.Cols[1].Type, ColTypeNumber,
		"Expected string type for \"x\". Actual: %s.", gdt.Cols[1].Type)
	gotestutil.AssertStringsEqual(t, gdt.Cols[1].Id, "y",
		"Expected id type for \"y\". Actual: %s.", gdt.Cols[1].Id)
}

func TestGoogleDataTable_AddRows(t *testing.T) {

	// Iterate through list of data items
	for _, v := range tData {
		gdt.AddRows(CreateRow(v[0], v[1]))
	}
	gotestutil.AssertEqual(t, len(gdt.Rows), 3, "Expected 3 rows. Actual: %d.", len(gdt.Rows))
}

func TestGoogleDataTable_ToJSON(t *testing.T) {
	j, err := gdt.ToJSON()
	if err != nil {
		t.Fatalf("Error on ToJSON: %s.\n", err)
	}
	log.Printf("JSON: %+v.\n", j)
}
