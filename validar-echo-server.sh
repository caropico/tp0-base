#!/bin/bash
mensaje="hola"
respuesta=$(echo "$mensaje" | docker run --rm -i --network tp0_testing_net alpine nc server 12345)

if [ "$mensaje" = "$respuesta" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi