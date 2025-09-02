import socket
import logging
import threading

from common.protocol import Protocol

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
        self._waiting_agencies_lock = threading.Lock()
        self._bets_storage_lock = threading.Lock()        
        self._barrier = threading.Barrier(self._num_clients, action=self.__do_sorteo_and_send_results)
        self._shutdown_event = threading.Event()
        self._active_threads = []
        self._threads_lock = threading.Lock()

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
            while self._is_running and not self._shutdown_event.is_set():
                client_sock = self.__accept_new_connection()
                if client_sock:
                    client_thread = threading.Thread(
                        target=self._handle_client,
                        args=(client_sock,)
                    )
                    client_thread.daemon = False
                with self._threads_lock:
                        self._active_threads.append(client_thread)
            
                client_thread.start()
                logging.info(f'action: thread_started | result: success | thread_id: {client_thread.ident} | is_alive: {client_thread.is_alive()}')
        except OSError:
            logging.info('action: server_loop_interrupted | result: success')
        finally:
            self._wait_for_threads()
            
    def _handle_client(self, client_sock):
        """Handle client connection and ensure thread cleanup"""
        try:
            self.__handle_client_connection(client_sock)
        finally:
            with self._threads_lock:
                current_thread = threading.current_thread()
                if current_thread in self._active_threads:
                    self._active_threads.remove(current_thread)
            

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        protocol = Protocol(client_sock)
        try:
            code = protocol.receive_code_message()
            if code == LOAD_BET_MESSAGE_CODE:
                self.__handle_bets_loads(protocol)
            elif code == CHECK_FOR_WINNERS_MESSAGE_CODE:
                self.__handle_check_for_winners(protocol)
        except OSError as e:
            logging.error(f"action: apuesta_recibida | result: fail")

            
    def __handle_bets_loads(self, protocol):
        """Maneja múltiples batches hasta recibir EOF"""
        keep_running = True
        while keep_running and not self._shutdown_event.is_set():
            try:
                msg, is_eof = protocol.receive_bet_message()
                bets_list = protocol.parse_message_to_bet(msg)
                with self._bets_storage_lock:
                    utils.store_bets(bets_list)
                logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets_list)}")
                protocol.send_ack_message()
                
                if is_eof:
                    logging.info("action: batch_session_completed | result: success")
                    keep_running = False
            except Exception as e:
                logging.error(f"action: batch_processing | result: fail | error: {e}")
                break
        protocol.close()
            
    def __handle_check_for_winners(self, protocol):
        try:
            if self._shutdown_event.is_set(): 
                return
            msg = protocol.receive_check_winners()
            agency_id = msg.strip()
            with self._waiting_agencies_lock:       
                self._waiting_agencies[agency_id] = protocol 
                if len(self._waiting_agencies) == self._num_clients:
                    self.__do_sorteo_and_send_results()
            try:
                self._barrier.wait(timeout=30.0)
            except threading.BrokenBarrierError:
                logging.error("action: barrier_broken | result: shutdown_in_progress")
            except Exception as e:
                logging.error(f"action: barrier_wait | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: handle_check_winners | result: fail | error: {e}")
            protocol.close()
        
    def __do_sorteo_and_send_results(self):    
        logging.info(f"action: sorteo | result: success")  
        
        all_bets = list(utils.load_bets())
        
        for agency_id, protocol in self._waiting_agencies.items():
            try:
                agency_bets = [bet for bet in all_bets if str(bet.agency) == agency_id]
                winners = [bet for bet in agency_bets if utils.has_won(bet)]
                winner_dnis = [str(winner.document) for winner in winners]
                
                protocol.send_winners_message(winner_dnis)
                logging.info(f"action: send_winners | result: success | agency_id: {agency_id} | winners: {len(winners)}")
                
            except Exception as e:
                logging.error(f"action: send_winners | result: fail | agency_id: {agency_id} | error: {e}")
            finally:
                protocol.close()
            
        self._waiting_agencies.clear()
        
        
    def _wait_for_threads(self):
        """Wait for all active threads to finish"""
        logging.info("action: waiting_for_threads | result: in_progress")
        try:
            if self._barrier:
                self._barrier.abort()
        except:
            pass
        
        with self._threads_lock:
            threads_to_wait = self._active_threads.copy()
        
        for thread in threads_to_wait:
            try:
                thread.join(timeout=5.0)
                if thread.is_alive():
                    logging.warning(f"action: thread_join | result: timeout | thread_id: {thread.ident}")
            except Exception as e:
                logging.error(f"action: thread_join | result: fail | error: {e}")
        
        logging.info("action: threads_cleanup | result: success")

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
        """
        Ensure graceful server shutdown
        """

        self._is_running = False
        self._shutdown_event.set()
        if self._server_socket:
            self._server_socket.close()
        logging.info('action: close_server_socket | result: success')