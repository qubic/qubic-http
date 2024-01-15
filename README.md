<br /><br />

qubic-http

<br /><br />

## Building the app binary:
```bash
$ go build -o "server" "./app/server"
```

## Running the server:
```bash
$ ./server --qubic-node-ips="65.21.10.217;148.251.184.163" --opensearch-host="http://93.190.139.223:9200"
2023/11/20 20:05:05 main: Config :
--web-host=0.0.0.0:8080
--web-read-timeout=5s
--web-write-timeout=5s
--web-shutdown-timeout=5s
--qubic-node-ips=[65.21.10.217 148.251.184.163]
--qubic-node-port=21841
--opensearch-host=http://93.190.139.223:9200
2023/11/20 20:05:05 main: API listening on 0.0.0.0:8080
```

## Docker usage:
```bash
$ docker build -t ghcr.io/qubic/qubic-http:latest .
$ docker run -p 8080:8080 -e QUBIC_API_SIDECAR_QUBIC_NODE_IPS="65.21.10.217;148.251.184.163" -e QUBIC_API_SIDECAR_OPENSEARCH_HOST="http://93.190.139.223:9200" ghcr.io/qubic/qubic-http:latest
```

## Getting identity info:
```bash
$ curl "localhost:8080/v1/identities/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF"
{
  "public_key": "4f27dc1b6a1a76d479833e5f1bed0d6d77c705a0290a632de94794dbee670dfa",
  "tick": 10894487,
  "balance": 10000,
  "incoming_amount": 1479299940,
  "outgoing_amount": 1479289940,
  "number_of_incoming_transfers": 125981,
  "number_of_outgoing_transfers": 2612,
  "latest_incoming_transfer_tick": 10894487,
  "latest_outgoing_transfer_tick": 11330319,
  "siblings": [
    "2a5af1c66af3ef4a294e09f27aed030d3faeffb9a1910012468b9c7f3e46bd9f",
    ...
  ]
}
```

## Getting tx by id:
```bash
$ curl "localhost:8080/v1/tx/qdnwqignkprgrbkdlbpsrprikjagdnzaafsbxngbtfczcfjnsvgelivfxoem"
{
  "bxid": "bacab66e7cbbb950c8b80facef94c5f0bf478ee3ab8ac270f541b1aba71cb8a7",
  "utime": "1704973083",
  "epoch": "91",
  "tick": "11963307",
  "type": "0",
  "src": "AMITKEPEADJMEBOKSXLFNNCQLQTCAAWIDNCKTDZGYEUCYCNOEEDJHMACPPTG",
  "dest": "DCHWKUZLNYTDYBFZTELHWEDQDQCBBKTKYUITIKCDKACHDUBFYRHGDZSFQLBG",
  "amount": "7357247517",
  "extra": "",
  "sig": "b6f330e9713047321de9ff3506078726a06119115ab27e8d805a9ebd2828a8189e96bd1a31d1a72702ba1ee5abebd130cf43e25ee68109f01c6f9cc42d100f00"
}
```

## Getting bx by id:
```bash
$ curl "localhost:8080/v1/bx/bacab66e7cbbb950c8b80facef94c5f0bf478ee3ab8ac270f541b1aba71cb8a7"
{
  "utime": "1704973097",
  "epoch": "91",
  "tick": "11963307",
  "type": "1",
  "src": "AMITKEPEADJMEBOKSXLFNNCQLQTCAAWIDNCKTDZGYEUCYCNOEEDJHMACPPTG",
  "dest": "DCHWKUZLNYTDYBFZTELHWEDQDQCBBKTKYUITIKCDKACHDUBFYRHGDZSFQLBG",
  "amount": "7357247517"
}
```

## Getting bx by id:
```bash
$ curl "localhost:8080/v1/status"
{
  "epoch": 91,
  "bxid": [
    91,
    12017682
  ],
  "txid": [
    91,
    12017336
  ],
  "quorum": [
    91,
    12017688
  ],
  "tick": [
    91,
    12017336
  ]
}
```

## Getting latest tick info:
```bash
$ curl "localhost:8080/v1/tick-info"
{
  "tick_duration": 2,
  "epoch": 86,
  "tick": 11368700,
  "number_of_aligned_votes": 0,
  "number_of_misaligned_votes": 0
}
```

## Getting tick data:
```bash
$ curl "localhost:8080/v1/tick-data/11356544"
{
  "computor_index": 420,
  "epoch": 86,
  "tick": 11356544,
  "millisecond": 0,
  "second": 4,
  "minute": 16,
  "hour": 15,
  "day": 12,
  "month": 12,
  "year": 23,
  "hex_union_data": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "hex_timelock": "7262190d095c8507d9077688d1892683236857c8c7bd556a1c26cc47f098d778",
  "transaction_digests": [
    "25a9a37da318564cae50aa67bf9af04948685736d2710338dc71467426c8a6ef",
    ...
  ],
  "contract_fees": null,
  "signature": "ae1e9cbbae7dc7895659d44c2828c77bca9e61b2de8c2a07e4df2e19dcad40b2fdf79b878c6decf8868d557105cc0a552688cc9e03f08e6b0d8e846844fa2700"
}
```

## Getting tick transactions:
```bash
$ curl "localhost:8080/v1/tick-transactions/11922277"
[
  {
    "source_public_key": "cf1fc50a6ecd7706aee39a90efb2e94b623714f715049419c940528dd939d9bd",
    "destination_public_key": "86ef6705ee15ca15915dff635a10a04b55614ee05f0108e63a94238b009bca7e",
    "amount": 200000,
    "tick": 11922277,
    "input_type": 0,
    "input_size": 0,
    "hash": "qttuixqwdopyggptwulggmltsikavawlsssmzrlyobrcenqzndjlhivbzdtn"
  },
  {
    "source_public_key": "3852a585cf4f1966d594a14d946b4976978889e86949822796a46acf9c915036",
    "destination_public_key": "9e1a100cfb556def7bcc6252e47ddf0985428637c3d1b3caa16f33fd98438d94",
    "amount": 0,
    "tick": 11356544,
    "input_type": 0,
    "input_size": 32,
    "hash": "piypdgfhmjbppdcezurkaiybuvtbitnpculrvgqdebkqgywbejxifuzenvxk"
  },
]
```

## Send raw tx:
```bash
$ curl "localhost:8080/v1/send-raw-tx" --json '{"hex_raw_tx": "C872E68E1C0ECCCE3BC6A87BC32E187C59BBA99AB81D7CC37E7D22F7423672A70E4EAF16A2218457BA8B46991B5CCA63E65AE65FF65C575A06743E40E8DA982A0100000000000000C60DAE0000000000A1C0B21A5C15D72275F7968D30A4F0520075F85A0232E180A5FC6C0137CC414F402404CF40773F444A25BCF30B6455B18A18FF7DD105F3223EECA8C566781A00"}'
```



