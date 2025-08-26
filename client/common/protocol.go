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
    message := fmt.Sprintf("%s|%s|%s|%d|%s|%d", 
        agencyId, bet.FirstName, bet.LastName, 
        bet.DNI, bet.Birthday, bet.BetNumber)

    if len(message) > 65535 {
        return fmt.Errorf("message too long")
    }
    
    messageSize := uint16(len(message))
    result := make([]byte, 2+len(message))
    binary.BigEndian.PutUint16(result[0:2], messageSize)
    copy(result[2:], []byte(message))

    totalWritten := 0
    for totalWritten < len(result) {
        n, err := conn.Write(result[totalWritten:])
        if err != nil {
            return err
        }
        totalWritten += n
    }
    return nil
}
