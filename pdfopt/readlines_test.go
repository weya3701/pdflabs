package pdfopt

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// func ReadLines(path string) ([]string, error) { ... }  // 你的 ReadLines 函式定義

// TestReadLines 測試 ReadLines 函式
func TestReadLines(t *testing.T) {
	// 建立測試情境

	// 1. 成功讀取檔案的測試案例
	t.Run("成功讀取檔案", func(t *testing.T) {
		// 建立一個臨時檔案
		content := "line1\nline2\nline3"
		tempFile, err := os.CreateTemp("", "testfile") // 使用CreateTemp建立臨時檔案，避免路徑衝突
		if err != nil {
			t.Fatalf("創建臨時檔案失敗: %v", err) // 使用Fatalf，測試失敗時會停止測試
		}
		defer os.Remove(tempFile.Name()) // 在測試結束後刪除臨時檔案
		if _, err := tempFile.WriteString(content); err != nil {
			t.Fatalf("寫入臨時檔案失敗: %v", err)
		}
		tempFile.Close() // 關閉檔案，確保內容被寫入

		// 呼叫 ReadLines 函式
		lines, err := ReadLines(tempFile.Name())
		if err != nil {
			t.Errorf("ReadLines 失敗: %v", err)
		}

		// 期望的結果
		expectedLines := strings.Split(content, "\n")

		// 檢查結果
		if !reflect.DeepEqual(lines, expectedLines) { // 使用DeepEqual比較slice
			t.Errorf("讀取到的行不符合預期.\nGot: %v\nWant: %v", lines, expectedLines)
		}
	})

	// 2. 檔案不存在的測試案例
	t.Run("檔案不存在", func(t *testing.T) {
		// 呼叫 ReadLines 函式
		lines, err := ReadLines("nonexistent_file.txt")

		// 檢查錯誤
		if err == nil {
			t.Errorf("預期會發生錯誤，但沒有錯誤")
		}

		// 檢查 lines 是否為 nil
		if lines != nil {
			t.Errorf("在發生錯誤時，lines 應該為 nil，但得到: %v", lines)
		}
	})

	// 3. 讀取空檔案的測試案例
	t.Run("讀取空檔案", func(t *testing.T) {
		// 建立一個空的臨時檔案
		tempFile, err := os.CreateTemp("", "emptyfile")
		if err != nil {
			t.Fatalf("創建臨時檔案失敗: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// 呼叫 ReadLines 函式
		lines, err := ReadLines(tempFile.Name())
		if err != nil {
			t.Errorf("ReadLines 失敗: %v", err)
		}

		// 期望的結果
		expectedLines := []string{}

		// 檢查結果
		if !reflect.DeepEqual(lines, expectedLines) {
			t.Errorf("讀取到的行不符合預期.\nGot: %v\nWant: %v", lines, expectedLines)
		}
	})

	// 4. 檔案包含空行的測試案例
	t.Run("檔案包含空行", func(t *testing.T) {
		content := "line1\n\nline3"
		tempFile, err := os.CreateTemp("", "emptylinefile")
		if err != nil {
			t.Fatalf("創建臨時檔案失敗: %v", err)
		}
		defer os.Remove(tempFile.Name())
		if _, err := tempFile.WriteString(content); err != nil {
			t.Fatalf("寫入臨時檔案失敗: %v", err)
		}
		tempFile.Close()

		lines, err := ReadLines(tempFile.Name())
		if err != nil {
			t.Errorf("ReadLines 失敗: %v", err)
		}

		expectedLines := strings.Split(content, "\n")
		if !reflect.DeepEqual(lines, expectedLines) {
			t.Errorf("讀取到的行不符合預期.\nGot: %v\nWant: %v", lines, expectedLines)
		}
	})

}
