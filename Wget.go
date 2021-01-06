package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func ShowLinks()  {
	fmt.Println(
		"\n--------------------------------------------------------------\n" +
		"Test links:\n" +
		" 10 Mbit file: http://ovh.net/files/10Mb.dat\n" +
		" 100 Mbit file: http://ovh.net/files/100Mb.dat\n" +
		" Jpeg file: https://i.ytimg.com/vi/3RTuFB2mq20/maxresdefault.jpg\n" +
		" Music file: http://commondatastorage.googleapis.com/codeskulptor-assets/week7-brrring.m4a\n" +
		"--------------------------------------------------------------")
}
func ShowBar(progress int)  {
	if progress>=90 	   {fmt.Print("\t [███████████████████████████░░░]\n\n")
	} else if progress>=80 {fmt.Print("\t [████████████████████████░░░░░░]\n\n")
	} else if progress>=70 {fmt.Print("\t [█████████████████████░░░░░░░░░]\n\n")
	} else if progress>=60 {fmt.Print("\t [██████████████████░░░░░░░░░░░░]\n\n")
	} else if progress>=50 {fmt.Print("\t [███████████████░░░░░░░░░░░░░░░]\n\n")
	} else if progress>=40 {fmt.Print("\t [████████████░░░░░░░░░░░░░░░░░░]\n\n")
	} else if progress>=30 {fmt.Print("\t [█████████░░░░░░░░░░░░░░░░░░░░░]\n\n")
	} else if progress>=20 {fmt.Print("\t [██████░░░░░░░░░░░░░░░░░░░░░░░░]\n\n")
	} else if progress>=10 {fmt.Print("\t [███░░░░░░░░░░░░░░░░░░░░░░░░░░░]\n\n")
	} else if progress>=0  {fmt.Print("\t [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]\n\n")
	}
}
func ClearTerminal()  {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
type Background struct {
	Total           uint64
	IsDownloaded    bool
	ShutdownChannel chan string
	Progress        int
	Speed           float64
}

func NewBackground() *Background {
	return &Background{
		Total:           0,
		IsDownloaded:    false,
		ShutdownChannel: make(chan string),
		Progress:        0,
		Speed:           0,
	}
}

func (bg *Background) Run(dataLength uint64) {
	// Loop that runs forever
	_count:=1
	fmt.Printf("\t\t↓ 0%%\t[%d bytes]  \t\t%.2f MByte/sec\t\t [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]\n\n",
		bg.Total, bg.Speed)
	for {
		select {
		case <-bg.ShutdownChannel:
			bg.ShutdownChannel <- "Stop"
			return
		case <-time.After(time.Second):
			break
		}
		bg.Progress=int(float32(bg.Total)/float32(dataLength) * float32(100))
		bg.Speed = float64(bg.Total)/float64(_count)/1048576
		fmt.Printf("\t\t↓ %d%%\t[%d bytes]  \t%.2f MByte/sec\t ", bg.Progress, bg.Total, bg.Speed)
		ShowBar(bg.Progress)
		_count++
	}
}

func (bg *Background) Shutdown() {
	bg.IsDownloaded = true

	bg.ShutdownChannel <- "Stop"
	<-bg.ShutdownChannel
	close(bg.ShutdownChannel)

	fmt.Printf("\t\t↓ 100%%\t[%d bytes]  \t%.2f MByte/sec\t\t [██████████████████████████████]\n",
		bg.Total, bg.Speed)
}

func (bg *Background) Write(p []byte) (int, error) {
	n := len(p)
	bg.Total += uint64(n)
	return n, nil
}

func DownloadFile(file string, url string) (error, uint64) {

	// Create the file
	out, err := os.Create(file + ".tmp")
	if err != nil {
		return err, 0
	}

	// Get the data
	data, err := http.Get(url)
	if err != nil {
		out.Close()
		return err, 0
	}
	defer data.Body.Close()

	dataLength:=uint64(data.ContentLength)
	timer := NewBackground()
	go timer.Run(dataLength)
	_,err = io.Copy(out, io.TeeReader(data.Body, timer))
	timer.Shutdown()
	if err != nil {
		out.Close()
		return err, 0
	}
	out.Close()

	err = os.Rename(file+".tmp", file)
	if  err != nil {
		return err, 0
	}
	return nil, dataLength
}


func main() {
	runtime.GOMAXPROCS(4)
	var URL, fileName, flag string
	for {
		for {
			ClearTerminal()
			ShowLinks()
			fmt.Println("\nPaste your download link: ")
			fmt.Scan(&URL)
			fileName = URL[strings.LastIndex(URL, "/")+1:]
			fmt.Printf("Your file: %s\nDOWNLOAD? [Y/N]\n", fileName)
			fmt.Scan(&flag)
			if flag == "y" || flag == "Y" {
				break
			}
		}
		fmt.Println("\t\tDOWNLOAD STARTED.")
		start:=time.Now()
		err, total := DownloadFile(fileName, URL)
		duration:=time.Since(start)
		if err != nil {
			fmt.Println("Error 0: Wrong link\n", err)
			time.Sleep(5*time.Second)
			continue
		}
		fmt.Printf("\t\tDOWNLOAD COMPLETE.\n\nDownloaded: %s (%d bytes) [%.2fs]\n",
			fileName, total, duration.Seconds())
		fmt.Println("DOWNLOAD MORE? [Y/N]")
		fmt.Scan(&flag)
		if flag == "n" || flag == "N" {
			break
		}
	}
}

// Links

// 10 Mbit file: http://ovh.net/files/10Mb.dat 100 Mbit file:
// http://ovh.net/files/100Mb.dat Jpeg file:
// https://i.ytimg.com/vi/3RTuFB2mq20/maxresdefault.jpg Music file:
// http://commondatastorage.googleapis.com/codeskulptor-assets/week7-brrring.m4a

// Исключения


