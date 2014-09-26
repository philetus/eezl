// public eezl interface

package eezl

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

