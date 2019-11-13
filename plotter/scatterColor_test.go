// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter_test

import (
	"testing"

	"github.com/coderme/plot/cmpimg"
)

func TestScatterColor(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatter_color, t, "scatterColor.png")
}
