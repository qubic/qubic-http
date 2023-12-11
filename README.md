<br /><br />

qubic-http

<br /><br />

## Building the app binary:
```bash
$ go build -o "server" "./app/server"
```

## Running the server:
```bash
$ ./server --qubic-node-ips="65.21.10.217;148.251.184.163"
2023/11/20 20:05:05 main: Config :
--web-host=0.0.0.0:8080
--web-read-timeout=5s
--web-write-timeout=5s
--web-shutdown-timeout=5s
--qubic-node-ips=[65.21.10.217 148.251.184.163]
--qubic-node-port=21841
2023/11/20 20:05:05 main: API listening on 0.0.0.0:8080
```

## Sending requests to the server:
```bash
$ curl "localhost:8080/v1/identities/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF"
{
  "public_key": "4f27dc1b6a1a76d479833e5f1bed0d6d77c705a0290a632de94794dbee670dfa",
  "incoming_amount": 1479299940,
  "outgoing_amount": 1479289940,
  "number_of_incoming_transfers": 125981,
  "number_of_outgoing_transfers": 2612,
  "latest_incoming_transfer_tick": 10894487,
  "latest_outgoing_transfer_tick": 11330319,
  "siblings": [
    "2a5af1c66af3ef4a294e09f27aed030d3faeffb9a1910012468b9c7f3e46bd9f",
    "1231fc9579e4a19c0568ea2ddcf5e6133397111f93a3910b6036dfc299580799",
    "17f6c2d15b2cdc00e0bb33e961c7c376b72d9f477c0e047bdd1524a0c4e771b1",
    "39217be72dad406e99910c70c082a32c59092b818d70c0aaab9b0bcc79a94f5c",
    "419c90c65a584a9a28b5bf0fdb1d43e7e62be60513177cccb4a590a14987324d",
    ...
  ]
}
```

## Docker usage:
```bash
$ docker build -t ghcr.io/qubic/qubic-http:latest .
$ docker run -p 8080:8080 -e QUBIC_API_SIDECAR_QUBIC_NODE_IPS="65.21.10.217;148.251.184.163" ghcr.io/qubic/qubic-http:latest
$ curl "localhost:8080/v1/identities/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF"
{
  "public_key": "4f27dc1b6a1a76d479833e5f1bed0d6d77c705a0290a632de94794dbee670dfa",
  "incoming_amount": 1479299940,
  "outgoing_amount": 1479289940,
  "number_of_incoming_transfers": 125981,
  "number_of_outgoing_transfers": 2612,
  "latest_incoming_transfer_tick": 10894487,
  "latest_outgoing_transfer_tick": 11330319,
  "siblings": [
        "2a5af1c66af3ef4a294e09f27aed030d3faeffb9a1910012468b9c7f3e46bd9f",
    "1231fc9579e4a19c0568ea2ddcf5e6133397111f93a3910b6036dfc299580799",
    "17f6c2d15b2cdc00e0bb33e961c7c376b72d9f477c0e047bdd1524a0c4e771b1",
    "39217be72dad406e99910c70c082a32c59092b818d70c0aaab9b0bcc79a94f5c",
    "419c90c65a584a9a28b5bf0fdb1d43e7e62be60513177cccb4a590a14987324d",
    ...
  ]
}
```
