package goui

import (
	"encoding/json"
	"log"
)

// Google Charts - Chart generation

type GoogleChart struct {
	ChartType    string
	Title        string
	Id           string
	ChartData    GoogleDataTable
	ChartOptions map[string]interface{}
}

var (
	chartInitTmpl = `
	{{/* the JSFunctions pipeline is an array of function names */}}
	<script>
	{{range JSFunctions}}
	 google.charts.setOnLoadCallback({{.}});
	{{end}}
	</script>`

	pieChartTmpl = `
	{{/* the ChartFunctions pipeline is an array of GoogleCharts */}}
	 <script type="text/javascript">
	 {{ range Charts}}
            function {{.Title}} () {
                var data = new google.visualization.DataTable({{.Data}});
                var options = {title: .Title};
                var chart = new google.visualization.{{.ChartType}}(document.getElementById({{.Id}}));
                chart.draw(data, options);
            }
         {{end}}
        </script>`
)

//
func NewGoogleChart(title, cType, id string) *GoogleChart {
	gc := &GoogleChart{
		Title: title,
		Id: id,
		ChartOptions: make(map[string]interface{}, 1),
	}
	switch cType {
	case "AreaChart": fallthrough
	case "BarChart": fallthrough
	case "BubbleChart": fallthrough
	case "Calendar": fallthrough
	case "ComboChart": fallthrough
	case "ColumnChart":fallthrough
	case "LineChart": fallthrough
	case "PieChart": fallthrough
	case "ScatterChart": fallthrough
	case "Table":
		gc.ChartType = cType
	default: gc.ChartType = "LineChart"
	}

	gc.ChartOptions["title"] = title

	return gc
}

func (gc *GoogleChart) SetData(gdt GoogleDataTable) *GoogleChart {
	gc.ChartData = gdt
	return gc
}

func (gc *GoogleChart) Data() string {
	jStr, err := gc.ChartData.ToJSON()
	if err != nil {
		log.Printf("Error converting chart data to JSON. %s", err)
		return `{"error": "1", "message": "Invalid chart data."}`
	}
	return jStr
}

func (gc *GoogleChart) SetOption(key string, value interface{}) *GoogleChart {
	gc.ChartOptions[key] = value
	return gc
}

func (gc *GoogleChart) Options() string {
	s, err := json.Marshal(gc.ChartOptions)
	if err != nil {
		log.Printf("GoogleChart.Options() JSON error: %s.\n", err)
		return ""
	}
	return string(s)
}

// Google Charts - data structure
// structures and methods for creating a google google.visualization.DataTable().
// The method ToJSON() renders JSON suitable for the DataTable contructor.
//
const (
	ColTYpeBoolean = "boolean"
	ColTypeString = "string"
	ColTypeNumber = "number"
	ColTypeDate = "date"
	ColTypeDateTime = "datetime"
	ColTypeTime = "timeofday"

	// @see: https://google-developers.appspot.com/chart/interactive/docs/roles
	RoleDomain = "domain"				// data point (hAxis)
	RoleData = "data"				// data point (yAxis)
	RoleAnnotation = "annotation"			// null or string
	RoleAnnotationText = "annotationText"		// null or string
	RoleCertainty = "certainty"
	RoleEmphasis = "emphasis"			// boolean
	RoleInterval = "interval"			// Must be in column pairs, with low/high values
	RoleScope = "scope"				// boolean
	RoleStyle = "style"				// null or string, CSS style format, single value or object
	RoleTooltip = "tooltip"				// string (default data point value)
)

// Reference: https://google-developers.appspot.com/chart/interactive/docs/reference#dataparam
type GoogleColumn struct {
	Id      string `json:"id,omitempty"`      //[Optional]
	Label   string `json:"label,omitempty"`   //[Optional]
	Type    string `json:"type"`              // values: boolean, number, string, date, datetime, timeofday or a role
	Pattern string `json:"pattern,omitempty"` //[Optional] String pattern that was used by a data source to format
	P       string `json:"p,omitempty"`       //[Optional] An object that is a map of custom values applied to the cell.
}
type GoogleDataItem struct {
	Value interface{} `json:"v"`           //The cell value. The data type should match the column data type.
	F     string      `json:"f,omitempty"` //[Optional] A string version of the v value, formatted for display.
	P     string      `json:"p,omitempty"` //[Optional] An object that is a map of custom values applied to the cell.
}
type GoogleRow struct {
	C []GoogleDataItem  `json:"c"`
	P map[string]string `json:"p,omitempty"` //[Optional] The table-level p property is a map of custom values applied to the whole DataTable.
}

type GoogleDataTable struct {
	Cols []GoogleColumn  `json:"cols"`
	Rows []GoogleRow     `json:"rows"`
}

// Create a new row
func NewGoogleDataTable() *GoogleDataTable {
	gdt := &GoogleDataTable{}
	return gdt
}

func (gdt *GoogleDataTable)AddColumns(c...GoogleColumn) *GoogleDataTable {
	for _, v := range c {
		gdt.Cols = append(gdt.Cols, v)
	}
	return gdt
}

func (gdt *GoogleDataTable)AddRows(r...GoogleRow) *GoogleDataTable {
	for _, v := range r {
		gdt.Rows = append(gdt.Rows, v)
	}
	return gdt
}

func (gdt *GoogleDataTable) ToJSON() (string, error) {
	j, err := json.Marshal(gdt)
	return string(j), err
}

func CreateRow(di...GoogleDataItem) GoogleRow {
	gr := GoogleRow{}
	for _, v := range di {
		gr.C = append(gr.C, v)
	}
	return gr
}


