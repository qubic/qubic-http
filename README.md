<br /><br />

<div align="center">
  <h1>qubic/api-sidecar</h1>
  <h4> REST service that makes it easy to interact with
    <a href="https://docs.qubic.world/">Qubic</a> network nodes.</h4>
  <h4>The sidecar api is used to interact with the Qubic network and it currently supports fetching information related to identity balances, with ongoing development to add additional functionalities exposed by Qubic nodes.</h4>
  <h4>It's using the <a href="https://github.com/0xluk/go-qubic"> go-qubic</a> sdk to interact with the Qubic network</h4>
</div>

<br /><br />

## Building the app binary:
```bash
$ go build -o "server" "./app/server"
```

## Running the server:
```bash
$ ./server --qubic-node-ip="65.21.10.217"
2023/11/20 20:05:05 main: Config :
--web-host=0.0.0.0:8080
--web-read-timeout=5s
--web-write-timeout=5s
--web-shutdown-timeout=5s
--qubic-node-ip=65.21.10.217
--qubic-node-port=21841
2023/11/20 20:05:05 main: API listening on 0.0.0.0:8080
```

## Sending requests to the server:
```bash
$ curl "localhost:8080/identities/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF"
{"Entity":{"PublicKey":[79,39,220,27,106,26,118,212,121,131,62,95,27,237,13,109,119,199,5,160,41,10,99,45,233,71,148,219,238,103,13,250],"IncomingAmount":10000,"OutgoingAmount":0,"NumberOfIncomingTransfers":125980,"NumberOfOutgoingTransfers":1047,"LatestIncomingTransferTick":10703131,"LatestOutgoingTransferTick":10701088},"Tick":10871666,"SpectrumIndex":14427983,"Siblings":[]}
```

## Docker usage:
```bash
$ docker build -t qubic-api-sidecar .
$ docker run -p 8080:8080 -e -e QUBIC_API_SIDECAR_QUBIC_NODE_IP=65.21.10.217 qubic-api-sidecar
$ curl "localhost:8080/identities/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF"
{"Entity":{"PublicKey":[79,39,220,27,106,26,118,212,121,131,62,95,27,237,13,109,119,199,5,160,41,10,99,45,233,71,148,219,238,103,13,250],"IncomingAmount":10000,"OutgoingAmount":0,"NumberOfIncomingTransfers":125980,"NumberOfOutgoingTransfers":1047,"LatestIncomingTransferTick":10703131,"LatestOutgoingTransferTick":10701088},"Tick":10871666,"SpectrumIndex":14427983,"Siblings":[]}
```
