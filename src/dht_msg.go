package dht

import (
	"encoding/json"
	"fmt"
	"net"
)

type DHTMsg struct {
	Key string `json:"key"`
	Src string `json:"src"`
	Dst string `json:"dst"`
	// Other?
	Req  string `json:"req"` //request type
	Opt  string `json:"opt"` //option --
	Data string `json:"data"`//data string [key:ip:port]
	//Bytes string `json:"bytes"`
}

type Transport struct {
	node       *DHTNode
	bindAdress string
	queue      chan *DHTMsg
}

//Initiates Transport
func CreateTransport(dhtNode *DHTNode, bindAdress string) *Transport {
	transport := &Transport{}
	transport.bindAdress = bindAdress
	transport.node = dhtNode
	transport.queue = make(chan *DHTMsg)
	go transport.handler()
	return transport
}

//Creates Message with all info
func CreateMsg(key string, src string, dst string, req string, opt string, data string) *DHTMsg {
	dhtMsg := &DHTMsg{}
	dhtMsg.Key = key
	dhtMsg.Src = src
	dhtMsg.Dst = dst
	dhtMsg.Req = req
	dhtMsg.Opt = opt
	dhtMsg.Data = data
	return dhtMsg
}

//Handles the different msg types from transport
func (transport *Transport) handler() {
	for {
		select {
		case msg := <-transport.queue:
			switch msg.Req {
			case "update":
				transport.node.updateNode(msg)
			case "join":
				transport.node.joinRing(msg)
			case "printring":
				if msg.Opt != transport.node.nodeId {
					transport.node.send("printring", transport.node.predecessor, msg.Opt, msg.Data+"->"+transport.node.nodeId)
				} else {
					fmt.Println(msg.Data)
				}
			case "printall":
				if msg.Opt != transport.node.nodeId {
					transport.node.send("printall", transport.node.successor, msg.Opt, msg.Data+"\n"+transport.node.nodeId+"\t"+transport.node.predecessor.nodeId+"\t"+transport.node.successor.nodeId)
				} else {
					fmt.Print("Node\tPre.\tSucc.")
					fmt.Println(msg.Data + "\n" + transport.node.nodeId + "\t" + transport.node.predecessor.nodeId + "\t" + transport.node.successor.nodeId)
				}
			}

		}
	}
}

//Listens for a connection
func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAdress)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	dec := json.NewDecoder(conn)
	fmt.Println("Started listening : " + transport.bindAdress)
	for {
		msg := DHTMsg{}
		err = dec.Decode(&msg)
		// We got something?
		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println(msg)
		//fmt.Println(transport.node.nodeId + ":> from:" + msg.Src + " to: " + msg.Dst)
		transport.queue <- &msg
		//fmt.Println(msg.Bytes + " " + "")
	}
}

//Sends a msg
func (dhtNode *DHTNode) send(req string, dstNode *DHTNode, opt, data string) {

	msg := CreateMsg(dhtNode.nodeId, dhtNode.transport.bindAdress, dstNode.transport.bindAdress, req, opt, data)
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	res, _ := json.Marshal(msg)
	_, err = conn.Write(res) // wat?

}
