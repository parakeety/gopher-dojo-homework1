/*
Package converter implements CLI that converts images to different formats.
Options
- input: input image format. defaults to jpg. i.e) -input=jpg
- output: output image format. defaults to png. i.e) output=png
*/
package converter

import (
	"flag"
	"strings"
	"path/filepath"
	"os"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/go-errors/errors"
)

var (
	inputExt = flag.String("input", "jpg", "input image format. i.e) -input=jpg")
	outputExt = flag.String("output", "png", "output image format. i.e) -output=png")
)

type CommandLine interface {
	Execute() error
}

type converter struct {
	InputExt string
	OutputExt string
}

func Converter() CommandLine {
	flag.Parse()

	inputImgExt := *inputExt
	if !strings.HasPrefix(inputImgExt, ".") {
		inputImgExt = "." + inputImgExt
	}

	outputImgExt := *outputExt
	if !strings.HasPrefix(outputImgExt, ".") {
		outputImgExt = "." + outputImgExt
	}

	return &converter{
		InputExt: inputImgExt,
		OutputExt: outputImgExt,
	}
}

func (c *converter) Execute() error {
	if len(flag.Args()) == 0 {
		return errors.New("please pass the directory where images you'd like exist")
	}

	dirRoot := flag.Args()[0]

	return filepath.Walk(dirRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != c.InputExt {
			return nil
		}

		inputFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		var img image.Image
		switch c.InputExt {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(inputFile)
		case ".png":
			img, err = png.Decode(inputFile)
		default:
			return errors.New("unsupported image extension")
		}

		if err != nil {
			return err
		}

		outputPath := strings.TrimSuffix(path, c.InputExt) + c.OutputExt
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		switch c.OutputExt {
		case ".jpg", ".jpeg":
			err = jpeg.Encode(outputFile, img, nil)
		case ".png":
			err = png.Encode(outputFile, img)
		default:
			return errors.New("unsupported image extension")
		}

		if err != nil {
			return err
		}

		return nil
	})
}
