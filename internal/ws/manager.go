package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

// ConnectionManager 线程安全的连接管理器
type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[uint]*Client
}

// NewConnectionManager 初始化
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint]*Client),
	}
}

// Store 用户上线
func (cm *ConnectionManager) Store(userID uint, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[userID] = &Client{conn: conn}
}

// Load 获取连接
func (cm *ConnectionManager) Load(userID uint) (*Client, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	client, ok := cm.connections[userID]
	return client, ok
}

// Delete 用户下线
func (cm *ConnectionManager) Delete(userID uint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, userID)
}

// SafeWriteJSON 线程安全写入
func (cm *ConnectionManager) SafeWriteJSON(userID uint, data interface{}) error {
	client, ok := cm.Load(userID)
	if !ok {
		return nil // 用户不在线，不是一个错误
	}

	client.writeMu.Lock()
	defer client.writeMu.Unlock()
	err := client.conn.WriteJSON(data)

	if err != nil {
		_ = client.conn.Close()
		cm.Delete(userID)
		return err
	}

	return nil

}

// Manager 全局单例
var Manager = NewConnectionManager()
