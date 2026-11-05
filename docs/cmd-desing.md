## Flock agent managment

### Start new flock agent 
flock agent start --service-name 'payment-service' --gossip-port 4200 --api-port 8081

### Start a new flock agent in deatch mode
flock agent start --service-name 'payment-service' --gossip-port 4200 --api-port 8081 -d 

### List all availbale agents in host
flock agent list (or ls)

### Stop flock agent
flock agent stop 324dfg3

### Show all agent specific inforation to the clinet
flock agent inspect 324dfg3

### Show logs from the specific agent
flock agent logs 324dfg3


## Flock cluster initializaiton
flock seed \
  --gossip-port 4202 \
  --api-port 8082 \
  --detach \
  192.168.1.10:4200 \
  192.168.1.11:4200

The following command initialize a new seed node in cluster that will be joined to other seed nodes and eventually discover all cluster topology
