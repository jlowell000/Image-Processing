package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/urfave/cli"
	"jlowell000.github.io/init/internal/kernalfunctions"
	"jlowell000.github.io/init/internal/pixalfunctions"
)

const (
	cliDelimitor = ","
)

var (
	availableFunctions = map[string]func(image.Image) image.Image{
		"greyscale":       pixalfunctions.GreyscaleHandle,
		"gaussian":        kernalfunctions.GaussianHandle,
		"sobel":           kernalfunctions.SobelHandle,
		"doubleThreshold": pixalfunctions.DoubleThresholdHandle,
		"fillInGaps":      pixalfunctions.FillInGapsHandle,
	}
)

func main() {

	var (
		functions, filePaths []string
		functionStr, errStr  string
		verboseOutput        bool
	)

	app := cli.NewApp()
	app.Name = "ImageProcessing"
	app.Usage = "Run functions over Image files\nSupports JPEG and PNG"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "output, o",
			Usage:       "makes verbose file output foreach step",
			Destination: &verboseOutput,
		},
		cli.StringFlag{
			Name:  "function, f",
			Value: "all",
			Usage: "Functions requested to run. \"all\" runs all functions for edge detection\n  Available functions: " +
				"[" + strings.Join(getAvailableFunctionsNames(), ", ") + "]",
			Destination: &functionStr,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			filePaths = strings.Split(c.Args()[0], cliDelimitor)
		}
		if len(functionStr) > 0 && functionStr != "all" {
			functions = strings.Split(functionStr, cliDelimitor)
			verifyFuncList := getAvailableFunctionsNames()
			for _, f := range functions {
				if !contains(verifyFuncList, f) {
					errStr = "Function " + f + " not in available"
					fmt.Printf("ERROR:  %s\n", errStr)
				}
			}
		} else {
			functions = getAvailableFunctionsNames()
		}
		if len(errStr) == 0 {
			if len(filePaths) > 0 && len(functions) > 0 {
				if verboseOutput {
					fmt.Printf("Funcions to run: %s\n", functions)
				}
				for _, file := range filePaths {
					fmt.Printf("Processing File: %s\n", file)
					loadedImage := readFileToImage(file)
					if loadedImage != nil {
						for _, f := range functions {
							loadedImage = availableFunctions[f](loadedImage)
							if verboseOutput {
								newFileName := createDerivedFileNames(file, f)
								fmt.Printf("Creating file: %s\n", newFileName)
								writeImageFile(newFileName, loadedImage)
							}
						}
						if !verboseOutput {
							newFileName := createDerivedFileNames(file, "output")
							fmt.Printf("Creating file: %s\n", newFileName)
							writeImageFile(newFileName, loadedImage)
						}
					} else {
						fmt.Printf("File at %s not found\n", file)
					}
				}
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func getAvailableFunctionsNames() []string {
	keys := reflect.ValueOf(availableFunctions).MapKeys()
	strkeys := make([]string, len(keys))
	for i, k := range keys {
		strkeys[i] = k.String()
	}
	return strkeys
}
func contains(list []string, e string) bool {
	for _, a := range list {
		if a == e {
			return true
		}
	}
	return false
}

func createDerivedFileNames(filePath string, mod string) string {
	return path.Dir(filePath) + string(os.PathSeparator) + mod + "_" + path.Base(filePath)
}

func readFileToImage(fileName string) image.Image {
	var loadedImage image.Image
	// Read image from inFile that already exists
	existingImageFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()
	fileExt := path.Ext(fileName)
	if fileExt == ".png" {
		loadedImage, err = png.Decode(existingImageFile)
	} else if fileExt == ".jpg" {
		loadedImage, err = jpeg.Decode(existingImageFile)
	}
	if err != nil {
		panic(err)
	}
	return loadedImage
}

func writeImageFile(fileName string, image image.Image) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileExt := path.Ext(fileName)
	if fileExt == ".png" {
		err = png.Encode(f, image)
	} else if fileExt == ".jpg" {
		err = jpeg.Encode(f, image, nil)
	}
	if err != nil {
		panic(err)
	}
}
