package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
)

func JSPConnection(w fyne.Window, url, password string) {
	reqBody := fmt.Sprintf("%s=%s", password, "echo success")
	fmt.Println(reqBody)
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

	if strings.TrimSpace(string(body)) == "success" {
		selectedType = "JSP"
		addConnection(url, password, selectedType)
		showJSPCommandPage(w, url, password)
	} else {
		dialog.ShowInformation("Connection Failed", fmt.Sprintf("Unexpected response: %s", string(body)), w)
	}
}

func showJSPCommandPage(w fyne.Window, url, password string) {
	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("Enter command...")

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord
	output.SetText("")

	scrollContainer := container.NewScroll(output)
	scrollContainer.SetMinSize(fyne.NewSize(600, 400))

	osLabel := widget.NewLabel("Detecting OS...")
	var detectedOS string
	go func() {
		detectedOS, _ = detectOs2(url, password)
		osLabel.SetText(fmt.Sprintf("OS: %s", detectedOS))
	}()

	executeCommand := func(url, password, command string) {
		wrappedCommand := fmt.Sprintf("%s", command)
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

	content := container.NewVBox(
		commandEntry,
		executeButton,
		scrollContainer,
		osLabel,
		backButton,
	)

	w.SetContent(content)
}

func detectOs2(url, password string) (string, error) {

	unameCommand := fmt.Sprintf("%s=uname", password)
	verCommand := fmt.Sprintf("%s=ver", password)

	fmt.Printf("Sending uname request: %s\n", unameCommand) // 打印请求体
	unameURL := fmt.Sprintf("%s?%s", url, unameCommand)

	resp, err := http.Get(unameURL)
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
	verURL := fmt.Sprintf("%s?%s", url, verCommand)
	resp, err = http.Get(verURL)
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
