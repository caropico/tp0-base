#!/bin/bash
mensaje="hola"
respuesta=$(docker run --rm --network tp0-base_testing_net alpine sh -c "echo '$mensaje' | nc server 12345")

if [ "$mensaje" = "$respuesta" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi