import sys


def generate_docker_compose(output_file, num_clients):
    compose_content = """name: tp0
services:
  server:
    container_name: server
    image: server:latest
    volumes:
      - ./server/config.ini:/config.ini 
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - NUM_CLIENTS=""" + str(num_clients) + """
    networks:
      - testing_net
"""
    for i in range(1, num_clients+1):
      compose_content += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{i}.csv:/agency.csv
    entrypoint: /client
    environment:
      - CLI_ID={i}
    networks:
      - testing_net
    depends_on:
      - server
"""

    compose_content += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""
    with open(output_file, 'w') as file:
        file.write(compose_content)

def main():
    output_file = sys.argv[1]
    input_clients = int(sys.argv[2])
    print(f"Creating file {output_file} with {input_clients} clients")

    generate_docker_compose(output_file, input_clients)
    
if __name__ == "__main__":
    main()