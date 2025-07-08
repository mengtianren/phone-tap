@echo off
chcp 65001 >nul
setlocal

REM === 配置路径 ===
set "OPENCV_DLL_DIR=C:\opencv\build\install\x64\mingw\bin"
set "DIST_DIR=dist"
set "EXE_NAME=main.exe"

echo.
echo [步骤1] 清理并创建 dist 目录...
if exist "%DIST_DIR%" (
    echo 正在删除旧的 dist 目录...
    rmdir /s /q "%DIST_DIR%"
)
mkdir "%DIST_DIR%"
mkdir "%DIST_DIR%\images"
mkdir "%DIST_DIR%\tools"
echo ✔️ 创建完成

echo.
echo [步骤2] 编译 Go 程序...
go build -o "%DIST_DIR%\%EXE_NAME%" main.go
if errorlevel 1 (
    echo ❌ 编译失败，请检查代码
    pause
    exit /b 1
)
echo ✔️ 编译成功

echo.
echo [步骤3] 拷贝 OpenCV DLL...
if not exist "%OPENCV_DLL_DIR%" (
    echo ❌ 未找到 DLL 路径: %OPENCV_DLL_DIR%
    pause
    exit /b 1
)
copy /y "%OPENCV_DLL_DIR%\*.dll" "%DIST_DIR%\" >nul
echo ✔️ DLL 拷贝完成

echo.
echo [步骤4] 拷贝资源文件...
if exist "images\close.png" (
    copy /y "images\close.png" "%DIST_DIR%\images\" >nul
) else (
    echo ⚠️ 缺少 close.png
)
if exist "images\success.png" (
    copy /y "images\success.png" "%DIST_DIR%\images\" >nul
) else (
    echo ⚠️ 缺少 success.png
)

if exist "tools\AdbWinApi.dll" (
    copy /y "tools\AdbWinApi.dll" "%DIST_DIR%\tools\" >nul
) else (
    echo ⚠️ 缺少 AdbWinApi.dll
)
if exist "tools\AdbWinUsbApi.dll" (
    copy /y "tools\AdbWinUsbApi.dll" "%DIST_DIR%\tools\" >nul
) else (
    echo ⚠️ 缺少 AdbWinUsbApi.dll
)

if exist "tools\adb.exe" (
    copy /y "tools\adb.exe" "%DIST_DIR%\tools\" >nul
    echo ✔️ adb.exe 拷贝完成
) else (
    echo ⚠️ 未找到 tools\adb.exe，跳过
)

echo.
echo ✅ 打包完成！可执行文件在：%DIST_DIR%\%EXE_NAME%
pause
endlocal
