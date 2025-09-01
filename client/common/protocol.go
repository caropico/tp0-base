package common

import (
    "fmt"
    "encoding/binary"
    "net"
    "strings"
)

type Protocol struct {
    conn net.Conn
}

const (
    EOF_FLAG_TRUE  = 0x01
    EOF_FLAG_FALSE = 0x00
)

func NewProtocol(conn net.Conn) *Protocol {
    return &Protocol{conn: conn}
}

func (p *Protocol) SendBetMessage(batchResult BatchResult, agencyId string) error {
    bet_msg := make([]string,0, len(batchResult.batch))
    for _, b := range batchResult.batch {
        bet_msg = append(bet_msg, fmt.Sprintf("%s;%s;%s;%d;%s;%d", 
            agencyId,b.FirstName, b.LastName, b.DNI, b.Birthday, b.BetNumber))
    }
    
    message := strings.Join(bet_msg, "\n")

    if len(message) > 65535 {
        return fmt.Errorf("message too long")
    }

    messageSize := uint16(len(message))
    result := make([]byte, 3+len(message))
    if batchResult.isEOF {
        result[0] = EOF_FLAG_TRUE
    } else {
        result[0] = EOF_FLAG_FALSE
    }
    binary.BigEndian.PutUint16(result[1:3], messageSize)
    copy(result[3:], []byte(message))

    return p.SendAll(result)
}

func (p *Protocol) SendAll(data []byte) error {
    totalWritten := 0
    for totalWritten < len(data) {
        n, err := p.conn.Write(data[totalWritten:])
        if err != nil {
            return err
        }
        totalWritten += n
    }
    return nil
}

func (p *Protocol) ReceiveAll() (byte,error) {
    ack := make([]byte,1)
    totalRead := 0
    for totalRead < 1 {
        n, err := p.conn.Read(ack[totalRead:])
        if err != nil {
            return 0,err
        }
        totalRead += n
    }

    return ack[0],nil
}

func (p *Protocol) Close() error {
    return p.conn.Close()
}