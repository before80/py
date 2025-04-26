package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"py/bs"
	"py/lg"
	"py/pg"
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
	var barMenuInfos []pg.BarMenuInfo
	barMenuInfos, err = pg.GetBarMenus(page, "https://docs.python.org/zh-cn/3.13/index.html")

	var secondMenuInfos []pg.MenuInfo
	for i, barMenuInfo := range barMenuInfos {
		err = pg.InitIndexMdFile(i, barMenuInfo)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}

		secondMenuInfos, err = pg.InsertBarMenuPageData(browserHwnd, barMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}

		if len(secondMenuInfos) <= 0 {
			continue
		}

		for j, menuInfo := range secondMenuInfos {
			err = pg.InitDetailPageMdFile(j, barMenuInfo, menuInfo)
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
				return
			}

			err = pg.InsertDetailPageData(browserHwnd, barMenuInfo, menuInfo, page)
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
	}
	time.Sleep(2000 * time.Second)
}
