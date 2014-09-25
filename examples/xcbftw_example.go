package main

import (
	"fmt"
	"github.com/philetus/eezl"
)

func main() {
	xscreen := eezl.Xconnect() // get connection to x server screen
	eezl := xscreen.NewEezl(300, 200) // open window
	fmt.Printf("created new window!\n")
	for {
		inpt := <-eezl.InputPipe
		fmt.Printf(".%d", inpt.Flavr)
	}
}
