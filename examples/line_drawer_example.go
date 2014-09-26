package main

import (
	"fmt"
	"github.com/philetus/eezl"
)

type point struct {
	y, x float64
}

type line struct {
	first, last point
}

func main() {
	xscreen := eezl.Xconnect() // get connection to x server screen
	ez := xscreen.NewEezl(600, 600) // open window
	fmt.Printf("created new eezl!\n")
	
	var pressed_flag bool = false
	var band line
	lines := make([]line, 8)
	
    var ln_thk float64 = 10.0
    ln_clr := [4]float64{1.0, 0.0, 0.0, 0.6} // slightly translucent red
    var bnd_thk float64 = 6.0
    bnd_clr := [4]float64{0.0, 0.0, 0.0, 0.4} // translucent gray
    bg_clr := [4]float64{1.0, 1.0, 1.0, 1.0} // opaque white

	for {
		select {
			
			case inpt := <-ez.InputPipe:
				//fmt.Printf(".%d", inpt.Flavr)
				switch inpt.Flavr {
				
				case eezl.PointerPress:
					pressed_flag = true
					band.first.y = float64(inpt.Y)
					band.first.x = float64(inpt.X)
					band.last.y = float64(inpt.Y)
					band.last.x = float64(inpt.X)

				case eezl.PointerRelease:
					lines = append(lines, band)
					pressed_flag = false
								   
					ez.Stain() // trigger eezl redraw
					
				case eezl.PointerMotion:
					if pressed_flag {
						band.last.y = float64(inpt.Y)
						band.last.x = float64(inpt.X)
						
						ez.Stain() // trigger eezl redraw
					}
				}
				
			case gel := <- ez.GelPipe:
			
				// fill background
				gel.SetColor(bg_clr[0], bg_clr[1], bg_clr[2], bg_clr[3])
				gel.Coat()
				
				// draw lines
				gel.SetColor(ln_clr[0], ln_clr[1], ln_clr[2], ln_clr[3])
				gel.SetWeight(ln_thk)
				for _, ln := range(lines) {
					gel.Jmto(ln.first.y, ln.first.x)
					gel.Rato(ln.last.y, ln.last.x)
					gel.Stroke()
					gel.Shake()
				}
				
				// draw band
				if pressed_flag {
					gel.SetColor(bnd_clr[0], bnd_clr[1], bnd_clr[2], bnd_clr[3])
					gel.SetWeight(bnd_thk)
					gel.Jmto(band.first.y, band.first.x)
					gel.Rato(band.last.y, band.last.x)
					gel.Stroke()
					gel.Shake()
				}
				
				// send trigger sig
				gel.Trigger()
		}
	}
}
