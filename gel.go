// +build !goci

// gel drawing context
//
// wraps cairo drawing functions

package eezl

/*
#cgo LDFLAGS: -lcairo
#include <cairo/cairo-xcb.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"

/*
import (
	"fmt"
)
*/

// gel has a cairo surface it wraps drawing functions of
type Gel struct {
	context *C.cairo_t
	trigger_pipe chan bool
	Height, Width int
}

// send signal that drawing is complete and gel should be rendered
func (self *Gel) Trigger() {
	self.trigger_pipe <- true
}

func (self *Gel) Jmto(y, x float64) {
	C.cairo_move_to(self.context, C.double(x), C.double(y))
}

func (self *Gel) Jmby(dy, dx float64) {
	C.cairo_rel_move_to(self.context, C.double(dx), C.double(dy))
}

func (self *Gel) Rato(y, x float64) {
	C.cairo_line_to(self.context, C.double(x), C.double(y))
}

func (self *Gel) Raby(dy, dx float64) {
	C.cairo_rel_line_to(self.context, C.double(dx), C.double(dy))
}

func (self *Gel) Beto(y, x, cy0, cx0, cy1, cx1 float64) {
	C.cairo_rel_curve_to(self.context, C.double(x), C.double(y), 
						 C.double(cx0), C.double(cy0),
						 C.double(cx1), C.double(cy1))
}

func (self *Gel) Beby(dy, dx, dcy0, dcx0, dcy1, dcx1 float64) {
	C.cairo_curve_to(self.context, C.double(dx), C.double(dy), 
					 C.double(dcx0), C.double(dcy0),
					 C.double(dcx1), C.double(dcy1))
}

// coat gel with current color
func (self *Gel) Coat() {
	C.cairo_paint(self.context)
}

// seal current gel subpath (from current position to last point after jump)
func (self *Gel) Seal() {
	C.cairo_close_path(self.context)
}

// clear current gel path
func (self *Gel) Shake() {
	C.cairo_new_path(self.context)
}

func (self *Gel) Stroke() {
	C.cairo_stroke_preserve(self.context)
}

func (self *Gel) Fill() {
	C.cairo_fill_preserve(self.context)
}

func (self *Gel) SetColor(r, g, b, a float64) {
	C.cairo_set_source_rgba(self.context, C.double(r), C.double(g), 
							C.double(b), C.double(a))
}

func (self *Gel) SetWeight(w float64) {
	C.cairo_set_line_width(self.context, C.double(w))
}

