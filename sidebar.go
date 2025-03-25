package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
)

var checkboxes []*widget.Check

func showSidebar(w fyne.Window) *fyne.Container {

	list := widget.NewList(
		func() int {
			return len(connections)
		},
		func() fyne.CanvasObject {

			hbox := container.NewHBox(
				widget.NewLabel("Template"),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {}),
				widget.NewButtonWithIcon("", theme.FolderIcon(), func() {}), // 新增按钮在第三个位置
				widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {}),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}),
			)
			return hbox
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			conn := connections[i]

			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s - %s", conn.Type, conn.URL))

			o.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				switch conn.Type {
				case "JSP":
					JSPConnection(w, conn.URL, conn.Password)
				case "ASP":
					ASPConnection(w, conn.URL, conn.Password)
				case "PHP":
					PHPConnection(w, conn.URL, conn.Password)
				case "ASPX":
					ASPXConnection(w, conn.URL, conn.Password)
				case "PHP+":
					PHPEConnection(w, conn.URL, conn.Password)
				case "JSP+":
					JSPEConnection(w, conn.URL, conn.Password)
				case "ASPX+":
					ASPXEConnection(w, conn.URL, conn.Password)
				}
			}

			o.(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				//for _, conn := range connections {
				//	fileManager(conn.URL, conn.Password)
				//}
				remoteURL := fmt.Sprintf("%s?%s=", conn.URL, conn.Password)
				fmt.Println(remoteURL)
				connType = conn.Type
				Osinfo := conn.OS
				fmt.Printf("1:%s\n", Osinfo)
				if connType == "PHP" || connType == "JSP" {
					if Osinfo == "Linux" {
						RunFileManager(w, remoteURL, conn.Type)
					} else {
						RunFileManager2(w, remoteURL, conn.Type)
					}
				}
				if connType == "ASP" || connType == "ASPX" {
					RunFileManager2(w, remoteURL, conn.Type)
				}
				if connType == "JSP+" {
					RunFileManager3(w, conn.URL, conn.Type, conn.Password)
				}
				if connType == "ASPX+" {
					RunFileManager4(w, conn.URL, conn.Type, conn.Password)
				}
				if connType == "PHP+" {
					Osinfo = detecteOSX(conn.URL, conn.Password)
					fmt.Printf("osinfo:%s", Osinfo)
					if Osinfo == "Windows" {
						fmt.Println("\nrun 4")
						RunFileManager4(w, conn.URL, conn.Type, conn.Password)
					} else {
						fmt.Println("\nrun 3")
						RunFileManager3(w, conn.URL, conn.Type, conn.Password)

					}

				}
			}

			o.(*fyne.Container).Objects[4].(*widget.Button).OnTapped = func() {
				showEditDialog(w, i)
			}

			o.(*fyne.Container).Objects[5].(*widget.Button).OnTapped = func() {
				connections = append(connections[:i], connections[i+1:]...)
				saveConnections()

				w.SetContent(firstPage())
			}
		},
	)

	scrollContainer := container.NewScroll(list)
	scrollContainer.SetMinSize(fyne.NewSize(250, 600))

	addButton := widget.NewButtonWithIcon("连接", theme.ContentAddIcon(), func() {
		w.SetContent(mainPage())
	})

	batchExecuteButton := widget.NewButtonWithIcon("批量执行", theme.ContentAddIcon(), func() {
		// 切换到批量执行模式
		showBatchExecuteMode(w)
	})

	shellExecuteButton := widget.NewButtonWithIcon("反弹shell", theme.ContentAddIcon(), func() {
		// 切换到批量执行模式
		showShellExecuteMode(w)
	})

	sidebar = container.NewVBox(
		addButton,
		batchExecuteButton,
		shellExecuteButton,
		scrollContainer,
		layout.NewSpacer(),
	)

	sidebar.Resize(fyne.NewSize(250, 1000))

	return sidebar
}

func showEditDialog(w fyne.Window, index int) {
	conn := connections[index]

	urlEntry := widget.NewEntry()
	urlEntry.SetText(conn.URL)

	passwordEntry := widget.NewEntry()
	passwordEntry.SetText(conn.Password)

	saveButton := widget.NewButton("保存", func() {
		conn.URL = urlEntry.Text
		conn.Password = passwordEntry.Text
		connections[index] = conn
		saveConnections()
		w.SetContent(firstPage())
		dialog.ShowInformation("Success", "修改已保存", w)
	})

	content := container.NewVBox(
		widget.NewLabel("编辑连接信息"),
		urlEntry,
		passwordEntry,
		container.NewHBox(saveButton),
	)

	dialogContent := container.NewPadded(content)
	dialogContent.Resize(fyne.NewSize(400, 300))

	editDialog := dialog.NewCustom("编辑连接信息", "关闭", dialogContent, w)
	editDialog.Resize(fyne.NewSize(400, 300))
	editDialog.Show()
}

func showBatchExecuteMode(w fyne.Window) {
	checkboxes = []*widget.Check{}

	list := widget.NewList(
		func() int {
			return len(connections)
		},
		func() fyne.CanvasObject {
			hbox := container.NewHBox(
				widget.NewCheck("", func(bool) {}),
				widget.NewLabel("Template"),
			)
			return hbox
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			conn := connections[i]

			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%s - %s", conn.Type, conn.URL))

			checkboxes = append(checkboxes, o.(*fyne.Container).Objects[0].(*widget.Check))
		},
	)

	scrollContainer := container.NewScroll(list)
	scrollContainer.SetMinSize(fyne.NewSize(250, 600))

	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("Enter you command")

	executeButton := widget.NewButtonWithIcon("确定", theme.ConfirmIcon(), func() {
		command := commandEntry.Text
		err := os.Truncate("output.txt", 0)
		if err != nil {
			fmt.Println("Error truncating file:", err)
		}
		executeSomeCommands(command)
	})

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(firstPage())
	})

	sidebar = container.NewVBox(
		commandEntry,
		executeButton,
		scrollContainer,
		backButton,
		layout.NewSpacer(), // 添加一些间距
	)

	sidebar.Resize(fyne.NewSize(250, 1000))
	w.SetContent(sidebar)
}
