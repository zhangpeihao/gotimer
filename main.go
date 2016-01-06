package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/zhangpeihao/gotimer/list"
)

const (
	programName = "gocmd"
	version     = "0.1"
)

var (
	bindAddress *string = flag.String("BindAddress", ":8001", "The bind address.")
	root        *string = flag.String("Root", "/root/src", "The root path.")
)

var g_list *list.List

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s version[%s]\r\nUsage: %s [OPTIONS]\r\n", programName, version, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	g_list = list.NewList()

	http.HandleFunc("/add", addHandler)

	go http.ListenAndServe(*bindAddress, nil)
	catchSignal()

}

func catchSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	os.Exit(0)
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	callback := req.FormValue("callback")
	time := req.FormValue("time")
	if callback == "" {
		log.Printf("callback or time is empty, callback: %s, time: %s\n", callback, time)
		w.Write([]byte(`{"ret":"1000","msg":"parameter error"}`))
		return
	}
	timeInt, err := strconv.Atoi(time)
	if err != nil {
		log.Printf("time isn't an integer, callback: %s, time: %s\n", callback, time)
		w.Write([]byte(`{"ret":"1000","msg":"parameter error"}`))
		return
	}
	g_list.AddTimer(&list.Timer{
		Time:     int64(timeInt),
		Callback: callback,
	})
}