package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// 获取env的内容(为适配Docker而做准备)
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 判断是否为手机端(用于后文是手机端返回手机端图片)
func isMobile(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	return strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone")
}

// 加载html页面
func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("static/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "html load failed", 500)
		return
	}
	t.Execute(w, nil)
}

// 一些页面
func cialloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ciallo～(∠・ω< )⌒☆")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// PC端随机图片
func PCRandomImageHendler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("image/pc")
	if err != nil || len(files) == 0 {
		http.Error(w, "images not found", 500)
		return
	}

	//随机选点图片
	rand.Seed(time.Now().UnixNano())
	file := files[rand.Intn(len(files))].Name()
	imagePath := filepath.Join("image/pc", file)

	switch filepath.Ext(file) {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	default:
		http.Error(w, "image type is not support", 415)
		return
	}

	http.ServeFile(w, r, imagePath)
}

// 移动端随机图片
func MobileRandomImageHendler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("image/mobile")
	if err != nil || len(files) == 0 {
		http.Error(w, "images not found", 500)
		return
	}

	//随机选点图片
	rand.Seed(time.Now().UnixNano())
	file := files[rand.Intn(len(files))].Name()
	imagePath := filepath.Join("image/mobile", file)

	switch filepath.Ext(file) {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	default:
		http.Error(w, "image type is not support", 415)
		return
	}

	http.ServeFile(w, r, imagePath)
}

func AutoRandomImageHendler(w http.ResponseWriter, r *http.Request) {
	ua := r.UserAgent()
	if isMobile(ua) {
		http.Redirect(w, r, "/mobile", 302)
	} else {
		http.Redirect(w, r, "/pc", 302)
	}
}

func main() {
	// 设置启动参数
	var (
		port = flag.String("port", getEnvString("port", "8000"), "Service Port")
	)
	flag.Parse()

	// 定义页面
	http.HandleFunc("/ciallo", cialloHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pc", PCRandomImageHendler)
	http.HandleFunc("/mobile", MobileRandomImageHendler)
	http.HandleFunc("/auto", AutoRandomImageHendler)

	// 启动服务
	log.Printf("Service is on port %s", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Panic(err)
	}
}
