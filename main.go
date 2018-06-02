package main

import (
	"fmt"
	"log"

	"flag"

	"github.com/parakeety/gopher-dojo-homework1/internal/converter"
)

var (
	inputExt  = flag.String("input", "jpg", "input image format. i.e) -input=jpg")
	outputExt = flag.String("output", "png", "output image format. i.e) -output=png")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		panic("no directory to parse specified") // should add as a testcase...?
	}

	if err := converter.Convert(flag.Args()[0], *inputExt, *outputExt); err != nil {
		log.Fatalf("failed converting image: %v", err)
	}
	fmt.Println("Image conversion complete")
}
