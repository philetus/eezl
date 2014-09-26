package main

import (
	"fmt"
	"github.com/philetus/eezl"
)

func main() {
	xscreen := eezl.Xconnect() // get connection to x server screen
	ez := xscreen.NewEezl(300, 200) // open window
	fmt.Printf("created new eezl!\n")
	for {
		select {
			
			case inpt := <-ez.InputPipe:
				fmt.Printf(".%d", inpt.Flavr)
				if inpt.Flavr == eezl.PointerPress {
					ez.Stain()
				}
				
			case gel := <- ez.GelPipe:
			
				// draw a red rectangle
				gel.SetColor(1.0, 0.0, 0.0, 1.0)
				gel.SetWeight(8.0)
				gel.Jmto(10.0, 10.0)
				gel.Raby(160.0, 0.0)
				gel.Raby(0.0, 80.0)
				gel.Raby(-160.0, 0.0)
				gel.Seal()
				gel.Stroke()
				gel.Shake()
				
				// send trigger sig
				gel.Trigger()
		}
	}
}
