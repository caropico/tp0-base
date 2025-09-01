package common

import (
	//"bufio"
	//"fmt"
	"net"
	"time"
	"io"
	"github.com/op/go-logging"

	"os"
)
var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxAmountBatch int
}

// ClientBet stores client bet data
type ClientBet struct {
	FirstName            string
	LastName string
	DNI    int
	Birthday    string
	BetNumber    int
}


// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	protocol *Protocol
	isRunning bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		isRunning: false,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {

	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.protocol = NewProtocol(conn)
	return nil
}

func (c *Client) receiveAck() error {
	ack, err := c.protocol.ReceiveAll()
		if err != nil {
			log.Errorf("action: receive_ack | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return err
		}

		if ack == 1{
			log.Infof("action: apuesta_enviada | result: success ")
		} else {
			log.Infof("action: apuesta_enviada | result: fail ")
		}
	return nil
}

func (c *Client) sendBetMessage(bets []ClientBet) error {
    err := c.protocol.SendBetMessage(bets, c.config.ID)
    if err != nil {
        log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
            c.config.ID, err)
    }
    return err
}

func (c *Client) closeConnection() {
    if c.protocol != nil {
        c.protocol.Close()
    }
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(signals chan os.Signal) {
	c.isRunning = true
	processor, err := createCSVProcessor("./agency.csv", c.config.MaxAmountBatch)
    if err != nil {
        log.Errorf("action: create_csv_processor | result: fail | client_id: %v | error: %v", c.config.ID, err)
        return
    }
    defer processor.Close() 
	for c.isRunning {
		select {
			case <- signals:
			log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
			c.isRunning = false
			return
		default: 
			dataBatch, err := processor.readNextBatch()
			if err != nil {
				if err == io.EOF {
					c.isRunning = false
					c.closeConnection()
					break
				}
			log.Errorf("action: read_batch | result: fail | error: %v", err)
			c.closeConnection()
			break
			}

			err = c.createClientSocket()
        	if err != nil {
            	log.Errorf("action: create_socket | result: fail | error: %v", err)
            	break
        	}
	
			err = c.sendBetMessage(dataBatch)
            if err != nil {
                log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", 
                    c.config.ID, err)
                c.closeConnection()
                break
            }
            
            err = c.receiveAck()
            if err != nil {
                log.Errorf("action: receive_batch_ack | result: fail | client_id: %v | error: %v", 
                    c.config.ID, err)
            }
            
            c.closeConnection()
		}
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
