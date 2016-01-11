package list

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type Timer struct {
	Time     int64  `json:"time"`
	Callback string `json:"callback"`
}

type List struct {
	savefile string
	list     []*Timer
	locker   sync.Mutex
	exitChan chan bool
}

func (l *List) Len() int {
	return len(l.list)
}

func (l *List) Less(i, j int) bool {
	return l.list[i].Time < l.list[j].Time
}

func (l *List) Swap(i, j int) {
	l.list[i], l.list[j] = l.list[j], l.list[i]
}

func NewList(savefile string) *List {
	l := &List{
		exitChan: make(chan bool),
		savefile: savefile,
	}
	file, err := os.Open(savefile)
	if err == nil {
		dec := json.NewDecoder(file)
		var timers []Timer
		err = dec.Decode(&timers)
		if err == nil {
			for _, t := range timers {
				l.list = append(l.list, &t)
			}
		}
		file.Close()
	}
	go l.Run()
	return l
}

func (l *List) Run() {
	for {
		select {
		case now := <-time.After(time.Second):
			l.onTimer(now.Unix())
		case <-l.exitChan:
			break
		}
	}
}

func (l *List) Save() {
	l.locker.Lock()
	defer l.locker.Unlock()
	// save
	fmt.Printf("Save, len: %d\n", len(l.list))
	if file, err := os.Create(l.savefile); err == nil {
		var timers []Timer
		for _, t := range l.list {
			if t != nil {
				timers = append(timers, *t)
			}
		}
		if len(timers) > 0 {
			enc := json.NewEncoder(file)
			enc.Encode(timers)
			file.Sync()
		}
		file.Close()
	}
}

func (l *List) AddTimer(timer *Timer) {
	log.Printf("add %+v\n", timer)
	l.locker.Lock()
	defer l.locker.Unlock()
	l.list = append(l.list, timer)
	sort.Sort(l)
}

func (l *List) Exit() {
	l.Save()
	l.exitChan <- true
}

func (l *List) onTimer(now int64) {
	l.locker.Lock()
	defer l.locker.Unlock()
	var index int
	var item *Timer
	for index, item = range l.list {
		if item.Time <= now {
			go callback(item)
		} else {
			if index == 0 {
				return
			}
			l.list = l.list[index:]
			return
		}
	}
	l.list = nil
}

func callback(item *Timer) {
	log.Printf("callback: %+v\n", item)
	resp, err := http.Get(item.Callback)
	if err != nil {
		return
	}
	resp.Body.Close()
}
