/*
Package converter implements CLI that converts images to different formats.
Options
- input: input image format. defaults to jpg. i.e) -input=jpg
- output: output image format. defaults to png. i.e) output=png
*/
package converter

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func Convert(dir, input, output string) error {
	// cleanse inputs
	dirRoot := strings.TrimSpace(dir)
	inputExt := strings.TrimSpace(input)
	outputExt := strings.TrimSpace(output)

	if err := validateInputs(dirRoot, inputExt, outputExt); err != nil {
		return err
	}

	if !strings.HasPrefix(inputExt, ".") {
		inputExt = "." + inputExt
	}

	if !strings.HasPrefix(outputExt, ".") {
		outputExt = "." + outputExt
	}

	imagePaths, err := findImagePaths(dirRoot, inputExt)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}
	for _, path := range imagePaths {
		eg.Go(func() error {
			return convert(path, inputExt, outputExt)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func validateInputs(src, input, output string) error {
	inputs := map[string]string{
		"src":    src,
		"input":  input,
		"output": output,
	}

	for key, value := range inputs {
		if len(value) == 0 {
			return fmt.Errorf("%s is empty", key)
		}
	}
	if !isSupportedFormat(input) {
		return errors.New("specified input extension is not supported")
	}

	if !isSupportedFormat(output) {
		return errors.New("specified output extension is not supported")
	}

	if input == output {
		return errors.New("please specify different extensions for input " +
			"& output")
	}
	return nil
}

func isSupportedFormat(ext string) bool {
	switch ext {
	case "jpeg", "jpg", "png":
		return true
	default:
		return false
	}
}
func findImagePaths(dirRoot, inputExt string) ([]string, error) {
	var imagePaths []string
	err := filepath.Walk(dirRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == inputExt {
			imagePaths = append(imagePaths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return imagePaths, nil
}

func convert(path, inputExt, outputExt string) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	var img image.Image
	switch inputExt {
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

	outputPath := strings.TrimSuffix(path, inputExt) + outputExt
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	switch outputExt {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outputFile, img, nil)
	case ".png":
		err = png.Encode(outputFile, img)
	default:
		return errors.New("unsupported image extension")
	}
	return err
}
