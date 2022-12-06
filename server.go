package main

import (
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

var WebsocketMessageType = map[int]string{
	websocket.BinaryMessage: "binary",
	websocket.TextMessage:   "text",
	websocket.CloseMessage:  "close",
	websocket.PingMessage:   "ping",
	websocket.PongMessage:   "pong",
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client connected")

	cmd := exec.Command("/bin/bash")
	tty, err := pty.Start(cmd)
	if err != nil {
		log.Println(err)
		return
	}

	var waiter sync.WaitGroup
	waiter.Add(1)

	// this is a keep-alive loop that ensures connection does not hang-up itself
	lastPongTime := time.Now()
	connection.SetPongHandler(func(msg string) error {
		lastPongTime = time.Now()
		return nil
	})
	go func() {
		for {
			if err := connection.WriteMessage(websocket.PingMessage, []byte("keepalive")); err != nil {
				return
			}
			keepalivePingTimeout := 5 * time.Second
			time.Sleep(keepalivePingTimeout / 2)
			if time.Since(lastPongTime) > keepalivePingTimeout {
				log.Printf("failed to get response from ping, triggering disconnect now...")
				waiter.Done()
				return
			}
			log.Print("received response from ping successfully")
		}
	}()

	// tty -> xterm
	go func() {
		for {
			buffer := make([]byte, 1024)
			readLength, err := tty.Read(buffer)
			if err != nil {
				log.Println("error reading from tty:", err)
				waiter.Done()
				return
			}
			if err := connection.WriteMessage(websocket.BinaryMessage, buffer[:readLength]); err != nil {
				log.Println("error writing to websocket:", err)
				waiter.Done()
				return
			}
			log.Printf("written %d bytes to websocket:", readLength)
		}
	}()

	// xterm -> tty
	go func() {
		for {
			_, buffer, err := connection.ReadMessage()
			if err != nil {
				return
			}

			bytesWritten, err := tty.Write(buffer)
			if err != nil {
				log.Println("error writing to tty:", err)
				return
			}

			log.Printf("wrote %d bytes to tty", bytesWritten)
		}
	}()

	waiter.Wait()
	log.Println("closing websocket connection")
}

func NewServer() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/xterm", handleWebsocket)
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlers.CORS()(router),
	}

	return server
}
