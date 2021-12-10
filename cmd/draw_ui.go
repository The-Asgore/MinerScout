package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func textList() (l *widgets.List) {
	l = widgets.NewList()
	l.Title = "Mining Sites"
	l.Rows = make([]string, 15)
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = true
	l.SetRect(0, 0, 50, 20)

	return
}

func barChart() (bc *widgets.BarChart) {
	bc = widgets.NewBarChart()
	bc.Data = []float64{}
	bc.Labels = []string{}
	bc.Title = "Statistics"
	bc.SetRect(50, 0, 130, 20)
	bc.BarWidth = 1
	bc.BarGap = 5
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

	return
}

func cpuCurve() (p *widgets.Plot) {
	p = widgets.NewPlot()
	p.Title = "Realtime CPU usage"
	p.Data = make([][]float64, 1)
	p.Data[0] = []float64{0, 0, 0}
	p.SetRect(0, 20, 65, 35)
	p.AxesColor = ui.ColorWhite
	p.LineColors[0] = ui.ColorGreen

	return
}

func memCurve() (p *widgets.Plot) {
	p = widgets.NewPlot()
	p.Title = "Realtime Memory usage"
	p.Data = make([][]float64, 1)
	p.Data[0] = []float64{0, 0, 0}
	p.SetRect(65, 20, 130, 35)
	p.AxesColor = ui.ColorWhite
	p.LineColors[0] = ui.ColorGreen

	return
}

func progressBar() (g *widgets.Gauge) {
	g = widgets.NewGauge()
	g.Title = "Scan progress"
	g.SetRect(0, 35, 65, 40)
	g.Percent = 0
	g.Label = ""
	g.BarColor = ui.ColorGreen
	g.LabelStyle = ui.NewStyle(ui.ColorYellow)

	return
}

func processedMessageBox() (m *widgets.Paragraph) {
	m = widgets.NewParagraph()
	m.Title = "Site Processing"
	m.Text = ""
	m.SetRect(0, 40, 65, 45)
	m.BorderStyle.Fg = ui.ColorYellow

	return
}

func errorMessageBox() (m *widgets.Paragraph) {
	m = widgets.NewParagraph()
	m.Title = "Error Message"
	m.Text = ""
	m.SetRect(65, 35, 130, 45)
	m.BorderStyle.Fg = ui.ColorYellow

	return
}
