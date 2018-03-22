# config

- `annie/config` 是一个目录
    - `api.go `是一个文件
    - `config.go` 是一个文件
    - `version.go` 是一个文件

一个目录三个文件，但是都在一个`代码包中- package config`

> 同一个package中的变量可以, 在同包不同文件中共用

---

## try

`go run main.go`

[`./examples/t1.go`](./examples/t1.go)

``` go
package examples

func hello(str string) string{ // 同包使用
	return str
}
```

[`./examples/t2.go`](./examples/t2.go)

``` go
package examples

import "fmt"

func Gethello(s string){ // 公用函数 被其他包使用
    fmt.Println(" package examples ")
    // hello 被使用
    fmt.Printf(hello("func hello")+
    ` in t1.go , t2.go can use func hello in Gethello `)
}
```

[./main.go](./main.go)

``` go
package main

import "github.com/chinanf-boy/annie-explain/examples"

func main(){
	examples.Gethello("local") // 
}
```
---

## annie/config/api.go

``` go
package config

// Bilibili
const (
	BILIBILI_API         string = "https://interface.bilibili.com/v2/playurl?"
	BILIBILI_BANGUMI_API string = "https://bangumi.bilibili.com/player/web_api/v2/playurl?"
	BILIBILI_TOKEN_API   string = "https://api.bilibili.com/x/player/playurl/token?"
)

```

## annie/config/config.go


``` go
package config

var (
	// Debug debug mode
	Debug bool
	// Version show version
	Version bool
	// InfoOnly Information only mode
	InfoOnly bool
	// Cookie http cookies
	Cookie string
	// Playlist download playlist
	Playlist bool
	// Refer use specified Referrer
	Refer string
	// Proxy HTTP proxy
	Proxy string
	// Socks5Proxy SOCKS5 proxy
	Socks5Proxy string
)

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Charset":  "UTF-8,*;q=0.5",
	"Accept-Encoding": "gzip,deflate,sdch",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36",
}

```

## annie/config/version.go

``` go
package config

// VERSION version of annie
const VERSION string = "0.5.0"

```