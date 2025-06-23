package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// getEnvString 获取环境变量值，如果不存在则返回默认值
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// renderTemplate 渲染HTML模板
func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("static/" + tmpl + ".html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Template load failed", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
	}
}

// getSupportedImageType 根据文件扩展名返回MIME类型
func getSupportedImageType(filename string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg", true
	case ".png":
		return "image/png", true
	case ".webp":
		return "image/webp", true
	case ".gif":
		return "image/gif", true
	default:
		return "", false
	}
}

// serveRandomImage 提供指定目录的随机图片
func serveRandomImage(w http.ResponseWriter, r *http.Request, imageDir string) {
	// 检查目录是否存在
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		log.Printf("Image directory does not exist: %s", imageDir)
		http.Error(w, "Image directory not found", http.StatusNotFound)
		return
	}

	// 读取目录
	files, err := os.ReadDir(imageDir)
	if err != nil {
		log.Printf("Error reading directory %s: %v", imageDir, err)
		http.Error(w, "Failed to read image directory", http.StatusInternalServerError)
		return
	}

	// 过滤出支持的图片文件
	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() {
			if _, supported := getSupportedImageType(file.Name()); supported {
				imageFiles = append(imageFiles, file.Name())
			}
		}
	}

	if len(imageFiles) == 0 {
		log.Printf("No supported images found in directory: %s", imageDir)
		http.Error(w, "No images found", http.StatusNotFound)
		return
	}

	// 随机选择图片
	rand.Seed(time.Now().UnixNano())
	selectedFile := imageFiles[rand.Intn(len(imageFiles))]
	imagePath := filepath.Join(imageDir, selectedFile)

	// 安全检查：确保路径在预期目录内
	absImagePath, err := filepath.Abs(imagePath)
	if err != nil {
		log.Printf("Error getting absolute path: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	absImageDir, err := filepath.Abs(imageDir)
	if err != nil {
		log.Printf("Error getting absolute directory path: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(absImagePath, absImageDir) {
		log.Printf("Security violation: path traversal attempt")
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}

	// 设置正确的Content-Type
	if contentType, supported := getSupportedImageType(selectedFile); supported {
		w.Header().Set("Content-Type", contentType)
		// 添加缓存控制头
		w.Header().Set("Cache-Control", "public, max-age=3600")
	} else {
		http.Error(w, "Unsupported image type", http.StatusUnsupportedMediaType)
		return
	}

	// 提供文件
	http.ServeFile(w, r, imagePath)
}

// detectUserAgent 检测用户代理以判断设备类型
func detectUserAgent(userAgent string) string {
	userAgent = strings.ToLower(userAgent)
	mobileKeywords := []string{"mobile", "android", "iphone", "ipad", "ipod", "blackberry", "windows phone"}
	
	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			return "mobile"
		}
	}
	return "pc"
}

// 页面处理函数
func cialloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ciallo～(∠・ω< )⌒☆")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// autoHandler 自动检测设备类型并返回对应图片
func autoHandler(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("User-Agent")
	deviceType := detectUserAgent(userAgent)
	
	var imageDir string
	if deviceType == "mobile" {
		imageDir = "image/mobile"
	} else {
		imageDir = "image/pc"
	}
	
	serveRandomImage(w, r, imageDir)
}

// pcRandomImageHandler PC端随机图片
func pcRandomImageHandler(w http.ResponseWriter, r *http.Request) {
	serveRandomImage(w, r, "image/pc")
}

// mobileRandomImageHandler 移动端随机图片
func mobileRandomImageHandler(w http.ResponseWriter, r *http.Request) {
	serveRandomImage(w, r, "image/mobile")
}

func main() {
	// 设置启动参数
	var (
		port = flag.String("port", getEnvString("port", "8000"), "Service Port")
	)
	flag.Parse()

	// 创建图片目录（如果不存在）
	os.MkdirAll("image/pc", 0755)
	os.MkdirAll("image/mobile", 0755)

	// 定义路由
	http.HandleFunc("/ciallo", cialloHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auto", autoHandler)
	http.HandleFunc("/pc", pcRandomImageHandler)
	http.HandleFunc("/mobile", mobileRandomImageHandler)

	// 启动服务
	log.Printf("Service is starting on port %s", *port)
	log.Printf("Available endpoints:")
	log.Printf("  - GET /        : Main page")
	log.Printf("  - GET /auto    : Auto-detect device and serve image")
	log.Printf("  - GET /pc      : Random PC image")
	log.Printf("  - GET /mobile  : Random mobile image")
	log.Printf("  - GET /ciallo  : Easter egg")
	
	server := &http.Server{
		Addr:         ":" + *port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server listening on :%s", *port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
