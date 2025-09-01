import socket
import logging

from common.protocol import Protocol
from common import utils



class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_running = False

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        self._is_running = True
        try: 
            while self._is_running:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    self.__handle_client_connection(client_sock)
        except OSError:
            logging.info('action: server_loop_interrupted | result: success')
            

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        bets_list = []
        protocol = Protocol(client_sock)
        keep_running = True
        while keep_running:
            try:
                msg, is_eof = protocol.receive_bet_message()
                addr = client_sock.getpeername()
                """logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')"""
                bets_list = protocol.parse_message_to_bet(msg)
                utils.store_bets(bets_list)
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets_list)}")
                protocol.send_ack_message()
                if is_eof:
                    logging.info(f"action: eof_received | result: success | ip: {addr[0]}")
                    keep_running = False
            except OSError as e:
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets_list)}")
        protocol.close()
            


    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    
    def shutdown(self):
        self._is_running = False
        if self._server_socket:
            self._server_socket.close()
        logging.info('action: close_server_socket | result: success')
