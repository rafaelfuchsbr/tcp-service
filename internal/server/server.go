package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/rafaelfuchsbr/tcp-service/internal/commands"
	"github.com/rafaelfuchsbr/tcp-service/internal/globals"
)

type Server struct {
	port     int
	serverId string
	protocol string
	listener net.Listener
}

func NewServer(port int) *Server {
	globals.ServerId = uuid.NewString()
	return &Server{
		port:     port,
		serverId: globals.ServerId,
		protocol: "tcp",
	}
}

func (s *Server) Serve() error {
	log.Println(fmt.Sprintf("Starting server on port %d", s.port))
	err := s.start()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) start() error {
	var err error
	s.listener, err = net.Listen(s.protocol, fmt.Sprintf(":%d", s.port))
	log.Println(fmt.Sprintf("Started listening port %d", s.port))
	if err != nil {
		return err
	}

	for {
		connection, err := s.listener.Accept()
		log.Println("Got connection")
		if err != nil {
			return err
		}
		globals.ConnectedClients.AddClient()
		log.Println("Connection accepted, calling handler")
		go func(connection net.Conn) {
			defer connection.Close()
			defer globals.ConnectedClients.RemoveClient()
			for {
				err := handleConnection(connection)
				if err == io.EOF {
					log.Println("Connection is closed")
					return
				}
				if err != nil {
					log.Println(fmt.Errorf("Error handling connection %w", err))
				}
			}
		}(connection)
	}
}

func handleConnection(connection net.Conn) error {
	log.Println("Reading input")
	input, err := readInput(connection)
	if err != nil {
		return err
	}
	log.Println("Processing command")
	command, err := commands.NewCommand(input)
	if err != nil {
		return err
	}

	response, err := command.Process()
	if err != nil {
		return err
	}

	connection.Write([]byte(fmt.Sprintf("%s\n", response)))

	log.Println("Command done")
	return nil
}

func readInput(connection net.Conn) (string, error) {
	input, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		return "", err
	}
	log.Println(fmt.Sprintf("Input=[%s]", input))
	return input, nil
}
