@echo off
setlocal

:: 设置Go程序的源代码路径
set SRC_DIR=main.go

:: 设置输出目录
set OUT_DIR=output

:: 设置程序名称
set APP_NAME=updav

:: 创建输出目录
if not exist %OUT_DIR% mkdir %OUT_DIR%

:: 编译Windows x86 64可执行文件
set GOOS=windows
set GOARCH=amd64
go build -o %OUT_DIR%\%APP_NAME%_windows_amd64.exe %SRC_DIR%

:: 编译Windows x86 32位可执行文件
set GOOS=windows
set GOARCH=386
go build -o %OUT_DIR%\%APP_NAME%_windows_386.exe %SRC_DIR%

:: 编译Windows ARM 64位可执行文件
set GOOS=windows
set GOARCH=arm64
go build -o %OUT_DIR%\%APP_NAME%_windows_arm64.exe %SRC_DIR%

:: 编译Linux x64位可执行文件
set GOOS=linux
set GOARCH=amd64
go build -o %OUT_DIR%\%APP_NAME%_linux_amd64 %SRC_DIR%

:: 编译Linux 32位可执行文件
set GOOS=linux
set GOARCH=386
go build -o %OUT_DIR%\%APP_NAME%_linux_386 %SRC_DIR%

:: 编译Linux ARM 32位可执行文件
set GOOS=linux
set GOARCH=arm
go build -o %OUT_DIR%\%APP_NAME%_linux_arm %SRC_DIR%

:: 编译Linux ARM 64位可执行文件
set GOOS=linux
set GOARCH=arm64
go build -o %OUT_DIR%\%APP_NAME%_linux_arm64 %SRC_DIR%

:: 编译Linux MIPS 32位可执行文件
set GOOS=linux
set GOARCH=mips
go build -o %OUT_DIR%\%APP_NAME%_linux_mips %SRC_DIR%

:: 编译Linux MIPS 64位可执行文件
set GOOS=linux
set GOARCH=mips64
go build -o %OUT_DIR%\%APP_NAME%_linux_mips64 %SRC_DIR%

:: 编译Linux MIPS 64位LE可执行文件
set GOOS=linux
set GOARCH=mips64le
go build -o %OUT_DIR%\%APP_NAME%_linux_mips64le %SRC_DIR%

:: 编译Linux MIPS 32位LE可执行文件
set GOOS=linux
set GOARCH=mipsle
go build -o %OUT_DIR%\%APP_NAME%_linux_mipsle %SRC_DIR%

:: 编译Linux RISC-V 64位可执行文件
set GOOS=linux
set GOARCH=riscv64
go build -o %OUT_DIR%\%APP_NAME%_linux_riscv64 %SRC_DIR%

:: 编译Linux LOONG64 64位可执行文件
set GOOS=linux
set GOARCH=loong64
go build -o %OUT_DIR%\%APP_NAME%_linux_loong64 %SRC_DIR%

:: 编译macOS ARM 64位可执行文件
set GOOS=darwin
set GOARCH=arm64
go build -o %OUT_DIR%\%APP_NAME%_darwin_arm64 %SRC_DIR%

:: 编译macOS x86 64位可执行文件
set GOOS=darwin
set GOARCH=amd64
go build -o %OUT_DIR%\%APP_NAME%_darwin_amd64 %SRC_DIR%

:: 编译完成,恢复环境变量
set GOOS=windows
set GOARCH=amd64
echo Compilation completed.
endlocal