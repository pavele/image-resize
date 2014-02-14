package main

import (
	opencv "github.com/lazywei/go-opencv/opencv"
	//"io/ioutil"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const srcDir = "res"
const destDir = "res/out"

var ratios = []float64{0.90, 0.80, 0.70, 0.60, 0.50, 0.40, 0.30, 0.10}

func main() {

	filepath.Walk(srcDir, resize)
}

func resize(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() {
		return nil
	}
	fmt.Printf("Processing: %s\n", path)

	srcImg := opencv.LoadImage(path)
	if srcImg == nil {
		return errors.New("Error loading image: " + path)
	}
	defer srcImg.Release()
	srcW := float64(srcImg.Width())
	srcH := float64(srcImg.Height())
	
	var fileName = fileInfo.Name()	
	var extension = filepath.Ext(fileName)
	var nameNoExt = strings.TrimRight(fileName, extension)
	
	for _, ratio := range ratios {
		thumbW := srcW * ratio
		thumbH := srcH * ratio
		thumbnail := opencv.Resize(srcImg, int(thumbW), int(thumbH), 1)
		name:= filepath.Join(destDir, fmt.Sprintf("%s__%f%s", nameNoExt, ratio * 10, extension))		
		
		opencv.SaveImage(name, thumbnail, 0)
	}
	fmt.Printf("Done processing: %s\n", path)
	return nil
}
