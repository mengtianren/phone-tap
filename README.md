# phone-tap



## 项目介绍
- 手机自动化相关，使用adb+go+gocv+opencv 实现手机自动化相关操作,项目中  gocv依赖opencv，opencv具体安装方式详见 [gocv官网](https://gocv.io/getting-started/windows/),
- 项目可配合 [Escrcpy](https://github.com/viarotel-org/escrcpy) 项目使用更佳
- 项目基于adb链接 可以使用有线及无线链接， 如本机没有安装adb 则需要切到本项目tools目录内使用adb命令
- 项目中使用的图片均为png格式，图片来源为手机截图，需要注意的是图片的像素比例需要和手机的像素比例保持一致，否则会点击失败
- 项目配置文件为config.yaml 具体可以看备注


## 技术选型

- go
- adb
- gocv
- opencv

## 项目运行

1. 克隆项目到本地
2. 进入项目目录，执行`go mod tidy`安装依赖
3. images文件夹需要存在close.png和success.png 以及liveClose.png文件，用来做匹配点击，可以使用截图工具进行图片切割，直接截图无效 需要保持原机像素比例
4. 执行`go run main.go`启动项目 本地开发
5. 双击`demo.bat`项目进行打包 会从（C:\opencv\build\install\x64\mingw\bin，C:\Program Files\mingw64\bin）文件夹复制对应的依赖，注意：需要mingw和opencv环境

## 项目结构

```
├── dist // 打包的文件夹
├── images // 根据相关图片进行的点击判断
├── tools // adb相关的调试
├── demo.bat // 打包的文件
├── go.mod // 项目的依赖
├── go.sum // 依赖的校验值
├── main.go  // 项目主入口
├── config   // 配置选项
├── config.yaml  // 项目的具体配置
├── README.md // 项目介绍
```

## 声明
- 本项目仅用于学习研究，不用于商业用途
- 本项目不承担因使用本项目而导致的任何损失或损害
- 本项目不对因使用本项目而导致的任何损失或损害承担责任
- 本项目不对因使用本项目的用途而导致的任何损失或损害承担责任
