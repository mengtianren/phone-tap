@echo off
chcp 65001 >nul
title ADB 无线连接工具 - 华为 P9 Plus
echo ==========================================
echo        一键 ADB 无线连接工具
echo        适用于：华为 P9 Plus
echo ==========================================
echo.

:: 输入手机 IP 地址
set /p ip=请输入手机的 IP 地址（如 192.168.1.88）:

:: 开始连接
echo 正在连接到 %ip%:5555 ...
adb connect %ip%:5555

:: 显示当前设备
echo.
adb devices

echo.
echo 如果列表中显示 device，说明连接成功！
pause
exit
