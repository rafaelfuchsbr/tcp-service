package globals

import (
	"log"
	"sync"
)

var ServerId string
var ConnectedClients Clients

type Clients struct {
	connected int
	sync.Mutex
}

func (c *Clients) Total() int {
	c.Lock()
	defer c.Unlock()
	return c.connected
}

func (c *Clients) AddClient() {
	log.Println("Adding client")
	c.Lock()
	defer c.Unlock()
	c.connected++
}

func (c *Clients) RemoveClient() {
	log.Println("Removing client")
	c.Lock()
	defer c.Unlock()
	c.connected--
}
