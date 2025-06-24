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

var allowedExtensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

// 获取env(适用于Docker)
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 判断是否为移动设备
func isMobile(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	return strings.Contains(ua, "mobile") ||
		strings.Contains(ua, "android") ||
		strings.Contains(ua, "iphone") ||
		strings.Contains(ua, "ipad") ||
		strings.Contains(ua, "ipod") ||
		strings.Contains(ua, "windows phone")
}

// 读取HTML页面
func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("static/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "html load failed", http.StatusInternalServerError)
		log.Println("Template load failed:", err)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "template execution failed", http.StatusInternalServerError)
		log.Println("Template exec failed:", err)
	}
}

// Ciallo～(∠・ω< )⌒☆
func cialloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ciallo～(∠・ω< )⌒☆")
}

// 主页
func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// 随机选取图片
func randomImageHandler(w http.ResponseWriter, r *http.Request, dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil || len(files) == 0 {
		http.Error(w, "images not found", http.StatusInternalServerError)
		log.Println("No images in", dir, ":", err)
		return
	}

	var imageFiles []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if _, ok := allowedExtensions[ext]; ok {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	if len(imageFiles) == 0 {
		http.Error(w, "no valid images found", http.StatusInternalServerError)
		log.Println("No valid images in", dir)
		return
	}

	file := imageFiles[rand.Intn(len(imageFiles))]
	imagePath := filepath.Join(dir, file)
	contentType := allowedExtensions[strings.ToLower(filepath.Ext(file))]

	// 缓存这种事情补药啊
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, imagePath)
}

// PC端图片
func PCRandomImageHandler(w http.ResponseWriter, r *http.Request) {
	randomImageHandler(w, r, "image/pc")
}

// 移动端图片
func MobileRandomImageHandler(w http.ResponseWriter, r *http.Request) {
	randomImageHandler(w, r, "image/mobile")
}

// 双端自动适配
func AutoRandomImageHandler(w http.ResponseWriter, r *http.Request) {
	ua := r.UserAgent()
	if isMobile(ua) {
		http.Redirect(w, r, "/mobile", http.StatusFound)
	} else {
		http.Redirect(w, r, "/pc", http.StatusFound)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	port := flag.String("port", getEnvString("port", "8000"), "Service Port")
	flag.Parse()

	http.HandleFunc("/ciallo", cialloHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pc", PCRandomImageHandler)
	http.HandleFunc("/mobile", MobileRandomImageHandler)
	http.HandleFunc("/auto", AutoRandomImageHandler)

	log.Printf("Service is on port %s", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Panic(err)
	}
}
