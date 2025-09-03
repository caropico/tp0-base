# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N.º 1:

#### Solución
Se creó el script generar-compose.sh, que permite crear automáticamente un archivo docker-compose.yaml con un servidor y una cantidad configurable de clientes (client1, client2, etc.). Al ejecutarlo, se llama a mi-generador.py, que construye el contenido del Docker Compose y lo guarda en el archivo indicado.

Se puede utilizar el generador de la siguiente manera:

```
./generar-compose.sh <output_file> <n_clients>
```

Para validar la correcta ejecución del script se puede ejecutar:

```
cat <output_file>
```

### Ejercicio N°2:

#### Solución
Se modificó la definición de docker-compose.yaml dentro de mi-generador.py y el Dockerfile del cliente para que los archivos de configuración no queden dentro de la imagen, sino que se inyecten como volúmenes al momento de levantar los contenedores.

De esta forma, al ejecutar el script generador, cualquier cambio en config.ini (servidor) o config.yaml (clientes) se aplica de inmediato sin necesidad de reconstruir la imagen.

Ejemplo de modificación del servicio server:
```
volumes:
  - ./server/config.ini:/config.ini
```

Para validar la correcta implementación se puede ejecutar el script:
```
./generar-compose.sh docker-compose-dev.yaml 1
```
Levantar los contenedores:
```
make docker-compose-up
```
Modificar la configuración ./server/config.ini o ./client/config.yaml

Verificar dentro del contenedor:
```
docker exec -it server cat /config.ini
```

### Ejercicio N°3:

#### Solución
Se creó el script en bash validar-echo-server.sh para comprobar automáticamente que el servidor funciona como un echo server.

El script levanta un contenedor temporal con alpine y netcat dentro de la misma red de Docker que el servidor (tp0_testing_net) y ejecuta el siguiente comando:
```
"echo '$mensaje' | nc server 12345"
```

Para validar el correcto funcionamiento ejecutar:
```
./generar-compose.sh docker-compose-dev.yaml 1
```
```
make docker-compose-up
```
```
./validar-echo-server.sh
```

### Ejercicio N°4:

#### Solución
Se modificaron el **cliente y el servidor** para que terminen de forma _graceful_ al recibir la señal `SIGTERM`. Esto implica que todos los recursos abiertos por la aplicación se cierren correctamente antes de que el proceso principal termine.

- **Servidor (Python)**:  
  - Se utiliza el módulo `signal` para capturar señales del sistema.  
  - Se define un manejador de señal (`signal_handler`) que se ejecuta al recibir `SIGTERM`.
  - Dentro del manejador:
    - Se llama a `server.shutdown()` para detener el loop de aceptación de conexiones.  
    - Se cierra el socket principal (`server._server_socket`).

- **Cliente (Go)**:  
  - Se usan los paquetes `os/signal` y `syscall` para capturar `SIGTERM`.
  - Se crea un canal de señales (signals) y se registra para recibir `syscall.SIGTERM`.
  - En cada iteración del loop de envío de mensajes, se verifica si llegó una señal:
    - Si hay señal, se sale del loop

Para validar el correcto funcionamiento ejecutar:
```
./generar-compose.sh docker-compose-dev.yaml 5
```
```
make docker-compose-up
```
```
docker compose -f docker-compose-dev.yaml stop -t 10
```
```
make docker-compose-logs
```


## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Solución
Cada cliente recibe sus datos de apuesta (NOMBRE, APELLIDO, DOCUMENTO, NACIMIENTO, NUMERO) como variables de entorno inyectadas por Docker Compose.
El cliente construye un ClientBet con esos valores y lo serializa usando el protocolo.

Luego, se implementó un protocolo de comunicación entre el cliente-servidor de la siguiente manera:
  - El cliente envia un mensaje con el siguiente formato. 2 bytes en big endian para el largo del mensaje que se enviara, y luego el mensaje, donde cada valor se encuentra delimitado por `;`.
```
<(2 bytes) largo del mensaje><(N bytes) agencyId;bet.FirstName;bet.LastName;bet.DNI;bet.Birthday; bet.BetNumber>
```

  - El servidor primero lee los bytes del largo del mensaje, y luego procede a leer el mensaje
  - Una vez que finalizo, responde un ACK de un byte para avisarle al cliente que la lectura fue correcta.
  - Una vez recibido el ACK por parte del cliente, se cierra la conexion.

Para evitar problemas de `short read` y `short write` se implementaron bucles de lectura y escritura en el protocolo:
  - Envío (cliente): se asegura que, si una operación de escritura no transmite todos los bytes, el proceso continúa hasta enviar el mensaje completo.

  - Recepción (servidor y cliente): se lee repetidamente hasta obtener:
    - Los 2 bytes iniciales que indican la longitud del mensaje.
    - Todos los bytes del mensaje según esa longitud.
    - El byte ACK de respuesta del servidor.

### Ejercicio N°6:

#### Solución
Se modificó la lógica de clientes y servidor para soportar el envío de apuestas en batchs (chunks de datos).

- Cliente
  - Cada cliente lee sus apuestas desde el archivo `agency-{N}.csv`, inyectado en el contenedor como volumen.
  - El procesamiento se realiza mediante un `CSVBatchProcessor`, que construye lotes de apuestas respetando dos restricciones:
    - Cantidad máxima de apuestas por batch, configurable en `config.yaml` bajo la clave `batch.maxAmount`.
    - Tamaño máximo del paquete de 8 KB, evitando que el mensaje supere este límite.
  - Cada batch se serializa y se envía al servidor. El protocolo se amplió para indicar si el batch enviado es el último:
    - 1 byte para el flag EOF (indica si es el último batch).
    - 2 bytes en big endian con la longitud del mensaje.
    - N bytes con las apuestas, delimitadas por saltos de línea `(\n)`.

- Servidor
  - Recibe mensajes de tamaño variable, parsea todas las apuestas contenidas en el batch y las almacena con store_bets.
  - Al recibir un batch con EOF flag, finaliza la comunicación con el cliente.

### Ejercicio N°7:

#### Solución
Se modificó el protocolo y la lógica de comunicación de la siguiente forma

- Cliente
  - Cada cliente lee sus apuestas desde el archivo `agency.csv` mediante `CSVBatchProcessor`, que arma lotes (BatchResult) de tamaño máximo configurable en `config.yaml`.
  - Al iniciar la conexión, el cliente envía un código de carga de apuestas `(0x02)`. Posteriormente, envía todos los batches con el siguiente formato:
```
<(1 byte) EOF flag><(2 bytes) length><(N bytes) apuestas>
```
  - Por cada batch, el servidor responde con un ACK `(0x01)`
  - Una vez finalizado el envío de apuestas, el cliente inicia una nueva conexión para notificar con el código check winners `(0x03)` y enviar su `agencyId`.
  - Tras el sorteo, recibe un mensaje con el código send winners `(0x04)` y la lista de DNIs de ganadores correspondientes a su agencia.

- Servidor
  - Carga de apuestas `(0x02)`: recibe múltiples batches hasta el EOF flag, donde parsea cada apuesta con `parse_message_to_bet`, y almacena usando `store_bets(...)`. Por cada batch responde con ACK `(0x01)` y cierra la conexión con el cliente.
  - Consulta de ganadores `(0x03)`: guarda la conexión asociada al `agencyId`.
  - Cuando se reciben consultas de todas las agencias, ejecuta el sorteo:
    - Se cargan todas las apuestas con `load_bets(...)`.
    - Se evalúan con `has_won(...)`.
    - Se filtran los ganadores por agencia y se responde a cada cliente con su lista de DNIs usando `send_winners_message`.
    - El formato de envío de ganadores 

    ```
    <(1 byte) código de mensaje><(2 bytes) longitud><(N bytes) DNIs ganadores>
    ```
    Donde los DNIs de los ganadores se encuentran delimitador por `;`. En caso de que no haya ganadores se envía `NO WINNERS`.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

#### Solución
Se modificó el servidor de la siguiente forma para el manejo de la concurrencia:

El servidor utiliza hilos `(threading.Thread)` para manejar múltiples clientes en paralelo.
  - Por cada conexión aceptada, se lanza un hilo que ejecuta `_handle_client`.
  - Se mantiene una lista de threads activos `(_active_threads)` protegida con `_threads_lock(threading.Lock)` para asegurar el manejo correcto de concurrencia.
  - Cada hilo procesa de manera independiente el envío de apuestas o la consulta de ganadores.

Para evitar condiciones de carrera se emplean locks `(threading.Lock)`:
  - `_waiting_agencies_lock`: protege el acceso al diccionario `_waiting_agencies`, donde se guardan las agencias que notificaron su intención de participar del sorteo.
  - `_bets_storage_lock`: pensado para proteger operaciones de almacenamiento de apuestas.

Se utiliza una barrera `(threading.Barrier)` inicializada con el número total de clientes.
  - Cada cliente que consulta ganadores se bloquea en la barrera hasta que todas las agencias lleguen.
  - Cuando la última agencia se registra, se dispara la acción `__do_sorteo_and_send_results`, que realiza el sorteo y envía resultados.

Se implementó un `_shutdown_event` para señalizar cuando el servidor debe detenerse.
  - El loop principal (run) se interrumpe si `_shutdown_event` está activo.
  - Antes de cerrar, se llama a `_wait_for_threads`:
    - Se aborta la barrera (para desbloquear hilos que esperaban).
    - Se realiza join con timeout a todos los hilos activos.

Se decidió utilizar multithreading debido a que las tareas del servidor son principalmente I/O. Usar procesos separados implicaría un overhead adicional de memoria y comunicación inter-procesos, innecesarios para este caso donde la carga de CPU es mínima. A su vez, los hilos son más livianos y nos permiten un graceful shutdown más simple.
