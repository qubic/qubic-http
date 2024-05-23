# The qubic http service

The `qubic-http` service acts as a bridge to the qubic network.  
Some of its features include:
- Checking balances
- Broadcasting transactions
- Tick information
- Block height


## Building from source

```shell
go build -o "server" "./app/grpc_server"
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

#### /balances/{identity}

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

### Asset related endpoints

#### /assets/{identity}/issued

```shell
curl localhost:8000/assets/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB/issued
```
```json
{
  "issuedAssets": [
    {
      "data": {
        "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
        "type": 1,
        "name": "RANDOM",
        "numberOfDecimalPlaces": 0,
        "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
      },
      "info": {
        "tick": 14057739,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
        "type": 1,
        "name": "QX",
        "numberOfDecimalPlaces": 0,
        "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
      },
      "info": {
        "tick": 14057739,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
        "type": 1,
        "name": "QTRY",
        "numberOfDecimalPlaces": 0,
        "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
      },
      "info": {
        "tick": 14057739,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
        "type": 1,
        "name": "QUTIL",
        "numberOfDecimalPlaces": 0,
        "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
      },
      "info": {
        "tick": 14057739,
        "universeIndex": 0
      }
    }
  ]
}
```

#### /assets/{identity}/owned

```shell
curl localhost:8000/assets/IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC/owned
```
```json
{
  "ownedAssets": [
    {
      "data": {
        "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 2,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 0,
        "numberOfUnits": "2",
        "issuedAsset": {
          "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
          "type": 1,
          "name": "RANDOM",
          "numberOfDecimalPlaces": 0,
          "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
        }
      },
      "info": {
        "tick": 14057652,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 2,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 1,
        "numberOfUnits": "1",
        "issuedAsset": {
          "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
          "type": 1,
          "name": "QX",
          "numberOfDecimalPlaces": 0,
          "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
        }
      },
      "info": {
        "tick": 14057652,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 2,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 9284980,
        "numberOfUnits": "186685601",
        "issuedAsset": {
          "issuerIdentity": "QWALLETSGQVAGBHUCVVXWZXMBKQBPQQSHRYKZGEJWFVNUFCEDDPRMKTAUVHA",
          "type": 1,
          "name": "QWALLET",
          "numberOfDecimalPlaces": 0,
          "unitOfMeasurement": [0, -48, 0, -48, 35, 24, 21]
        }
      },
      "info": {
        "tick": 14057652,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 2,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 2741973,
        "numberOfUnits": "10",
        "issuedAsset": {
          "issuerIdentity": "TFUYVBXYIYBVTEMJHAJGEJOOZHJBQFVQLTBBKMEHPEVIZFXZRPEYFUWGTIWG",
          "type": 1,
          "name": "QFT",
          "numberOfDecimalPlaces": 0,
          "unitOfMeasurement": [84, 79, 75, 69, 78, 0, 0]
        }
      },
      "info": {
        "tick": 14057652,
        "universeIndex": 0
      }
    }
  ]
}
```

#### /assets/{identity}/possessed

```shell
curl localhost:8000/assets/IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC/possessed
```
```json
{
  "possessedAssets": [
    {
      "data": {
        "possessorIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 3,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 9707976,
        "numberOfUnits": "2",
        "ownedAsset": {
          "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
          "type": 3,
          "padding": 0,
          "managingContractIndex": 1,
          "issuanceIndex": 9707976,
          "numberOfUnits": "2",
          "issuedAsset": {
            "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
            "type": 1,
            "name": "RANDOM",
            "numberOfDecimalPlaces": 0,
            "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
          }
        }
      },
      "info": {
        "tick": 14057759,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "possessorIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 3,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 9707978,
        "numberOfUnits": "1",
        "ownedAsset": {
          "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
          "type": 3,
          "padding": 0,
          "managingContractIndex": 1,
          "issuanceIndex": 9707978,
          "numberOfUnits": "1",
          "issuedAsset": {
            "issuerIdentity": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB",
            "type": 1,
            "name": "QX",
            "numberOfDecimalPlaces": 0,
            "unitOfMeasurement": [0, 0, 0, 0, 0, 0, 0]
          }
        }
      },
      "info": {
        "tick": 14057759,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "possessorIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 3,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 9707980,
        "numberOfUnits": "186685601",
        "ownedAsset": {
          "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
          "type": 3,
          "padding": 0,
          "managingContractIndex": 1,
          "issuanceIndex": 9707980,
          "numberOfUnits": "186685601",
          "issuedAsset": {
            "issuerIdentity": "QWALLETSGQVAGBHUCVVXWZXMBKQBPQQSHRYKZGEJWFVNUFCEDDPRMKTAUVHA",
            "type": 1,
            "name": "QWALLET",
            "numberOfDecimalPlaces": 0,
            "unitOfMeasurement": [0, -48, 0, -48, 35, 24, 21]
          }
        }
      },
      "info": {
        "tick": 14057759,
        "universeIndex": 0
      }
    },
    {
      "data": {
        "possessorIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
        "type": 3,
        "padding": 0,
        "managingContractIndex": 1,
        "issuanceIndex": 9707982,
        "numberOfUnits": "10",
        "ownedAsset": {
          "ownerIdentity": "IGJQYTMFLVNIMEAKLANHKGNGZPFCFJGSMVOWMNGLWCZWKFHANHGCBYODMKBC",
          "type": 3,
          "padding": 0,
          "managingContractIndex": 1,
          "issuanceIndex": 9707982,
          "numberOfUnits": "10",
          "issuedAsset": {
            "issuerIdentity": "TFUYVBXYIYBVTEMJHAJGEJOOZHJBQFVQLTBBKMEHPEVIZFXZRPEYFUWGTIWG",
            "type": 1,
            "name": "QFT",
            "numberOfDecimalPlaces": 0,
            "unitOfMeasurement": [84, 79, 75, 69, 78, 0, 0]
          }
        }
      },
      "info": {
        "tick": 14057759,
        "universeIndex": 0
      }
    }
  ]
}
```