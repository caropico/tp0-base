package common

import (
    "fmt"
    "encoding/binary"
    "net"
)

type Bet struct {
    Agency    string
    FirsName      string
    Last   string
    DNI  int
    Birthday string
    Number    int
}


func SendBetMessage(conn net.Conn, bet ClientBet, agencyId string) error {
    message := fmt.Sprintf("%s;%s;%s;%d;%s;%d", 
        agencyId, bet.FirstName, bet.LastName, 
        bet.DNI, bet.Birthday, bet.BetNumber)

    if len(message) > 65535 {
        return fmt.Errorf("message too long")
    }
    
    messageSize := uint16(len(message))
    result := make([]byte, 2+len(message))
    binary.BigEndian.PutUint16(result[0:2], messageSize)
    copy(result[2:], []byte(message))

    err := SendAll(conn, result)

    return err
}

func SendAll(conn net.Conn, data []byte) error {
    totalWritten := 0
    for totalWritten < len(data) {
        n, err := conn.Write(data[totalWritten:])
        if err != nil {
            return err
        }
        totalWritten += n
    }
    return nil
}

func ReceiveAll(conn net.Conn) (byte,error) {
    ack := make([]byte,1)
    totalRead := 0
    for totalRead < 1 {
        n, err := conn.Read(ack[totalRead:])
        if err != nil {
            return 0,err
        }
        totalRead += n
    }

    return ack[0],nil
}
