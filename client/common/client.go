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
	ack, err := c.protocol.ReceiveAck()
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

func (c *Client) sendBetMessage(batchResult BatchResult) error {
    err := c.protocol.SendBetMessage(batchResult, c.config.ID)
    if err != nil {
        log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
            c.config.ID, err)
    }
    return err
}

func (c *Client) CheckForWinners() error {
    err := c.protocol.SendCheckForWinnersMessage(c.config.ID)
    if err != nil {
        log.Errorf("action: send_check_winners | result: fail | client_id: %v | error: %v",
            c.config.ID, err)
        return err
    }
    return nil
}

func (c *Client) ReceiveWinners() ([]string, error) {
    msg, err := c.protocol.ReceiveWinnersMessage()
    if err != nil {
        log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
            c.config.ID, err)
			return []string{}, err
    }
    return msg, nil
}

func (c *Client) closeConnection() {
    if c.protocol != nil {
        c.protocol.Close()
    }
}

func (c *Client) sendAllBatches(signals chan os.Signal) error {
    processor, err := createCSVProcessor("./agency.csv", c.config.MaxAmountBatch)
    if err != nil {
        return err
    }
    defer processor.Close()
    
    err = c.createClientSocket()
    if err != nil {
        return err
    }
    defer c.closeConnection()

	err = c.protocol.SendLoadBetCode()
    if err != nil {
        return err
    }
    
    for c.isRunning && processor.HasMoreBatches() {
        select {
        case <-signals:
            return nil
        default:
            dataBatch, err := processor.readNextBatch(c.config.ID)
            if err != nil {
                if err == io.EOF {
                    break
                }
                return err
            }
            
            err = c.sendBetMessage(dataBatch)
            if err != nil {
                return err
            }
            
            err = c.receiveAck()
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func (c *Client) checkForWinners(signals chan os.Signal) error {
    err := c.createClientSocket()
    if err != nil {
        return err
    }
    defer c.closeConnection()
    
    err = c.protocol.SendCheckForWinnersMessage(c.config.ID)
    if err != nil {
        return err
    }
    
    winners, err := c.protocol.ReceiveWinnersMessage()
    if err != nil {
        return err
    }
    
    log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
    return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(signals chan os.Signal) {
	c.isRunning = true

	err := c.sendAllBatches(signals)
    if err != nil {
        log.Errorf("action: send_batches | result: fail | error: %v", err)
        return
    }
    
    err = c.checkForWinners(signals)
    if err != nil {
        log.Errorf("action: check_winners | result: fail | error: %v", err)
    }

	time.Sleep(c.config.LoopPeriod)
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}