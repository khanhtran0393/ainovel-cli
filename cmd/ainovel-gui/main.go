package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
	"github.com/voocel/ainovel-cli/internal/entry/webgui"
	"github.com/voocel/ainovel-cli/internal/rules"
)

func init() {
	// Pure-go .env loader
	data, err := os.ReadFile(".env")
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
				if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
					val = val[1 : len(val)-1]
				} else if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
					val = val[1 : len(val)-1]
				}
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func main() {
	rules.EnsureHomeRulesDir()

	customPort := ""
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--port" && i+1 < len(os.Args) {
			customPort = os.Args[i+1]
		}
	}

	// Tải cấu hình, KHÔNG THOÁT NGAY NẾU LỖI ĐỂ TRÁNH APP IM LÌM
	cfg, err := bootstrap.LoadConfig("")
	if err != nil {
		os.WriteFile("ainovel-gui-crash.log", []byte(fmt.Sprintf("Cảnh báo tải config: %v\nTiếp tục với cấu hình mặc định để duy trì Web GUI.", err)), 0644)
		cfg = bootstrap.Config{Style: "default"}
	}

	bundle := assets.Load(cfg.Style)

	// Port Auto-Hopping: Nếu không chỉ định --port, tự động quét tìm cổng rảnh từ 8080 đến 8099
	port := customPort
	if port == "" {
		for p := 8080; p <= 8099; p++ {
			pStr := fmt.Sprintf("%d", p)
			ln, err := net.Listen("tcp", ":"+pStr)
			if err == nil {
				ln.Close()
				port = pStr
				break
			}
		}
		if port == "" {
			port = "8080" // Fallback
		}
	}

	// Tự động mở trình duyệt sau 500ms
	go func() {
		time.Sleep(500 * time.Millisecond)
		url := "http://localhost:" + port
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		}
		if err != nil {
			os.WriteFile("ainovel-gui-crash.log", []byte(fmt.Sprintf("Không thể tự động mở trình duyệt: %v\nVui lòng mở thủ công tại %s", err, url)), 0644)
		}
	}()

	if err := webgui.Run(cfg, bundle, "", port); err != nil {
		os.WriteFile("ainovel-gui-crash.log", []byte(fmt.Sprintf("Lỗi khởi động webgui: %v", err)), 0644)
		os.Exit(1)
	}
}
