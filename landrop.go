package main

import "fmt"
import "networking"
import "files"
import "lds"
import "github.com/howeyc/fsnotify"
import "os"


func main() {
	fmt.Println("Start...")
	fmt.Println("Making chains...")
	fromFStoUDPfevents := make(chan *fsnotify.FileEvent, 100)
	fromUDPtoEventHandlerEvents := make(chan lds.NBEvent, 100)
	fetchFileQueue := make(chan lds.NBEvent, 100)
	fsQueue := make(chan lds.BFile, 100)
	fmt.Println("Creating StatusServer...")
	go networking.CreateStatusServer(fromUDPtoEventHandlerEvents)
	go networking.CreateFileServer()
	
	go networking.FileFetcher(fetchFileQueue, fsQueue)
	go networking.BroadcastEvents("localhost:42000", fromFStoUDPfevents)
	
	//go networking.SendStatus("localhost:42000", "hello")
	//go networking.SendStatus("localhost:42000", "new")
	fmt.Println("Watching...")
	go files.WatchAndAlert("./share/", fromFStoUDPfevents)
	go files.FileHandler(fsQueue)
	go EventHandler(fromUDPtoEventHandlerEvents, fetchFileQueue, fsQueue)

}

func EventHandler(udpqueue chan lds.NBEvent, tcpqueue chan lds.NBEvent, fsqueue chan lds.BFile) {
	// udp -> tcpqueue
	// udp -> fsqueue
}