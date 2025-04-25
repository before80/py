package main

import (
	"cppreference/bs"
	"cppreference/exc"
	"cppreference/lg"
	"cppreference/pg"
	_ "embed"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"os"
	"slices"
	"sort"
	"strconv"
	"time"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("%v", err))
		}
	}()

	_ = err
	var browser *rod.Browser
	var page *rod.Page
	var browserHwnd win.HWND

	// 打开浏览器
	browser, err = bs.GetBrowser(strconv.Itoa(0))
	defer browser.MustClose()
	// 创建新页面
	page = browser.MustPage()
	browserHwnd = robotgo.GetHWND()

	hn2HInfo, _ := pg.GetAllHeaderInfo(page)

	// 创建output文件夹
	err = os.MkdirAll("output/std", 0777)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("无法创建%s目录：%v\n", "output/std", err), 3)
		return
	}

	// 给hn排序
	hns := make([]string, 0, len(hn2HInfo))
	for hn := range hn2HInfo {
		hns = append(hns, hn)
	}
	sort.Strings(hns)
	lg.InfoToFileAndStdOut(fmt.Sprintf("hns=%v\n", hns))

	for index, hn := range hns {
		hInfo := hn2HInfo[hn]
		_ = hn
		if !(hn == "stdio") {
			continue
		}

		//if hn != "stdckdint" {
		//	continue
		//}

		err = pg.InitSpecialHeaderMdFile(index, hInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}

		idInfos, err := pg.GetSomeoneHeaderAllIdentifierInfo(hInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}

		//fmt.Printf("idInfos=%v\n", idInfos)
		//time.Sleep(1000 * time.Second)

		if len(idInfos) > 0 {
			for _, idInfo := range idInfos {
				if idInfo.Url != "" && !slices.Contains(exc.ExcludeHeaderIdentifierUrl, idInfo.Url) {
					err = pg.GetIdentifierData(browserHwnd, hn, idInfo, page)
					if err != nil {
						lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v\n", err), 3)
					}
				}
			}
		}
	}
	time.Sleep(2000 * time.Second)
}
