import struct
import socket
import logging
from common import utils


MAX_MESSAGE_FIELDS = 6


class Protocol:
    def __init__(self, client_socket):
        self.client_sock = client_socket
        
    def receive_bet_message(self):
        """
        Receive messages from the client socket and return it as a string
        """
        eof_flag = self._receive_eof_flag()
        message_length = self._receive_bytes_length()
            
        message = self._receive_message(message_length)
        
        is_eof = (eof_flag == 0x01)
        return (message, is_eof)
    
    def _receive_eof_flag(self):
        """
        Receive the EOF flag (1 byte)
        """
        try:
            eof_bytes = b''
            while len(eof_bytes) < 1:
                chunk = self.client_sock.recv(1 - len(eof_bytes))
                if not chunk:
                    raise ConnectionError("Connection closed while reading EOF flag")
                eof_bytes += chunk
            return eof_bytes[0]
        except Exception as e:
            raise Exception(f"Error receiving EOF flag: {e}")
            
        
    def _receive_bytes_length(self):
        """
        Receive bytes length of the message from the client
        """
        try:
            length_bytes = b''
            while len(length_bytes) < 2:
                chunk = self.client_sock.recv(2-len(length_bytes))
                if not chunk:
                    raise ConnectionError("Connection closed while reading length")
                length_bytes += chunk
            return struct.unpack('>H', length_bytes)[0]
        except Exception as e:
            raise Exception(f"Error receiving length: {e}")
        
        
    def _receive_message(self, message_length):
        """
        Receive the message of bets from the client
        """
        try:
            message_bytes = b''
            while len(message_bytes) < message_length:
                chunk = self.client_sock.recv(message_length - len(message_bytes))
                if not chunk:
                    raise ConnectionError("Connection closed while reading message")
                message_bytes += chunk
            return message_bytes.decode('utf-8')
        except Exception as e:
            raise Exception(f"Error receiving message: {e}")
        
        
    def send_ack_message(self):
        """
        Send an acknowledgment byte to the client
        """
        try:
            ack_byte = b'\x01'
            total_bytes_sent = 0
            while total_bytes_sent < 1:
                bytes_sent = self.client_sock.send(ack_byte[total_bytes_sent:])
                if bytes_sent == 0:
                    raise ConnectionError("Connection closed while sending ACK")
                total_bytes_sent += bytes_sent
        except Exception as e:
            raise Exception(f"Error sending ACK: {e}")
        
    def parse_message_to_bet(self,msg: str) -> list[utils.Bet]:
        """
        Parse a message string to a Bet object and store it
        """
        try:
            lines = msg.strip().split('\n')
            if len(lines) < 1:
                raise ValueError("Empty message")
            bets = []
            for i in range(0, len(lines)):
                line = lines[i].strip()
                if not line:
                    continue
                    
                fields = line.split(';')
                if len(fields) != MAX_MESSAGE_FIELDS:
                    raise ValueError(f"Line {i}: Expected {MAX_MESSAGE_FIELDS} fields, got {len(fields)}")
                
                bet = utils.Bet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5])
                bets.append(bet)
            return bets
        except Exception as e:
            raise Exception(f"Error parsing bets: {e}")
        
    def close(self):
        if self.client_sock:
            self.client_sock.close()
    