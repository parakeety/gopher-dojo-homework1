package main

import (
	"log"
	"fmt"

	"github.com/parakeety/gopher-dojo-homework1/internal/converter"
)

func main() {
	if err := converter.Converter().Execute(); err != nil {
		log.Fatalf("failed converting image: %v", err)
	}
	fmt.Println("Image conversion complete")
}
