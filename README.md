# gotesseract
tesseract的一个小应用

![demo.gif](https://github.com/guoruibiao/gotesseract/raw/master/demo.gif)

测试代码
```go
package main

import (
	"net/http"
	"io"
		"fmt"
	"github.com/guoruibiao/gotesseract"
	"encoding/base64"
	"io/ioutil"
	"regexp"
		"os"
			"log"
)


var template string = `
<head>
<meta charset="UTF-8">
<title>Document</title>
<style>
body {
    display: -webkit-flex;
    display: flex;
    -webkit-justify-content: center;
    justify-content: center;
}
#tar_box {
    width: 500px;
    height: 500px;
    border: 1px solid red;
}
</style>
<body>
    <div id="container"> </div>
<script src="https://cdn.bootcss.com/jquery/1.12.0/jquery.js"></script>

<br><hr>
    <div>
        <h3>识别结果</h3>
        <p id="inspectresult"></p>
    </div>
</body>
<script>
document.addEventListener('paste', function (event) {
    console.log(event)
    var isChrome = false;
    if ( event.clipboardData || event.originalEvent ) {
        //not for ie11   某些chrome版本使用的是event.originalEvent
        var clipboardData = (event.clipboardData || event.originalEvent.clipboardData);
        if ( clipboardData.items ) {
            // for chrome
            var    items = clipboardData.items,
                len = items.length,
                blob = null;
            isChrome = true;
            //items.length比较有意思，初步判断是根据mime类型来的，即有几种mime类型，长度就是几（待验证）
            //如果粘贴纯文本，那么len=1，如果粘贴网页图片，len=2, items[0].type = 'text/plain', items[1].type = 'image/*'
            //如果使用截图工具粘贴图片，len=1, items[0].type = 'image/png'
            //如果粘贴纯文本+HTML，len=2, items[0].type = 'text/plain', items[1].type = 'text/html'
            // console.log('len:' + len);
            // console.log(items[0]);
            // console.log(items[1]);
            // console.log( 'items[0] kind:', items[0].kind );
            // console.log( 'items[0] MIME type:', items[0].type );
            // console.log( 'items[1] kind:', items[1].kind );
            // console.log( 'items[1] MIME type:', items[1].type );

            //阻止默认行为即不让剪贴板内容在div中显示出来
            event.preventDefault();

            //在items里找粘贴的image,据上面分析,需要循环
            for (var i = 0; i < len; i++) {
                if (items[i].type.indexOf("image") !== -1) {
                    // console.log(items[i]);
                    // console.log( typeof (items[i]));

                    //getAsFile()  此方法只是living standard  firefox ie11 并不支持
                    blob = items[i].getAsFile();
                }
            }
            if ( blob !== null ) {
                var reader = new FileReader();
                reader.onload = function (event) {
                    // event.target.result 即为图片的Base64编码字符串
                    var base64_str = event.target.result
                    //可以在这里写上传逻辑 直接将base64编码的字符串上传（可以尝试传入blob对象，看看后台程序能否解析）
console.log(base64_str);
var image = document.createElement('img');
image.id = "imagecontainer";
image.src= base64_str;
document.getElementById("container").appendChild(image);
                    uploadImgFromPaste(base64_str, 'paste', isChrome);
                }
                reader.readAsDataURL(blob);
            }
        } else {
            //for firefox
            setTimeout(function () {
                //设置setTimeout的原因是为了保证图片先插入到div里，然后去获取值
                var imgList = document.querySelectorAll('#tar_box img'),
                    len = imgList.length,
                    src_str = '',
                    i;
                for ( i = 0; i < len; i ++ ) {
                    if ( imgList[i].className !== 'my_img' ) {
                        //如果是截图那么src_str就是base64 如果是复制的其他网页图片那么src_str就是此图片在别人服务器的地址
                        src_str = imgList[i].src;
                    }
                }
                jquerywayupload(src_str);
            }, 1);
        }
    } else {
        //for ie11
        setTimeout(function () {
            var imgList = document.querySelectorAll('#tar_box img'),
                len = imgList.length,
                src_str = '',
                i;
            for ( i = 0; i < len; i ++ ) {
                if ( imgList[i].className !== 'my_img' ) {
                    src_str = imgList[i].src;
                }
            }
			jquerywayupload(src_str)
            // uploadImgFromPaste(src_str, 'paste', isChrome);
        }, 1);
    }
})

function jquerywayupload(file){
var formData = new FormData();
formData.append("file", $("#imagecontainer").val());
console.log(data);
	$.ajax({
		type : "POST",
		url : "/upload",
		dataType : "json",
		data : formData,
		success : function(data) {
			var res = data;
console.log(data);
			if (!res) {
				alert("上传失败！");
			} else {
				$("#inspectresult").html(res);
				alert("上传成功!");
			}
		},
		error : function(err) {
console.log(err);
			alert("由于网络原因，上传失败。");
		}
	});
}

function uploadImgFromPaste (file, type, isChrome) {
    var formData = new FormData();
    formData.append('file', file);
    formData.append('submission-type', type);

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/upload');
    xhr.onload = function () {
        if ( xhr.readyState === 4 ) {
            if ( xhr.status === 200 ) {
console.log("responseData:", xhr.responseText);
$("#inspectresult").html(xhr.responseText);
                // var data = JSON.parse( xhr.responseText ),
                //     tarBox = document.getElementById('tar_box');
                // if ( isChrome ) {
                //     var img = document.createElement('img');
                //     img.className = 'my_img';
                //     img.src = data.store_path;
                //     tarBox.appendChild(img);
                // } else {
                //     var imgList = document.querySelectorAll('#tar_box img'),
                //         len = imgList.length,
                //         i;
                //     for ( i = 0; i < len; i ++) {
                //         if ( imgList[i].className !== 'my_img' ) {
                //             imgList[i].className = 'my_img';
                //             imgList[i].src = data.store_path;
                //         }
                //     }
                // }

            } else {
				$("#inspectresult").html(xhr.statusText);
                console.log( xhr.statusText );
            }
        };
    };
    xhr.onerror = function (e) {
        console.log( xhr.statusText );
    }
    xhr.send(formData);
}
</script>
`

var tesseract2 *gotesseract.TesseractOCR

func picindex2(response http.ResponseWriter, request *http.Request) {
	io.WriteString(response, template)
}

func picupload2(w http.ResponseWriter, r *http.Request) {
	// 参考： https://studygolang.com/articles/5171
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		fmt.Println(r.PostFormValue("file"))
		// ddd, _ := base64.StdEncoding.DecodeString(r.PostFormValue("file")) //成图片文件并把文件写入到buffer
		filename := "/tmp/output.png"
		WriteFile(filename, r.PostFormValue("file"))
		fmt.Println(filename + "文件已保存到本地！")

		// 识别并进行展示
		result , err := tesseract2.Inspect(filename)
		fmt.Println(result, err)
		if err != nil {
			io.WriteString(w, err.Error())
		}else{
			io.WriteString(w, result)
		}
	}
}

//写入文件,保存
func WriteFile(path string, base64_image_content string) bool {

	b, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64_image_content)
	if !b {
		return false
	}

	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(base64_image_content), 2)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取
	fmt.Println("fileType=" + fileType)
	base64Str := re.ReplaceAllString(base64_image_content, "")


	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	err := ioutil.WriteFile(path, byte, 0666)
	if err != nil {
		log.Println(err)
	}

	return false
}

//判断文件是否存在

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true

}


func main() {
	tesseract2 = gotesseract.NewTesseractOCR()
	http.HandleFunc("/", picindex2)
	http.HandleFunc("/upload", picupload2)
	fmt.Println("图片识别服务运行中 http://localhost:9999/ ")
	http.ListenAndServe(":9999", nil)
}
```