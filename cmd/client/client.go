package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {

	var wg sync.WaitGroup

	clientsPtr := flag.Int("clients", 1, "")
	interationsPtr := flag.Int("iterations", 1, "")
	commandsPtr := flag.String("commands", "WHERE,WHO,WHY", "")
	hostPtr := flag.String("host", "127.0.0.1", "")
	portPtr := flag.Int("port", 12345, "")

	flag.Parse()

	nrClients := *clientsPtr
	wg.Add(nrClients)

	log.Println(fmt.Sprintf("Running %d clients, %d iterations each with commands %s", nrClients, *interationsPtr, *commandsPtr))

	for i := 0; i < nrClients; i++ {
		go func(i int) {
			defer wg.Done()
			log.Println(fmt.Sprintf("[client-%d] Running client %d", i, i))
			client := Client{
				host:       *hostPtr,
				port:       *portPtr,
				commands:   *commandsPtr,
				iterations: *interationsPtr,
				clientId:   i,
			}
			err := client.run()
			if err != nil {
				log.Println(fmt.Errorf("Error running client [%d]: %w", i, err))
			}
			log.Println(fmt.Sprintf("[client-%d] Exiting client %d", i, i))
		}(i)
	}
	wg.Wait()
}

type Client struct {
	connection net.Conn
	commands   string
	host       string
	port       int
	iterations int
	clientId   int
}

func (c *Client) run() error {

	connectionString := fmt.Sprintf("%s:%d", c.host, c.port)

	log.Println(fmt.Sprintf("[client-%d] Trying to connect to %s", c.clientId, connectionString))
	connection, err := net.Dial("tcp", connectionString)
	if err != nil {
		return fmt.Errorf("[client-%d] Error connecting to server", c.clientId)
	}

	log.Println(fmt.Sprintf("[client-%d] Connected to %s", c.clientId, connectionString))
	c.connection = connection
	defer c.connection.Close()

	for i := 0; i < c.iterations; i++ {
		err = c.processCommandList()
		if (err == io.EOF) || (errors.Is(err, syscall.EPIPE)) || (errors.Is(err, syscall.ECONNRESET)) {
			log.Println(fmt.Sprintf("[client-%d] Remote connection closed", c.clientId))
			return nil
		}
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}

	return nil
}

func (c *Client) processCommandList() error {
	for _, command := range strings.Split(c.commands, ",") {
		err := c.sendCommand(command)
		if err != nil {
			return err
		}

		reponse, err := c.getResponse()
		if err != nil {
			return err
		}
		log.Println(fmt.Sprintf("[client-%d] Response from server: %s", c.clientId, strings.Trim(reponse, "\n")))
	}
	return nil
}

func (c *Client) sendCommand(command string) error {
	_, err := c.connection.Write([]byte(command + "\n"))
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("[client-%d] Command sent: %s", c.clientId, command))
	return nil
}

func (c *Client) getResponse() (string, error) {
	response, err := bufio.NewReader(c.connection).ReadString('\n')
	if err != nil {
		return "", err
	}
	return response, nil
}
