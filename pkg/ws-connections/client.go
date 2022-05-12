package wsconnections

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WebsocketClient struct {
	// Поточне з'єднання
	conn *websocket.Conn

	// Налаштування підключення
	ConnectionString string

	// Канал відправки даних
	send chan []byte

	// Канал отримання даних
	recive chan []byte

	// Статус відкритий/закритий канал
	close chan bool
}

func New() *WebsocketClient {
	host := os.Getenv("SERVER_SOCKET_HOST")
	if host == "" {
		host = "195.95.233.38:8182"
	}
	deviceID := os.Getenv("DEVICE_ID")
	if deviceID == "" {
		deviceID = "2"
	}

	client := new(WebsocketClient)

	u := url.URL{Scheme: "ws", Host: host, Path: "device/" + deviceID}
	client.ConnectionString = u.String()

	client.send = make(chan []byte)
	client.recive = make(chan []byte)

	return client
}

// Підключаємось до сокета
func (c *WebsocketClient) Connect() error {
	conn, err := c.connect()

	if err == nil {
		c.conn = conn

		go c.watchSend()

		go c.watchRecive()

		go c.ping()

		return nil
	}
	return err
}

func (c *WebsocketClient) connect() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.ConnectionString, nil)
	return conn, err
}

// Пишемо в канал для відправки по сокету
func (c *WebsocketClient) Write(message interface{}) {
	if c.conn == nil {
		// Якщо коннект не встановлений щось будем робити
		return
	}

	if data, err := json.Marshal(message); err == nil {
		c.send <- data
	} else {
		log.Println("Error send message via socket", err)
		// Тут ще буде додатковий обработчкик
	}
}

// Відправляємо повідомлення з каналу на сервер
func (c *WebsocketClient) watchSend() {
	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error send message", err)
				// Обработчик помилки
			}
		}
	}
}

// Отримуємо повідомлення від сервера
func (c *WebsocketClient) watchRecive() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("Error recive message", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				conn, err := c.connect()
				if err == nil {
					c.conn = conn
					log.Println("Reconnected")
				} else {
					c.conn = nil
				}
			}
			// Обработчик помилки
			continue
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		if string(message) != "pong" {
			c.recive <- message
		}
	}
}

// Повертаємо канал через який нам приходять повідомлення
func (c *WebsocketClient) Recive() chan []byte {
	return c.recive
}

// Ping - pong
func (c *WebsocketClient) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2)); err != nil {
				log.Println("Error pong", err)
			}
		}
	}
}
