## bilibili-extractors

> 关于 bilibili 的 解析

---

那么我们重新开始一段旅程吧, 如果你看了⬆️[小节](./readme.md), 就会知道上小节描述>知道-哪个网站, 然后把`用户输入-url-` 给到 对应的`extractors`

---

`annie/main.go`

``` go
    extractors.Bilibili(videoURL)
```

---

本次粒子🌰

``` bash
annie -c cookies.txt https://www.bilibili.com/video/av20203945/
```

> 哔哩哔哩 来了

---


---

## 1. Bilibili

`annie/extractors/bilibili.go`


<details>

``` go
type bilibiliOptions struct {
	Bangumi  bool
	Subtitle string
	Aid      string
	Cid      string
	HTML     string
}

// Bilibili download function
func Bilibili(url string) {
	var options bilibiliOptions
	if strings.Contains(url, "bangumi") { // bangumi 是 哔哩哔哩-视频列表-页 有多个视频文件 
		options.Bangumi = true
	}
	html := request.Get(url) // 获取 <html> .. </html> : string
	if !config.Playlist { // 用户命令行 `annie -p` 设置 Download playlist 
		options.HTML = html
		data, err := getMultiPageData(html) // 从 html 中 获取视频列表 data: multiPage
        if err == nil && !options.Bangumi { // 并不在 哔哩哔哩-视频列表-页 单视频
// 但本次例子，有 options.Bangumi == true

			// handle URL that has a playlist, mainly for unified titles
			// bangumi doesn't need this
			pageString := utils.MatchOneOf(url, `\?p=(\d+)`) // 哔哩哔哩-有时还分 1-part/2-part 就下一part
			var p int
			if pageString == nil {
				// https://www.bilibili.com/video/av20827366/
				p = 1
			} else {
                // https://www.bilibili.com/video/av20827366/?p=2
				p, _ = strconv.Atoi(pageString[1])// 将字符串转换为十进制整数，即：ParseInt(s, 10, 0) 的简写）
//参考 http://www.cnblogs.com/golove/p/3262925.html 
			}
			page := data.VideoData.Pages[p-1]
			options.Aid = data.Aid
			options.Cid = strconv.Itoa(page.Cid) // 将整数转换为十进制字符串形式（即：FormatInt(i, 10) 的简写
			options.Subtitle = page.Part
		}
		bilibiliDownload(url, options) // 单一视频网址
		return
	}
	if options.Bangumi {
		dataString := utils.MatchOneOf(html, `window.__INITIAL_STATE__=(.+?);`)[1]
		var data bangumiData
		json.Unmarshal([]byte(dataString), &data)
		for _, u := range data.EpList { // 循环-单一视频网址
			bilibiliDownload(
				fmt.Sprintf("https://www.bilibili.com/bangumi/play/ep%d", u.EpID), options,
			)
		}
	} else { // 如果输入了 `annie -p` 但 却不在 哔哩哔哩-视频列表-页
		data, err := getMultiPageData(html)
		if err != nil {
			// this page has no playlist
			options.HTML = html
			bilibiliDownload(url, options)
			return
		}
		// https://www.bilibili.com/video/av20827366/?p=1
		for _, u := range data.VideoData.Pages {  // 哔哩哔哩-有时还分 1-part/2-part 全-part下
			options.Aid = data.Aid
			options.Cid = strconv.Itoa(u.Cid)
			options.Subtitle = u.Part
			bilibiliDownload(url, options)
		}
	}
}
```

- `bilibiliDownload`

> 那么可以看出，最重要也还是，找链接，给到[bilibiliDownload(url, options)](#bilibilidownload)

- [`getMultiPageData`](#getmultipagedata)

> 找到匹配项, 然后`string`转成->`json`储存

</details>

---

## 2. bilibiliDownload

<details>

``` go
func bilibiliDownload(url string, options bilibiliOptions) downloader.VideoData {
	var (
		aid, cid, html string
	)
	if options.HTML != "" {
		// reuse html string, but this can't be reused in case of playlist
		html = options.HTML
	} else {
		html = request.Get(url)
	}
	if options.Aid != "" && options.Cid != "" { // 可以看出来 b站 主要靠 这两🆔 视频运输很重要
		aid = options.Aid
		cid = options.Cid
	} else {
		if options.Bangumi {
			cid = utils.MatchOneOf(html, `"cid":(\d+)`)[1]
			aid = utils.MatchOneOf(html, `"aid":(\d+)`)[1]
		} else {
			cid = utils.MatchOneOf(html, `cid=(\d+)`)[1]
			aid = utils.MatchOneOf(url, `av(\d+)`)[1]
		}
	}
    api := genAPI(aid, cid, options.Bangumi) 
    // 配合 在 config.BILIBILI_TOKEN_API : "https://api.bilibili.com/x/player/playurl/token?"
    // 和 两🆔 还有 bilibiliPlayer.min.js 之类 组合出 api链接, 天啊好麻烦阿 😱

	apiData := request.Get(api)
	var dataDict bilibiliData
	json.Unmarshal([]byte(apiData), &dataDict) // 主要是视频链接和质量

	// get the title
	doc := parser.GetDoc(html)
    title := parser.Title(doc)
	if options.Subtitle != "" {
		title = fmt.Sprintf("%s %s", title, options.Subtitle) // 拿到标题
	}

    urls, size := genURL(dataDict.DURL) // 全部url, 总大小 
    // ⚠️bilibili用的播放器是 flv 
    // https://github.com/Bilibili/flv.js

	data := downloader.VideoData{ // 定义好视频数据快
		Site:    "哔哩哔哩 bilibili.com",
		Title:   title,
		URLs:    urls,
		Type:    "video",
		Size:    size,
		Quality: quality[dataDict.Quality],
	}
	data.Download(url) // 通用下载
	return data
}
```

- [genApi](#genapi)

> 从一个有验证的api中, 获取真实的视频链接 验证这一步很重要, 要`黑客`噢😯

- [genURL](#genurl)

> 总结-从genApi-获取的真实的视频块
</details>

---

## genApi

> 从一个有验证的api中, 获取真实的视频链接

<details>

``` js
const (
	// BiliBili blocks keys from time to time.
	// You can extract from the Android client or bilibiliPlayer.min.js
	appKey string = "84956560bc028eb7" 
// 看到没有，密钥阿！
	secKey string = "94aba54af9065f71de72f5508f1cd42e"
)

func genAPI(aid, cid string, bangumi bool) string {
	var (
		baseAPIURL string
		params     string
	)
	utoken := ""
	if config.Cookie != "" {
		utoken = request.Get(fmt.Sprintf(
			"%said=%s&cid=%s", config.BILIBILI_TOKEN_API, aid, cid,
		))
		var t token
		json.Unmarshal([]byte(utoken), &t)
		if t.Code != 0 {
			log.Println(config.Cookie)
			log.Fatal("Cookie error: ", t.Message)
		}
		utoken = t.Data.Token
	}
	if bangumi {
		// The parameters need to be sorted by name
		// qn=0 flag makes the CDN address different every time
		// quality=116(1080P 60) is the highest quality so far
		params = fmt.Sprintf(
			"appkey=%s&cid=%s&module=bangumi&otype=json&qn=116&quality=116&season_type=4&type=&utoken=%s",
			appKey, cid, utoken,
		)
		baseAPIURL = config.BILIBILI_BANGUMI_API
	} else {
		params = fmt.Sprintf(
			"appkey=%s&cid=%s&otype=json&qn=116&quality=116&type=",
			appKey, cid,
		)
		baseAPIURL = config.BILIBILI_API
	}
	// bangumi utoken also need to put in params to sign, but the ordinary video doesn't need
	api := fmt.Sprintf(
		"%s%s&sign=%s", baseAPIURL, params, utils.Md5(params+secKey),
	)
	if !bangumi && utoken != "" {
		api = fmt.Sprintf("%s&utoken=%s", api, utoken)
	}
	return api
}
```
</details>

## genURL

> 总结-从genApi-获取的真实的视频块

<details>

``` go
func genURL(durl []dURLData) ([]downloader.URLData, int64) {
	var (
		urls []downloader.URLData
		size int64
	)
	for _, data := range durl {
		size += data.Size
		url := downloader.URLData{
			URL:  data.URL,
			Size: data.Size,
			Ext:  "flv",
		}
		urls = append(urls, url)
	}
	return urls, size
}
```
</details>

## getMultiPageData

> 找到匹配项, 然后`string`转成->`json`储存

<details>

``` go
func getMultiPageData(html string) (multiPage, error) {
	var data multiPage
	multiPageDataString := utils.MatchOneOf(
		html, `window.__INITIAL_STATE__=(.+?);\(function`, //找到匹配项
	)
	if multiPageDataString == nil {
		return data, errors.New("This page has no playlist")
	}
	json.Unmarshal([]byte(multiPageDataString[1]), &data) // 然后`string`转成->`json`储存
	return data, nil
}
```


</details>