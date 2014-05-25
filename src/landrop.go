package landrop


import "fmt"
import "landrop_go/networking/udp"


func main() {
	go udp.CreateStatusServer();
	go udp.SendStatus("localhost:42000", "hello");
	go udp.SendStatus("localhost:42000", "new");
}