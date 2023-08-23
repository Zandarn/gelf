package gelf

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Jeffail/gabs"
)

var hostname string = "localhost"

const (
	defaultEndpoint     = "127.0.0.1:12201"
	defaultMaxChunkSize = 1420
)

const (
	EMERGENCY int = 0
	ALERT     int = 1
	CRITICAL  int = 2
	ERROR     int = 3
	WARNING   int = 4
	NOTICE    int = 5
	INFO      int = 6
	DEBUG     int = 7
)

const (
	LOCAL0 string = "local0"
	LOCAL1 string = "local1"
	LOCAL2 string = "local2"
	LOCAL3 string = "local3"
	LOCAL4 string = "local4"
	LOCAL5 string = "local5"
	LOCAL6 string = "local6"
	LOCAL7 string = "local7"
)

var LogFacility = facility{}
var LogLevel = level{}

type facility struct {
	LOCAL0 string
	LOCAL1 string
	LOCAL2 string
	LOCAL3 string
	LOCAL4 string
	LOCAL5 string
	LOCAL6 string
	LOCAL7 string
}

type level struct {
	EMERGENCY int
	ALERT     int
	CRITICAL  int
	ERROR     int
	WARNING   int
	NOTICE    int
	INFO      int
	DEBUG     int
}

type Message struct {
	Version   string            `json:"version"`
	Host      string            `json:"host"`
	Message   string            `json:"message"`
	Timestamp float64           `json:"timestamp"`
	Level     int               `json:"level"`
	Facility  string            `json:"facility"`
	Extra     map[string]string `json:"-"`
}

type Config struct {
	Endpoint     string
	MaxChunkSize int
}

type Client struct {
	Config *Config
}

func GelfClient() *Client {
	client := &Client{}
	client.setDefaults()

	return client
}

func (g *Client) setDefaults() {
	g.Config = &Config{
		Endpoint:     defaultEndpoint,
		MaxChunkSize: defaultMaxChunkSize,
	}

	LogLevel = level{
		EMERGENCY: EMERGENCY,
		ALERT:     ALERT,
		CRITICAL:  CRITICAL,
		ERROR:     ERROR,
		WARNING:   WARNING,
		NOTICE:    NOTICE,
		INFO:      INFO,
		DEBUG:     DEBUG,
	}

	LogFacility = facility{
		LOCAL0: LOCAL0,
		LOCAL1: LOCAL1,
		LOCAL2: LOCAL2,
		LOCAL3: LOCAL3,
		LOCAL4: LOCAL4,
		LOCAL5: LOCAL5,
		LOCAL6: LOCAL6,
		LOCAL7: LOCAL7,
	}

	hostname, _ = os.Hostname()
}

func (g *Client) SendMessage(message Message) (n int, err error) {
	data, err := g.prepare(message)
	if err != nil {
		log.Printf("Error preparing message: %s", err)
		return 0, err
	}

	if _, err := g.udp(data); err != nil {
		return 0, err
	}

	return len(data), nil
}

func (g *Client) prepare(message Message) ([]byte, error) {
	message.setDefaults()
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error json.Marshal: %s", err)
		return []byte{}, err
	}

	c, err := gabs.ParseJSON(jsonMessage)
	if err != nil {
		log.Printf("Error gabs.ParseJSON: %s", err)
		return nil, err
	}

	for key, value := range message.Extra {
		_, err = c.Set(value, fmt.Sprintf("_%s", key))
		if err != nil {
			log.Printf("Error adding extra fields: %s", err)
			return nil, err
		}
	}

	// end of message
	data := append(c.Bytes(), '\n', 0)

	return data, nil
}

func (message *Message) setDefaults() {
	if message.Version == "" {
		message.Version = "1.1"
	}

	if message.Host == "" {
		message.Host = hostname
	}

	if message.Timestamp == 0 {
		message.Timestamp = float64(time.Now().Unix())
	}
}

func (g *Client) udp(message []byte) (n int, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", g.Config.Endpoint)
	if err != nil {
		log.Printf("Error resolving UDP address: %s", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("Error dialing UDP: %s", err)
		return
	}
	defer conn.Close()

	n, err = conn.Write(message)
	if err != nil {
		log.Printf("Error writing to UDP connection: %s", err)
		return 0, err
	}

	return n, nil
}
