package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"golang.org/x/net/webdav"
)

const port = ":8080"
const wwwroot = "www"
const uploadDir = "temp"

// 嵌入www目录
var (
	//go:embed www/*
	webDir embed.FS
)

// webdav的配置
var (
	addr         = ":8848"
	path         = "/" // WebDAV 服务的根目录
	webDAVServer *http.Server
	once         sync.Once
	mu           sync.Mutex
)

var (
	c   *exec.Cmd
	f   interface{}
	err error
)

func init() {
}

func main() {
	// c = exec.Command("cmd")
	// f, err = conpty.Start(c.Path)
	c = exec.Command("/bin/sh")
	f, err = pty.Start(c)
	if err != nil {
		panic(err)
	}

	m := melody.New()

	go func() {
		for {
			buf := make([]byte, 1024)
			var read int
			// 使用ConPTY的读取方法
			// read, err = f.(*conpty.ConPty).Read(buf)
			read, err = f.(*os.File).Read(buf)
			if err != nil {
				break
			}
			m.Broadcast(buf[:read])
		}
	}()

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		// 使用ConPTY的写入方法
		// f.(*conpty.ConPty).Write(msg)
		// 使用pty的写入方法
		f.(*os.File).Write(msg)
	})
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/upload", PostApiUpload) // 处理上传文件的路由
	r.GET("/webdav", getWebdav)      // 开启webdav服务
	r.GET("/upgrade", getUpgrade)    // 处理升级
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	// r.NoRoute(httpWeb)
	staticFp, _ := fs.Sub(webDir, wwwroot)
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(staticFp))))
	r.Run(port)
}

func httpWeb(c *gin.Context) {
	path := c.Request.URL.Path // 获取请求路径
	if path == "/" {           // 根目录重定向到 index.html
		path = "/index.html"
	}
	filePath := filepath.Join(wwwroot, path) // 构造完整文件路径
	_, err := os.Stat(filePath)              // 检查文件是否存在
	if err == nil {
		http.ServeFile(c.Writer, c.Request, filePath) // 文件存在，则返回文件
	} else {
		c.String(http.StatusNotFound, "404 Not Found") // 文件不存在，返回 404
	}
}

func PostApiUpload(c *gin.Context) {
	// 从表单中获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取文件"})
		return
	}
	defer file.Close()

	// 创建 temp 目录（如果不存在）
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0755)
	}

	// 创建目标文件
	filePath := "temp/" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建文件"})
		return
	}
	defer out.Close()

	// 流式写入文件
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功", "filePath": filePath})
}

func getWebdav(c *gin.Context) {
	// 获取cmd参数
	cmd := c.Query("cmd")
	if cmd == "enable" {
		mu.Lock()
		defer mu.Unlock()
		once.Do(func() { go startWebDAV() })
		c.JSON(http.StatusOK, gin.H{
			"msg": "enabled",
		})
	} else if cmd == "disable" {
		mu.Lock()
		defer mu.Unlock()
		if webDAVServer != nil {
			// 关闭 WebDAV 服务
			err := webDAVServer.Shutdown(context.Background())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "error",
				})
				return
			}
			webDAVServer = nil
			once = sync.Once{} // 重置 sync.Once，以便可以重新启动服务
			fmt.Println("[WebDAV] 服务已关闭")
			c.JSON(http.StatusOK, gin.H{
				"msg": "disabled",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": "error",
			})
		}
	} else if cmd == "status" {
		if webDAVServer != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "enabled",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "disabled",
			})
		}
	}
}

func startWebDAV() {
	fmt.Println("addr=", addr, ", path=", path) // 在控制台输出配置

	// 创建 WebDAV 服务
	webDAVServer = &http.Server{
		Addr: addr,
		Handler: &webdav.Handler{
			FileSystem: webdav.Dir(path),
			LockSystem: webdav.NewMemLS(),
		},
	}

	fmt.Println("[WebDAV] 服务正在启动，监听端口:", addr)
	if err := webDAVServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("无法开启 WebDAV 服务:", err)
	}
}

func getUpgrade(c *gin.Context) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "start", "upgrade.bat")
		err := cmd.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "升级失败",
			})
			return
		}
	} else {
		cmd := exec.Command("sh", "upgrade.sh")
		err := cmd.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "升级失败",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "升级成功",
	})
}
