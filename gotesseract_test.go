package gotesseract

import (
	"testing"
			)

func TestTesseractOCR_Detect(t *testing.T) {
	tesseract := NewTesseractOCR()
	filename := "/tmp/pic.jpg"
	result, err := tesseract.Inspect(filename)
	if err != nil {
		t.Error(err)
	}else{
		t.Log(result)
	}
	
}