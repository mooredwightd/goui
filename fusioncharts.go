package goui

type FusionDataItem struct {
	Label string        `json:"label"`
	Value string        `json:"value"`
}


type FusionChartData struct {
	// DOM or element id
	RenderAt   string                `json:"renderAt"`
	// Format values: "json","jsonurl", "csv", "xml", "xmlurl"
	DataFormat string                `json:"dataFormat"`
	DataSource string                `json:"dataSource"`
	// Used to render DataSource to JSON
	data       []FusionDataItem
}
