# phone-tap



## 项目介绍
- 手机自动化相关，使用adb+go+gocv+opencv 实现手机自动化相关操作,项目中gocv依赖opencv，opencv具体安装方式详见gocv官网

## 技术选型

- go
- adb
- gocv
- opencv

## 项目运行

1. 克隆项目到本地
2. 进入项目目录，执行`go mod tidy`安装依赖
3. 执行`go run main.go`启动项目 本地开发
4. 双击`demo.bat`项目进行打包 会从C:\opencv\build\install\x64\mingw\bin文件夹复制对应的opencv依赖

## 项目结构

```
├── dist // 打包的文件夹
├── images // 根据相关图片进行的点击判断
├── tools // adb相关的调试
├── demo.bat // 打包的文件
├── go.mod // 项目的依赖
├── go.sum // 依赖的校验值
├── main.go  // 项目主入口
├── README.md // 项目介绍