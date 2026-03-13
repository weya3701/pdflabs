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

// 請用golang實作一個名為ReadLines讀取檔案路徑內容並回傳
// 回傳格式請使用[]string
// 回傳資料[]string, error
// 函式引入參數請使用path string型態
// **程式碼說明:**
//
// 1.  **`package main`**:  聲明程式屬於 `main` package。
// 2.  **`import ("bufio", "os")`**:  導入所需的套件：
//   - `bufio`:  用於緩衝 I/O 操作，方便逐行讀取檔案。
//   - `os`:  提供檔案操作相關的功能，例如開啟、關閉檔案。
//
// 3.  **`func ReadLines(path string) ([]string, error)`**:  定義 `ReadLines` 函數，它接受一個檔案路徑 `path` (字串) 作為輸入，並回傳兩個值：
//   - `[]string`:  一個字串 slice，包含檔案中的每一行。
//   - `error`:  一個 error 物件，如果沒有錯誤則為 `nil`。
//
// 4.  **`var lines []string`**:  宣告一個字串 slice `lines`，用於儲存檔案中的每一行。
// 5.  **`file, err := os.Open(path)`**:  嘗試開啟指定路徑的檔案。如果發生錯誤（例如：檔案不存在、權限不足等），`err` 會被賦值。
// 6.  **`if err != nil { return nil, err }`**:  如果開啟檔案時發生錯誤，立即回傳 `nil` (空 slice) 和 `err`。
// 7.  **`defer file.Close()`**:  使用 `defer` 語句確保在函數結束時關閉檔案，即使發生錯誤也是如此。  這是一種良好的程式設計習慣，可以避免資源洩漏。
// 8.  **`scanner := bufio.NewScanner(file)`**:  建立一個 `bufio.Scanner`，用於逐行讀取檔案內容。
// 9.  **`for scanner.Scan() { ... }`**:  迴圈遍歷檔案的每一行：
//   - `scanner.Scan()`:  讀取下一行。如果成功，`scanner.Text()` 可以取得該行內容。
//   - `lines = append(lines, scanner.Text())`:  將讀取的行添加到 `lines` slice 中。
//
// 10. **`if err := scanner.Err(); err != nil { return nil, err }`**:  在迴圈結束後，檢查 `scanner` 過程中是否發生錯誤 (例如：I/O 錯誤)。 如果發生錯誤，則返回 `nil` 和 `err`。
// 11. **`return lines, nil`**:  如果一切順利，則回傳儲存了檔案所有行的 `lines` slice 和 `nil` (表示沒有錯誤)。
// ************************************測試階段*************************************************
// 為Readliness函式撰寫一個測試函式
func ReadLines(path string) ([]string, error) {
	// 1. 打開檔案
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// 確保函式結束時關閉檔案資源
	defer file.Close()

	var lines []string = []string{}
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
