@echo off
:: 设置代码页为 UTF-8
chcp 65001 >nul
setlocal
title SupperSystem Build Tool

:: --- 自动定位路径 ---
cd /d "%~dp0.."
set "PROJECT_ROOT=%cd%"

set "APP_NAME=SupperSystem"
set "BUILD_DIR=dist_release"
set "GO_ENTRY=./cmd/suppersystem/main.go"

echo ========================================
echo Project Root: %PROJECT_ROOT%
echo ========================================

echo.
echo [1/3] Cleaning old build files...
if exist "%BUILD_DIR%" rd /s /q "%BUILD_DIR%"
mkdir "%BUILD_DIR%"

echo.
echo [2/3] Building Frontend (Vue 3)...
cd /d "%PROJECT_ROOT%\web"
if not exist node_modules (
    echo Installing dependencies...
    call npm install
)
call npm run build

echo.
echo [3/3] Compiling Backend (Go)...
cd /d "%PROJECT_ROOT%"
:: 编译入口路径
go build -ldflags="-s -w" -o "%BUILD_DIR%/%APP_NAME%.exe" "./cmd/suppersystem"

if %errorlevel% neq 0 (
    echo.
    echo ERROR: Go compilation failed.
    pause
    exit /b %errorlevel%
)

echo.
echo ========================================
echo BUILD SUCCESS!
echo Output Directory: %PROJECT_ROOT%\%BUILD_DIR%
echo ========================================

:: 拷贝配置文件
if not exist "%BUILD_DIR%\configs" mkdir "%BUILD_DIR%\configs"
copy /y "configs\config.json" "%BUILD_DIR%\configs\config.json" >nul 2>&1

pause