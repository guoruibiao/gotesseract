package main

/**
 * 客服群老是有这样的反馈：用户A充值了，却没到账。然后给个截图，有一串交易单号。
 *     但是这种一般都是用户A将钱充值到了自己的小号，然后拿大号过来询问，说没到账。
 *
 *  每次都要手动输入那串交易单号，很长而且容易出错。于是有了下面的代码，截个图自动识别下。
 */
import (
	"github.com/guoruibiao/commands"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"log"
	"github.com/fsnotify/fsnotify"
)

// 借助fsnotify过滤出新增的文件内容
// 借助tesseract 识别图片截图内容
// 利用pbcopy将文件内容拷贝到系统剪切板
// do anything else.


type TesseractHelper struct {
	commander *commands.Commands
	WatchPath string
	Result string
}

func NewTesseractHelper(watchpath string) (*TesseractHelper){
	return &TesseractHelper{
		commander: commands.New(),
		WatchPath:watchpath,
	}
}

// 懒得改名了 其实也没必要有返回值
func (this *TesseractHelper) getNewlyFilename() (string, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 为了让主进程一直等待，实现go协程的运行流程。
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				//if event.Op&fsnotify.Write == fsnotify.Write {
				//	log.Println("modified file: ", event.Name)
				//}
			    if event.Op & fsnotify.Create == fsnotify.Create {
			    	log.Println("newly add file: ", event.Name)
			    	// 识别处理
			    	filename := event.Name[8:]
				    // 可以考虑用 go func(){}()的形式，不过好像没啥必要
			    	this.doTesseractInspect(filename)
			    	fmt.Println("识别结果为: ", this.output())
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)
			}
		}
	}()

	err = watcher.Add(this.WatchPath)
	if err != nil {
		log.Fatal(err)
	}
	// 如上，让主进程一直等待，不至于运行完就退出了。
	<-done
	return "/tmp/Xnip2019-11-01_08-26-29.jpg", nil
}

func (this *TesseractHelper) doTesseractInspect(newlyfile string) {
	// 生成一个临时文件
	tempDetectFile := "/tmp/tesseract-detect-result"
	defer func(){
		// 删除临时文件
		this.commander.Run("rm", tempDetectFile + ".txt")
	}()
	// 调用tesseract 识别
    this.commander.Run("tesseract", newlyfile, tempDetectFile)
	file, err := os.Open(tempDetectFile+".txt")
	if err != nil {
		fmt.Println("open detected file failed, ", err.Error())
		os.Exit(0)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read detected files failed, ", err.Error())
		os.Exit(0)
	}
	this.Result = strings.Replace(string(bytes), "\n", "", -1)
}

func (this *TesseractHelper) output() (string) {
    // 为了自用 所以随意添加一些东西
	return fmt.Sprintf(`select * from payment_order_detail where transaction_id='%s';`, this.Result)
}

func main() {
	th := NewTesseractHelper("/tmp")
	th.getNewlyFilename()
	//th.doTesseractInspect()
    // tesseract-detect-result.txt.txt
    //fmt.Println("识别结果为：", th.output())
}
