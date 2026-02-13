package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	// 允许跨域
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Connect(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, ok := uid.(uint)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("WebSocket Upgrade Failed", zap.Error(err))
		return
	}

	Manager.Store(userID, conn)
	zap.L().Info("User Connected", zap.Uint("uid", userID))

	//  清理
	defer func() {
		Manager.Delete(userID)
		_ = conn.Close()
		zap.L().Info("User Disconnected", zap.Uint("uid", userID))
	}()

	// 保持连接活跃
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
