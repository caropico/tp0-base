import struct
import socket

def receive_bet_message(client_sock):
    try:
        length_bytes = b''
        while len(length_bytes) < 2:
            chunk = client_sock.recv(2 - len(length_bytes))
            if not chunk:
                raise ConnectionError("Connection closed while reading length")
            length_bytes += chunk
        
        message_length = struct.unpack('>H', length_bytes)[0]
        
        message_bytes = b''
        while len(message_bytes) < message_length:
            chunk = client_sock.recv(message_length - len(message_bytes))
            if not chunk:
                raise ConnectionError("Connection closed while reading message")
            message_bytes += chunk
        
        message = message_bytes.decode('utf-8')
        
        return message
        
    except Exception as e:
        raise Exception(f"Error receiving message: {e}")