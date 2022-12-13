package ipc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type IPC struct {
	sync.RWMutex
	Send    chan *Message
	Message chan *Message
	Done    chan bool
}

func New() *IPC {
	return &IPC{
		Send:    make(chan *Message),
		Message: make(chan *Message),
		Done:    make(chan bool),
	}
}

func (ipc *IPC) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	go ipc.reader(ctx)
	go ipc.writer(ctx)
	go ipc.watch(ctx)
	ipc.Send <- &Message{Op: IpcReady}
	<-ipc.Done
	cancel()
}

func (ipc *IPC) Close() {
	close(ipc.Done)
	os.Exit(0)
}

func (ipc *IPC) watch(c context.Context) {
	for {
		select {
		case data, ok := <-ipc.Message:
			if ok {
				if data.Op == IpcDispatch {
					// handle data.T and data.D.....
					ipc.Send <- &Message{IpcDispatch, "TEST_EVENT", nil}
				} else if data.Op == IpcExit {
					ipc.Close()
					return
				} else {
					log.Println("Invalid Opcode.")
					continue
				}
			} else {
				ipc.Close()
				return
			}
		case <-c.Done():
			return
		}
	}
}

func (ipc *IPC) writer(c context.Context) {
	defer close(ipc.Send)
	for {
		select {
		case msg, ok := <-ipc.Send:
			if ok {
				data, err := MarshalJSON(msg)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Print(data + "\\n")
				}
			} else {
				return
			}
		case <-c.Done():
			return
		}
	}
}

func (ipc *IPC) reader(c context.Context) {
	defer close(ipc.Message)
	input := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-time.Tick(10 * time.Millisecond):
			data, err := input.ReadBytes('\n')
			if err != nil {
				log.Println(err)
				continue
			}
			if len(data) > 0 {
				data = bytes.TrimSuffix(data, []byte("\n"))
				if len(data) > 0 {
					var msg Message
					if err := json.Unmarshal(data, &msg); err != nil {
						log.Println(err)
						continue
					}
					ipc.Message <- &msg
				}
			}
		case <-c.Done():
			return
		}
	}
}
