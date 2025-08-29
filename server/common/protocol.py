import struct
import socket
import logging
from common import utils


MAX_MESSAGE_FIELDS = 6


def receive_bet_message(client_sock):
    
    message_length = receive_bytes_length(client_sock)
        
    message = receive_message(client_sock, message_length)
    
    return message
        
    
def receive_bytes_length(client_sock):
    try:
        length_bytes = b''
        while len(length_bytes) < 2:
            chunk = client_sock.recv(2-len(length_bytes))
            if not chunk:
                raise ConnectionError("Connection closed while reading length")
            length_bytes += chunk
        return struct.unpack('>H', length_bytes)[0]
    except Exception as e:
        raise Exception(f"Error receiving length: {e}")
    
    
def receive_message(client_sock, message_length):
    try:
        message_bytes = b''
        while len(message_bytes) < message_length:
            chunk = client_sock.recv(message_length - len(message_bytes))
            if not chunk:
                raise ConnectionError("Connection closed while reading message")
            message_bytes += chunk
        return message_bytes.decode('utf-8')
    except Exception as e:
        raise Exception(f"Error receiving message: {e}")
    
    
def send_ack_message(client_sock):
    try:
        ack_byte = b'\x01'
        total_bytes_sent = 0
        while total_bytes_sent < 1:
            bytes_sent = client_sock.send(ack_byte[total_bytes_sent:])
            if bytes_sent == 0:
                raise ConnectionError("Connection closed while sending ACK")
            total_bytes_sent += bytes_sent
    except Exception as e:
        raise Exception(f"Error sending ACK: {e}")
    
def parse_message_to_bet(msg: str) -> list[utils.Bet]:
    """
    Parse a message string to a Bet object and store it
    """
    try:
        lines = msg.strip().split('\n')
        if len(lines) < 1:
            raise ValueError("Empty message")
        bets = []
        for i in range(1, len(lines)):
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
    