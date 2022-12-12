package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TTYSize struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}

func NewServer(ctx context.Context) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		cmd := exec.Command("bash")
		tty, err := pty.Start(cmd)
		if err != nil {
			runtime.LogWarningf(ctx, "Error starting pty: %s", err.Error())
			return
		}

		waiter := sync.WaitGroup{}
		waiter.Add(1)

		var resizeMsg TTYSize

		// tty -> xterm
		go func() {
			for {
				buffer := make([]byte, 1024)
				readLength, err := tty.Read(buffer)
				if err != nil {
					runtime.LogDebugf(ctx, "Error reading from pty (%s), restarting...", err.Error())
					cmd := exec.Command("bash")
					tty, err = pty.Start(cmd)
					pty.Setsize(tty, &pty.Winsize{
						Rows: resizeMsg.Rows,
						Cols: resizeMsg.Cols,
					})
					runtime.WindowHide(ctx)
					if err != nil {
						runtime.LogWarningf(ctx, "Error restarting pty: %s", err.Error())
					}
					continue
				}
				if err := connection.WriteMessage(websocket.BinaryMessage, buffer[:readLength]); err != nil {
					runtime.LogWarningf(ctx, "Error writing to websocket: %s", err.Error())
				}
			}
		}()

		// xterm -> tty
		go func() {
			for {
				var err error
				messageType, buffer, err := connection.ReadMessage()

				if messageType == websocket.BinaryMessage {
					runtime.LogDebugf(ctx, "Received binary message: %s", string(buffer))
					err = json.Unmarshal(buffer, &resizeMsg)
					if err != nil {
						runtime.LogWarningf(ctx, "Error unmarshalling resize message: %s", err.Error())
						continue
					}

					if err := pty.Setsize(tty, &pty.Winsize{
						Rows: resizeMsg.Rows,
						Cols: resizeMsg.Cols,
					}); err != nil {
						runtime.LogWarningf(ctx, "Error resizing pty: %s", err.Error())
						continue
					}

					continue
				}
				if err != nil {
					runtime.LogWarningf(ctx, "Error reading from websocket: %s", err.Error())
					continue
				}
				if _, err := tty.Write(buffer); err != nil {
					runtime.LogWarningf(ctx, "Error writing to pty: %s", err.Error())
				}
			}
		}()

		waiter.Wait()

	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
