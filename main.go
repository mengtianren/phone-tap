package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"phone/config"
	"strings"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

const closeImg = "./images/close.png"
const closeLiveImg = "./images/liveClose.png"
const successImg = "./images/success.png"
const saveFold = "./images"

func init() {
	// 自动设置 PATH 环境变量，确保 DLL 能加载
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("❌ 获取可执行文件路径失败:", err)
		return
	}
	exeDir := filepath.Dir(exePath)
	os.Setenv("PATH", os.Getenv("PATH")+";"+exeDir)
	rand.Seed(time.Now().UnixNano())

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

	// 如果是识别 closeImg 此处可修改为所有对比都使用白色对比，则启用颜色过滤
	if templatePath == closeImg {
		// 设定HSV白色范围（饱和度小、亮度高）
		lower := gocv.NewScalar(0, 0, 200, 0)    // S=0, V=200
		upper := gocv.NewScalar(180, 40, 255, 0) // S=40, V=255

		filterHSV := func(img gocv.Mat) gocv.Mat {
			hsv := gocv.NewMat()
			gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)
			mask := gocv.NewMat()
			gocv.InRangeWithScalar(hsv, lower, upper, &mask)
			hsv.Close()
			return mask
		}

		// 过滤截图和模板图中的白色区域
		filteredSrc := filterHSV(src)
		defer filteredSrc.Close()
		filteredTpl := filterHSV(tpl)
		defer filteredTpl.Close()

		// 可选调试保存：
		// gocv.IMWrite("./filtered_screen.png", filteredSrc)
		// gocv.IMWrite("./filtered_template.png", filteredTpl)

		resultCols := filteredSrc.Cols() - filteredTpl.Cols() + 1
		resultRows := filteredSrc.Rows() - filteredTpl.Rows() + 1
		results := gocv.NewMatWithSize(resultRows, resultCols, gocv.MatTypeCV32F)
		defer results.Close()

		gocv.MatchTemplate(filteredSrc, filteredTpl, &result, gocv.TmCcoeffNormed, gocv.NewMat())
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
		if maxVal < 0.7 {
			return 0, 0, fmt.Errorf("未找到匹配区域，置信度不足: %.2f", maxVal)
		}

		return maxLoc.X, maxLoc.Y, nil
	}

	// 非 closeImg 的默认流程（原图直接匹配）

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
	screenFile := fmt.Sprintf("%s/phone_%d_1.png", saveFold, i)
	err := takeScreenshot(dev, screenFile)
	if err != nil {
		fmt.Println("❌ 截图失败:", err)

		return
	}

	x, y, err := matchImage(screenFile, closeImg)
	if err == nil {

		offsetX := x + config.Cfg.CloseOffsetX + rand.Intn(21) - 10 // [-10, +10]
		offsetY := y + config.Cfg.CloseOffsetY + rand.Intn(21) - 10 // [-10, +10]

		fmt.Printf("✅ 设备 %s ,下标 %d 找到【取消按钮】匹配: 点击 (%d,%d)\n", dev, i, offsetX, offsetY)
		tap(dev, offsetX, offsetY)
	} else {
		fmt.Printf("❌ 设备 %s 下标 %d 未匹配【取消按钮】图块: %v\n", dev, i, err)
		if config.Cfg.LiveCloseStart {
			fmt.Println("⚠️ 兼容直播间")
			x3, y3, err3 := matchImage(screenFile, closeLiveImg)
			if err3 == nil {
				offsetX := x3 + config.Cfg.LiveCloseOffsetX + rand.Intn(21) - 10 // [-10, +10]
				offsetY := y3 + config.Cfg.LiveCloseOffsetY + rand.Intn(21) - 10 // [-10, +10]

				time.Sleep(time.Duration(config.Cfg.LiveCloseTime) * time.Second)
				delay := rand.Float64()*float64(5) + 1
				time.Sleep(time.Duration(delay * float64(time.Second)))
				fmt.Printf("⚠️ 兼容直播间 设备 %s ,下标 %d 找到【取消按钮】匹配: 点击 (%d,%d)\n", dev, i, offsetX, offsetY)
				tap(dev, offsetX, offsetY)
			} else {
				fmt.Printf("❌ 设备 %s 下标 %d 进程非直播间\n", dev, i)
			}

		}

	}
	delay := rand.Float64()*float64(config.Cfg.AwaitTime) + 1
	fmt.Printf("等待 %.2f 秒\n", delay)
	time.Sleep(time.Duration(delay * float64(time.Second)))
	fmt.Printf("等待 %.2f 秒完成\n", delay)

	screenFile2 := fmt.Sprintf("%s/phone_%d_2.png", saveFold, i)
	err2 := takeScreenshot(dev, screenFile2)
	if err2 != nil {
		fmt.Println("❌ 截图失败222:", err2)
		return
	}

	x1, y1, err1 := matchImage(screenFile2, successImg)
	if err1 == nil {
		offsetX := x1 + config.Cfg.SuccessOffsetX + rand.Intn(21) - 10 // [-10, +10]
		offsetY := y1 + config.Cfg.SuccessOffsetY + rand.Intn(21) - 10 // [-10, +10]
		fmt.Printf("✅ 设备 %s 下标 %d 找到【完成按钮】匹配: 点击 (%d,%d)\n ", dev, i, offsetX, offsetY)
		tap(dev, offsetX, offsetY)
	} else {
		fmt.Printf("❌ 设备 %s 下标 %d 未匹配【完成按钮】图块: %v\n", dev, i, err1)
	}
}

func filterByColor(img gocv.Mat, lower, upper gocv.Scalar) gocv.Mat {
	hsv := gocv.NewMat()
	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)
	mask := gocv.NewMat()
	gocv.InRangeWithScalar(hsv, lower, upper, &mask)
	hsv.Close()
	return mask
}

func main() {
	index := 1
	fmt.Println(`
	1. 请确保已打开USB调试
	2. 请确保已打开USB调试（安全设置）
	3. 请确保手机已连接PC
	4. 第一个区域文件名为close.png
	5. 第二个区域文件名为success.png
	6. 请确保文件与可执行文件在同一目录下
	7. 请确保文件路径中不包含中文
	`)

	// var name string
	// fmt.Print("请输入你要进行的步骤1.养号，2.操作：")
	// fmt.Scanln(&name)
	// fmt.Println("你好，", name)
	// return

	devices, err := getDevices()
	if err != nil || len(devices) == 0 {
		fmt.Println("❌ 未发现设备:", err)
		return
	}
	// 使用go协程 实现步骤一致
	for {
		fmt.Printf("--------第 %d 轮开始------\n", index)
		var wg sync.WaitGroup
		for i, dev := range devices {
			wg.Add(1)
			go start(i, dev, &wg)
		}
		wg.Wait()

		fmt.Printf("--------第 %d 轮结束，等待 %d 秒后再次启动------\n", index, config.Cfg.EndTime)
		time.Sleep(time.Duration(config.Cfg.EndTime) * time.Second)
		index++
	}
}
