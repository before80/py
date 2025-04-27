package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-vgo/robotgo"
	"github.com/tailscale/win"
	"py/bs"
	"py/lg"
	"py/pg"
	"slices"
	"strconv"
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
	var thirdMenuInfos []pg.MenuInfo
	var fourthMenuInfos []pg.MenuInfo
	for i, barMenuInfo := range barMenuInfos {
		if slices.Contains([]string{"tutorial", "whatsnew_3_13"}, barMenuInfo.Filename) {
			continue
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("bar正要处理 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))
		err = pg.InitIndexMdFile(i, barMenuInfo)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("bar初始化完成 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))

		secondMenuInfos, err = pg.InsertBarMenuPageData(browserHwnd, barMenuInfo, page)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
			return
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("bar插入数据完成 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))

		if len(secondMenuInfos) <= 0 {
			continue
		}

		lg.InfoToFileAndStdOut(fmt.Sprintf("second处理二级菜单中 -> file=%s, menu=%s\n", barMenuInfo.Filename, barMenuInfo.MenuName))

		for j, secondMenuInfo := range secondMenuInfos {
			if barMenuInfo.Filename == "library" &&
				slices.Contains([]string{"constants", "allos", "binary", "crypto", "datatypes", "fileformats", "filesys", "functional", "numeric", "persistence", "text", "constants", "exceptions", "functions", "intro"}, secondMenuInfo.Filename) {
				continue
			}

			thirdMenuInfos, err = pg.GetThirdLevelMenu(secondMenuInfo, page)
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
				return
			}
			lg.InfoToFileAndStdOut(fmt.Sprintf("second获取第三级菜单完成 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
				barMenuInfo.Filename, barMenuInfo.MenuName,
				secondMenuInfo.Filename, secondMenuInfo.MenuName))

			// 存在第三级菜单的情况
			if len(thirdMenuInfos) > 0 {
				err = pg.InitSecondIndexMdFile(j, barMenuInfo, secondMenuInfo)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				lg.InfoToFileAndStdOut(fmt.Sprintf("second初始化完成1 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
					barMenuInfo.Filename, barMenuInfo.MenuName,
					secondMenuInfo.Filename, secondMenuInfo.MenuName))

				err = pg.InsertSecondMenuPageData(browserHwnd, barMenuInfo, secondMenuInfo, page)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}

				lg.InfoToFileAndStdOut(fmt.Sprintf("second插入数据完成1 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
					barMenuInfo.Filename, barMenuInfo.MenuName,
					secondMenuInfo.Filename, secondMenuInfo.MenuName))

				for k, thirdMenuInfo := range thirdMenuInfos {
					fourthMenuInfos, err = pg.GetFourthLevelMenu(thirdMenuInfo, page)
					if err != nil {
						lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
						return
					}
					lg.InfoToFileAndStdOut(fmt.Sprintf("third获取第四级菜单完成 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s\n",
						barMenuInfo.Filename, barMenuInfo.MenuName,
						secondMenuInfo.Filename, secondMenuInfo.MenuName,
						thirdMenuInfo.Filename, thirdMenuInfo.MenuName,
					))

					//存在第四级菜单的情况
					if len(fourthMenuInfos) > 0 {
						err = pg.InitThirdIndexMdFile(k, barMenuInfo, secondMenuInfo, thirdMenuInfo)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						lg.InfoToFileAndStdOut(fmt.Sprintf("third初始化完成1 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s\n",
							barMenuInfo.Filename, barMenuInfo.MenuName,
							secondMenuInfo.Filename, secondMenuInfo.MenuName,
							thirdMenuInfo.Filename, thirdMenuInfo.MenuName))

						err = pg.InsertThirdMenuPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, page)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}

						lg.InfoToFileAndStdOut(fmt.Sprintf("third插入数据完成1 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s\n",
							barMenuInfo.Filename, barMenuInfo.MenuName,
							secondMenuInfo.Filename, secondMenuInfo.MenuName,
							thirdMenuInfo.Filename, thirdMenuInfo.MenuName))

						for l, fourthMenuInfo := range fourthMenuInfos {
							err = pg.InitFourthDetailPageMdFile(l, barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo)
							if err != nil {
								lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
								return
							}
							lg.InfoToFileAndStdOut(fmt.Sprintf("fourth初始化完成1 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s fourthfile=%s fourthmenu=%s\n",
								barMenuInfo.Filename, barMenuInfo.MenuName,
								secondMenuInfo.Filename, secondMenuInfo.MenuName,
								thirdMenuInfo.Filename, thirdMenuInfo.MenuName,
								fourthMenuInfo.Filename, fourthMenuInfo.MenuName,
							))

							err = pg.InsertFourthDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, fourthMenuInfo, page)
							if err != nil {
								lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
								return
							}
							lg.InfoToFileAndStdOut(fmt.Sprintf("fourth插入数据完成1 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s fourthfile=%s fourthmenu=%s\n",
								barMenuInfo.Filename, barMenuInfo.MenuName,
								secondMenuInfo.Filename, secondMenuInfo.MenuName,
								thirdMenuInfo.Filename, thirdMenuInfo.MenuName,
								fourthMenuInfo.Filename, fourthMenuInfo.MenuName,
							))
						}
					} else {
						//不存在第四级菜单的情况
						err = pg.InitThirdDetailPageMdFile(k, barMenuInfo, secondMenuInfo, thirdMenuInfo)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						lg.InfoToFileAndStdOut(fmt.Sprintf("third初始化完成2 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s\n",
							barMenuInfo.Filename, barMenuInfo.MenuName,
							secondMenuInfo.Filename, secondMenuInfo.MenuName,
							thirdMenuInfo.Filename, thirdMenuInfo.MenuName))

						err = pg.InsertThirdDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, thirdMenuInfo, page)
						if err != nil {
							lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
							return
						}
						lg.InfoToFileAndStdOut(fmt.Sprintf("third插入数据完成 -> file=%s, menu=%s secondfile=%s secondmenu=%s thirdfile=%s thirdmenu=%s\n",
							barMenuInfo.Filename, barMenuInfo.MenuName,
							secondMenuInfo.Filename, secondMenuInfo.MenuName,
							thirdMenuInfo.Filename, thirdMenuInfo.MenuName))
					}
				}
			} else {
				// 不存在第三级菜单的情况
				err = pg.InitDetailPageMdFile(j, barMenuInfo, secondMenuInfo)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				lg.InfoToFileAndStdOut(fmt.Sprintf("second初始化完成2 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
					barMenuInfo.Filename, barMenuInfo.MenuName,
					secondMenuInfo.Filename, secondMenuInfo.MenuName))

				err = pg.InsertDetailPageData(browserHwnd, barMenuInfo, secondMenuInfo, page)
				if err != nil {
					lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("%v", err), 3)
					return
				}
				lg.InfoToFileAndStdOut(fmt.Sprintf("second插入数据完成2 -> file=%s, menu=%s secondfile=%s secondmenu=%s\n",
					barMenuInfo.Filename, barMenuInfo.MenuName,
					secondMenuInfo.Filename, secondMenuInfo.MenuName))
			}
		}
	}
	lg.InfoToFileAndStdOut("已全部完成")
	_ = browser.Close()
}
