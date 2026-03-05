package pdfopt

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dslipak/pdf"
	"github.com/signintech/gopdf"
)

type PdfOptParameter struct {
	OriginFilename   string
	SourceFilename   string
	OutputFilename   string
	Content          string
	MultiLineContent []string
	ContentFromFile  string
	PageNum          int
	X                *float64
	Y                *float64
	PdfCoordinate    Coordinate
}

type Coordinate struct {
	MarginLeft   float64
	MarginTop    float64
	MarginBottom float64
	PageHeight   float64
	LineHeight   float64
	ContentWidth float64
}

type tytbPDF struct {
	PDFInstance gopdf.GoPdf
}

var (
	instance *tytbPDF
	once     sync.Once
)

// ReadLines 讀取檔案路徑並回傳字串切片
func ReadLines(path string) ([]string, error) {
	// 1. 打開檔案
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// 確保函式結束時關閉檔案資源
	defer file.Close()

	var lines []string
	// 2. 使用 Scanner 逐行掃描
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 將每一行的文字加入切片
		lines = append(lines, scanner.Text())
	}

	// 3. 檢查掃描過程中是否有錯誤
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func GetCurrentFiles(root string) []string {

	var vfiles []string = []string{}
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		vfiles = append(vfiles, path)
		return nil
	})

	return vfiles
}

func getPDFInstance() *tytbPDF {
	once.Do(func() {
		instance = &tytbPDF{
			PDFInstance: gopdf.GoPdf{},
		}
	})
	instance.PDFInstance.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	return instance
}

func (tp *tytbPDF) SetPageSize(cfg gopdf.Config) bool {
	var rsp bool = true
	tp.PDFInstance.Start(cfg)
	return rsp
}

// Get source pdf pages number
func getPDFpages(filename string) int {
	var rsp int = 0
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	r, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		fmt.Printf("解析失敗: %v\n", err)
		return rsp
	}

	numPages := r.NumPage()
	return numPages
}

// Merge 2 pdf file.
func Merge2file(params PdfOptParameter) {

	files := []string{
		params.OriginFilename,
		params.SourceFilename,
	}

	writer := getPDFInstance()
	writer.PDFInstance.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	defer writer.PDFInstance.Close()
	err := writer.PDFInstance.AddTTFFont("myfont", "./fonts/ArialBlack.ttf")
	if err != nil {
		fmt.Printf("Can not add font %v", err)
	}

	for _, file := range files {
		if file != "" {
			sourcePages := getPDFpages(file)
			for i := 1; i <= sourcePages; i++ {
				tpl := writer.PDFInstance.ImportPage(file, i, "/MediaBox")
				writer.PDFInstance.AddPage()
				writer.PDFInstance.UseImportedTemplate(tpl, 0, 0, 595.28, 841.89)
			}
		}
	}

	err = writer.PDFInstance.WritePdf(params.OutputFilename)
	if err != nil {
		log.Printf("Can not wirte pdf file: %v", err)
	}
}

func StickTags(params PdfOptParameter) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err := pdf.AddTTFFont("myfont", "./fonts/tahoma.ttf")
	defer pdf.Close()
	if err != nil {
		fmt.Printf("Can not add font %v", err)
	}

	sourcePages := getPDFpages(params.OriginFilename)
	for i := 1; i <= sourcePages; i++ {
		tpl := pdf.ImportPage(params.OriginFilename, i, "/MediaBox")
		pdf.AddPage()
		pdf.UseImportedTemplate(tpl, 0, 0, 595.28, 841.89)
		if i == params.PageNum {
			pdf.SetFont("myfont", "", 12)
			pdf.SetXY(*params.X, *params.Y)
			pdf.Text(params.Content)
		}
	}
	pdf.WritePdf(params.OutputFilename)
}

func WriteContentFromFile(params PdfOptParameter) {

	filecontent, _ := ReadLines(params.ContentFromFile)
	params.MultiLineContent = filecontent
	WriteContent(params)
}

func WriteContent(params PdfOptParameter) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err := pdf.AddTTFFont("myfont", "./fonts/tahoma.ttf")
	if err != nil {
		fmt.Printf("Can not add font %v", err)
	}

	pdf.AddPage()

	currY := params.PdfCoordinate.MarginTop

	for _, paragraph := range params.MultiLineContent {

		fmt.Println("paragraph: ", paragraph)
		pdf.SetFont("myfont", "", 12)
		// 將段落依照寬度切成多行
		lines, err := pdf.SplitText(paragraph, params.PdfCoordinate.ContentWidth)
		fmt.Println("err: ", err)

		for _, line := range lines {
			// 檢查是否超過頁面底部
			if currY+params.PdfCoordinate.LineHeight > (params.PdfCoordinate.PageHeight - params.PdfCoordinate.MarginBottom) {
				pdf.AddPage()                          // 換頁
				currY = params.PdfCoordinate.MarginTop // 重置 Y 座標到頁首
			}

			// 寫入文字
			pdf.SetXY(params.PdfCoordinate.MarginLeft, currY)
			pdf.Text(line)

			// 增加 Y 座標 (準備寫下一行)
			currY += params.PdfCoordinate.LineHeight
		}

		// 段落與段落之間可以額外增加一點間距
		currY += 10.0
	}
	pdf.WritePdf(params.OutputFilename)
	defer pdf.Close()
}

// Read pdf content.
// For Test.
func ReadPdfText(path string) string {
	f, err := pdf.Open(path)

	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	b, err := f.GetPlainText()
	if err != nil {
		return ""
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		return ""
	}

	return buf.String()

}
