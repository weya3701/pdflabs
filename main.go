package main

import (
	"flag"
	"pdflabs/pdfopt"
)

func main() {

	action := flag.String("action", "", "Command action name")
	originfilename := flag.String("origin", "", "Origin File Name")
	sourcefilename := flag.String("source", "", "Source File Name")
	outputfilename := flag.String("output", "", "Output File Name")
	tagContent := flag.String("content", "", "Content")
	contentfromfile := flag.String("cf", "", "Content From File")
	pagenum := flag.Int("pagenum", 1, "Page Number")
	f_x := flag.Float64("x", 0, "X")
	f_y := flag.Float64("y", 0, "Y")

	flag.Parse()

	var params pdfopt.PdfOptParameter

	params = pdfopt.PdfOptParameter{
		OriginFilename:  *originfilename,
		SourceFilename:  *sourcefilename,
		OutputFilename:  *outputfilename,
		Content:         *tagContent,
		ContentFromFile: *contentfromfile,
		PageNum:         *pagenum,
		X:               f_x,
		Y:               f_y,
		PdfCoordinate: pdfopt.Coordinate{
			MarginLeft:   50.0,
			MarginTop:    50.0,
			MarginBottom: 50.0,
			PageHeight:   841.89,
			LineHeight:   20.0,
			ContentWidth: 500.0,
			Position:     "bl",
			Scale:        0.3,
			Fontname:     "Helvetica",
			Points:       12,
			Rot:          0,
		},
	}

	if *action == "Merge" {
		// pdfopt.Merge2file(params)
		pdfopt.Merge2fileOptimized(params)
		// pdfopt.Merge2file(params)
	}
	if *action == "WriteContentFromFile" {
		pdfopt.WriteContentFromFile(params)
	}
	if *action == "StickTags" {
		pdfopt.StickTagsTest(params)
	}
}
