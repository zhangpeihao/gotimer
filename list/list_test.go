package list

import (
	"testing"
	"time"
)

func TestTimerList(t *testing.T) {
	list := NewList()
	for i := 0; i < 10; i++ {
		list.AddTimer(&Timer{
			time:     time.Now().Unix() + int64(i),
			callback: "http://www.baidu.com",
		})
	}
	time.Sleep(time.Second * time.Duration(13))
}
