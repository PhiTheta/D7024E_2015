package dht

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNet1(t *testing.T) {
	var wg sync.WaitGroup
	Notice("\n### TestNet1 Started\n")
	node := [8]*DHTNode{}
	fmt.Println("Legend:")
	Error("Error ")
	Notice("Notice ")
	Info("Info ")
	Warn("Warn \n")

	fmt.Print("")
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	node[0] = makeDHTNode(&id0, "localhost:1110")
	node[1] = makeDHTNode(&id1, "localhost:1111")
	node[2] = makeDHTNode(&id2, "localhost:1112")
	node[3] = makeDHTNode(&id3, "localhost:1113")
	node[4] = makeDHTNode(&id4, "localhost:1114")
	node[5] = makeDHTNode(&id5, "localhost:1115")
	node[6] = makeDHTNode(&id6, "localhost:1116")
	node[7] = makeDHTNode(&id7, "localhost:1117")
	wg.Add(8)
	go node[0].startServer(&wg)
	go node[1].startServer(&wg)
	go node[2].startServer(&wg)
	go node[3].startServer(&wg)
	go node[4].startServer(&wg)
	go node[5].startServer(&wg)
	go node[6].startServer(&wg)
	go node[7].startServer(&wg)
	wg.Wait()

	go node[1].send("join", node[0].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	go node[2].send("join", node[1].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	go node[3].send("join", node[1].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	go node[4].send("join", node[2].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	go node[5].send("join", node[2].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	go node[6].send("join", node[3].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	go node[7].send("join", node[3].bindAddress, "", "", "")
	time.Sleep(200 * time.Millisecond)
	// go node[1].printAll()
	// time.Sleep(200 * time.Millisecond)

	// Iterates and updates each nodes fingers.
	fmt.Println("Setup fingers!")
	for i := 0; i < 8; i++ {
		go node[i].setupFingers()
	}
	time.Sleep(100 * time.Millisecond)
	/*for i := 0; i < 8; i++ {
		Infoln("Node " + node[i].nodeId + ":" + node[i].FingersToString())
	}*/
	time.Sleep(50 * time.Millisecond)
	go node[0].printAll()
	time.Sleep(50 * time.Millisecond)
	node[1].send("lookup", node[1].successor.bindAddress, "", "10", "")

	time.Sleep(5000 * time.Millisecond)
	// go node[1].printAll()
	time.Sleep(5000 * time.Millisecond)

	//key string, src string, dst string, bytes string
}
