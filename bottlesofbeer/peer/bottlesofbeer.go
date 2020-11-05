package main
import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"time"

	//"time"
	"bottlesofbeer/stubs"
)

var nextAddr string

type Bottles struct {}

func makeRequest(num int, nextAddr string) {
	request := stubs.Request{Number: num}
	response := new(stubs.Response)
	client, _ := rpc.Dial("tcp", nextAddr)
	defer client.Close()
	done := make(chan *rpc.Call, 1)
	client.Go(stubs.Sing, request, response, done)
	fmt.Println(((<-done).Reply).(*stubs.Response).Message)
}

func (s *Bottles) Sing(req stubs.Request, res *stubs.Response) (err error) {
	if req.Number > 0 {
		numString := strconv.Itoa(req.Number)
		plural := "bottles"
		if req.Number == 1 {
			plural = "bottle"
		}
		if req.Number > 1 {
			res.Message = numString + " " + plural + " of beer on the wall, " + numString + " " + plural + " of beer. Take one down, pass it around..."
		} else if req.Number == 1 {
			res.Message = numString + " bottle of beer on the wall, " + numString + " bottle of beer. Take one down, pass it around..."
		}
		time.Sleep(1 * time.Second)
		go makeRequest(req.Number - 1, nextAddr)
	} else {
		time.Sleep(1 * time.Second)
		res.Message = "No more bottles of beer on the wall. Goodbye."
	}
	return
}

func main(){
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	bottles := flag.Int("n",-1, "Bottles of Beer (launches song if not 0)")
	flag.Parse()
	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.
	if *bottles >= 0 {
		go makeRequest(*bottles, nextAddr)
	}
	rpc.Register(&Bottles{})
	listener, _ := net.Listen("tcp", ":" + *thisPort)
	defer listener.Close()
	rpc.Accept(listener)
}
