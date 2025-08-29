package common

import (
    "fmt"
    "encoding/binary"
    "net"
    "strings"
)

func SendBetMessage(conn net.Conn, bet []ClientBet, agencyId string) error {
    bet_msg := make([]string,0, len(bet))
    for _, b := range bet {
        bet_msg = append(bet_msg, fmt.Sprintf("%s;%s;%s;%d;%s;%d", 
            agencyId,b.FirstName, b.LastName, b.DNI, b.Birthday, b.BetNumber))
    }
    
    message := strings.Join(bet_msg, "\n")

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
