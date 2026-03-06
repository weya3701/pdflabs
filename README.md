### PDF加工工具

#### 安裝

* git clone git@github.com:weya3701/pdflabs.git

* go build -o pdflab

#### 操作功能說明

* 合併模式(Merge)

    * 參數：

        * action: Merge
        
        * origin: 原始PDF檔

        * source: 合併來源PDF檔

        * output: 輸出PDF檔名

            pdflab --action=Merge --origin=<file.pdf> --source=<file.pdf> --output=<file.pdf>


* 疊加寫入、貼標模式(StickTags)

    * 參數：

        * action: StickTags

        * origin: 原始PDF檔

        * pagenum: 張貼頁數

        * x: 張貼座標X軸

        * y: 張貼座標Y軸

        * content: 張貼內容

            pdflab --action=StickTags --origin=<file.pdf> --pagenum=<1> --x=500 --y=500 --content=<Tag, content>

* 文字檔內容轉PDF(WriteContentFromFile)

    * 參數：

        * action: WriteContentFromFile

        * cf: 內容文字檔

        * output: 輸出PDF檔

            pdflab --action=WriteContentFromFile --cf=<content.txt> --output=<output.pdf>
