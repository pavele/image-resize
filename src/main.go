package main

import (
	opencv "github.com/lazywei/go-opencv/opencv"
	"io/ioutil"
	"runtime"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const srcDir = "images"
const destDir = "result"

const numberOfSamples = 100

var ratios = []float64{0.90, 0.80, 0.70, 0.60, 0.50, 0.40, 0.30, 0.10}
type resolution struct{
	width, height int
}
var resolutions = []resolution{
	{ 640, 480},
	{ 800, 600},
	{ 1024, 768},
	{ 1152, 864},
	{ 1280, 720},
	{ 1360, 768},
	{ 1400, 1050},
	{ 1600, 980},
	{ 1600, 1200},
	{ 1900, 1080} }

func main() {
	runtime.GOMAXPROCS(4)
	//filepath.Walk(srcDir, resize)
	var files,_ = ioutil.ReadDir(srcDir)
	var samples [numberOfSamples]os.FileInfo
	for i:=0; i< numberOfSamples; i++ {
		samples[i] = files[i% len(files)]
	}
	
	var stime = time.Now()
	
	//for i:=0; i< numberOfSamples; i++ {
		//fmt.Printf("%s\n", samples[i].Name())
		//resize(srcDir + "/" + samples[i].Name(), samples[i], nil)
	//}
	
	semaphore := make(chan int,4)
	go func(quit chan int) {
		fmt.Printf("First thread started\n")
		for i:=0; i< 0.25*numberOfSamples;i++ {
			resize(srcDir + "/" + samples[i].Name(), samples[i], nil)
		}
		fmt.Printf("First thread finished\n")
		quit <- 1
	}(semaphore)
	
	go func(quit chan int) {
		fmt.Printf("Second thread started\n")
		for i:= int(0.25*numberOfSamples); i< 0.5*numberOfSamples;i++ {
			resize(srcDir + "/" + samples[i].Name(), samples[i], nil)
		}
		fmt.Printf("Second thread finished\n")
		quit<-1
	}(semaphore)	
	
	/*go func(quit chan int) {
		fmt.Printf("Third thread started\n")
		for i:= int(0.5*numberOfSamples); i< 0.75*numberOfSamples;i++ {
			resize(srcDir + "/" + samples[i].Name(), samples[i], nil)
		}
		fmt.Printf("Third thread finished\n")
		quit<-1
	}(semaphore)	
	
	go func(quit chan int) {
		fmt.Printf("Fourth thread started\n")
		for i:=int(0.75*numberOfSamples); i< numberOfSamples;i++ {
			resize(srcDir + "/" + samples[i].Name(), samples[i], nil)
		}
		fmt.Printf("Fourth thread finished\n")
		quit<-1
	}(semaphore)
	
	<-semaphore
	<-semaphore*/
	<-semaphore
	<-semaphore
	
	fmt.Printf("Started at: %s\n", stime)
	fmt.Printf("Finished at: %s\n", time.Now())
}

func resize(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() {
		return nil
	}
	//fmt.Printf("Processing: %s\n", path)

	srcImg := opencv.LoadImage(path)
	if srcImg == nil {
		return errors.New("Error loading image: " + path)
	}
	defer srcImg.Release()
	//srcW := float64(srcImg.Width())
	//srcH := float64(srcImg.Height())
	
	var fileName = fileInfo.Name()	
	var extension = filepath.Ext(fileName)
	var nameNoExt = strings.TrimRight(fileName, extension)
	
	for _, res := range resolutions {
		//thumbW := srcW * ratio
		//thumbH := srcH * ratio
		//fmt.Printf("Resolution %dx%d\n", res.width, res.height)
		thumbW := res.width
		thumbH := res.height
		thumbnail := opencv.Resize(srcImg, int(thumbW), int(thumbH), 1)
		name:= filepath.Join(destDir, fmt.Sprintf("%s__%dx%d_%s", nameNoExt, res.width, res.height, extension))		
		
		opencv.SaveImage(name, thumbnail, 0)
	}
	//fmt.Printf("Done processing: %s\n", path)
	return nil
}
