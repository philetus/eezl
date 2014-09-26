// public eezl interface

package eezl

import (
	//"fmt"
	"github.com/philetus/eezl/keys"
)

// mark eezl as stained to trigger new gel to be sent down gel pipe
func (self *Eezl) Stain() {

	// wrap sending dirty signal in select with default to make it non-blocking
	// -- if there is already a new gel pending just return
	select {
		case self.stainPipe <- 
			&gelStain{Resize: false, Height: self.height, Width: self.width}:
		default:
	}
}

// input flavrs
const (
	PointerMotion int = iota
	PointerPress
	PointerRelease
	KeyPress
	KeyRelease
)

// represents pointer motion and press events and keyboard events
type Input struct {
	Flavr int
	Timestamp int
	Y, X int
	Stroke *keys.Key
}


