package gotesseract

import (
	"os"
	"github.com/guoruibiao/commands"
	"io/ioutil"
		)

/**
 * 对外提供"代理接口"
 */

type TesseractOCR struct {
	commander commands.Commands
}

func NewTesseractOCR() *TesseractOCR {
	commander := commands.New()
	return &TesseractOCR{
		commander: *commander,
	}
}

func (this *TesseractOCR) Inspect(filename string) (result string, err error) {
	_, err = os.Stat(filename)
	if err != nil {
		return
	}
	
	outputfile := "/tmp/hello"
	// TODO 可能需要做下健壮性判断
	params := []string{filename, outputfile}
	this.commander.OuterRun("tesseract", params...)
	
	file, err := os.Open(outputfile+".txt")
	if err != nil {
		return
	}
	
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	
	os.Remove(outputfile + ".txt")
	return string(bytes), nil
}
