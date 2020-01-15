package main

import (
	"fmt"
)

func main() {
	config := initConfig()
	fmt.Printf("%+v\n ", config)
	startController(config)
}
