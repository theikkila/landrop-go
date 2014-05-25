
package networking

import (
	"fmt"
	"net"
	"os"
	"github.com/howeyc/fsnotify"
	"encoding/json"
	"encoding/gob"
	"lds"
)


func BroadcastEvents(service string, fevents chan *fsnotify.FileEvent) {
	for {
		event := <-fevents
		var e lds.BEvent
		e.Fname = event.Name
		if event.IsCreate(){
			e.Etype = 1
			e.Etypestring = "CREATE"
		}
		if event.IsDelete(){
			e.Etype = 2
			e.Etypestring = "DELETE"
		}

		if event.IsModify(){
			e.Etype = 3
			e.Etypestring = "MODIFY"
		}
		if event.IsRename(){
			e.Etype = 4
			e.Etypestring = "RENAME"
		}

		b, err := json.Marshal(e)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(e.Etypestring);
		go  SendStatus(service, b)
	}
}

func SendStatus(service string, msg []byte) {

	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)

	conn, err := net.DialUDP("udp4", nil, udpAddr)
	checkError(err)
	fmt.Println("Sending status...");
	_, err = conn.Write([]byte(msg))
	checkError(err)

	var buf [512]byte
	_, err = conn.Read(buf[0:])
	
	checkError(err)

}

func FileFetcher(queue chan lds.NBEvent, files chan lds.BFile){
	for{
		event := <-queue
		// Fetch file
		fmt.Println("Fetching file...")
		var file lds.BFile;
		file.Name = event.Fname;
		files<-file;
	}
}

func CreateStatusServer(queue chan lds.NBEvent) {
	fmt.Println("creating...")
	service := ":41000"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	fmt.Println(udpAddr)
	conn, err := net.ListenUDP("udp4", udpAddr)
	checkError(err)
	for {
		handleUDPClient(conn, queue)
	}
}

func CreateFileServer() {
	service := ":42000"
	ln, err := net.Listen("tcp", service)
	if err != nil {
		return
	}

	for{
		c, err := ln.Accept()
		if err != nil{
			continue
		}
		go handleFileRequest(c)
	}
}

func handleFileRequest(c net.Conn) {
	var f lds.FileRequest
	err := gob.NewDecoder(c).Decode(&f)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("received", f.Fname)
		file, err := os.Open(f.Fname)
		if err != nil {
			fmt.Println(err)
		}
		err = gob.NewEncoder(c).Encode(file);
		if err != nil {
			fmt.Println(err)
		}
	}
	c.Close()
}

func SendRequest(address string, filename string) {
	// connect to the server
    c, err := net.Dial("tcp", address)
    if err != nil {
        fmt.Println(err)
        return
    }

    // send the message
    var fr lds.FileRequest
    fr.Fname = filename
    fmt.Println("Sending", fr.Fname)
    err = gob.NewEncoder(c).Encode(fr)
    if err != nil {
        fmt.Println(err)
    }
    c.Close()
}

func handleUDPClient(conn *net.UDPConn, queue chan lds.NBEvent) {
	fmt.Println("Handling new client...")
	var buf []byte = make([]byte, 10)

	s, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return
	}
	go handlePacket(buf, s, addr.IP.String(), queue);
}

func handlePacket(buf []byte, s int, ip string, queue chan lds.NBEvent){
	var msg = string( buf[0:s] );
	//addr.IP.String()
	fmt.Println(msg);
	var m lds.BEvent
	err := json.Unmarshal(buf, m)
	if err != nil {
        fmt.Println(err)
    }
	var e lds.NBEvent
	e.Fname = m.Fname
	e.Etype = m.Etype
	e.Etypestring = m.Etypestring
	e.Remote = ip
	queue<-e
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
