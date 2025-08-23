#!/bin/bash
respuesta=$(echo "hola" | docker run --rm -i --network tp0-base_testing_net alpine nc server 12345)

if [ "$mensaje" = "$respuesta" ]; then
    echo "Echo server funciona correctamente"
else
    echo "Error: enviado '$mensaje', recibido '$respuesta'"
fi