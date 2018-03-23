# Annie

「 👾 Annie is a fast, simple and clean video downloader built with Go.

 」

[![explain](http://llever.com/explain.svg)](https://github.com/chinanf-boy/Source-Explain)
    
Explanation

> "version": "1.0.0"

[github source](https://github.com/iawia002/annie)

~~[english](./README.en.md)~~

---

⏬下载啦下载啦

- 本次从小示例出发 比如「抖音」

``` js
$ annie https://www.douyin.com/share/video/6509219899754155272

 Site:    抖音 douyin.com
Title:    好冷  逢考必过
 Type:    mp4
 Size:    2.63 MiB (2762719 Bytes)

 741.70 KiB / 2.63 MiB [=========>--------------------------]  27.49% 1.98 MiB/s
```

> ⏰本项目-`代码示例`运行，请正确-clone本项目到-GOPATH

---

本目录



---

## main

> 作为`go`项目，要想作为命令行，被打包📦，被编译，需要具备

``` go
package main // 作为 编译入口
```

`annie/main.go`

``` go
package main

import (
	"flag" // 作为 go 内置的命令解析
	"fmt" // 终端-输出
	"net/url" // 网络/网址

	"github.com/iawia002/annie/config" // 作者-默认配置
	"github.com/iawia002/annie/extractors" // 作者-对应网站解析
	"github.com/iawia002/annie/utils" // 作者-工具集
)

func init() { // 初始化 运行顺序派在 func main() 前

// 命令行解析大概描述一下 主要 4个变量
    // flag.***( 1: 默认配置变量, 2: 命令行选项{-p -i 之类}, 3: 默认值, 4: 说明描述 )
	flag.BoolVar(&config.Debug, "d", false, "Debug mode") // 
	flag.BoolVar(&config.Version, "v", false, "Show version")
	flag.BoolVar(&config.InfoOnly, "i", false, "Information only")
	flag.StringVar(&config.Cookie, "c", "", "Cookie")
	flag.BoolVar(&config.Playlist, "p", false, "Download playlist")
	flag.StringVar(&config.Refer, "r", "", "Use specified Referrer")
	flag.StringVar(&config.Proxy, "x", "", "HTTP proxy")
    flag.StringVar(&config.Socks5Proxy, "s", "", "SOCKS5 proxy")

// 其实 在每个 package 包里面，最先运行的函数都是自定义的 func init()
}

func main() {
// > 记得我们的小示例吗 `annie https://www.douyin.com/share/video/6509219899754155272` 
	flag.Parse() // 命令行选项定义后，要开始解析 Parse 启动解析
	args := flag.Args() // 除开-定义选项 其他-用户输入命令选项
	if config.Version { // 如果 你选项中具备 -v 那么它
		fmt.Printf(
			"annie: version %s, A simple and clean video downloader.\n", config.VERSION,
		)
		return // 只到这里就结束了
	}
	if len(args) < 1 { // 没有url-输出错误，但其实应该显示一下粒子
		fmt.Println("error")
		return
	}
	videoURL := args[0] // 拿到url
	u, err := url.ParseRequestURI(videoURL) // 内置库 url解析
	if err != nil {
		fmt.Println(err)
		return
    }
    

// 下载抖音的视频

    domain := utils.Domain(u.Host) // 拿到-域

// 在这里-往上⬆️, 做了两件事
// 1. 控制命令行选项，达到获取配置-config
// 2. 确认域名

	switch domain {
	case "douyin":
		extractors.Douyin(videoURL) // 完整网址给到-抖音解析器
	// case "bilibili":
	// 	extractors.Bilibili(videoURL)
    // 。。。

```

- [utils.Domain](#domain)

> 从给予的URL中获得域名

- [config explain ](./config.explain.md)

> 有关-`"github.com/iawia002/annie/config"` 的声明与定义

- [extractors explain ](./extractors.explain.md)

> 提取库-定义不同网址的解析

---

## utils

下面工具函数-一般不细看，只要知道什么功能就行

<details>


``` go
package utils

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/iawia002/annie/request"
)

```

- [MatchOneOf](#matchoneof)

- [MatchAll](#matchall)

- [FileSize](#filesize)

- [Domain](#domain)

- [FileName](#filename)

- [FilePath](#filepath)

- [StringInSlice](#stringinslice)

- [GetNameAndExt](#getnameandext)

- [Md5](#md5)

- [M3u8URLs](#m3u8urls)


``` go
// MatchOneOf match one of the patterns 匹配一个就返回
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

```

### matchall

``` go
// MatchAll return all matching results 匹配所有
func MatchAll(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern) // 内置
	value := re.FindAllStringSubmatch(text, -1)
	return value
}

```

### filesize

``` go
// FileSize return the file size of the specified path file 返回指定路径文件的文件大小
func FileSize(filePath string) int64 {
	file, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return 0
	}
	return file.Size()
}

```

### domain

``` go
// Domain get the domain of given URL 从给予的URL中获得域名
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)[1]
	return domain
}

```

### filename

``` go
// FileName Converts a string to a valid filename 将字符串转换为有效的文件名
func FileName(name string) string {
	// FIXME(iawia002) file name can't have /
	name = strings.Replace(name, "/", " ", -1)
	name = strings.Replace(name, "|", "-", -1)
	name = strings.Replace(name, ":", "：", -1)
	if runtime.GOOS == "windows" {
		winSymbols := []string{
			"\"", "?", "*", "\\", "<", ">",
		}
		for _, symbol := range winSymbols {
			name = strings.Replace(name, symbol, " ", -1)
		}
	}
	return name
}

```

### filepath

``` go
// FilePath gen valid filename 生成有效的文件名
func FilePath(name, ext string, escape bool) string {
	fileName := fmt.Sprintf("%s.%s", name, ext)
	if escape {
		fileName = FileName(fileName)
	}
	return fileName
}

```

### stringinslice

``` go
// StringInSlice if a string is in the list 如果一个字符串在列表中
func StringInSlice(str string, list []string) bool {
	for _, a := range list {
		if a == str {
			return true
		}
	}
	return false
}

```

### getnameandext

``` go
// GetNameAndExt return the name and ext of the URL 返回URL的名称和分机号
// https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg ->
// 1f5a87801a0711e898b12b640777720f, jpg
func GetNameAndExt(uri string) (string, string) {
	u, _ := url.ParseRequestURI(uri)
	s := strings.Split(u.Path, "/")
	filename := strings.Split(s[len(s)-1], ".")
	if len(filename) > 1 {
		return filename[0], filename[1]
	}
	// Image url like this
	// https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg/w650
	// has no suffix
	contentType := request.ContentType(uri, uri)
	return filename[0], strings.Split(contentType, "/")[1]
}

```

### md5

``` go
// Md5 md5 hash 哈希
func Md5(text string) string {
	sign := md5.New()
	sign.Write([]byte(text))
	return fmt.Sprintf("%x", sign.Sum(nil))
}

```

### m3u8urls

``` go
// M3u8URLs get all urls from m3u8 url 从m3u8网址获取所有网址
func M3u8URLs(uri string) []string {
	html := request.Get(uri)
	lines := strings.Split(html, "\n")
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "http") {
				urls = append(urls, line)
			} else {
				base, _ := url.Parse(uri)
				u, _ := url.Parse(line)
				urls = append(urls, fmt.Sprintf("%s", base.ResolveReference(u)))
			}
		}
	}
	return urls
}
```
</details>