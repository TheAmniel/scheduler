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

	"github.com/theamniel/scheduler/types"
	"gorm.io/gorm"
)

type IPC struct {
	sync.RWMutex
	Send    chan *Message
	Message chan *Message
	Done    chan bool
	db      *gorm.DB
}

func New() *IPC {
	return &IPC{
		Send:    make(chan *Message),
		Message: make(chan *Message),
		Done:    make(chan bool),
	}
}

func (ipc *IPC) GetDatabase() *gorm.DB {
	ipc.RLock()
	defer ipc.RUnlock()
	return ipc.db
}

func (ipc *IPC) SetDatabase(db *gorm.DB) {
	ipc.Lock()
	ipc.db = db
	ipc.Unlock()
}

func (ipc *IPC) HasDatabase() bool {
	ipc.RLock()
	defer ipc.RUnlock()
	return ipc.db != nil
}

func (ipc *IPC) ToSchedule(payload any) *types.Schedule {
	item := payload.(map[string]any)
	s := &types.Schedule{
		ID:        item["id"].(string),
		CreatedAt: int64(item["created_at"].(float64)),
		ExpiresAt: int64(item["expires_at"].(float64)),
	}
	if item["content"] != nil {
		s.Content = item["content"].(string)
	}
	return s
}

func (ipc *IPC) Schedule(s *types.Schedule) {
	if ipc.HasDatabase() {
		db := ipc.GetDatabase()
		if err := db.Where("id = ?", s.ID).Select("id").First(&types.Schedule{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err = db.Create(s).Error; err != nil {
					log.Printf("Unable to create schedule: %v\n", err)
					return
				}
			} else {
				log.Printf("Unable to find or create schedule: %v\n", err)
				return
			}
		}
		go ipc.WatchSchedule(s)
	}
}

func (ipc *IPC) WatchSchedule(s *types.Schedule) {
	select {
	case <-time.After(time.Duration(s.ExpiresAt-s.CreatedAt) * time.Millisecond):
		db := ipc.GetDatabase()
		if rs := db.Select("id").Where("id = ?", s.ID).Find(&types.Schedule{}).RowsAffected; rs > 0 {
			if err := db.Unscoped().Delete(&s).Error; err != nil {
				log.Printf("Unable to delete schedule from db: %v\n", err)
			}
			ipc.Send <- &Message{"schedule:done", s}
		}
	}
}

func (ipc *IPC) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	go ipc.reader(ctx)
	go ipc.writer(ctx)
	go ipc.watch(ctx)

	var pendings []*types.Schedule
	db := ipc.GetDatabase()
	if err := db.Find(&pendings).Error; err != nil {
		log.Println(err)
	} else {
		for _, sch := range pendings {
			ipc.Schedule(sch)
		}
	}

	ipc.Send <- &Message{Event: "schedule:ready"}
	<-ipc.Done
	cancel()
}

func (ipc *IPC) Close() {
	close(ipc.Done)
	os.Exit(0)
}

func (ipc *IPC) watch(c context.Context) {
	db := ipc.GetDatabase()
	for {
		select {
		case data, ok := <-ipc.Message:
			if ok {
				switch data.Event {
				case "schedule:add":
					{
						ipc.Schedule(ipc.ToSchedule(data.Data))
						ipc.Send <- &Message{"schedule:added", true}
					}
				case "schedule:exists":
					{
						rs := db.Select("id").Where("id = ?", data.Data.(string)).Find(&types.Schedule{}).RowsAffected
						ipc.Send <- &Message{"schedule:exists", rs > 0}
					}
				case "schedule:delete":
					{
						err := db.Unscoped().Where("id = ?", data.Data.(string)).Delete(&types.Schedule{}).Error
						if err != nil {
							log.Printf("Unable to delete schedule: %v\n", err)
							continue
						}
						ipc.Send <- &Message{"schedule:delete", true}
					}
				case "exit":
					ipc.Close()
					return
				default:
					log.Println("Unknown event")
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
				ipc.RLock()
				data, err := MarshalJSON(msg)
				ipc.RLock()
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
			ipc.RLock()
			data, err := input.ReadBytes('\n')
			ipc.RUnlock()
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
