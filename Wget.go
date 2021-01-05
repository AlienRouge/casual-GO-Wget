package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func ShowLinks()  {
	fmt.Println(
		"--------------------------------------------------------------\n" +
		"Test links:\n" +
		" 10 Mbit file: http://ovh.net/files/10Mb.dat\n" +
		" 100 Mbit file: http://ovh.net/files/100Mb.dat\n" +
		" Jpeg file: https://i.ytimg.com/vi/3RTuFB2mq20/maxresdefault.jpg\n" +
		" Music file: http://commondatastorage.googleapis.com/codeskulptor-assets/week7-brrring.m4a\n" +
		"--------------------------------------------------------------")
}

type Counter struct {
	Total uint64
	IsDownloaded bool
	ShutdownChannel chan string
}

func NewCounter() *Counter {
	return &Counter{
		Total: 0,
		IsDownloaded:         false,
		ShutdownChannel: make(chan string),
		//Interval:        time.Second,
		//period:          interval,
	}
}

func (wc *Counter) Run() {
	log.Println("Download started.\n Downloaded:",wc.Total, "bytes")
	// Loop that runs forever
	for {
		select {
		case <-wc.ShutdownChannel:
			wc.ShutdownChannel <- "Stop"
			return
		case <-time.After(time.Second):
			break
		}
		log.Println("Downloaded:", wc.Total, "bytes")
		//started := time.Now()
		//wc.Action()
		//finished := time.Now()

		//duration := finished.Sub(started)
		//wc.period = wc.Interval - duration

	}
}

func (wc *Counter) Shutdown() {
	wc.IsDownloaded = true

	wc.ShutdownChannel <- "Stop"
	<-wc.ShutdownChannel

	close(wc.ShutdownChannel)
	log.Println("Download complete\n","Downloaded:", wc.Total, "bytes")
}

func (wc *Counter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	return n, nil
}


func DownloadFile(file string, url string) error {

	// Create the file
	out, err := os.Create(file + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	data, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer data.Body.Close()

	timer := NewCounter()
	go timer.Run()
	_,err = io.Copy(out, io.TeeReader(data.Body, timer))
	timer.Shutdown()
	if err != nil {
		out.Close()
		return err
	}
	out.Close()

	err = os.Rename(file+".tmp", file)
	if  err != nil {
		return err
	}
	return nil
}


func main() {
	var URL, fileName string
	for {
		ShowLinks()
		fmt.Println("\nPaste your download link: ")
		fmt.Scan(&URL)
		fileName = URL[strings.LastIndex(URL, "/")+1:]
		fmt.Printf("Your file: %s\nDOWNLOAD? [Y/N]\n", fileName)
		f := "N"
		fmt.Scan(&f)
		if f == "y" || f == "Y" {
			break
		}
	}
	err := DownloadFile(fileName, URL)
	if err != nil {
		panic(err)
	}
}

// Links

// 10 Mbit file: http://ovh.net/files/10Mb.dat
// 100 Mbit file: http://ovh.net/files/100Mb.dat
// Jpeg file: https://i.ytimg.com/vi/3RTuFB2mq20/maxresdefault.jpg
// Music file: http://commondatastorage.googleapis.com/codeskulptor-assets/week7-brrring.m4a

// Исключения
// Проценты
// Рефактор