package common

import (
    "fmt"
    "encoding/binary"
    "net"
    "strings"
    "io"
)


const (
    MAX_BET_MESSAGE_SIZE = 65535
    LOAD_BET_MESSAGE_CODE = 0x02
    CHECK_FOR_WINNERS_MESSAGE_CODE = 0x03
    SEND_WINNERS_MESSAGE_CODE = 0x04 
    EOF_FLAG_TRUE  = 0x01
    EOF_FLAG_FALSE = 0x00
)

type Protocol struct {
    conn net.Conn
}

func NewProtocol(conn net.Conn) *Protocol {
    return &Protocol{conn: conn}
}

func (p *Protocol) SendLoadBetCode() error {
    result := make([]byte, 1)
    result[0] = LOAD_BET_MESSAGE_CODE
    return p.SendAll(result)
}

func (p *Protocol) SendBetMessage(batchResult BatchResult, agencyId string) error {
    bet_msg := make([]string,0, len(batchResult.batch))
    for _, b := range batchResult.batch {
        bet_msg = append(bet_msg, fmt.Sprintf("%s;%s;%s;%d;%s;%d", 
            agencyId,b.FirstName, b.LastName, b.DNI, b.Birthday, b.BetNumber))
    }
    
    message := strings.Join(bet_msg, "\n")

    if len(message) > MAX_BET_MESSAGE_SIZE {
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

    err := p.SendAll(result)

    return err
}

func (p *Protocol) SendCheckForWinnersMessage(agencyId string) error {
    message := agencyId
    
    if len(message) > MAX_BET_MESSAGE_SIZE {
        return fmt.Errorf("agency ID too long")
    }
    
    messageSize := uint16(len(message))
    result := make([]byte, 3+len(message))
    result[0] = CHECK_FOR_WINNERS_MESSAGE_CODE
    binary.BigEndian.PutUint16(result[1:3], messageSize)
    copy(result[3:], []byte(message))

    err := p.SendAll(result)
    return err
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

func (p *Protocol) ReceiveAck() (byte,error) {
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

func (p *Protocol)  ReceiveWinnersMessage() ([]string, error) {
    codeBytes := make([]byte, 1)
    _, err := io.ReadFull(p.conn, codeBytes)
    if err != nil {
        return nil, fmt.Errorf("error receiving message code: %w", err)
    }
    
    if codeBytes[0] != SEND_WINNERS_MESSAGE_CODE {
        return nil, fmt.Errorf("unexpected message code: %d", codeBytes[0])
    }
    
    sizeBytes := make([]byte, 2)
    _, err = io.ReadFull(p.conn, sizeBytes)
    if err != nil {
        return nil, fmt.Errorf("error receiving message size: %w", err)
    }
    
    messageSize := binary.BigEndian.Uint16(sizeBytes)
    
    messageBytes := make([]byte, messageSize)
    _, err = io.ReadFull(p.conn, messageBytes)
    if err != nil {
        return nil, fmt.Errorf("error receiving message: %w", err)
    }
    
    message := string(messageBytes)
    if message == "NO_WINNERS" {
        return []string{}, nil
    }
    
    return strings.Split(message, ";"), nil
}

func (p *Protocol) Close() error {
    return p.conn.Close()
}