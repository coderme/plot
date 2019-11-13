// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"fmt"
	"log"
	"os"

	"github.com/coderme/gonum/mat"
	"github.com/coderme/plot" 
	"github.com/coderme/plot/palette"
	"github.com/coderme/plot/plotter"
	"github.com/coderme/plot/vg"
	"github.com/coderme/plot/vg/draw"
	"github.com/coderme/plot/vg/vgimg"
)

func ExampleHeatMap() {
	m := offsetUnitGrid{
		XOffset: -2,
		YOffset: -1,
		Data: mat.NewDense(3, 4, []float64{
			1, 2, 3, 4,
			5, 6, 7, 8,
			9, 10, 11, 12,
		})}
	pal := palette.Heat(12, 1)
	h := plotter.NewHeatMap(m, pal)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Heat map"

	p.X.Tick.Marker = integerTicks{}
	p.Y.Tick.Marker = integerTicks{}

	p.Add(h)

	// Create a legend.
	l, err := plot.NewLegend()
	if err != nil {
		log.Panic(err)
	}
	thumbs := plotter.PaletteThumbnailers(pal)
	for i := len(thumbs) - 1; i >= 0; i-- {
		t := thumbs[i]
		if i != 0 && i != len(thumbs)-1 {
			l.Add("", t)
			continue
		}
		var val float64
		switch i {
		case 0:
			val = h.Min
		case len(thumbs) - 1:
			val = h.Max
		}
		l.Add(fmt.Sprintf("%.2g", val), t)
	}

	p.X.Padding = 0
	p.Y.Padding = 0
	p.X.Max = 1.5
	p.Y.Max = 1.5

	img := vgimg.New(250, 175)
	dc := draw.New(img)

	l.Top = true
	// Calculate the width of the legend.
	r := l.Rectangle(dc)
	legendWidth := r.Max.X - r.Min.X
	l.YOffs = -p.Title.Font.Extents().Height // Adjust the legend down a little.

	l.Draw(dc)
	dc = draw.Crop(dc, 0, -legendWidth-vg.Millimeter, 0, 0) // Make space for the legend.
	p.Draw(dc)
	w, err := os.Create("testdata/heatMap.png")
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}
