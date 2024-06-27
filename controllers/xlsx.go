package controllers

import (
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func (idb InDb) setStyle(xlsx *excelize.File, sheetName string, letters []string) {
	styles, err := xlsx.NewStyle(title_style)
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	xlsx.SetCellStyle(sheetName, "A1", letters[0], styles)

	styles, err = xlsx.NewStyle(header_style)
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	xlsx.SetCellStyle(sheetName, "A2", letters[1], styles)

	styles, err = xlsx.NewStyle(cell_style)
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	xlsx.SetCellStyle(sheetName, "A3", letters[2], styles)
}

var title_style = `{
	"font": {
		"size": 20,
		"color": "#000000"
	},
	"border":  [{
		"type": "left",
		"style": 2,
		"color": "#000000"
	}, {
		"type": "right",
		"style": 2,
		"color": "#000000"
	}, {
		"type": "top",
		"style": 2,
		"color": "#000000"
	}, {
		"type": "bottom",
		"style": 2,
		"color": "#000000"
	}],
	"alignment": {
		"horizontal": "center",
		"vertical": "center"
	}
}`

var header_style = `{
	"font": {
		"size": 12,
		"color": "#000000"
	},
	"border":  [{
		"type": "left",
		"style": 2,
		"color": "#000000"
	}, {
		"type": "right",
		"style": 2,
		"color": "#000000"
	}, {
		"type": "top",
		"style": 2,
		"color": "#000000"
	}, {
		"type": "bottom",
		"style": 2,
		"color": "#000000"
	}],
	"alignment": {
		"horizontal": "center",
		"vertical": "center"
	}
}`

var cell_style = `{
    "fill": {
        "type": "pattern",
        "color": ["#E0EBF5"]
    },
	"alignment": {
		"horizontal": "left",
		"vertical": "center"
	}
}`
