package pdfopt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dslipak/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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

// desc := "pos:bl, scale:0.3, fontname:Helvetica, points:12, c:1 0 0, rot:0"
// desc = fmt.Sprintf("%s, off: %f %f", desc, *params.Y, *params.X)

type Coordinate struct {
	MarginLeft   float64
	MarginTop    float64
	MarginBottom float64
	PageHeight   float64
	LineHeight   float64
	ContentWidth float64
	Position     string
	Scale        float64
	Fontname     string
	Points       int
	Rot          int
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

// 切割來源檔，依頁數儲存
func splitSourceFile(tmpPrefix string, fileIdx int, filename string, splitConds []int) string {

	var writerIO io.Writer

	pdffile, _ := os.Open(filename)
	defer pdffile.Close()

	// 定義一個接受 io.ReadSeeker 的變數
	// var seeker io.ReadSeeker = pdffile

	var writer gopdf.GoPdf
	var err error = nil
	var outputFilename string
	fmt.Println(tmpPrefix, filename, splitConds)
	outputFilename = fmt.Sprintf("%s/%s%d.pdf", "pdfsplittmp", tmpPrefix, fileIdx)

	// 建立暫存PDF檔
	writerIO, err = os.Create(outputFilename)
	if err != nil {
		fmt.Println("Can not create file: %w", err)
	}

	// init pdf writer
	writer = gopdf.GoPdf{}
	writer.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = writer.AddTTFFont("myfont", "./fonts/ArialBlack.ttf")

	if err != nil {
		fmt.Printf("Can not add font %v", err)
	}

	for i := splitConds[0]; i <= splitConds[1]; i++ {

		// tpl := writer.ImportPageStream(&seeker, i, "/MediaBox")
		tpl := writer.ImportPage(filename, i, "/MediaBox")
		writer.AddPage()
		writer.UseImportedTemplate(tpl, 0, 0, 595.28, 841.89)
	}

	_, err = writer.WriteTo(writerIO)
	if err != nil {
		fmt.Println("Can not write pdf file: %w", err)
	}

	defer writer.Close()
	return outputFilename
}

// 檔案切割功能，將一個PDF檔依頁數切割成數個暫存檔，輸入:檔名(string)，輸出:檔名列表([]string)
func splitPDFfile(filename string) []string {
	var rsp []string = []string{""}
	var fileCount [][]int
	fmt.Println(filename)

	fileCount, err := splitCount(filename, fileCount)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fileCount)

	return rsp
}

func splitCount(filename string, fileCount [][]int) ([][]int, error) {

	// 使用getPDFpages函式取得檔案總頁數，如果總頁數超過100頁時，用以下方式處理
	// 1）以每個檔頁50頁為前題，計算需要分割為多少檔案
	// 2) 將切割資訊用[][]int格式儲存[[1, 50], [51, 100], [101, 150]....]
	// var fileCount [][]int # FIXME
	var err error = nil

	fmt.Println("正在處理檔案：", filename)

	// 1. 使用 getPDFpages 函式取得檔案總頁數
	// 假設 getPDFpages 的定義為 func getPDFpages(string) (int, error)
	totalPages := getPDFpages(filename)

	// 2. 判斷邏輯：如果總頁數超過 100 頁
	if totalPages > 100 {
		pageSize := 100
		for start := 1; start <= totalPages; start += pageSize {
			end := start + pageSize - 1

			// 確保最後一個分段不會超過總頁數
			if end > totalPages {
				end = totalPages
			}

			// 將範圍 [開始頁, 結束頁] 加入結果
			fileCount = append(fileCount, []int{start, end})
		}
	} else {
		// 如果不滿 100 頁，通常會視為一個分段或不處理，這裡預設為單一分段
		fileCount = append(fileCount, []int{1, totalPages})
	}

	return fileCount, err
}

func mergeFiles(files []string, outputFile string) {

	var err error
	var writerIO io.Writer
	fmt.Println(writerIO)
	writerIO, err = os.Create(outputFile)
	if err != nil {
		fmt.Println("Can not create file %w", err)
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = pdf.AddTTFFont("myfont", "./fonts/ArialBlack.ttf")
	if err != nil {
		fmt.Println("Can not add fonts: ", err)
	}

	for _, file := range files {
		if file != "" {
			// 取得頁數
			pageNums := getPDFpages(file)
			for i := 1; i <= pageNums; i++ {
				tpl := pdf.ImportPage(file, i, "/MediaBox")
				pdf.AddPage()
				pdf.UseImportedTemplate(tpl, 0, 0, 595.28, 841.89)
			}

			// err := pdf.ImportPagesFromSource(file, "/MediaBox")
			// fmt.Println("Can not write pdf file: ", err)
		}
	}

	_, err = pdf.WriteTo(writerIO)
	if err != nil {
		fmt.Println("Can not write pdf file: %w", err)
	}
	// err = pdf.WritePdf(outputFile)
	// if err != nil {
	// 	log.Printf("Can not wirte pdf file: %v", err)
	// }

	defer pdf.Close()
}

var (
	fontData []byte
	fontOnce sync.Once
)

func loadFont() {
	data, err := os.ReadFile("./fonts/ArialBlack.ttf")
	if err == nil {
		fontData = data
	}
}

func SecureMerge(files []string, outputPath string) error {
	// pdfcpu 會直接在底層處理交叉引用表（XRef Table）的合併
	// 這比 gopdf 的 ImportPage 穩定得多，且不會有字型亂碼問題
	err := api.MergeCreateFile(files, outputPath, false, nil)
	if err != nil {
		fmt.Printf("合併失敗: %v\n", err)
	}
	return err
}

func Merge2fileOptimized(params PdfOptParameter) error {
	var err error = nil

	files := []string{params.OriginFilename, params.SourceFilename}
	err = SecureMerge(files, params.OutputFilename)

	return err
}

// Merge 2 pdf file.
// FIXME. 改用Merge2fileOptimized
func Merge2file(params PdfOptParameter) {

	var tmpCnt int = 1
	var tmpFiles []string

	files := []string{
		params.OriginFilename,
		params.SourceFilename,
	}

	var fileCount [][]int
	var err error = nil

	writer := gopdf.GoPdf{}
	writer.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = writer.AddTTFFont("myfont", "./fonts/ArialBlack.ttf")
	defer writer.Close()
	if err != nil {
		fmt.Printf("Can not add font %v", err)
	}

	for _, file := range files {
		fileCount, _ := splitCount(file, fileCount)
		fmt.Println("fileCount: ", fileCount)
		for _, splitConds := range fileCount {
			fname := splitSourceFile("testPdf", tmpCnt, file, splitConds)
			tmpFiles = append(tmpFiles, fname)
			tmpCnt = tmpCnt + 1
		}
	}
	fmt.Println(tmpFiles)

	mergeFiles(tmpFiles, params.OutputFilename)
}

func StickTagsTest(params PdfOptParameter) {
	// 設定輸入與輸出路徑
	inFile := params.OriginFilename
	outFile := params.OutputFilename

	// 1. 設定文字內容
	text := params.Content

	// 2. 設定參數描述字串 (Description)
	// 這是 pdfcpu 的核心設計，透過字串控制位置、字型、顏色、旋轉等
	// pos:c (置中), s:1.0 (縮放), f:Helvetica (字型), points:48 (大小), c:1 0 0 (紅色)
	desc := fmt.Sprintf(
		"pos:%s, scale:%f, fontname:%s, points:%d, c:%d %d %d, rot:%d",
		params.PdfCoordinate.Position,
		params.PdfCoordinate.Scale,
		params.PdfCoordinate.Fontname,
		params.PdfCoordinate.Points,
		0,
		0,
		0,
		params.PdfCoordinate.Rot,
	)
	desc = fmt.Sprintf(
		"%s, off: %f %f",
		desc,
		*params.Y,
		*params.X,
	)

	// 3. 設定要插入的頁面 (nil 代表所有頁面)
	selectedPages := []string{"1"}

	// 4. 執行插入 (onTop = true 代表當作圖章，覆蓋在內容之上)
	onTop := true

	// 使用預設配置
	conf := model.NewDefaultConfiguration()

	err := api.AddTextWatermarksFile(inFile, outFile, selectedPages, onTop, text, desc, conf)
	if err != nil {
		log.Fatalf("插入文字失敗: %v", err)
	}

	fmt.Println("成功生成:", outFile)
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
		// tpl := pdf.ImportPage(params.OriginFilename, i, "/MediaBox")
		// pdf.AddPage()
		// pdf.UseImportedTemplate(tpl, 0, 0, 595.28, 841.89)
		if i == params.PageNum {
			tpl := pdf.ImportPage(params.OriginFilename, i, "/MediaBox")
			pdf.AddPage()
			pdf.UseImportedTemplate(tpl, 0, 0, 595.28, 841.89)
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
