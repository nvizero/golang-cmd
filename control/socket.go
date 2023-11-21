package control

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	connections map[*websocket.Conn]chan string
	lock        sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[*websocket.Conn]chan string),
	}
}

func (manager *ConnectionManager) Add(conn *websocket.Conn) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.connections[conn] = make(chan string, 20) // 为每个连接创建一个消息通道
}

func (manager *ConnectionManager) Remove(conn *websocket.Conn) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	close(manager.connections[conn]) // 关闭消息通道
	delete(manager.connections, conn)
}

func (manager *ConnectionManager) SendToConnection(conn *websocket.Conn, message string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if channel, ok := manager.connections[conn]; ok {
		channel <- message
	}
}
func (manager *ConnectionManager) SendToAll(message string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	for conn := range manager.connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			// 可以選擇在這裡移除連接
			fmt.Printf("Error sending message: %v", err)
		}
	}
}
