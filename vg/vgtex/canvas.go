// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vgtex provides a vg.Canvas implementation for LaTeX, targeted at
// the TikZ/PGF LaTeX package: https://sourceforge.net/projects/pgf
//
// vgtex generates PGF instructions that will be interpreted and rendered by LaTeX.
// vgtex allows to put any valid LaTeX notation inside plot's strings.
package vgtex // import "github.com/coderme/plot/vg/vgtex"

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/coderme/plot/vg"
)

const degPerRadian = 180 / math.Pi

const (
	defaultHeader = `%%%%%% generated by gonum/plot %%%%%%
\documentclass{standalone}
\usepackage{pgf}
\begin{document}
`
	defaultFooter = "\\end{document}\n"
)

// Canvas implements the vg.Canvas interface, translating drawing
// primitives from gonum/plot to PGF.
type Canvas struct {
	buf   *bytes.Buffer
	w, h  vg.Length
	stack []context

	// If document is true, Canvas.WriteTo will generate a standalone
	// .tex file that can be fed to, e.g., pdflatex.
	document bool
	id       int64 // id is a unique identifier for this canvas
}

type context struct {
	color      color.Color
	dashArray  []vg.Length
	dashOffset vg.Length
	linew      vg.Length
}

// New returns a new LaTeX canvas.
func New(w, h vg.Length) *Canvas {
	return newCanvas(w, h, false)
}

// NewDocument returns a new LaTeX canvas that can be readily
// compiled into a standalone document.
func NewDocument(w, h vg.Length) *Canvas {
	return newCanvas(w, h, true)
}

func newCanvas(w, h vg.Length, document bool) *Canvas {
	c := &Canvas{
		buf:      new(bytes.Buffer),
		w:        w,
		h:        h,
		document: document,
		id:       time.Now().UnixNano(),
	}
	if !document {
		c.wtex(`%%%% gonum/plot created for LaTeX/pgf`)
		c.wtex(`%%%% you need to add:`)
		c.wtex(`%%%%   \usepackage{pgf}`)
		c.wtex(`%%%% to your LaTeX document`)
	}
	c.wtex("")
	c.wtex(`\begin{pgfpicture}`)
	c.stack = make([]context, 1)
	vg.Initialize(c)
	return c
}

func (c *Canvas) context() *context {
	return &c.stack[len(c.stack)-1]
}

// Size returns the width and height of the canvas.
func (c *Canvas) Size() (w, h vg.Length) {
	return c.w, c.h
}

// SetLineWidth implements the vg.Canvas.SetLineWidth method.
func (c *Canvas) SetLineWidth(w vg.Length) {
	c.context().linew = w
}

// SetLineDash implements the vg.Canvas.SetLineDash method.
func (c *Canvas) SetLineDash(pattern []vg.Length, offset vg.Length) {
	c.context().dashArray = pattern
	c.context().dashOffset = offset
}

// SetColor implements the vg.Canvas.SetColor method.
func (c *Canvas) SetColor(clr color.Color) {
	c.context().color = clr
}

// Rotate implements the vg.Canvas.Rotate method.
func (c *Canvas) Rotate(rad float64) {
	c.wtex(`\pgftransformrotate{%g}`, rad*degPerRadian)
}

// Translate implements the vg.Canvas.Translate method.
func (c *Canvas) Translate(pt vg.Point) {
	c.wtex(`\pgftransformshift{\pgfpoint{%gpt}{%gpt}}`, pt.X, pt.Y)
}

// Scale implements the vg.Canvas.Scale method.
func (c *Canvas) Scale(x, y float64) {
	c.wtex(`\pgftransformxscale{%g}`, x)
	c.wtex(`\pgftransformyscale{%g}`, y)
}

// Push implements the vg.Canvas.Push method.
func (c *Canvas) Push() {
	c.wtex(`\begin{pgfscope}`)
	c.stack = append(c.stack, *c.context())
}

// Pop implements the vg.Canvas.Pop method.
func (c *Canvas) Pop() {
	c.stack = c.stack[:len(c.stack)-1]
	c.wtex(`\end{pgfscope}`)
	c.wtex("")
}

// Stroke implements the vg.Canvas.Stroke method.
func (c *Canvas) Stroke(p vg.Path) {
	if c.context().linew <= 0 {
		return
	}
	c.wstyle()
	c.wpath(p)
	c.wtex(`\pgfusepath{stroke}`)
	c.wtex("")
}

// Fill implements the vg.Canvas.Fill method.
func (c *Canvas) Fill(p vg.Path) {
	c.wstyle()
	c.wpath(p)
	c.wtex(`\pgfusepath{fill, stroke}`)
	c.wtex("")
}

// FillString implements the vg.Canvas.FillString method.
func (c *Canvas) FillString(f vg.Font, pt vg.Point, text string) {
	c.wcolor()
	pt.X += 0.5 * f.Width(text)
	c.wtex(`\pgftext[base,at={\pgfpoint{%gpt}{%gpt}}]{%s}`, pt.X, pt.Y, text)
}

// DrawImage implements the vg.Canvas.DrawImage method.
// DrawImage will first save the image inside a PNG file and have the
// generated LaTeX reference that file.
// The file name will be "gonum-pgf-image-<canvas-id>-<time.Now()>.png
func (c *Canvas) DrawImage(rect vg.Rectangle, img image.Image) {
	fname := fmt.Sprintf("gonum-pgf-image-%v-%v.png", c.id, time.Now().UnixNano())
	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(fmt.Errorf("vgtex: error encoding image to PNG: %v", err))
	}

	var (
		xmin   = rect.Min.X
		ymin   = rect.Min.Y
		width  = rect.Size().X
		height = rect.Size().Y
	)
	c.wtex(`\pgftext[base,left,at=\pgfpoint{%gpt}{%gpt}]{\pgfimage[height=%gpt,width=%gpt]{%s}}`, xmin, ymin, height, width, fname)
}

func (c *Canvas) indent(s string) string {
	return strings.Repeat(s, len(c.stack))
}

func (c *Canvas) wtex(s string, args ...interface{}) {
	fmt.Fprintf(c.buf, c.indent("  ")+s+"\n", args...)
}

func (c *Canvas) wstyle() {
	c.wdash()
	c.wlineWidth()
	c.wcolor()
}

func (c *Canvas) wdash() {
	if len(c.context().dashArray) == 0 {
		return
	}
	str := `\pgfsetdash{`
	for _, d := range c.context().dashArray {
		str += fmt.Sprintf("{%gpt}", d)
	}
	str += fmt.Sprintf("}{%gpt}", c.context().dashOffset)
	c.wtex(str)
}

func (c *Canvas) wlineWidth() {
	c.wtex(`\pgfsetlinewidth{%gpt}`, c.context().linew)
}

func (c *Canvas) wcolor() {
	col := c.context().color
	if col == nil {
		col = color.Black
	}
	r, g, b, a := col.RGBA()
	alpha := 255.0 / float64(a)
	// FIXME(sbinet) \color will last until the end of the current TeX group
	// use \pgfsetcolor and \pgfsetstrokecolor instead.
	// it needs a named color: define it on the fly (storing it at the beginning
	// of the document.)
	c.wtex(
		`\color[rgb]{%g,%g,%g}`,
		float64(r)*alpha/255.0,
		float64(g)*alpha/255.0,
		float64(b)*alpha/255.0,
	)

	opacity := float64(a) / math.MaxUint16
	c.wtex(`\pgfsetstrokeopacity{%g}`, opacity)
	c.wtex(`\pgfsetfillopacity{%g}`, opacity)
}

func (c *Canvas) wpath(p vg.Path) {
	for _, comp := range p {
		switch comp.Type {
		case vg.MoveComp:
			c.wtex(`\pgfpathmoveto{\pgfpoint{%gpt}{%gpt}}`, comp.Pos.X, comp.Pos.Y)
		case vg.LineComp:
			c.wtex(`\pgflineto{\pgfpoint{%gpt}{%gpt}}`, comp.Pos.X, comp.Pos.Y)
		case vg.ArcComp:
			start := comp.Start * degPerRadian
			angle := comp.Angle * degPerRadian
			r := comp.Radius
			c.wtex(`\pgfpatharc{%g}{%g}{%gpt}`, start, angle, r)
		case vg.CurveComp:
			var a, b vg.Point
			switch len(comp.Control) {
			case 1:
				a = comp.Control[0]
				b = a
			case 2:
				a = comp.Control[0]
				b = comp.Control[1]
			default:
				panic("vgtex: invalid number of control points")
			}
			c.wtex(`\pgfcurveto{\pgfpoint{%gpt}{%gpt}}{\pgfpoint{%gpt}{%gpt}}{\pgfpoint{%gpt}{%gpt}}`,
				a.X, a.Y, b.X, b.Y, comp.Pos.X, comp.Pos.Y)
		case vg.CloseComp:
			c.wtex("%% path-close")
		default:
			panic(fmt.Errorf("vgtex: unknown path component type: %v\n", comp.Type))
		}
	}
}

// WriteTo implements the io.WriterTo interface, writing a LaTeX/pgf plot.
func (c *Canvas) WriteTo(w io.Writer) (int64, error) {
	var (
		n   int64
		nn  int
		err error
	)
	b := bufio.NewWriter(w)
	if c.document {
		nn, err = b.Write([]byte(defaultHeader))
		n += int64(nn)
		if err != nil {
			return n, err
		}
	}
	m, err := c.buf.WriteTo(b)
	n += m
	if err != nil {
		return n, err
	}
	nn, err = fmt.Fprintf(b, "\\end{pgfpicture}\n")
	n += int64(nn)
	if err != nil {
		return n, err
	}

	if c.document {
		nn, err = b.Write([]byte(defaultFooter))
		n += int64(nn)
		if err != nil {
			return n, err
		}
	}
	return n, b.Flush()
}
