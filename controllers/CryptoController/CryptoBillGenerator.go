package controllers

import (
	"fmt"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func CryptoBill() {
	m := pdf.NewMaroto(consts.Portrait, consts.A5)
	m.SetPageMargins(40, 10, 40)

	buildingHeading(m)
	buildingList(m)
	m.OutputFileAndClose("/home/umut/goKripto/Pdfs/deneme.pdf")
	fmt.Println("a")
}

func buildingHeading(m pdf.Maroto) {
	m.RegisterHeader(func() {
		m.Row(25, func() {
			m.Col(12, func() {
				m.FileImage("Pdfs/Images/proxolab.jpg", props.Rect{
					Center: true,
				})
			})
		})
	})

	m.Row(0, func() {
		m.Col(0, func() {
			m.Text("Deneme123+%&/", props.Text{
				Top:   0,
				Style: consts.Bold,
				Align: consts.Center,
				Color: color.NewBlack(),
			})
		})
	})
}

func buildingList(m pdf.Maroto) {
	tableHeadings := []string{"Crypto Name", "Amount", "Price", "Deneme"}
	contents := [][]string{
		{"Bitcoin", "10", "$50000", "Sample1"},
		{"Ethereum", "5", "$3000", "Sample2"},
		{"Litecoin", "20", "$150", "Sample3"},
	}

	lightPurpleColor := getLightDark()
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Denemeeee", props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Courier,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	m.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{4, 7, 2, 4}, // Burada sütun genişliklerini ayarlayın
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{4, 7, 2, 4}, // Burada sütun genişliklerini ayarlayın
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

}

func getDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}

func getLightDark() color.Color {
	return color.Color{
		Red:   153,
		Green: 153,
		Blue:  153,
	}
}
