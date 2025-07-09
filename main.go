package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

const closeImg = "./images/close.png"
const successImg = "./images/success.png"

func init() {
	// 自动设置 PATH 环境变量，确保 DLL 能加载
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径失败:", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	os.Setenv("PATH", os.Getenv("PATH")+";"+exeDir)

}

func getDevices() ([]string, error) {
	cmd := exec.Command("./tools/adb.exe", "devices")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var devices []string
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "List of devices") {
			continue
		}
		if strings.HasSuffix(line, "\tdevice") {
			fields := strings.Fields(line)
			devices = append(devices, fields[0])
		}
	}
	return devices, nil
}

func takeScreenshot(deviceID, savePath string) error {
	cmd := exec.Command("./tools/adb.exe", "-s", deviceID, "exec-out", "screencap", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.WriteFile(savePath, out.Bytes(), 0644)
}

var adbMu sync.Mutex

func tap(deviceID string, x, y int) error {
	// 防止adb不可靠导致点击不一致情况
	adbMu.Lock()
	defer adbMu.Unlock()
	cmd := exec.Command("./tools/adb.exe", "-s", deviceID, "shell", "input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))
	return cmd.Run()
}

func matchImage(screenPath, templatePath string) (int, int, error) {
	src := gocv.IMRead(screenPath, gocv.IMReadColor)
	if src.Empty() {
		return 0, 0, fmt.Errorf("读取截图失败: %s", screenPath)
	}
	defer src.Close()

	tpl := gocv.IMRead(templatePath, gocv.IMReadColor)
	if tpl.Empty() {
		return 0, 0, fmt.Errorf("读取模板失败: %s", templatePath)
	}
	defer tpl.Close()

	resultCols := src.Cols() - tpl.Cols() + 1
	resultRows := src.Rows() - tpl.Rows() + 1
	result := gocv.NewMatWithSize(resultRows, resultCols, gocv.MatTypeCV32F)
	defer result.Close()

	mask := gocv.NewMat() // 空 mask 符合要求
	defer mask.Close()

	// 修复关键点：添加 mask 参数
	gocv.MatchTemplate(src, tpl, &result, gocv.TmCcoeffNormed, mask)
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	if maxVal < 0.8 {
		return 0, 0, errors.New("未找到匹配区域，置信度不足")
	}

	return maxLoc.X, maxLoc.Y, nil
}

func start(i int, dev string, wg *sync.WaitGroup) {
	// 关闭
	defer wg.Done()
	screenFile := fmt.Sprintf("phone_%d_1.png", i)
	err := takeScreenshot(dev, screenFile)
	if err != nil {
		fmt.Println("截图失败:", err)
		return
	}

	x, y, err := matchImage(screenFile, closeImg)
	if err == nil {
		offsetX := x + 400 + rand.Intn(21) - 10 // [-10, +10]
		offsetY := y + rand.Intn(21) - 10       // [-10, +10]

		fmt.Printf("设备 %s ,下标 %d 找到【取消按钮】匹配: 点击 (%d,%d)\n", dev, i, offsetX, offsetY)
		tap(dev, offsetX, offsetY)
	} else {
		fmt.Printf("设备 %s 下标 %d 未匹配【取消按钮】图块: %v\n", dev, i, err)

	}
	fmt.Println("等待3秒")
	time.Sleep(3 * time.Second)
	fmt.Println("等待3秒完成")

	screenFile2 := fmt.Sprintf("phone_%d_2.png", i)
	err2 := takeScreenshot(dev, screenFile2)
	if err2 != nil {
		fmt.Println("截图失败222:", err2)
		return
	}

	x1, y1, err1 := matchImage(screenFile2, successImg)
	if err1 == nil {
		offsetX := x1 + rand.Intn(21) - 10 // [-10, +10]
		offsetY := y1 + rand.Intn(21) - 10 // [-10, +10]
		fmt.Printf("设备 %s 下标 %d 找到【完成按钮】匹配: 点击 (%d,%d)\n ", dev, i, offsetX, offsetY)
		tap(dev, offsetX, offsetY)
	} else {
		fmt.Printf("设备 %s 下标 %d 未匹配【完成按钮】图块: %v\n", dev, i, err1)
	}
}

func main() {
	fmt.Println(`
	1. 请确保已打开USB调试
	2. 请确保已打开USB调试（安全设置）
	3. 请确保手机已连接PC
	4. 第一个区域文件名为close.png
	5. 第二个区域文件名为success.png
	6. 请确保文件与可执行文件在同一目录下
	7. 请确保文件路径中不包含中文
	`)
	devices, err := getDevices()
	if err != nil || len(devices) == 0 {
		fmt.Println("未发现设备:", err)
		return
	}
	// 使用go协程 实现步骤一致
	for {
		var wg sync.WaitGroup
		for i, dev := range devices {
			wg.Add(1)
			go start(i, dev, &wg)
		}
		wg.Wait()
		time.Sleep(10 * time.Second)
	}
}
