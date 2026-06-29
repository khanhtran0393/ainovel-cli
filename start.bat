@echo off
chcp 65001 >nul
title Khoi dong ainovel-cli

echo ==========================================
echo    Khoi dong ainovel-cli (Ban tieng Viet)
echo ==========================================
echo.

:: 1. Kiem tra thu muc can thiet theo README
if not exist "config" mkdir config
if not exist "workspace" mkdir workspace

:: 2. Kiem tra file config.json de fix loi setup: could not open a new TTY tren Docker
if not exist "config\config.json" (
    echo [*] Khong tim thay config.json. Dang tao cau hinh mac dinh ^(Ollama^)...
    (
        echo {
        echo   "provider": "ollama",
        echo   "model": "gemma4:12b",
        echo   "providers": {
        echo     "ollama": {
        echo       "base_url": "http://host.docker.internal:11434/v1",
        echo       "models": ["gemma4:12b", "qwen3.5:27b"]
        echo     }
        echo   }
        echo }
    ) > config\config.json
    echo [*] Da tao xong config\config.json
    echo.
)

:menu
cls
echo ==========================================
echo               MENU CHINH
echo ==========================================
echo 1. Chay app bang Docker (Khuyen nghi, khong can cai dat Go)
echo 2. Chay app bang Go (Build tu source code, can Go ^>= 1.21)
echo.
echo --- TIEN ICH ^& SUA LOI (Theo huong dan README) ---
echo 3. Tai model AI cho Ollama (gemma4:12b) - Dung cho cau hinh mac dinh
echo 4. Xoa du lieu truyen dang viet do de bat dau cuon moi (Sua loi App loop)
echo 0. Thoat
echo ==========================================
set /p choice="Nhap lua chon cua ban: "

if "%choice%"=="1" goto rundocker
if "%choice%"=="2" goto rungo
if "%choice%"=="3" goto pullmodel
if "%choice%"=="4" goto clean
if "%choice%"=="0" exit /b

echo [!] Lua chon khong hop le.
timeout /t 2 >nul
goto menu

:rundocker
echo.
echo [*] Kiem tra moi truong Docker...
docker --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [!] Khong tim thay Docker tren may. Yeu cau: Docker Desktop ^>= 24.
    echo [*] Se tu dong mo trang tai Docker Desktop trong giay lat...
    timeout /t 3 >nul
    start https://www.docker.com/products/docker-desktop/
    pause
    goto menu
)

:check_docker_running
docker info >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [!] Docker chua duoc bat hoac dang khoi dong.
    echo [*] Vui long mo app "Docker Desktop", doi bieu tuong ca voi chuyen xanh roi nhan phim bat ky de tiep tuc...
    pause >nul
    goto check_docker_running
)

echo [*] Dang build Docker image (ainovel-cli-vi)...
docker build -t ainovel-cli-vi .
if %ERRORLEVEL% neq 0 (
    echo [!] Build Docker that bai.
    pause
    goto menu
)
echo.
echo [*] Dang khoi dong ung dung... (Nhan Ctrl+C de thoat)
set ENV_ARGS=
if exist ".env" (
    set ENV_ARGS=--env-file .env
)
docker run --rm -it %ENV_ARGS% -v "%CD%\config:/root/.ainovel" -v "%CD%\workspace:/workspace" -e TERM=xterm-256color ainovel-cli-vi
pause
goto menu

:rungo
echo.
echo [*] Kiem tra moi truong Go...
go version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [!] Khong tim thay Go tren may. Yeu cau: Go ^>= 1.21.
    echo [*] Se tu dong mo trang tai Go lang trong giay lat...
    timeout /t 3 >nul
    start https://go.dev/dl/
    pause
    goto menu
)

echo [*] Kiem tra va tai cac thu vien...
go mod tidy
if %ERRORLEVEL% neq 0 (
    echo [!] Loi khi tai thu vien.
    pause
    goto menu
)

echo [*] Dang build ung dung bang Go...
go build -o ainovel-cli.exe ./cmd/ainovel-cli/
if %ERRORLEVEL% neq 0 (
    echo [!] Build Go that bai. Hay dam bao ban dang dung ban Go ^>= 1.21.
    pause
    goto menu
)
echo.
echo [*] Dang khoi dong ung dung... (Nhan Ctrl+C de thoat)
ainovel-cli.exe
pause
goto menu

:pullmodel
echo.
echo [*] Kiem tra moi truong Ollama...
ollama --version >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [!] Khong tim thay Ollama.
    echo [*] Se tu dong mo trang tai Ollama trong giay lat...
    timeout /t 3 >nul
    start https://ollama.com/
    pause
    goto menu
)
echo [*] Dang tai model gemma4:12b... (Qua trinh nay phu thuoc vao toc do mang cua ban)
ollama pull gemma4:12b
echo [*] Hoan tat!
pause
goto menu

:clean
echo.
echo [!] CANH BAO: Hanh dong nay se xoa toan bo tien do viet truyen hien tai trong thu muc workspace/output.
echo Truong hop su dung:
echo - Ban muon bat dau viet mot tieu thuyet hoan toan moi.
echo - App bi loi lap vo han ma khong dung lai duoc (Theo Troubleshooting trong README).
echo.
set /p confirm="Ban co chac chan muon xoa khong? (y/n): "
if /i "%confirm%"=="y" (
    rmdir /s /q "workspace\output" 2>nul
    echo [*] Da xoa workspace/output thanh cong.
) else (
    echo [*] Da huy thao tac xoa.
)
pause
goto menu
