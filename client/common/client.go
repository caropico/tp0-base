package common

import (
	//"bufio"
	//"fmt"
	"net"
	"time"

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
	bet   ClientBet
	protocol *Protocol
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet ClientBet) *Client {
	client := &Client{
		config: config,
		bet: bet,
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

func (c *Client) sendBetMessage() error {
	err := c.protocol.SendBetMessage(c.bet, c.config.ID)
		if err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
		}
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
			log.Infof("action: apuesta_enviada | result: success | dni: ${%d} | numero: ${%d}", 
    			c.bet.DNI, c.bet.BetNumber)
		} else {
			log.Infof("action: apuesta_enviada | result: fail | dni: ${%d} | numero: ${%d}", 
    			c.bet.DNI, c.bet.BetNumber)
		}
	return nil
}

func (c *Client) closeConnection() {
    if c.protocol != nil {
        c.protocol.Close()
    }
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(signals chan os.Signal) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	select {
	case <- signals:
		log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
		return
	default: 
		c.createClientSocket()

		c.sendBetMessage()

		c.receiveAck()

		c.closeConnection()

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
