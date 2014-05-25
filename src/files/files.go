package files

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/howeyc/fsnotify"
	"log"
	"lds"
)

func GetFileTree(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
}

func FileSystemWorker(queue chan lds.BEvent) {
	for{
		fmt.Println("FileSystemWorker called")
		event := <-queue
		fmt.Println(event.Fname);
	}
}

func FileHandler(queue chan lds.BFile){
	for{
		file := <-queue
		fmt.Print("FileHandler ")
		fmt.Println(file.Name)
	}
}

func WatchAndAlert(path string, fevents chan *fsnotify.FileEvent) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for{
			select{
			case ev := <-watcher.Event:
				fevents<-ev
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}

	<-done
	watcher.Close()
}
