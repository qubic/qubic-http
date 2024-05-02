# The qubic http service

The `qubic-http` service acts as a bridge to the qubic network.  
Some of its features include:
- Checking balances
- Broadcasting transactions
- Tick information
- Block height


## Building from source

```shell
go build -o "server" "./app/server"
```

## Running the service

`qubic-http` requires an instance of [qubic-nodes](https://github.com/qubic/go-qubic-nodes) for fetching reliable node addresses.

### Configuration

The service can be configured either by CLI parameters or environment variables.  
The current configuration will be printed at startup.

#### Required parameters

```shell
QUBIC_API_SIDECAR_SERVER_HTTPS_HOST:          "0.0.0.0:8000"
QUBIC_API_SIDECAR_SERVER_GRPC_HOST:           "0.0.0.0:8001"
QUBIC_API_SIDECAR_SERVER_MAX_TICK_FETCH_URL:  "http://127.0.0.1:8080/max-tick"

QUBIC_API_SIDECAR_POOL_NODE_FETCHER_URL:      "http://127.0.0.1:8080/status"
```

#### Optional parameters
```shell
QUBIC_API_SIDECAR_SERVER_READ_TIMEOUT:        (default: 5s)
QUBIC_API_SIDECAR_SERVER_WRITE_TIMEOUT:       (default: 5s)
QUBIC_API_SIDECAR_SERVER_SHUTDOWN_TIMEOUT:    (default: 5s)

QUBIC_API_SIDECAR_POOL_NODE_FETCHER_TIMEOUT:  (default: 2s)
QUBIC_API_SIDECAR_POOL_INITIAL_CAP:           (default: 5)
QUBIC_API_SIDECAR_POOL_MAX_IDLE:              (default: 20)
QUBIC_API_SIDECAR_POOL_MAX_CAP:               (default: 30)
QUBIC_API_SIDECAR_POOL_IDLE_TIMEOUT:          (default: 15s)
```

### Docker (recommended)
A `docker-compose.yml` file is provided in this repository. You can run it as-is using `docker compose up -d`.  
This will start a `qubic-nodes` service alongside `qubic-http`.
> It may be necessary to configure `qubic-nodes` with an up-to-date list of peers if the service fails to start.

### Standalone
Setting up a `qubic-nodes` instance is outside the scope of this documentation, but you can learn about it on the`qubic-nodes` [GitHub](https://github.com/qubic/go-qubic-nodes).  
We assume that `qubic-nodes` is running locally on port 8080.  

Export the required environment variables:
```shell
export QUBIC_API_SIDECAR_SERVER_HTTPS_HOST="0.0.0.0:8000"
export QUBIC_API_SIDECAR_SERVER_GRPC_HOST="0.0.0.0:8001"
export QUBIC_API_SIDECAR_SERVER_MAX_TICK_FETCH_URL="http://127.0.0.1:8080/max-tick"
export QUBIC_API_SIDECAR_POOL_NODE_FETCHER_URL="http://127.0.0.1:8080/status"
```
You can now start the service:
```shell
./server
```


## Available endpoints

### Information related endpoints

#### /tick-info

```shell
curl http://localhost:8000/tick-info
```
```json
{
  "tickInfo": {
    "tick": 13692242,
    "duration": 7,
    "epoch": 107,
    "initialTick": 13680000
  }
}
```

#### /block-height

```shell
curl http://localhost:8000/block-height
```
```json
{
  "blockHeight": {
    "tick": 13692267,
    "duration": 0,
    "epoch": 107,
    "initialTick": 13680000
  }
}
```

#### /balances/{id}

```shell
curl http://localhost:8000/balances/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF
```
```json
{
  "balance": {
    "id": "PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF",
    "balance": "0",
    "validForTick": 13692275,
    "latestIncomingTransferTick": 0,
    "latestOutgoingTransferTick": 0
  }
}
```

### Transaction related endpoints

#### /broadcast-transaction

```shell
curl -X POST http://localhost:8000/broadcast-transaction -d '{"encodedTransaction": "..."}'
```
```json
{
    "peersBroadcasted": 3,
    "encodedTransaction": "...",
    "transactionId": "oxmqdbynwbisqcjgphyaexhlknmaanipiyxpulatkdjxpdqsqtiovjhcxqkd"
}
```