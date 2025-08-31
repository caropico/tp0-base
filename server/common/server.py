import socket
import logging

from common import protocol
from common import utils

LOAD_BET_MESSAGE_CODE = 0x02
CHECK_FOR_WINNERS_MESSAGE_CODE = 0x03

class Server:
    def __init__(self, port, listen_backlog,num_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._is_running = False
        self._waiting_agencies = {}
        self._num_clients = num_clients

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
        try:
            code = protocol.receive_code_message(client_sock)
            if code == LOAD_BET_MESSAGE_CODE:
                self.__handle_bets_loads(client_sock)
            elif code == CHECK_FOR_WINNERS_MESSAGE_CODE:
                self.__handle_check_for_winners(client_sock)
        except OSError as e:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets_list)}")
            
    def __handle_bets_loads(self, client_sock):
        try:
            msg = protocol.receive_bet_message(client_sock)
            addr = client_sock.getpeername()
            """logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')"""
            bets_list = protocol.parse_message_to_bet(msg)
            utils.store_bets(bets_list)
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets_list)}")
        except OSError as e:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets_list)}")
        finally:
            protocol.send_ack_message(client_sock)
            client_sock.close()
            
    def __handle_check_for_winners(self, client_sock):
        try:
            msg = protocol.receive_bet_message(client_sock)
            agency_id = msg.strip()
            addr = client_sock.getpeername()
            logging.info(f'action: check_winners | result: success | ip: {addr[0]} | agency_id: {agency_id}')            
            self._waiting_agencies[agency_id] = client_sock
            
            self.__validate_amount_of_agencies_ready()
            
        except Exception as e:
            logging.error(f"action: handle_check_winners | result: fail | error: {e}")
            client_sock.close()
        
    def __validate_amount_of_agencies_ready(self):
        agencies_ready = len(self._waiting_agencies)
        
        if agencies_ready == self._num_clients:
            logging.info(f"action: sorteo | result: success")  
            
            all_bets = list(utils.load_bets())
            
            for agency_id, client_sock in self._waiting_agencies.items():
                try:
                    agency_bets = [bet for bet in all_bets if str(bet.agency) == agency_id]
                    
                    winners = [bet for bet in agency_bets if utils.has_won(bet)]
                    winner_dnis = [str(winner.document) for winner in winners]
                                
                    protocol.send_winners_message(client_sock, winner_dnis)
                    
                    logging.info(f"action: send_winners | result: success | agency_id: {agency_id} | winners: {len(winners)}")
                    
                except Exception as e:
                    logging.error(f"action: send_winners | result: fail | agency_id: {agency_id} | error: {e}")
                finally:
                    client_sock.close()
            
            self._waiting_agencies.clear()
            
        else:
            logging.info(f"action: waiting_agencies | result: in_progress | agencies_ready: {agencies_ready}/{self._num_clients}")

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
