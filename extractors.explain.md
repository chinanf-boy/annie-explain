# annie-extractors

[annie 支持的网站](https://github.com/iawia002/annie#supported-sites)

---

- [douyin](#douyin)

---

类型和结构定义，是 go 的特性，这些特性对于-程序的稳定和递进有影响

但，要适度。一般的类型也就那几个，string, number 之类

其中`json`之类的内置解析，go 会教给你一些招式，让你可以快速击倒 -`json🐶`

> try `go run main.go json`

然后继续往下看

---

## extractors-Douyin

<details>

``` go
package extractors

import (
	"encoding/json"

	"github.com/iawia002/annie/downloader"
	"github.com/iawia002/annie/request"
	"github.com/iawia002/annie/utils"
)

// json 数据
type douyinVideoURLData struct {
	URLList []string `json:"url_list"`
}

type douyinVideoData struct {
	PlayAddr     douyinVideoURLData `json:"play_addr"`
	RealPlayAddr string             `json:"real_play_addr"`
}

type douyinData struct {
	Video douyinVideoData `json:"video"`
	Desc  string          `json:"desc"`
}
// 对应下来-就是 这样的数据结构
// {
//     desc:"",
//     video:{
//         play_addr:{
//             url_list:"",
//         }
//         real_play_addr:"",
//     }
// }



// Douyin download function
func Douyin(url string) downloader.VideoData {
	html := request.Get(url) // 从url-获取完整的 <html>..html内容..</html>
    vData := utils.MatchOneOf(html, `var data = \[(.*?)\];`)[1] // 从 html 找出 匹配
    // 所谓的 命令行选项 
    // annie url
    // 就能下载所想要的视频或其他，都是对网页内容的解析 获取真实的下载地址
    // 从而进行下载

	var dataDict douyinData 

	json.Unmarshal([]byte(vData), &dataDict) // 解析从 html:string -> json 获得 的匹配项

    size := request.Size(dataDict.Video.RealPlayAddr, url) 
    // 从 要下载的 网址 head.Content-Length 知道下载的 文件大小

	urlData := downloader.URLData{
		URL:  dataDict.Video.RealPlayAddr,
		Size: size,
		Ext:  "mp4",
	} // 下载器 定义的 数据块
	data := downloader.VideoData{ 
		Site:  "抖音 douyin.com",
		Title: utils.FileName(dataDict.Desc),
		Type:  "video",
		URLs:  []downloader.URLData{urlData}, // 有时候 是 一串视频流
		Size:  size,
    } // 下载器 定义的 视频数据 
    
	data.Download(url) // 开始下载
	return data
}
```

- `json.Unmarshal([]byte(vData), &dataDict)`

> 可以试试 `go run main.go json` [`./examples/t4-json.go`](./examples/t4-json.go) 解析从 string -> json 获得 的匹配项



</details>

### Download

`annie/downloader/downloader.go`

代码 121-139

<details>

``` go
func (data VideoData) Download(refer string) {
	if data.Size == 0 {
		data.calculateTotalSize()
	}
	data.printInfo()
	if config.InfoOnly {
		return
    }
    // pb 是 进度条库
	bar := pb.New64(data.Size).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.SetMaxWidth(1000)
	bar.Start()
	if len(data.URLs) == 1 { // 本次的例子只有一个视频流
		// only one fragment
		data.urlSave(data.URLs[0], refer, data.Title, bar)
		bar.Finish()
    }
    // 。。
}
```

- `pb` - [ 可以看看 》》github source](https://github.com/cheggaaa/pb)

> 作为下载状态的进度条

- urlSave

> urlSave( 真实视频网址, 用户输入网址, 视频名, 进度条实例 )
</details>

---

### urlSave

`annie/downloader/downloader.go`

代码 67-19

<details>


``` go
// urlSave save url file
func (data VideoData) urlSave(
	urlData URLData, refer, fileName string, bar *pb.ProgressBar,
) {
	filePath := utils.FilePath(fileName, urlData.Ext, false) // 组合-本地下载路径-文件名
	fileSize := utils.FileSize(filePath) //  文件大小
	// TODO: Live video URLs will not return the size // 直播不会返回大小
	if fileSize == urlData.Size { // 如果相等 自然下载完
		fmt.Printf("%s: file already exists, skipping\n", filePath)
		bar.Add64(fileSize)
		return
	}
	tempFilePath := filePath + ".download"
	tempFileSize := utils.FileSize(tempFilePath)
	headers := map[string]string{
		"Referer": refer, // 用户输入网址
	}
	var file *os.File
    if tempFileSize > 0 { // 还是
        //状态-显示
		// range start from zero
		headers["Range"] = fmt.Sprintf("bytes=%d-", tempFileSize)
		file, _ = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644) 
		bar.Add64(tempFileSize)
	} else {
        // 新建文件
		file, _ = os.Create(tempFilePath)
	}

	// close and rename temp file at the end of this function
	// must be done here to avoid the following request error to cause the file can't close properly
	defer func() { 
        // 在结束本函数 时 defer 后面的函数 「注意⚠️是函数运行 不仅仅是定义/声明」 都会运行，所以一般用来关闭 文件 数据库 连接 的关闭工作
		file.Close()
		// must close the file before rename or it will cause `The process cannot access the file because it is being used by another process.` error.
		err := os.Rename(tempFilePath, filePath)
		if err != nil {
			log.Fatal(err)
		}
	}() // <--- 运行

	res := request.Request("GET", urlData.URL, nil, headers)
	if res.StatusCode >= 400 {
        // color 是 颜色库 帮-显示信息-加颜色
		red := color.New(color.FgRed)
		log.Print(urlData.URL)
		log.Fatal(red.Sprintf("HTTP error: %d", res.StatusCode))
	}
	defer res.Body.Close()
    writer := io.MultiWriter(file, bar)
    // go语言中 - 同时输出到文件和控制台(命令行）

    // 请注意，io.Copy从输入读取32kb（最大值）并将它们写入输出，然后重复。 也就是说-一步到位，不用管了
	_, copyErr := io.Copy(writer, res.Body) // res.Body 即是视频流本身 复制给文件 和 进度条
	if copyErr != nil { // 错误
		log.Fatal(fmt.Sprintf("Error while downloading: %s, %s", urlData.URL, copyErr))
	}

```

- `request.Request("GET", urlData.URL, nil, headers)`

> 重中之重, 在这步之后我们就拿到-真实数据和状态了

> Request( 网页请求方式, 真实网址, io.Reader ?? ,请求头)

- `io.MultiWriter(file, bar)` - `io.Copy(writer, res.Body)`

> 你可以试试 `go run main.go pb` 查看[相关代码](./examples/t3-pb.go)

- `color` 颜色库

> [github source ](https://github.com/fatih/color)


</details>


### request-Request

`annie/request/request.go`

代码 24-107

<details>

``` go
// Request base request
func Request(
	method, url string, body io.Reader, headers map[string]string,
) *http.Response {
	transport := &http.Transport{
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	if config.Proxy != "" {
    // 添加代理
		var httpProxy, err = netURL.Parse(config.Proxy)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(httpProxy)
	}
	if config.Socks5Proxy != "" {
    // socks-代理
		dialer, err := proxy.SOCKS5(
			"tcp",
			config.Socks5Proxy,
			nil,
			&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			panic(err)
		}
		transport.Dial = dialer.Dial
    }
    // 请求客户端 - 使用 Client.Do请求
	client := &http.Client{
		Timeout:   time.Second * 100,
		Transport: transport,
    }
// 定义好- 网址的请求信息🆕
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Print(url)
		panic(err)
	}
	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Referer", url)
	if config.Cookie != "" {
		var cookie string
		if _, fileErr := os.Stat(config.Cookie); fileErr == nil {
			// Cookie is a file
			data, _ := ioutil.ReadFile(config.Cookie)
			cookie = string(data)
		} else {
			// Just strings
			cookie = config.Cookie
		}
		req.Header.Set("Cookie", cookie)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if config.Refer != "" {
		req.Header.Set("Referer", config.Refer)
    }
// ———————— 定义完成✅

    // 使用 Client.Do请求
	res, err := client.Do(req)
	if err != nil {
		log.Print(url)
		panic(err)
    }
    // 调试时-状态显示
	if config.Debug {
		blue := color.New(color.FgBlue)
		fmt.Println()
		blue.Printf("URL:         ")
		fmt.Printf("%s\n", url)
		blue.Printf("Method:      ")
		fmt.Printf("%s\n", method)
		blue.Printf("Headers:     ")
		pretty.Printf("%# v\n", req.Header)
		blue.Printf("Status Code: ")
		if res.StatusCode >= 400 {
			color.Red("%d", res.StatusCode)
		} else {
			color.Green("%d", res.StatusCode)
		}
    }
    // 返回请求结果 
	return res

	    // _, copyErr := io.Copy(writer, res.Body) // res.Body 即是视频流本身 复制给文件 和 进度条
// 上小节的
```


</details>