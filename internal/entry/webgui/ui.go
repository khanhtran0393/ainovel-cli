package webgui

// WebUIHTML chứa giao diện HTML/CSS/JS đỉnh cao (Glassmorphism, Dark Mode sang trọng, TailwindCSS, Font Google, Lucide Icons, SSE Realtime streaming).
const WebUIHTML = `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Novel Studio — Enterprise Authoring System</title>
    <!-- Tailwind CSS CDN -->
    <script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
    <!-- Google Fonts: Inter & Outfit -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Outfit:wght@400;500;600;700&display=swap" rel="stylesheet">
    <!-- Lucide Icons -->
    <script src="https://unpkg.com/lucide@latest"></script>
    <!-- Vue.js CDN -->
    <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
    <style>
        body {
            font-family: 'Inter', sans-serif;
            background-color: #0b0f19;
            color: #f8fafc;
        }
        h1, h2, h3, h4, h5, h6, .brand {
            font-family: 'Outfit', sans-serif;
        }
        ::-webkit-scrollbar {
            width: 6px;
            height: 6px;
        }
        ::-webkit-scrollbar-track {
            background: #111827;
        }
        ::-webkit-scrollbar-thumb {
            background: #374151;
            border-radius: 3px;
        }
        ::-webkit-scrollbar-thumb:hover {
            background: #4b5563;
        }
        .glass-panel {
            background: rgba(17, 24, 39, 0.75);
            backdrop-filter: blur(16px);
            border: 1px solid rgba(255, 255, 255, 0.08);
            box-shadow: 0 10px 30px 0 rgba(0, 0, 0, 0.4);
        }
        .glass-header {
            background: rgba(11, 15, 25, 0.85);
            backdrop-filter: blur(20px);
            border-bottom: 1px solid rgba(255, 255, 255, 0.08);
        }
        .pulse-border {
            box-shadow: 0 0 25px rgba(99, 102, 241, 0.25);
            border-color: rgba(99, 102, 241, 0.4);
            animation: pulse-indigo 3s infinite;
        }
        @keyframes pulse-indigo {
            0% { box-shadow: 0 0 15px rgba(99, 102, 241, 0.2); }
            50% { box-shadow: 0 0 35px rgba(99, 102, 241, 0.4); }
            100% { box-shadow: 0 0 15px rgba(99, 102, 241, 0.2); }
        }
        .toast {
            position: fixed;
            bottom: 24px;
            right: 24px;
            z-index: 9999;
            display: flex;
            align-items: center;
            padding: 12px 24px;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.6);
            transition: all 0.3s ease;
        }
        /* Hiệu ứng Radar Scanner */
        .radar-line {
            width: 100%;
            height: 2px;
            background: linear-gradient(90deg, transparent, rgba(99, 102, 241, 0.8), transparent);
            animation: scan 2.5s linear infinite;
        }
        @keyframes scan {
            0% { transform: translateY(0px); opacity: 0.2; }
            50% { opacity: 1; }
            100% { transform: translateY(180px); opacity: 0.2; }
        }
    </style>
</head>
<body class="min-h-screen bg-slate-950 selection:bg-indigo-500 selection:text-white">
    <div id="app" class="min-h-screen flex flex-col">
        <!-- Toast Notification -->
        <div v-if="toast.show" :class="['toast text-xs font-medium text-white border flex items-center space-x-2', toast.type === 'success' ? 'bg-emerald-600 border-emerald-500' : 'bg-rose-600 border-rose-500']">
            <i v-if="toast.type === 'success'" data-lucide="check-circle" class="w-4 h-4"></i>
            <i v-else data-lucide="alert-circle" class="w-4 h-4"></i>
            <span>{{ toast.message }}</span>
        </div>

        <!-- Header / Navbar -->
        <header class="glass-header sticky top-0 z-50 px-6 py-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
                <div class="p-2.5 bg-gradient-to-tr from-indigo-600 via-indigo-500 to-violet-500 rounded-xl shadow-lg shadow-indigo-500/30">
                    <i data-lucide="cpu" class="w-6 h-6 text-white"></i>
                </div>
                <div>
                    <h1 class="text-xl font-bold tracking-tight bg-gradient-to-r from-indigo-200 via-white to-indigo-200 bg-clip-text text-transparent">AI Novel Studio</h1>
                    <p class="text-xs text-indigo-400 font-medium flex items-center space-x-1">
                        <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 animate-ping"></span>
                        <span>Enterprise Master Orchestrator v3.0</span>
                    </p>
                </div>
            </div>
            
            <!-- NAVBAR CHUYỂN TAB KÈM TRẠNG THÁI ACTIVE -->
            <nav class="flex items-center space-x-1.5 bg-slate-900/90 p-1.5 rounded-xl border border-slate-800 shadow-inner">
                <button @click="switchTab('compose')" :class="['px-4 py-2 rounded-lg font-medium text-sm transition flex items-center space-x-2', currentTab === 'compose' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/40 border border-indigo-500/30' : 'text-slate-400 hover:text-white hover:bg-slate-800/60']">
                    <i data-lucide="pen-tool" class="w-4 h-4"></i>
                    <span>Sáng tác</span>
                    <span v-if="isRunning && currentTab !== 'compose'" class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                </button>
                <button @click="switchTab('chapters')" :class="['px-4 py-2 rounded-lg font-medium text-sm transition flex items-center space-x-2', currentTab === 'chapters' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/40 border border-indigo-500/30' : 'text-slate-400 hover:text-white hover:bg-slate-800/60']">
                    <i data-lucide="library" class="w-4 h-4"></i>
                    <span>Danh sách chương</span>
                    <span class="text-[10px] px-1.5 py-0.5 bg-slate-950 text-indigo-300 rounded-full font-mono">{{ chaptersList.length }}</span>
                </button>
                <button @click="switchTab('rules')" :class="['px-4 py-2 rounded-lg font-medium text-sm transition flex items-center space-x-2', currentTab === 'rules' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/40 border border-indigo-500/30' : 'text-slate-400 hover:text-white hover:bg-slate-800/60']">
                    <i data-lucide="shield-check" class="w-4 h-4"></i>
                    <span>Quy tắc & Chẩn đoán</span>
                    <span class="w-1.5 h-1.5 rounded-full bg-amber-400 animate-ping"></span>
                </button>
                <button @click="switchTab('config')" :class="['px-4 py-2 rounded-lg font-medium text-sm transition flex items-center space-x-2', currentTab === 'config' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/40 border border-indigo-500/30' : 'text-slate-400 hover:text-white hover:bg-slate-800/60']">
                    <i data-lucide="sliders" class="w-4 h-4"></i>
                    <span>Cấu hình & API Keys</span>
                </button>
            </nav>

            <div class="flex items-center space-x-4">
                <!-- Hiệu ứng Engine Status Trực quan -->
                <div :class="['flex items-center space-x-2.5 px-4 py-2 rounded-xl border transition duration-300', isRunning ? 'bg-emerald-500/10 border-emerald-500/40 text-emerald-300 shadow-lg shadow-emerald-500/10' : 'bg-slate-900 border-slate-800 text-slate-300']">
                    <span :class="['w-2.5 h-2.5 rounded-full', isRunning ? 'bg-emerald-400 animate-ping' : 'bg-amber-500']"></span>
                    <span class="text-xs font-bold uppercase tracking-wider font-mono">{{ isRunning ? 'Engine: ĐANG HOẠT ĐỘNG' : 'Engine: Tạm dừng' }}</span>
                </div>
                
                <!-- Nút Tắt App -->
                <button @click="shutdownApp" class="px-3 py-2 bg-rose-600/20 hover:bg-rose-600/30 text-rose-400 hover:text-rose-300 border border-rose-500/30 rounded-xl font-medium text-xs transition flex items-center space-x-1.5 shadow-lg shadow-rose-500/10">
                    <i data-lucide="power" class="w-3.5 h-3.5"></i>
                    <span>Tắt App</span>
                </button>
            </div>
        </header>

        <!-- Main content -->
        <main class="flex-1 max-w-7xl w-full mx-auto px-6 py-8 flex flex-col">
            
            <!-- MÀN HÌNH LOADING KHI CHUYỂN TAB (SKELETON ANIMATION) -->
            <div v-if="isLoadingTab" class="flex-1 flex flex-col items-center justify-center space-y-4 my-24">
                <div class="p-4 bg-indigo-600/10 border border-indigo-500/20 rounded-2xl animate-pulse">
                    <i data-lucide="loader-2" class="w-10 h-10 text-indigo-400 animate-spin"></i>
                </div>
                <h3 class="text-sm font-semibold text-slate-300">Đang đồng bộ hóa dữ liệu Tab với máy chủ...</h3>
                <p class="text-xs text-slate-500 font-mono">Quét bộ nhớ RAM & Kiểm toán tệp tin hệ thống</p>
            </div>

            <!-- TAB 1: SÁNG TÁC -->
            <div v-show="!isLoadingTab && currentTab === 'compose'" class="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-8">
                <!-- Vùng Trái: Bảng điều khiển khởi tạo -->
                <div class="lg:col-span-5 flex flex-col space-y-6">
                    <div class="glass-panel p-6 rounded-2xl flex flex-col space-y-5">
                        <div class="flex items-center justify-between border-b border-slate-800 pb-4">
                            <h2 class="text-lg font-semibold text-white flex items-center space-x-2">
                                <i data-lucide="sparkles" class="w-5 h-5 text-indigo-400"></i>
                                <span>Khởi tạo Tiểu thuyết</span>
                            </h2>
                            <span class="text-xs px-2.5 py-1 bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded-full font-mono font-medium">v3.0 AI Mode</span>
                        </div>

                        <!-- Chế độ viết -->
                        <div class="flex space-x-2 p-1 bg-slate-900 rounded-xl border border-slate-800">
                            <button @click="composeMode = 'new'" :class="['flex-1 py-2 text-xs font-medium rounded-lg transition', composeMode === 'new' ? 'bg-indigo-600 text-white shadow' : 'text-slate-400 hover:text-white']">
                                Viết truyện mới
                            </button>
                            <button @click="composeMode = 'resume'" :class="['flex-1 py-2 text-xs font-medium rounded-lg transition', composeMode === 'resume' ? 'bg-indigo-600 text-white shadow' : 'text-slate-400 hover:text-white']">
                                Tiếp tục phiên cũ
                            </button>
                        </div>

                        <!-- Nhập Prompt & AI Assist cho truyện mới -->
                        <div v-if="composeMode === 'new'" class="flex flex-col space-y-4">
                            <!-- THÔNG BÁO DỌN DẸP TRUYỆN CŨ -->
                            <div v-if="chaptersList.length > 0 && !isRunning" class="p-3 bg-rose-500/10 border border-rose-500/20 rounded-xl flex items-start space-x-2.5 text-xs text-rose-200/90 leading-relaxed">
                                <i data-lucide="info" class="w-4 h-4 text-rose-400 shrink-0 mt-0.5"></i>
                                <div>
                                    <span>Hệ thống phát hiện có <b>{{ chaptersList.length }} chương</b> từ phiên làm việc cũ. Khi bạn bấm <b>'Bắt đầu viết mới'</b>, động cơ sẽ tự động sao lưu & dọn sạch bộ nhớ để xuất phát trang mới từ Chương 1!</span>
                                </div>
                            </div>

                            <!-- Tuỳ biến số chương Hỗ trợ 500+ Chương -->
                            <div class="flex flex-col space-y-3 p-4 bg-slate-900/60 border border-slate-800 rounded-xl shadow-inner">
                                <div class="flex items-center justify-between">
                                    <span class="text-xs font-medium text-slate-300 flex items-center space-x-1.5">
                                        <i data-lucide="layers" class="w-3.5 h-3.5 text-indigo-400"></i>
                                        <span>Quy mô tác phẩm (Số chương)</span>
                                    </span>
                                    <div class="flex items-center space-x-1.5">
                                        <input type="number" v-model="chapterCount" min="1" max="1000" class="w-16 p-1 bg-slate-950 border border-indigo-500/40 rounded-lg text-xs text-indigo-300 font-mono font-bold text-center focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                                        <span class="text-xs text-slate-400 font-medium">chương</span>
                                    </div>
                                </div>
                                <input type="range" v-model="chapterCount" min="1" max="1000" class="w-full accent-indigo-500 bg-slate-800 h-1.5 rounded-lg appearance-none cursor-pointer" />
                                <div class="flex justify-between text-[10px] text-slate-500 font-mono">
                                    <span>1 chương (Truyện ngắn)</span>
                                    <span>500+ chương (Siêu Trường Thiên)</span>
                                </div>
                            </div>

                            <!-- Trợ lý AI tạo Prompt -->
                            <div class="flex flex-col space-y-2 p-4 bg-indigo-950/30 border border-indigo-800/30 rounded-xl shadow-inner">
                                <label class="text-xs font-medium text-indigo-300 flex items-center justify-between">
                                    <span class="flex items-center space-x-1.5">
                                        <i data-lucide="wand-2" class="w-3.5 h-3.5 text-indigo-400"></i>
                                        <span>AI Hỗ trợ Tạo ý tưởng sáng tác</span>
                                    </span>
                                </label>
                                <div class="flex space-x-2">
                                    <input v-model="ideaKeywords" type="text" placeholder="Gõ vài từ khóa: án mạng, biệt thự cổ..." class="flex-1 p-2.5 bg-slate-900 border border-slate-700 rounded-lg text-xs text-slate-200 placeholder:text-slate-600 focus:outline-none focus:ring-1 focus:ring-indigo-500" />
                                    <button @click="enhancePrompt" :disabled="isEnhancing" class="px-3 py-2 bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-medium rounded-lg transition flex items-center space-x-1 shadow-md shrink-0">
                                        <i v-if="!isEnhancing" data-lucide="zap" class="w-3.5 h-3.5 fill-current text-amber-300"></i>
                                        <i v-else data-lucide="loader-2" class="w-3.5 h-3.5 animate-spin"></i>
                                        <span>{{ isEnhancing ? 'Đang nghĩ...' : 'Phát triển' }}</span>
                                    </button>
                                </div>
                                <p class="text-[11px] text-slate-400 leading-normal mt-1">Chỉ cần gõ vài từ khóa, AI sẽ tự động kiến trúc bối cảnh và mâu thuẫn chính siêu chi tiết.</p>
                            </div>

                            <div class="flex flex-col space-y-2">
                                <label class="text-xs font-medium text-slate-300 flex items-center justify-between">
                                    <span>Nội dung Gợi ý sáng tác (Prompt)</span>
                                    <span class="text-slate-500">Bắt buộc</span>
                                </label>
                                <textarea v-model="promptText" rows="6" placeholder="Nội dung miêu tả chi tiết bối cảnh, nhân vật và phong cách văn học..." class="w-full p-4 bg-slate-900/80 border border-slate-800 rounded-xl text-xs text-slate-100 placeholder:text-slate-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition leading-relaxed"></textarea>
                            </div>
                        </div>

                        <!-- Lời nhắc tiếp tục phiên cũ -->
                        <div v-if="composeMode === 'resume'" class="p-4 bg-amber-500/10 border border-amber-500/20 rounded-xl flex items-start space-x-3">
                            <i data-lucide="history" class="w-5 h-5 text-amber-400 shrink-0 mt-0.5"></i>
                            <p class="text-xs text-amber-200/90 leading-relaxed">
                                Động cơ sẽ tự động khôi phục ngữ cảnh (Context Engine), tóm tắt (Store Summary) và các chương đã lưu trong thư mục <code class="bg-slate-900 px-1 py-0.5 rounded text-amber-300">output/novel</code> để tiếp tục sáng tác liền mạch.
                            </p>
                        </div>

                        <button @click="startWriting" :disabled="isRunning" :class="['w-full py-3.5 px-4 rounded-xl font-bold text-sm flex items-center justify-center space-x-2 shadow-xl transition uppercase tracking-wider', isRunning ? 'bg-slate-800 text-slate-500 cursor-not-allowed border border-slate-700' : 'bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 text-white shadow-indigo-500/25 active:scale-[0.98]']">
                            <i v-if="!isRunning" data-lucide="play" class="w-4 h-4 fill-current"></i>
                            <i v-else data-lucide="loader-2" class="w-4 h-4 animate-spin"></i>
                            <span>{{ isRunning ? 'Động cơ đang sáng tác...' : (composeMode === 'new' ? 'Bắt đầu viết mới' : 'Khôi phục & Viết tiếp') }}</span>
                        </button>
                    </div>

                    <!-- Bảng trạng thái & thống kê nhanh ĐÃ SỬA LOGIC SO VỚI TRUYỆN CŨ -->
                    <div class="glass-panel p-6 rounded-2xl flex flex-col space-y-4">
                        <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider flex items-center justify-between">
                            <span class="flex items-center space-x-2">
                                <i data-lucide="activity" class="w-4 h-4 text-emerald-400"></i>
                                <span>Tiến trình Động cơ</span>
                            </span>
                            <span class="text-[10px] text-emerald-400 font-mono flex items-center space-x-1">
                                <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
                                <span>Live Status</span>
                            </span>
                        </h3>
                        <div class="grid grid-cols-2 gap-4">
                            <div class="p-4 bg-slate-900/60 border border-slate-800/80 rounded-xl flex flex-col">
                                <span class="text-xs text-slate-400">Số chương phiên này</span>
                                <span class="text-2xl font-bold text-white mt-1 flex items-baseline space-x-1">
                                    <span>{{ isRunning ? chaptersList.length : (composeMode === 'new' ? '0' : chaptersList.length) }}</span>
                                    <span v-if="composeMode === 'new' && !isRunning" class="text-xs font-normal text-slate-500 font-sans">(Chờ chạy)</span>
                                    <span v-else class="text-xs font-normal text-emerald-400 font-sans">chương</span>
                                </span>
                            </div>
                            <div class="p-4 bg-slate-900/60 border border-slate-800/80 rounded-xl flex flex-col">
                                <span class="text-xs text-slate-400">Thao tác cuối</span>
                                <span class="text-xs font-medium text-indigo-400 mt-2 truncate">{{ lastAction || 'Ready / Idle' }}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Vùng Phải: Terminal Live Streaming -->
                <div class="lg:col-span-7 flex flex-col space-y-4">
                    <div class="flex items-center justify-between px-2">
                        <h3 class="text-sm font-semibold text-slate-300 flex items-center space-x-2">
                            <i data-lucide="terminal" class="w-4 h-4 text-indigo-400"></i>
                            <span>Luồng xuất văn bản & Logs Trực tiếp (SSE Streaming)</span>
                        </h3>
                        <div class="flex items-center space-x-2">
                            <span v-if="isRunning" class="text-xs font-mono text-indigo-400 flex items-center space-x-1 bg-indigo-500/10 px-2.5 py-1 rounded-lg border border-indigo-500/20">
                                <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 animate-ping"></span>
                                <span>Đang truyền luồng...</span>
                            </span>
                            <button @click="clearLog" class="text-xs text-slate-400 hover:text-slate-200 bg-slate-900 px-3 py-1.5 rounded-lg border border-slate-800 transition flex items-center space-x-1">
                                <i data-lucide="trash-2" class="w-3.5 h-3.5"></i>
                                <span>Xóa log</span>
                            </button>
                        </div>
                    </div>

                    <!-- Vùng hiển thị Streaming kèm hiệu ứng phát sáng khi hoạt động -->
                    <div :class="['flex-1 min-h-[600px] max-h-[750px] glass-panel rounded-2xl p-6 flex flex-col font-mono text-xs overflow-y-auto border transition duration-500 bg-slate-950/90 relative', isRunning ? 'pulse-border' : 'border-slate-800 shadow-inner']" id="stream-terminal">
                        <!-- Hiệu ứng Radar Line khi engine chạy -->
                        <div v-if="isRunning" class="absolute top-0 left-0 right-0 radar-line pointer-events-none"></div>

                        <div v-if="streamLines.length === 0" class="flex-1 flex flex-col items-center justify-center text-slate-600 space-y-3 my-12">
                            <i data-lucide="cpu" class="w-12 h-12 text-slate-700 stroke-1"></i>
                            <p class="text-xs font-sans text-center max-w-md">Động cơ đang trong trạng thái chờ. Hãy thiết lập thông số bên trái và nhấn nút Bắt đầu để quan sát toàn bộ quá trình Tác nhân suy luận và viết truyện trực tiếp...</p>
                        </div>
                        <div v-for="(line, idx) in streamLines" :key="idx" class="py-1 whitespace-pre-wrap leading-relaxed break-words" :class="getLineClass(line)">{{ line.text }}</div>
                    </div>
                </div>
            </div>

            <!-- TAB 2: DANH SÁCH CHƯƠNG -->
            <div v-show="!isLoadingTab && currentTab === 'chapters'" class="flex-1 flex flex-col space-y-6">
                <!-- Tiêu đề & Công dụng của Tab -->
                <div class="p-6 bg-slate-900/90 border border-slate-800 rounded-2xl flex items-start space-x-4 shadow-inner">
                    <div class="p-3 bg-indigo-500/10 text-indigo-400 rounded-xl border border-indigo-500/20 shrink-0">
                        <i data-lucide="library" class="w-6 h-6"></i>
                    </div>
                    <div class="flex flex-col space-y-2 w-full">
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-white">Quản lý, Đọc & Chỉnh sửa Các Chương Đã Viết</h3>
                            <span class="text-xs font-mono text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-lg border border-emerald-500/20 flex items-center space-x-1">
                                <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span>
                                <span>Thư mục đồng bộ: output/novel/chapters</span>
                            </span>
                        </div>
                        <p class="text-xs text-slate-400 leading-relaxed">
                            Toàn bộ các chương tiểu thuyết do AI sáng tác được tự động lưu trữ tại thư mục hệ thống. Tab này giúp bạn theo dõi tiến trình sáng tác theo thời gian thực, đọc nội dung trọn vẹn, chỉnh sửa câu chữ trực tiếp theo ý muốn và tải xuống trọn bộ tiểu thuyết thành 1 file hoàn chỉnh.
                        </p>
                    </div>
                </div>

                <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
                    <!-- Danh sách chương -->
                    <div class="lg:col-span-4 flex flex-col space-y-4">
                        <div class="flex items-center justify-between px-2">
                            <h3 class="text-sm font-semibold text-slate-300 flex items-center space-x-2">
                                <i data-lucide="book" class="w-4 h-4 text-indigo-400"></i>
                                <span>Danh sách chương đã viết</span>
                            </h3>
                            <button @click="loadChapters(true)" class="text-xs text-indigo-400 hover:text-indigo-300 flex items-center space-x-1 bg-indigo-500/10 px-2.5 py-1 rounded-lg border border-indigo-500/20 transition shadow">
                                <i data-lucide="refresh-cw" class="w-3 h-3"></i>
                                <span>Làm mới bộ nhớ</span>
                            </button>
                        </div>
                        <div class="glass-panel p-4 rounded-2xl flex flex-col space-y-2 max-h-[650px] overflow-y-auto">
                            <div v-if="chaptersList.length === 0" class="p-8 text-center text-xs text-slate-500 font-sans">
                                Chưa có chương truyện nào trong thư mục <code class="bg-slate-900 px-1 py-0.5 rounded text-slate-400">output/novel/chapters</code>
                            </div>
                            <!-- ĐÃ SỬA CÚ PHÁP OPTIONAL CHAINING THÀNH LOGIC AN TOÀN -->
                            <button v-for="ch in chaptersList" :key="ch.path" @click="selectChapter(ch)" :class="['w-full p-4 rounded-xl text-left border transition flex flex-col space-y-2.5', (selectedChapter && selectedChapter.path === ch.path) ? 'bg-indigo-600/20 border-indigo-500/50 text-white shadow-lg shadow-indigo-500/10' : 'bg-slate-900/60 border-slate-800/80 text-slate-300 hover:border-slate-700']">
                                <div class="flex items-center justify-between w-full">
                                    <span class="font-semibold text-sm truncate text-white">{{ ch.title || ch.name }}</span>
                                    <span class="text-xs px-2.5 py-0.5 bg-slate-800 text-indigo-400 border border-slate-700 rounded-full shrink-0 font-mono">{{ ch.size }} KB</span>
                                </div>
                                <div class="text-xs text-slate-400 line-clamp-2 leading-relaxed font-sans">{{ ch.excerpt }}</div>
                            </button>
                        </div>
                        
                        <!-- Download Full Novel -->
                        <a href="/api/download-all" target="_blank" class="w-full py-3 bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 text-white text-xs font-medium rounded-xl shadow-lg shadow-emerald-600/20 transition flex items-center justify-center space-x-2">
                            <i data-lucide="package" class="w-4 h-4"></i>
                            <span>📦 Tải xuống Trọn bộ Tiểu thuyết (.md)</span>
                        </a>
                    </div>

                    <!-- Vùng đọc & chỉnh sửa truyện -->
                    <div class="lg:col-span-8 flex flex-col space-y-4">
                        <div class="flex items-center justify-between px-2">
                            <div class="flex items-center space-x-3">
                                <h3 class="text-sm font-semibold text-slate-300 flex items-center space-x-2">
                                    <i data-lucide="eye" class="w-4 h-4 text-indigo-400"></i>
                                    <span>Chế độ Đọc & Chỉnh sửa</span>
                                </h3>
                                <span v-if="selectedChapter" class="text-xs text-slate-400 font-mono">Tệp: {{ selectedChapter.name }}</span>
                            </div>
                            
                            <div v-if="selectedChapter" class="flex items-center space-x-2">
                                <button @click="isEditing = !isEditing" :class="['px-3 py-1.5 rounded-lg text-xs font-medium border transition flex items-center space-x-1.5', isEditing ? 'bg-amber-500/20 text-amber-300 border-amber-500/40' : 'bg-slate-900 text-slate-300 border-slate-800 hover:bg-slate-800']">
                                    <i data-lucide="edit-3" class="w-3.5 h-3.5"></i>
                                    <span>{{ isEditing ? 'Xem văn bản' : '✏️ Chỉnh sửa' }}</span>
                                </button>
                                <button v-if="isEditing" @click="saveChapterContent" class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-lg text-xs font-medium shadow-md transition flex items-center space-x-1.5">
                                    <i data-lucide="save" class="w-3.5 h-3.5"></i>
                                    <span>Lưu thay đổi</span>
                                </button>
                                <button @click="downloadCurrentChapter" class="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-lg text-xs font-medium border border-slate-700 transition flex items-center space-x-1.5">
                                    <i data-lucide="download" class="w-3.5 h-3.5"></i>
                                    <span>Tải file này</span>
                                </button>
                            </div>
                        </div>

                        <div class="flex-1 min-h-[600px] max-h-[750px] glass-panel rounded-2xl p-8 flex flex-col overflow-y-auto border border-slate-800">
                            <div v-if="!selectedChapter" class="flex-1 flex flex-col items-center justify-center text-slate-600 space-y-3 my-12 font-sans">
                                <i data-lucide="file-text" class="w-12 h-12 text-slate-700 stroke-1"></i>
                                <p class="text-xs">Chọn một chương truyện ở danh sách bên trái để đọc toàn bộ nội dung hoặc chỉnh sửa trực tiếp...</p>
                            </div>
                            <div v-else-if="!isEditing" class="text-slate-100 font-serif text-lg leading-relaxed whitespace-pre-wrap max-w-3xl mx-auto w-full selection:bg-indigo-500 selection:text-white">
                                {{ chapterContent }}
                            </div>
                            <textarea v-else v-model="chapterContent" class="flex-1 w-full p-4 bg-slate-955 border border-indigo-500/30 rounded-xl font-mono text-sm text-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition leading-relaxed"></textarea>
                        </div>
                    </div>
                </div>
            </div>

            <!-- TAB 3: QUY TẮC & CHẨN ĐOÁN -->
            <div v-show="!isLoadingTab && currentTab === 'rules'" class="flex-1 flex flex-col space-y-8 max-w-6xl mx-auto w-full">
                <!-- Tiêu đề & Công dụng của Tab -->
                <div class="p-6 bg-slate-900/90 border border-slate-800 rounded-2xl flex items-start space-x-4 shadow-inner">
                    <div class="p-3 bg-indigo-500/10 text-indigo-400 rounded-xl border border-indigo-500/20 shrink-0">
                        <i data-lucide="shield-check" class="w-6 h-6"></i>
                    </div>
                    <div class="flex flex-col space-y-2 w-full">
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-white">Quản trị Quy tắc Sáng tác & Chẩn đoán Báo cáo Hệ thống</h3>
                            <span class="text-xs font-mono text-amber-400 bg-amber-500/10 px-2.5 py-1 rounded-lg border border-amber-500/20 flex items-center space-x-1">
                                <span class="w-1.5 h-1.5 rounded-full bg-amber-400 animate-ping"></span>
                                <span>Trình Quét Cú Pháp: Hoạt Động</span>
                            </span>
                        </div>
                        <p class="text-xs text-slate-400 leading-relaxed">
                            Tab này cung cấp hai tính năng cấp cao: (1) Quản trị bộ quy tắc kiểm duyệt nội dung của Tác nhân (như quy định số từ, cụm từ cấm, văn phong tùy chọn) tại <code class="text-indigo-300 bg-slate-950 px-1.5 py-0.5 rounded">./.ainovel/rules</code>. (2) Chạy báo cáo chẩn đoán luồng sáng tác (Runtime Diagnostics) để kiểm toán trạng thái bộ nhớ, tránh tình trạng kẹt chu trình.
                        </p>
                    </div>
                </div>

                <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
                    <!-- KHU VỰC QUY TẮC SÁNG TÁC (COL-SPAN-7) -->
                    <div class="lg:col-span-7 flex flex-col space-y-6">
                        <div class="glass-panel p-6 rounded-2xl flex flex-col space-y-4">
                            <div class="flex items-center justify-between border-b border-slate-800 pb-4">
                                <h3 class="text-sm font-semibold text-white flex items-center space-x-2">
                                    <i data-lucide="book-open-check" class="w-4 h-4 text-indigo-400"></i>
                                    <span>Quy tắc Dự án hiện tại (project_style.md)</span>
                                </h3>
                                <button @click="saveRules" class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white font-medium text-xs rounded-xl shadow-lg shadow-indigo-600/20 transition flex items-center space-x-1.5">
                                    <i data-lucide="save" class="w-3.5 h-3.5"></i>
                                    <span>Lưu Quy tắc</span>
                                </button>
                            </div>
                            <textarea v-model="projectRules" rows="12" class="w-full p-4 bg-slate-950 border border-slate-800 rounded-xl font-mono text-xs text-indigo-200 placeholder:text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition leading-relaxed" placeholder="Nhập quy tắc kiểm duyệt của bạn tại đây..."></textarea>
                            
                            <!-- Hướng dẫn nhanh Template YAML -->
                            <div class="p-4 bg-slate-900/60 border border-slate-800 rounded-xl flex flex-col space-y-2 text-[11px] text-slate-400">
                                <span class="font-semibold text-slate-300">💡 Hướng dẫn cấu hình YAML (Tùy chọn):</span>
                                <p>Bạn có thể đặt khối YAML ở đầu file để thiết lập các ràng buộc cứng cho Tác nhân kiểm tra (commit_chapter):</p>
                                <code class="bg-slate-950 p-2.5 rounded text-amber-300 font-mono leading-normal">---
chapter_words: 2500-5000               # Ràng buộc số từ mỗi chương
forbidden_phrases: ["bất chợt", "có thể nói"] # Từ cấm
fatigue_words: {không khỏi: 1}         # Cảnh báo từ sáo rỗng
---</code>
                            </div>
                        </div>

                        <div class="glass-panel p-6 rounded-2xl flex flex-col space-y-4">
                            <div class="flex items-center justify-between border-b border-slate-800 pb-4">
                                <h3 class="text-sm font-semibold text-white flex items-center space-x-2">
                                    <i data-lucide="globe" class="w-4 h-4 text-amber-400"></i>
                                    <span>Quy tắc Toàn cục (~/.ainovel/rules/global_style.md)</span>
                                </h3>
                                <button @click="saveRules" class="px-4 py-1.5 bg-amber-600 hover:bg-amber-500 text-white font-medium text-xs rounded-xl shadow-lg shadow-amber-600/20 transition flex items-center space-x-1.5">
                                    <i data-lucide="save" class="w-3.5 h-3.5"></i>
                                    <span>Lưu Toàn cục</span>
                                </button>
                            </div>
                            <textarea v-model="globalRules" rows="8" class="w-full p-4 bg-slate-950 border border-slate-800 rounded-xl font-mono text-xs text-amber-200 placeholder:text-slate-700 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-transparent transition leading-relaxed" placeholder="Quy tắc toàn cục có hiệu lực trên mọi cuốn sách..."></textarea>
                        </div>
                    </div>

                    <!-- KHU VỰC BÁO CÁO CHẨN ĐOÁN (COL-SPAN-5) -->
                    <div class="lg:col-span-5 flex flex-col space-y-4">
                        <div class="flex items-center justify-between px-2">
                            <h3 class="text-sm font-semibold text-slate-300 flex items-center space-x-2">
                                <i data-lucide="activity" class="w-4 h-4 text-emerald-400"></i>
                                <span>Báo cáo Chẩn đoán Hệ thống</span>
                            </h3>
                            <button @click="runDiag" :disabled="isDiagRunning" class="px-4 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white font-medium text-xs rounded-xl shadow-lg shadow-emerald-600/20 transition flex items-center space-x-1.5">
                                <i v-if="!isDiagRunning" data-lucide="refresh-cw" class="w-3.5 h-3.5"></i>
                                <i v-else data-lucide="loader-2" class="w-3.5 h-3.5 animate-spin"></i>
                                <span>{{ isDiagRunning ? 'Đang chẩn đoán...' : 'Chạy Chẩn đoán' }}</span>
                            </button>
                        </div>

                        <div class="flex-1 min-h-[550px] glass-panel rounded-2xl p-6 flex flex-col font-mono text-xs overflow-y-auto border border-slate-800 bg-slate-950/90 relative">
                            <div v-if="isDiagRunning" class="absolute inset-0 bg-slate-950/80 backdrop-blur-sm flex flex-col items-center justify-center space-y-3 z-10">
                                <i data-lucide="loader-2" class="w-10 h-10 text-emerald-400 animate-spin"></i>
                                <span class="text-xs font-semibold text-emerald-300">Đang kiểm toán tệp tin & vùng nhớ RAM...</span>
                            </div>

                            <div v-if="!diagContent" class="flex-1 flex flex-col items-center justify-center text-slate-600 space-y-3 my-12 font-sans text-center">
                                <i data-lucide="stethoscope" class="w-12 h-12 text-slate-700 stroke-1"></i>
                                <p class="text-xs max-w-sm">Nhấn nút 'Chạy Chẩn đoán' để tổng hợp tín hiệu runtime, kiểm tra độ kẹt step, tình trạng cấp phát model và bằng chứng bất thường...</p>
                            </div>
                            <div v-else class="text-slate-200 whitespace-pre-wrap leading-relaxed selection:bg-indigo-500 selection:text-white">
                                {{ diagContent }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- TAB 4: CẤU HÌNH & API KEY -->
            <div v-show="!isLoadingTab && currentTab === 'config'" class="flex-1 flex flex-col space-y-8 max-w-6xl mx-auto w-full">
                <!-- Tiêu đề & Công dụng của Tab -->
                <div class="p-6 bg-slate-900/90 border border-slate-800 rounded-2xl flex items-start space-x-4 shadow-inner">
                    <div class="p-3 bg-indigo-500/10 text-indigo-400 rounded-xl border border-indigo-500/20 shrink-0">
                        <i data-lucide="key" class="w-6 h-6"></i>
                    </div>
                    <div class="flex flex-col space-y-2 w-full">
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-white">Quản lý API Key & Cấu hình Động cơ (Config.json / .env)</h3>
                            <div class="flex items-center space-x-2">
                                <button @click="pingKeys" :disabled="isPinging" class="px-3 py-1 bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 text-white text-xs font-medium rounded-lg shadow transition flex items-center space-x-1.5">
                                    <i v-if="!isPinging" data-lucide="radio" class="w-3.5 h-3.5 text-amber-300"></i>
                                    <i v-else data-lucide="loader-2" class="w-3.5 h-3.5 animate-spin"></i>
                                    <span>{{ isPinging ? 'Đang test API Keys...' : '⚡ Kiểm tra Quota API Keys' }}</span>
                                </button>
                                <span class="text-xs font-mono text-indigo-300 bg-indigo-500/10 px-2.5 py-1 rounded-lg border border-indigo-500/20 flex items-center space-x-1">
                                    <span class="w-1.5 h-1.5 rounded-full bg-indigo-400 animate-ping"></span>
                                    <span>Mật mã 256-bit</span>
                                </span>
                            </div>
                        </div>
                        <p class="text-xs text-slate-400 leading-relaxed">
                            Tab này đóng vai trò là Trung tâm Kiểm soát Truy cập & Cấu hình (Enterprise Key Management). Tại đây bạn có thể trực tiếp dán danh sách 9+ API Key của mình vào tệp <code class="text-amber-300 bg-slate-950 px-1.5 py-0.5 rounded">.env</code> hoặc tinh chỉnh tham số xoay vòng mô hình trong <code class="text-indigo-300 bg-slate-950 px-1.5 py-0.5 rounded">config.json</code>. Hệ thống sẽ lưu và áp dụng ngay lập tức vào bộ nhớ mà không cần khởi động lại.
                        </p>
                    </div>
                </div>

                <!-- BẢNG KẾT QUẢ PING TEST API KEYS -->
                <div v-if="pingResults" class="p-6 bg-slate-900/95 border border-indigo-500/40 rounded-2xl flex flex-col space-y-3 shadow-2xl">
                    <div class="flex items-center justify-between border-b border-slate-800 pb-3">
                        <h4 class="text-sm font-semibold text-white flex items-center space-x-2">
                            <i data-lucide="activity" class="w-4 h-4 text-emerald-400"></i>
                            <span>Báo cáo Kiểm toán Hạn mức API Keys (Live Quota Audit)</span>
                        </h4>
                        <button @click="pingResults = null" class="text-xs text-slate-400 hover:text-white transition">Đóng bảng</button>
                    </div>
                    <div class="text-xs font-mono whitespace-pre-wrap leading-relaxed" :class="pingResults.success ? 'text-emerald-300' : 'text-rose-300'">
                        {{ pingResults.message }}
                    </div>
                </div>

                <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    <!-- .ENV CONFIG -->
                    <div class="glass-panel p-6 rounded-2xl flex flex-col space-y-4">
                        <div class="flex items-center justify-between border-b border-slate-800 pb-4">
                            <h3 class="text-sm font-semibold text-white flex items-center space-x-2">
                                <i data-lucide="file-key" class="w-4 h-4 text-amber-400"></i>
                                <span>Cấu hình API Keys (.env)</span>
                            </h3>
                            <button @click="saveConfig('env')" class="px-4 py-1.5 bg-amber-600 hover:bg-amber-500 text-white font-medium text-xs rounded-xl shadow-lg shadow-amber-600/20 transition flex items-center space-x-1.5">
                                <i data-lucide="save" class="w-3.5 h-3.5"></i>
                                <span>Lưu .env</span>
                            </button>
                        </div>
                        <textarea v-model="envContent" rows="16" class="w-full p-4 bg-slate-950 border border-slate-800 rounded-xl font-mono text-xs text-amber-200 placeholder:text-slate-700 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-transparent transition" placeholder="GEMINI_API_KEY_1=AIzaSy...\nGEMINI_API_KEY_2=AIzaSy..."></textarea>
                    </div>

                    <!-- CONFIG.JSON -->
                    <div class="glass-panel p-6 rounded-2xl flex flex-col space-y-4">
                        <div class="flex items-center justify-between border-b border-slate-800 pb-4">
                            <h3 class="text-sm font-semibold text-white flex items-center space-x-2">
                                <i data-lucide="file-code" class="w-4 h-4 text-indigo-400"></i>
                                <span>Cấu hình Hệ thống (config.json)</span>
                            </h3>
                            <button @click="saveConfig('config')" class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white font-medium text-xs rounded-xl shadow-lg shadow-indigo-600/20 transition flex items-center space-x-1.5">
                                <i data-lucide="save" class="w-3.5 h-3.5"></i>
                                <span>Lưu config.json</span>
                            </button>
                        </div>
                        <textarea v-model="configContent" rows="16" class="w-full p-4 bg-slate-950 border border-slate-800 rounded-xl font-mono text-xs text-indigo-200 placeholder:text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition" placeholder="{...}"></textarea>
                    </div>
                </div>
            </div>
        </main>

        <!-- Footer -->
        <footer class="glass-header border-t border-slate-800/80 mt-auto py-4 px-6 flex items-center justify-between text-xs text-slate-500">
            <p>© 2026 AI Novel Studio. Created for Enterprise IAM & Advanced Algorithmic Creation.</p>
            <div class="flex items-center space-x-4 text-slate-400">
                <span>Core: Go / AgentCore</span>
                <span>•</span>
                <span>GUI: Glassmorphism Ultra</span>
            </div>
        </footer>
    </div>

    <script>
        const { createApp } = Vue

        createApp({
            data() {
                return {
                    currentTab: 'compose',
                    isLoadingTab: false,
                    composeMode: 'new',
                    promptText: '',
                    chapterCount: 10,
                    ideaKeywords: '',
                    isEnhancing: false,
                    isRunning: false,
                    lastAction: 'Ready / Idle',
                    streamLines: [],
                    chaptersList: [],
                    selectedChapter: null,
                    chapterContent: '',
                    isEditing: false,
                    envContent: '',
                    configContent: '',
                    projectRules: '',
                    globalRules: '',
                    diagContent: '',
                    isDiagRunning: false,
                    isPinging: false,
                    pingResults: null,
                    eventSource: null,
                    toast: { show: false, message: '', type: 'success' },
                }
            },
            mounted() {
                lucide.createIcons();
                this.checkStatus();
                this.loadChapters();
                this.loadConfig();
                this.loadRules();
                setInterval(this.checkStatus, 3000);
            },
            updated() {
                lucide.createIcons();
            },
            methods: {
                showToast(msg, type = 'success') {
                    this.toast.message = msg;
                    this.toast.type = type;
                    this.toast.show = true;
                    setTimeout(() => { this.toast.show = false; }, 4000);
                },
                switchTab(tab) {
                    if (this.currentTab === tab) return;
                    this.isLoadingTab = true;
                    this.currentTab = tab;
                    
                    // Kích hoạt đồng bộ hóa ngầm để mang lại cảm giác ứng dụng làm việc mượt mà
                    if (tab === 'chapters') this.loadChapters();
                    if (tab === 'rules') this.loadRules();
                    if (tab === 'config') this.loadConfig();

                    setTimeout(() => { 
                        this.isLoadingTab = false;
                        setTimeout(() => { lucide.createIcons(); }, 50);
                    }, 300);
                },
                getLineClass(line) {
                    if (line.type === 'error') return 'text-rose-400 font-semibold';
                    if (line.type === 'event') return 'text-indigo-300 font-semibold';
                    if (line.type === 'system') return 'text-amber-400 font-medium';
                    return 'text-slate-200';
                },
                clearLog() {
                    this.streamLines = [];
                    this.showToast('Đã xóa bộ nhớ log hiển thị.', 'success');
                },
                async checkStatus() {
                    try {
                        const res = await fetch('/api/status');
                        const data = await res.json();
                        this.isRunning = data.is_running;
                        if (data.last_action) this.lastAction = data.last_action;
                        if (this.isRunning && !this.eventSource) {
                            this.connectSSE();
                        } else if (!this.isRunning && this.eventSource) {
                            this.eventSource.close();
                            this.eventSource = null;
                        }
                    } catch (e) {
                        console.error("Status check error:", e);
                    }
                },
                async enhancePrompt() {
                    if (!this.ideaKeywords.trim()) {
                        this.showToast('Vui lòng gõ vài từ khóa để AI phát triển ý tưởng!', 'error');
                        return;
                    }
                    this.isEnhancing = true;
                    this.showToast('Đang kết nối AI để kiến trúc tình tiết truyện...', 'success');
                    try {
                        const res = await fetch('/api/enhance-prompt', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ keywords: this.ideaKeywords })
                        });
                        const data = await res.json();
                        if (data.success) {
                            this.promptText = data.enhanced_prompt;
                            this.showToast('✨ Phát triển ý tưởng thành công! Đã dán vào Prompt.', 'success');
                        } else {
                            this.showToast('Lỗi phát triển ý tưởng: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.showToast('Không thể kết nối tới AI: ' + e.message, 'error');
                    } finally {
                        this.isEnhancing = false;
                    }
                },
                async startWriting() {
                    if (this.isRunning) return;
                    this.clearLog();
                    this.streamLines.push({ text: '[SYSTEM] Đang khởi tạo tiến trình sáng tác...', type: 'system' });
                    this.showToast('Động cơ đang khởi động tiến trình sáng tác...', 'success');
                    
                    const endpoint = this.composeMode === 'new' ? '/api/start' : '/api/resume';
                    const body = this.composeMode === 'new' ? JSON.stringify({ prompt: this.promptText, chapter_count: parseInt(this.chapterCount) }) : JSON.stringify({});
                    
                    try {
                        const res = await fetch(endpoint, {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: body
                        });
                        const data = await res.json();
                        if (data.success) {
                            this.isRunning = true;
                            this.connectSSE();
                            if (this.composeMode === 'new') {
                                this.chaptersList = []; // Xóa danh sách truyện cũ trên UI để đếm sạch từ 0
                            }
                            this.showToast('🚀 Động cơ khởi động thành công! Vui lòng quan sát Logs.', 'success');
                        } else {
                            this.streamLines.push({ text: '[ERROR] Khởi tạo thất bại: ' + data.error, type: 'error' });
                            this.showToast('Khởi tạo thất bại: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.streamLines.push({ text: '[ERROR] Không thể kết nối tới máy chủ: ' + e.message, type: 'error' });
                        this.showToast('Lỗi kết nối máy chủ.', 'error');
                    }
                },
                connectSSE() {
                    if (this.eventSource) this.eventSource.close();
                    this.eventSource = new EventSource('/api/stream');
                    this.eventSource.onmessage = (event) => {
                        try {
                            const data = JSON.parse(event.data);
                            if (data.text) {
                                this.streamLines.push({ text: data.text, type: data.type || 'text' });
                                this.scrollToBottom();
                                // Refresh danh sách chương ngầm nếu có chương mới
                                if (data.text.includes('Hoàn thành chương') || data.text.includes('commit')) {
                                    this.loadChapters();
                                }
                            }
                        } catch(e) {
                            this.streamLines.push({ text: event.data, type: 'text' });
                            this.scrollToBottom();
                        }
                    };
                    this.eventSource.onerror = () => {
                        this.eventSource.close();
                        this.eventSource = null;
                    };
                },
                scrollToBottom() {
                    setTimeout(() => {
                        const container = document.getElementById('stream-terminal');
                        if (container) container.scrollTop = container.scrollHeight;
                    }, 50);
                },
                async loadChapters(showNotify = false) {
                    try {
                        const res = await fetch('/api/chapters');
                        const data = await res.json();
                        this.chaptersList = data.chapters || [];
                        if (showNotify) this.showToast('🔄 Đã đồng bộ hóa ' + this.chaptersList.length + ' chương từ bộ nhớ.', 'success');
                    } catch (e) {
                        console.error("Load chapters error:", e);
                    }
                },
                async selectChapter(ch) {
                    this.selectedChapter = ch;
                    this.isEditing = false;
                    try {
                        const res = await fetch('/api/chapters/' + encodeURIComponent(ch.name));
                        const data = await res.json();
                        this.chapterContent = data.content || 'Không thể đọc nội dung tệp...';
                    } catch (e) {
                        this.chapterContent = 'Lỗi kết nối khi tải tệp...';
                    }
                },
                async saveChapterContent() {
                    if (!this.selectedChapter) return;
                    try {
                        const res = await fetch('/api/chapters/' + encodeURIComponent(this.selectedChapter.name), {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ content: this.chapterContent })
                        });
                        const data = await res.json();
                        if (data.success) {
                            this.showToast('💾 Lưu nội dung chương thành công!', 'success');
                            this.isEditing = false;
                            this.loadChapters();
                        } else {
                            this.showToast('Lưu thất bại: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.showToast('Lỗi kết nối khi lưu chương.', 'error');
                    }
                },
                downloadCurrentChapter() {
                    if (!this.selectedChapter) return;
                    const blob = new Blob([this.chapterContent], { type: 'text/markdown;charset=utf-8;' });
                    const url = URL.createObjectURL(blob);
                    const link = document.createElement('a');
                    link.setAttribute('href', url);
                    link.setAttribute('download', this.selectedChapter.name);
                    link.style.visibility = 'hidden';
                    document.body.appendChild(link);
                    link.click();
                    document.body.removeChild(link);
                    this.showToast('📥 Đã tải xuống tệp ' + this.selectedChapter.name, 'success');
                },
                async loadConfig() {
                    try {
                        const res = await fetch('/api/config');
                        const data = await res.json();
                        this.envContent = data.env || '';
                        this.configContent = data.config || '';
                    } catch (e) {
                        console.error("Load config error:", e);
                    }
                },
                async saveConfig(type) {
                    try {
                        const body = type === 'env' ? { env: this.envContent } : { config: this.configContent };
                        const res = await fetch('/api/config', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(body)
                        });
                        const data = await res.json();
                        if (data.success) {
                            this.showToast('💾 Lưu cấu hình ' + type + ' thành công! Máy chủ đã ghi nhận.', 'success');
                        } else {
                            this.showToast('Lưu thất bại: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.showToast('Lỗi kết nối khi lưu cấu hình.', 'error');
                    }
                },
                async pingKeys() {
                    this.isPinging = true;
                    this.pingResults = { success: true, message: 'Đang gửi ping xác thực tới máy chủ Google Gemini API trên toàn bộ Key...' };
                    this.showToast('Đang kiểm tra Hạn mức (Quota) của toàn bộ API Keys...', 'success');
                    try {
                        const res = await fetch('/api/ping-keys', { method: 'POST' });
                        const data = await res.json();
                        this.pingResults = {
                            success: data.success,
                            message: data.message || data.error
                        };
                        if (data.success) {
                            this.showToast('⚡ Kiểm tra API Key hoàn tất! Đã có Key sẵn sàng hoạt động.', 'success');
                        } else {
                            this.showToast('Cảnh báo: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.pingResults = { success: false, message: 'Không thể kết nối tới máy chủ để test Key: ' + e.message };
                        this.showToast('Lỗi kết nối máy chủ.', 'error');
                    } finally {
                        this.isPinging = false;
                    }
                },
                async loadRules() {
                    try {
                        const res = await fetch('/api/rules');
                        const data = await res.json();
                        this.projectRules = data.project_rules || '';
                        this.globalRules = data.global_rules || '';
                    } catch (e) {
                        console.error("Load rules error:", e);
                    }
                },
                async saveRules() {
                    try {
                        const res = await fetch('/api/rules', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ project_rules: this.projectRules, global_rules: this.globalRules })
                        });
                        const data = await res.json();
                        if (data.success) {
                            this.showToast('💾 Lưu quy tắc sáng tác thành công! Hệ thống sẽ tự động gộp vào engine.', 'success');
                        } else {
                            this.showToast('Lưu quy tắc thất bại: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.showToast('Lỗi kết nối khi lưu quy tắc.', 'error');
                    }
                },
                async runDiag() {
                    this.isDiagRunning = true;
                    this.diagContent = 'Đang thu thập dữ liệu runtime, kiểm toán RAM và tổng hợp bằng chứng chẩn đoán...';
                    this.showToast('Đang tổng hợp dữ liệu chẩn đoán hệ thống...', 'success');
                    try {
                        const res = await fetch('/api/diag');
                        const data = await res.json();
                        if (data.success) {
                            this.diagContent = data.diag_content;
                            this.showToast('🔍 Đã xuất Báo cáo Chẩn đoán thành công!', 'success');
                        } else {
                            this.diagContent = 'Lỗi xuất chẩn đoán: ' + data.error;
                            this.showToast('Lỗi xuất chẩn đoán: ' + data.error, 'error');
                        }
                    } catch (e) {
                        this.diagContent = 'Không thể kết nối để chạy chẩn đoán: ' + e.message;
                        this.showToast('Lỗi kết nối khi chạy chẩn đoán.', 'error');
                    } finally {
                        this.isDiagRunning = false;
                    }
                },
                async shutdownApp() {
                    if (!confirm('Bạn có chắc chắn muốn tắt hoàn toàn máy chủ Web GUI ngầm?')) return;
                    try {
                        await fetch('/api/shutdown', { method: 'POST' });
                        this.showToast('🔴 Máy chủ đã được tắt thành công. Bạn có thể đóng trình duyệt này.', 'success');
                        setTimeout(() => { window.close(); }, 1000);
                    } catch (e) {
                        alert('Đã gửi lệnh tắt máy chủ.');
                    }
                }
            }
        }).mount('#app')
    </script>
</body>
</html>`
