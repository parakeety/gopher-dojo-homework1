package converter

import (
	"errors"
	"os"
	"testing"
)

const testdataDir = "testdata/"

func TestConvert(t *testing.T) {
	errorCases := []struct {
		input    string
		output   string
		expected error
	}{
		{"", "jpeg", errors.New("input is empty")},
		{"png", "", errors.New("output is empty")},
		{"png", "png", errors.New("please specify different extensions for input & output")},
		{"gif", "jpg", errors.New("specified input extension is not supported")},
	}

	for _, tc := range errorCases {
		t.Run(tc.input+" -> "+tc.output, func(t *testing.T) {
			err := Convert(testdataDir, tc.input, tc.output)
			if !compareErrors(err, tc.expected) {
				t.Errorf("actual: %v, expected: %v", err, tc.expected)
			}
		})
	}

	successCases := []struct {
		input      string
		output     string
		outputPath string
	}{
		{"jpg", "png", testdataDir + "inputjpg.png"},
		{"png", "jpg", testdataDir + "inputpng.jpg"},
	}

	for _, tc := range successCases {
		t.Run(tc.input+" -> "+tc.output, func(t *testing.T) {
			testConvert(t, tc.input, tc.output, tc.outputPath)
		})
	}
}

func compareErrors(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1.Error() == err2.Error() {
		return true
	}

	return false
}

func testConvert(t *testing.T, input, output, outputPath string) {
	t.Helper()

	err := Convert(testdataDir, input, output)
	if err != nil {
		t.Fatalf("expected error to be nil: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("expected output file at path: %s", outputPath)
	}
	if err := os.Remove(outputPath); err != nil {
		t.Fatalf("failed remove output file: %v", err)
	}
}
