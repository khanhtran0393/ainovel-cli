package webgui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/voocel/ainovel-cli/assets"
	"github.com/voocel/ainovel-cli/internal/bootstrap"
	"github.com/voocel/ainovel-cli/internal/diag"
	"github.com/voocel/ainovel-cli/internal/entry/startup"
	"github.com/voocel/ainovel-cli/internal/host"
	"github.com/voocel/ainovel-cli/internal/logger"
	"github.com/voocel/ainovel-cli/internal/rules"
	"github.com/voocel/ainovel-cli/internal/store"
)

type WebServer struct {
	cfg        bootstrap.Config
	bundle     assets.Bundle
	configPath string

	eng        *host.Host
	isRunning  bool
	lastAction string
	clients    map[chan string]bool
	mu         sync.Mutex
}

// panicRecoveryMiddleware bọc bảo vệ toàn bộ máy chủ, tránh sập tiến trình ainovel-gui.exe khi xảy ra sự cố ngầm
func panicRecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("HTTP Panic Recovered", "error", err, "path", r.URL.Path)
				http.Error(w, fmt.Sprintf("Máy chủ Web GUI đã tự động chống chọi lỗi ngầm (Panic Shield): %v", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func Run(cfg bootstrap.Config, bundle assets.Bundle, configPath, port string) error {
	if port == "" {
		port = "8080"
	}
	server := &WebServer{
		cfg:        cfg,
		bundle:     bundle,
		configPath: configPath,
		clients:    make(map[chan string]bool),
		lastAction: "Ready / Idle",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", panicRecoveryMiddleware(server.handleIndex))
	mux.HandleFunc("/api/status", panicRecoveryMiddleware(server.handleStatus))
	mux.HandleFunc("/api/start", panicRecoveryMiddleware(server.handleStart))
	mux.HandleFunc("/api/resume", panicRecoveryMiddleware(server.handleResume))
	mux.HandleFunc("/api/stream", panicRecoveryMiddleware(server.handleStream))
	mux.HandleFunc("/api/chapters", panicRecoveryMiddleware(server.handleChaptersList))
	mux.HandleFunc("/api/chapters/", panicRecoveryMiddleware(server.handleChapterDetail))
	mux.HandleFunc("/api/config", panicRecoveryMiddleware(server.handleConfig))
	mux.HandleFunc("/api/shutdown", panicRecoveryMiddleware(server.handleShutdown))
	mux.HandleFunc("/api/enhance-prompt", panicRecoveryMiddleware(server.handleEnhancePrompt))
	mux.HandleFunc("/api/download-all", panicRecoveryMiddleware(server.handleDownloadAll))
	mux.HandleFunc("/api/rules", panicRecoveryMiddleware(server.handleRules))
	mux.HandleFunc("/api/diag", panicRecoveryMiddleware(server.handleDiag))
	mux.HandleFunc("/api/ping-keys", panicRecoveryMiddleware(server.handlePingKeys))

	addr := ":" + port
	fmt.Fprintf(os.Stderr, "\n[GUI] Khởi động Web GUI Server thành công tại http://localhost:%s\n", port)
	return http.ListenAndServe(addr, mux)
}

func (s *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(WebUIHTML))
}

func (s *WebServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	isRunning := s.isRunning
	lastAction := s.lastAction
	s.mu.Unlock()

	// Thu thập thống kê bộ nhớ & số lượng Goroutines để đo đạc độ ổn định
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	data := map[string]any{
		"is_running":  isRunning,
		"last_action": lastAction,
		"goroutines":  runtime.NumGoroutine(),
		"sys_mb":      m.Sys / 1024 / 1024,
		"alloc_mb":    m.Alloc / 1024 / 1024,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *WebServer) broadcast(msg map[string]any) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	str := string(data)

	// GIẢI PHÁP TỐI ƯU CẠNH TRANH KHÓA (Lock Contention Relief)
	// Sao chép nhanh danh sách client channel ra mảng tĩnh, giải phóng mutex ngay lập tức
	s.mu.Lock()
	activeClients := make([]chan string, 0, len(s.clients))
	for client := range s.clients {
		activeClients = append(activeClients, client)
	}
	s.mu.Unlock()

	for _, client := range activeClients {
		select {
		case client <- str:
		default:
			// Bỏ qua client bị nghẽn để không làm kẹt luồng stream AI
		}
	}
}

func (s *WebServer) handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Engine đang hoạt động"})
		return
	}
	s.isRunning = true
	s.lastAction = "Cleaning old session & Starting new project"
	s.mu.Unlock()

	var body struct {
		Prompt       string `json:"prompt"`
		ChapterCount int    `json:"chapter_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Prompt) == "" {
		s.mu.Lock()
		s.isRunning = false
		s.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Thiếu prompt khởi tạo"})
		return
	}

	oldDir := "output/novel/chapters"
	if _, err := os.Stat(oldDir); err == nil {
		backupDir := fmt.Sprintf("output/novel/chapters_backup_%s", time.Now().Format("20060102_150405"))
		_ = os.Rename(oldDir, backupDir)
		_ = os.MkdirAll(oldDir, 0755)
	}

	finalPrompt := strings.TrimSpace(body.Prompt)
	if body.ChapterCount > 0 {
		finalPrompt = fmt.Sprintf("Quy mô tác phẩm: Dự kiến %d chương. %s", body.ChapterCount, finalPrompt)
	}

	go s.runEngine(finalPrompt, false)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (s *WebServer) handleResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Engine đang hoạt động"})
		return
	}
	s.isRunning = true
	s.lastAction = "Resuming session"
	s.mu.Unlock()

	go s.runEngine("", true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

func (s *WebServer) runEngine(prompt string, isResume bool) {
	eng, err := host.New(s.cfg, s.bundle)
	if err != nil {
		s.mu.Lock()
		s.isRunning = false
		s.lastAction = "Error: " + err.Error()
		s.mu.Unlock()
		s.broadcast(map[string]any{"text": "[ERROR] Khởi tạo host thất bại: " + err.Error(), "type": "error"})
		return
	}

	s.mu.Lock()
	s.eng = eng
	s.mu.Unlock()

	_ = os.MkdirAll(eng.Dir(), 0755)

	cleanup := logger.SetupFile(eng.Dir(), "webgui.log", false)
	defer cleanup()
	defer eng.Close()
	defer func() {
		s.mu.Lock()
		s.isRunning = false
		s.eng = nil
		s.mu.Unlock()
		_, _ = diag.Export(store.NewStore(eng.Dir()))
	}()

	if !isResume {
		plan, err := startup.PrepareQuick(startup.Request{
			Mode:        startup.ModeQuick,
			UserPrompt:  prompt,
			OutputDir:   eng.Dir(),
			Interactive: true,
		})
		if err != nil {
			s.broadcast(map[string]any{"text": "[ERROR] Chuẩn bị thất bại: " + err.Error(), "type": "error"})
			return
		}
		s.broadcast(map[string]any{"text": fmt.Sprintf("[SYSTEM] Khởi động phiên sáng tác mới tại: %s", eng.Dir()), "type": "system"})
		if err := eng.StartPrepared(plan.StartPrompt); err != nil {
			s.broadcast(map[string]any{"text": "[ERROR] Khởi động động cơ thất bại: " + err.Error(), "type": "error"})
			return
		}
	} else {
		_, err := eng.ReplayQueue(0)
		if err != nil {
			s.broadcast(map[string]any{"text": "[ERROR] Replay queue thất bại: " + err.Error(), "type": "error"})
			return
		}
		label, err := eng.Resume()
		if err != nil {
			s.broadcast(map[string]any{"text": "[ERROR] Khôi phục thất bại: " + err.Error(), "type": "error"})
			return
		}
		if label == "" {
			s.broadcast(map[string]any{"text": "[ERROR] Không tìm thấy phiên làm việc có thể khôi phục trong thư mục đầu ra", "type": "error"})
			return
		}
		s.broadcast(map[string]any{"text": fmt.Sprintf("[SYSTEM] Khôi phục phiên làm việc: %s (%s)", eng.Dir(), label), "type": "system"})
	}

	s.consumeEngine(eng)
}

func (s *WebServer) consumeEngine(eng *host.Host) {
	for {
		select {
		case ev, ok := <-eng.Events():
			if !ok {
				s.broadcast(map[string]any{"text": "[SYSTEM] Hoàn tất luồng xử lý sự kiện.", "type": "system"})
				return
			}
			if strings.TrimSpace(ev.Summary) != "" {
				s.mu.Lock()
				s.lastAction = fmt.Sprintf("[%s] %s", ev.Category, ev.Summary)
				s.mu.Unlock()
				ts := ev.Time.Format("15:04:05")
				msg := fmt.Sprintf("[%s] [%s] %s", ts, ev.Category, ev.Summary)
				slog.Info("UI Event", "category", ev.Category, "summary", ev.Summary)
				s.broadcast(map[string]any{"text": msg, "type": "event"})
			}
		case delta, ok := <-eng.Stream():
			if !ok {
				continue
			}
			if delta == host.StreamClearSentinel {
				s.broadcast(map[string]any{"text": "\n\n", "type": "text"})
				continue
			}
			if delta == "" {
				continue
			}
			s.broadcast(map[string]any{"text": delta, "type": "text"})
		case _, ok := <-eng.Done():
			if !ok {
				s.broadcast(map[string]any{"text": "[SYSTEM] Phiên làm việc đã kết thúc.", "type": "system"})
				return
			}
			s.broadcast(map[string]any{"text": "[SYSTEM] Động cơ báo hiệu hoàn thành tiến trình.", "type": "system"})
			return
		}
	}
}

func (s *WebServer) handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Tăng buffer lên 300 để đảm bảo không bị mất log khi AI stream ở tốc độ cực cao
	ch := make(chan string, 300)
	s.mu.Lock()
	s.clients[ch] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, ch)
		s.mu.Unlock()
		close(ch)
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *WebServer) handleChaptersList(w http.ResponseWriter, r *http.Request) {
	chapters := make([]map[string]any, 0)
	outputDir := "output/novel/chapters"
	files, err := os.ReadDir(outputDir)
	if err == nil {
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			info, err := f.Info()
			if err != nil {
				continue
			}
			path := filepath.Join(outputDir, f.Name())
			content, _ := os.ReadFile(path)
			contentStr := string(content)
			lines := strings.Split(contentStr, "\n")
			title := f.Name()
			excerpt := ""
			for _, line := range lines {
				if strings.HasPrefix(line, "# ") {
					title = strings.TrimPrefix(line, "# ")
				} else if strings.TrimSpace(line) != "" && excerpt == "" && !strings.HasPrefix(line, "#") {
					excerpt = line
					if len(excerpt) > 120 {
						excerpt = excerpt[:120] + "..."
					}
				}
			}
			chapters = append(chapters, map[string]any{
				"name":    f.Name(),
				"title":   title,
				"path":    path,
				"size":    info.Size() / 1024,
				"excerpt": excerpt,
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"chapters": chapters})
}

func (s *WebServer) handleChapterDetail(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/chapters/")
	if name == "" || strings.Contains(name, "..") {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join("output/novel/chapters", name)

	if r.Method == http.MethodGet {
		content, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"content": string(content)})
		return
	}

	if r.Method == http.MethodPost {
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
			return
		}
		_ = os.MkdirAll("output/novel/chapters", 0755)
		if err := os.WriteFile(path, []byte(body.Content), 0644); err != nil {
			json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"success": true})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *WebServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		envContent, _ := os.ReadFile(".env")
		cfgPath := s.configPath
		if cfgPath == "" {
			cfgPath = "config.json"
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				cfgPath = "config/config.json"
			}
		}
		cfgContent, _ := os.ReadFile(cfgPath)

		json.NewEncoder(w).Encode(map[string]any{
			"env":    string(envContent),
			"config": string(cfgContent),
		})
		return
	}

	if r.Method == http.MethodPost {
		var body struct {
			Env    *string `json:"env"`
			Config *string `json:"config"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
			return
		}

		if body.Env != nil {
			os.WriteFile(".env", []byte(*body.Env), 0644)
			for _, line := range strings.Split(*body.Env, "\n") {
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
					os.Setenv(key, val)
				}
			}
		}

		if body.Config != nil {
			cfgPath := s.configPath
			if cfgPath == "" {
				cfgPath = "config.json"
				if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
					cfgPath = "config/config.json"
				}
			}
			if dir := filepath.Dir(cfgPath); dir != "." {
				_ = os.MkdirAll(dir, 0755)
			}
			os.WriteFile(cfgPath, []byte(*body.Config), 0644)
			newCfg, err := bootstrap.LoadConfig(cfgPath)
			if err == nil {
				s.mu.Lock()
				s.cfg = newCfg
				s.mu.Unlock()
			}
		}

		json.NewEncoder(w).Encode(map[string]any{"success": true})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *WebServer) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "message": "Đang tắt máy chủ..."})
	go func() {
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}

func (s *WebServer) handleDownloadAll(w http.ResponseWriter, r *http.Request) {
	outputDir := "output/novel/chapters"
	files, err := os.ReadDir(outputDir)
	if err != nil || len(files) == 0 {
		http.Error(w, "Không tìm thấy chương truyện nào trong hệ thống", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	buf.WriteString("# TRỌN BỘ TIỂU THUYẾT\n\nĐược sáng tác bởi AI Novel Studio.\n\n---\n\n")

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		path := filepath.Join(outputDir, f.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		buf.WriteString(string(content))
		buf.WriteString("\n\n---\n\n")
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"Full_Novel_AI_Studio.md\"")
	w.Write(buf.Bytes())
}

func (s *WebServer) handlePingKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var allKeys []string
	seenKeys := make(map[string]bool)
	addKey := func(k string) {
		k = strings.TrimSpace(k)
		if k != "" && !seenKeys[k] {
			allKeys = append(allKeys, k)
			seenKeys[k] = true
		}
	}

	envData, err := os.ReadFile(".env")
	if err == nil {
		for _, line := range strings.Split(string(envData), "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "GEMINI") && strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				val := strings.TrimSpace(parts[1])
				if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
					val = val[1 : len(val)-1]
				} else if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
					val = val[1 : len(val)-1]
				}
				addKey(val)
			}
		}
	}
	addKey(os.Getenv("GEMINI_API_KEY"))
	for i := 1; i <= 20; i++ {
		addKey(os.Getenv(fmt.Sprintf("GEMINI_API_KEY_%d", i)))
	}

	if len(allKeys) == 0 {
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Không tìm thấy API Key nào trong tệp .env hoặc biến môi trường hệ thống. Vui lòng dán Key vào ô bên dưới và nhấn 'Lưu .env' trước!",
		})
		return
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("=== BÁO CÁO KIỂM TOÁN QUOTA %d API KEYS ===\nThời gian: %s\n\n", len(allKeys), time.Now().Format("15:04:05 02/01/2006")))

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": "ping"},
				},
			},
		},
	}
	reqBytes, _ := json.Marshal(reqBody)

	successCount := 0
	for idx, apiKey := range allKeys {
		masked := apiKey
		if len(apiKey) > 12 {
			masked = apiKey[:6] + "..." + apiKey[len(apiKey)-4:]
		}
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
		
		// BẢO VỆ CHỐNG TREO BĂNG KHI PING (TIMEOUT HARDENING)
		client := &http.Client{Timeout: 6 * time.Second}
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			buf.WriteString(fmt.Sprintf("[Key #%d] %s ➔ ❌ Lỗi kết nối (Timeout/DNS): %v\n", idx+1, masked, err))
			continue
		}
		if resp.StatusCode == http.StatusOK {
			buf.WriteString(fmt.Sprintf("[Key #%d] %s ➔ ✔️ HOẠT ĐỘNG TỐT (HTTP 200 OK - Còn Quota)\n", idx+1, masked))
			successCount++
		} else {
			buf.WriteString(fmt.Sprintf("[Key #%d] %s ➔ ⚠️ HẾT QUOTA / KHÔNG HỢP LỆ (HTTP %d)\n", idx+1, masked, resp.StatusCode))
		}
		resp.Body.Close()
	}

	buf.WriteString(fmt.Sprintf("\nTổng kết: %d/%d Key hoạt động bình thường. Hệ thống sẵn sàng xoay vòng tự động!", successCount, len(allKeys)))

	json.NewEncoder(w).Encode(map[string]any{
		"success": successCount > 0,
		"message": buf.String(),
	})
}

func (s *WebServer) handleEnhancePrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Keywords string `json:"keywords"`
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Keywords) == "" {
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Thiếu từ khóa"})
		return
	}

	var allKeys []string
	seenKeys := make(map[string]bool)
	addKey := func(k string) {
		k = strings.TrimSpace(k)
		if k != "" && !seenKeys[k] {
			allKeys = append(allKeys, k)
			seenKeys[k] = true
		}
	}

	envData, err := os.ReadFile(".env")
	if err == nil {
		for _, line := range strings.Split(string(envData), "\n") {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "GEMINI") && strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				val := strings.TrimSpace(parts[1])
				if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
					val = val[1 : len(val)-1]
				} else if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
					val = val[1 : len(val)-1]
				}
				addKey(val)
			}
		}
	}
	addKey(os.Getenv("GEMINI_API_KEY"))
	for i := 1; i <= 20; i++ {
		addKey(os.Getenv(fmt.Sprintf("GEMINI_API_KEY_%d", i)))
	}

	if len(allKeys) == 0 {
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Không tìm thấy API Key nào trong .env hoặc hệ thống để gọi AI trợ giúp"})
		return
	}

	models := []string{"gemini-1.5-pro", "gemini-1.5-flash", "gemini-2.5-flash", "gemini-2.0-flash-exp"}

	promptPayload := fmt.Sprintf(`Bạn là một trợ lý sáng tác tiểu thuyết hàng đầu. Hãy dựa vào các từ khóa sau của người dùng để mở rộng thành một Gợi ý sáng tác (Prompt) chi tiết, đầy đủ bối cảnh, nhân vật, mâu thuẫn chính và không khí truyện.
Từ khóa của người dùng: %s

Hãy trả về duy nhất nội dung Gợi ý sáng tác hoàn chỉnh (không giải thích dòng ngoài).`, body.Keywords)

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{"text": promptPayload},
				},
			},
		},
	}
	reqBytes, _ := json.Marshal(reqBody)

	// BẢO VỆ CHỐNG TREO BĂNG (TIMEOUT HARDENING) CHO KẾT NỐI GEMINI API
	client := &http.Client{Timeout: 10 * time.Second}
	var lastErr string
	for _, modelName := range models {
		for keyIdx, apiKey := range allKeys {
			url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", modelName, apiKey)
			resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBytes))
			if err != nil {
				lastErr = fmt.Sprintf("Lỗi kết nối (%s/Key_%d): %v", modelName, keyIdx+1, err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				lastErr = fmt.Sprintf("Hết Quota hoặc bị chặn (%s/Key_%d): HTTP %d", modelName, keyIdx+1, resp.StatusCode)
				continue
			}

			var geminiResp struct {
				Candidates []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				} `json:"candidates"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
				resp.Body.Close()
				lastErr = fmt.Sprintf("Lỗi giải mã phản hồi (%s/Key_%d)", modelName, keyIdx+1)
				continue
			}
			resp.Body.Close()

			if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
				enhanced := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
				fmt.Fprintf(os.Stderr, "\n[AI ASSIST] Phát triển ý tưởng thành công với %s (Key #%d)\n", modelName, keyIdx+1)
				json.NewEncoder(w).Encode(map[string]any{"success": true, "enhanced_prompt": enhanced})
				return
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   fmt.Sprintf("Toàn bộ %d API Keys đã hết Quota trên cả 4 mốc mô hình (Pro ➔ Flash). Lỗi cuối: %s", len(allKeys), lastErr),
	})
}

func (s *WebServer) handleRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cwd, _ := os.Getwd()
	projRulesDir := rules.DefaultProjectRulesDir(cwd)
	projFile := filepath.Join(projRulesDir, "project_style.md")

	homeRulesDir := rules.DefaultHomeRulesDir()
	globalFile := filepath.Join(homeRulesDir, "global_style.md")

	if r.Method == http.MethodGet {
		projData, _ := os.ReadFile(projFile)
		globalData, _ := os.ReadFile(globalFile)
		json.NewEncoder(w).Encode(map[string]any{
			"project_rules": string(projData),
			"global_rules":  string(globalData),
		})
		return
	}

	if r.Method == http.MethodPost {
		var body struct {
			ProjectRules *string `json:"project_rules"`
			GlobalRules  *string `json:"global_rules"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
			return
		}

		if body.ProjectRules != nil {
			_ = os.MkdirAll(projRulesDir, 0o755)
			_ = os.WriteFile(projFile, []byte(*body.ProjectRules), 0o644)
		}
		if body.GlobalRules != nil {
			_ = os.MkdirAll(homeRulesDir, 0o755)
			_ = os.WriteFile(globalFile, []byte(*body.GlobalRules), 0o644)
		}
		json.NewEncoder(w).Encode(map[string]any{"success": true})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *WebServer) handleDiag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_ = os.MkdirAll("output/novel", 0755)
	storeInst := store.NewStore("output/novel")
	exportPath, err := diag.Export(storeInst)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	data, err := os.ReadFile(exportPath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Đã xuất chẩn đoán nhưng không thể đọc file: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success":      true,
		"diag_content": string(data),
		"export_path":  exportPath,
	})
}
