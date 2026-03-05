package main

import (
	"flag"
	"pdflabs/pdfopt"
)

func main() {

	action := flag.String("action", "", "Command action name")
	originfilename := flag.String("of", "", "Origin File Name")
	sourcefilename := flag.String("sf", "", "Source File Name")
	outputfilename := flag.String("opf", "", "Output File Name")
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
		},
	}

	if *action == "Merge" {
		pdfopt.Merge2file(params)
	}
	if *action == "WriteContentFromFile" {
		pdfopt.WriteContentFromFile(params)
	}
	if *action == "StickTags" {
		pdfopt.StickTags(params)
	}
}

// 1) 合併模式：輸入：原始檔、來源檔，輸出：檔案
// 2) 疊加寫入模式： 輸入：原始檔、來源內容(txt格式)，輸出：檔案
// 3) 註解模式： 輸入：原始檔、頁數、內容區塊，輸出：檔案
