# Plot 


`plot` is a fork  of [Gonum plot](https://github.com/gonum/plot) which is the new, official fork of code.google.com/p/plotinum.
It provides an API for building and drawing plots in Go.

*Note* that this new API is still in flux and may change.
See the wiki for some [example plots](http://github.com/gonum/plot/wiki/Example-plots).

For additional Plotters, see the [Community Plotters](https://github.com/gonum/plot/wiki/Community-Plotters) Wiki page.

There is a discussion list on Google Groups: gonum-dev@googlegroups.com.

`plot` is split into a few packages:

* The `plot` package provides simple interface for laying out a plot and provides primitives for drawing to it.
* The `plotter` package provides a standard set of `Plotter`s which use the primitives provided by the `plot` package for drawing lines, scatter plots, box plots, error bars, etc. to a plot. You do not need to use the `plotter` package to make use of `gonum/plot`, however: see the wiki for a tutorial on making your own custom plotters.
* The `plotutil` package contains a few routines that allow some common plot types to be made very easily. This package is quite new so it is not as well tested as the others and it is bound to change.
* The `vg` package provides a generic vector graphics API that sits on top of other vector graphics back-ends such as a custom EPS back-end, draw2d, SVGo, X-Window and gopdf.

## Documentation

Documentation is available at:

  https://godoc.org/gonum.org/v1/plot

## Installation

You can get `plot` using go get:

`go get github.com/coderme/plot/...`

