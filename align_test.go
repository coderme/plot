// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot_test

import (
	"testing"

	"github.com/coderme/plot/cmpimg"
)

func TestAlign(t *testing.T) {
	cmpimg.CheckPlot(ExampleAlign, t, "align.png")
}
