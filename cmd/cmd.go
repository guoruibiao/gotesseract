package main

import (
	"flag"
	"github.com/guoruibiao/gotesseract"
	"fmt"
)

func main() {
	input := flag.String("i", "/tmp/input", "待检测图片完整路径")
	flag.Parse()
	
	tesseract := gotesseract.NewTesseractOCR()
	result, err := tesseract.Inspect(*input)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(result)
	}
}
