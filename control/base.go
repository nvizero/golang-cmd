package control

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var statusChan = make(chan string, 20)
var params []string
var show_cate string

type MyForm struct {
	Host     string `form:"host"`
	Date     string `form:"date"`
	Category string `form:"category"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var set_type = []string{"get_date", "stop_server", "start_server", "status"}

const (
	GET_DATE     = 0
	STOP_SERVER  = 1
	START_SERVER = 2
	STATUS       = 3
)

func SetServer(ctx context.Context, host, cate string, wg *sync.WaitGroup) {
	defer wg.Done()
	var message string
	if cate == set_type[START_SERVER] {
		message = fmt.Sprintf("<b style='color:red;'>執行%s .... 會比較久 大概三分鐘</b>\n", cate)
	} else {
		message = fmt.Sprintf("<b style='color:blue;'>執行%s<b>\n", cate)
	}
	statusChan <- message

	cmd := "ssh"
	sshhost := fmt.Sprintf("root@%s", host)

	switch cate {
	case set_type[GET_DATE]:
		params = []string{sshhost, "-tt", "sh /root/get_date.sh -u"}
	case set_type[STOP_SERVER]:
		params = []string{sshhost, "-tt", "sh", "/root/stop.sh -u"}
	case set_type[START_SERVER]:
		params = []string{sshhost, "-tt", "sh", "/root/start.sh -u"}
	case set_type[STATUS]:
		params = []string{sshhost, "-tt", "sh", "/root/status.sh -u"}
	default:
		params = []string{sshhost, "-tt", "sh", "/root/get_date.sh -u"}
	}
	cmdStr := fmt.Sprintf("%s %s", cmd, strings.Join(params, " "))
	ExecCmd(cmdStr)
}

func SetDate(sdate, host string, wg *sync.WaitGroup) {
	defer wg.Done()
	message := fmt.Sprintf("<b style='color:blue;'>執行set date<b>\n")
	statusChan <- message
	cmd := "ssh"
	sshhost := fmt.Sprintf("root@%s", host)
	params := []string{sshhost, "sh", "/root/set_date.sh", fmt.Sprintf("'%s'", sdate)}

	cmdStr := fmt.Sprintf("%s %s", cmd, strings.Join(params, " "))
	ExecCmd(cmdStr)
}

func ExecCmd(cmdStr string) {
	command := exec.Command("/bin/sh", "-c", cmdStr)
	stdout, err := command.StdoutPipe()
	ChkErr(err)
	stderr, err := command.StderrPipe()
	ChkErr(err)

	if err := command.Start(); err != nil {
		log.Fatal("--", err)
	}

	go func() {
		defer command.Wait()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			statusChan <- scanner.Text()
		}
	}()

	go func() {
		defer command.Wait()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			statusChan <- scanner.Text()
		}
	}()

	// 等待命令执行完成
	if err := command.Wait(); err != nil {
		fmt.Println("命令执行失败:", err)
	}

	//statusChan <- string(out)
}

func CHttp() {
	r := gin.Default()
	r.Static("/static", "templates/static")
	r.LoadHTMLGlob("templates/*.html") // Load HTML templates from the "templates" directory

	r.GET("/", func(c *gin.Context) {
		data := gin.H{
			"Title": "慶餘年更新日期",
		}
		c.HTML(http.StatusOK, "index.html", data) // Render the HTML template
	})

	// Handle the form submission
	r.POST("/submit", func(c *gin.Context) {
		var wg sync.WaitGroup
		var form MyForm
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		server_date := form.Date
		server_host := form.Host
		category := form.Category
		done := make(chan struct{})

		if category == "set_date" {

			go func() {
				wg.Add(1)
				SetServer(ctx, server_host, set_type[STOP_SERVER], &wg)
				done <- struct{}{} // 发送完成信号
			}()
			<-done // 等待第一个 goroutine 完成

			go func() {
				wg.Add(1)
				SetDate(server_date, server_host, &wg)
				done <- struct{}{} // 发送完成信号
			}()
			<-done // 等待第二个 goroutine 完成

			go func() {
				wg.Add(1)
				SetServer(ctx, server_host, set_type[START_SERVER], &wg)
				done <- struct{}{} // 发送完成信号
			}()
			<-done // 等待第一个 goroutine 完成
		} else if category == "status" {
			go func() {
				wg.Add(1)
				SetServer(ctx, server_host, set_type[STATUS], &wg)
				done <- struct{}{} // 发送完成信号
			}()
			<-done // 等待第一个 goroutine 完成

		} else {
			go func() {
				wg.Add(1)
				SetServer(ctx, server_host, set_type[GET_DATE], &wg)
				done <- struct{}{} // 发送完成信号
			}()
			<-done // 等待第一个 goroutine 完成
		}
		// 等待所有 goroutine 完成
		wg.Wait()
		c.JSON(http.StatusOK, gin.H{
			"Host": form.Host,
			"Date": form.Date,
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		// 處理 WebSocket 連接
		handleWebSocket(c.Writer, c.Request)
	})

	r.Run(":8080")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	go func() {
		for {
			select {
			case status := <-statusChan:
				// 在这里将状态信息发送到WebSocket连接
				if err := conn.WriteMessage(websocket.TextMessage, []byte(status)); err != nil {
					return
				}
			}
		}
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println("hi", messageType)
		fmt.Println("p=", p)
		// 这里可以处理来自WebSocket客户端的消息
	}
}

func ChkErr(err error) {
	if err != nil {
		log.Fatal("------error-------   \n", err)
	}
}
