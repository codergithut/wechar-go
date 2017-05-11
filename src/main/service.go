package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"strconv"
	"regexp"
	"time"
	"io"
	"crypto/md5"
	"os"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	//fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		//路径要改成相对路径才行
		t, _ := template.ParseFiles("D:\\mybatis-plugins\\wechar-go\\src\\template\\login.gtpl")
		fmt.Print("haha")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func check(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		//路径要改成相对路径才行
		t, _ := template.ParseFiles("D:\\mybatis-plugins\\wechar-go\\src\\template\\check.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println(checkNull(r,"username"))
		fmt.Println(checkAge(r,"age"))
		fmt.Println(checkHight(r,"height"));
	}
}

func checkXSS (w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		//路径要改成相对路径才行
		t, _ := template.ParseFiles("D:\\mybatis-plugins\\wechar-go\\src\\template\\xss.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
	}
}

func checkRepSubmit(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("D:\\mybatis-plugins\\wechar-go\\src\\template\\rep.gtpl")
		t.Execute(w, token)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			fmt.Println(token + "第一次")
		} else {
			//不存在token报错
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
	}

}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("D:\\mybatis-plugins\\wechar-go\\src\\template\\upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		fmt.Println(handler.Filename)
		f, err := os.OpenFile("D:\\abc\\"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	http.HandleFunc("/login", login)         //设置访问的路由
	http.HandleFunc("/check", check)
	http.HandleFunc("/xss", checkXSS)
	http.HandleFunc("/rep", checkRepSubmit)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func checkNull(r *http.Request, checkParam string)(bool) {
	if len(r.Form[checkParam][0])==0{
		//为空的处理
		return false
	}
	return true;
}

func checkAge(r *http.Request, checkParam string) bool {
	getint,err:=strconv.Atoi(r.Form.Get(checkParam))
	if err!=nil{
		return false
	}
	//接下来就可以判断这个数字的大小范围了
	if getint >100 {
		return false
	}
	return true;
}

func checkHight(r *http.Request, checkParam string) bool {
	if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get(checkParam)); !m {
		return false
	}
	return true;
}

func checkZh(r *http.Request, checkParam string) bool {
	if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
		return false
	}
	return true;
}

func checkEn(r *http.Request, checkParam string) bool {
	if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("engname")); !m {
		return false
	}
	return true;
}

func checkEmail(r *http.Request, checkParam string) bool {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, r.Form.Get("email")); !m {
		return false;
	}
	return true;
}

func checkPhone(r *http.Request, checkParam string) bool {
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
		return false
	}
	return true;
}

func checkSelectSingle(r *http.Request, checkParam string, slice []string) bool {
	//slice:=[]string{"apple","pear","banane"}
	//
	//v := r.Form.Get("fruit")
	//for _, item range slice {
	 //   if item == v {
		//return true
	 //   }
        //}
	//
        //return false
	return false;
}

func timeFormate(r *http.Request, checkParam string) {
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("Go launched at %s\n", t.Local())
}

func checkPersonId(r *http.Request, checkParam string) bool {
	//验证15位身份证，15位的是全部数字
	if m, _ := regexp.MatchString(`^(\d{15})$`, r.Form.Get("usercard")); !m {
		return false
	} else {
		return true;
	}

	//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
	if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, r.Form.Get("usercard")); !m {
		return false
	} else {
		return true;
	}

}
