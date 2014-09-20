package main

import (
	"fmt"
	"github.com/philetus/eezl"
)

func main() {
	xscreen := eezl.Xconnect() // get connection to x server screen
	xscreen.NewXwin(300, 200) // open window
	fmt.Printf("created new window!\n")
	for {} // loop
}
