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
$ curl "localhost:8080/v1/address/PKXGRCNOEEDLEGTLAZOSXMEYZIEDLGMSPNTJJJBHIBJISHFFYBBFDVGHRJQF"
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

## Getting status:
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
  "signature": "ae1e9cbbae7dc7895659d44c2828c77bca9e61b2de8c2a07e4df2e19dcad40b2fdf79b878c6decf8868d557105cc0a552688cc9e03f08e6b0d8e846844fa2700",
  "potentialBx": [
    {
      "index": 35,
      "dest": "YOANRGAZNOPBCBYTYQFBFYEARPKDZEEIBYZORWQTTDOKCIEGLZZMACGDOEKA",
      "amount": "1203390386"
    }
  ]
}
```

## Getting quorum data by tick
```bash
$ curl "localhost:8080/v1/quorum/12023688"
{
  "computor": 542,
  "epoch": 91,
  "tick": 12023688,
  "time": [
    0,
    26,
    5,
    7,
    16,
    1,
    24
  ],
  "prevRTD": "12250904539969688494",
  "saltRTD": "5465134535200475419",
  "digests": [
    "a2cecdc3110116a22c654501e23261480ccb07f52e34f5ab19861da21bb658f1",
    "34a98e98d893e3a22db055ded34e71702d1f6b38bf43627dc592a199fe819ef6",
    "c9df65d2c2053559c6259fc9f58e88086ad84aedacd9d5bd6981b2dba09ad769",
    "f7e2fd47adb07b334c2c30ee6fb5d844ee4d707f17e8bc3b84d960447b96f654",
    "e217aa9e67ee3e5319a8feb5004d0ebc26a61cc2434b87d5810ea59dc18512c7",
    "d1216643b9279c9d0b6278d05c22e1bd95616c283e4a3f85eb6d9a8e6a6cdd5a",
    "07d621e3bc092473eb0e06c7fed4b3cb7f82cb77cfa3372002252f131c01f68f",
    "9a4076ba1e30e76a5e7ecdcf1f58bb8ca8ec4a62251a006e0ad463bd2e104c26"
  ],
  "sig": "56ea4b6f4fe2cb2e1728a3f6aa9b6b21343b17fdd42b046e67e0a7f9ade479079925e53f5f880d93e123539a8124d8c603d54e7d5c21f7278320bfa31b7c1600",
  "diffs": [
    {
      "computor": 3,
      "saltRTD": "13747147356733467208",
      "digests": [
        "47da2aeebb85fade0c4abb953d56eebee0f294b99e14f5b9efb4714219fca992",
        "33c8368920dd6d59c436a09ff88e06b98320676db6be94e0d1f2e870a980ab08",
        "fe1744396133c5dc6a02ad9170db69d6cd404ba7a7381ba595e6e5a71242ed82",
        "9a4076ba1e30e76a5e7ecdcf1f58bb8ca8ec4a62251a006e0ad463bd2e104c26"
      ],
      "sig": "a0f99dd2ae753da54edacfb4d6dc5022f1e7e33c44e1ee51eb0714781dc1dc6b6a5c7c560c813f98173de48a09f58c03e1276238ca8fa149ca3b162c9db50600"
    },
    {
      "computor": 4,
      "saltRTD": "16517841960992476753",
      "digests": [
        "3a9d7318e20b91bc6e82d5b15349ecb683c47e58caecd59ca9849d70c1b910cf",
        "7ccbf8bfa4c1ce6d6147ebe5d03b16b8a5b783f032104e7259a04dcae0e36567",
        "fb21ff0b75f90fe1aa81ce54cce91b64eee439ed87316ab64a7816a935a0c710",
        "9a4076ba1e30e76a5e7ecdcf1f58bb8ca8ec4a62251a006e0ad463bd2e104c26"
      ],
      "sig": "514e5699ddf35c02351d6643b540573db103135b2f64a606bfddad5985f48d43756ebe9eba92a7cbe7d01a6134397aba5a3ca18793386301bfb09ac6e5fe1300"
    }
  ],
  "numvotes": 603
}
```

## Getting list computors per epoch
```bash
$ curl "localhost:8080/v1/computors/90"
{
  "epoch": "90",
  "pubkeys": [
    "LVMMXNAABKUVYGFUSPUTVPVABARALIUZGNLQJZLXEEOAIKBNGDKUIZICAUFE",
    "WJWPMUQJURGNNBWDGXJSJFPSHCTCGPWFHJKBKICWHBBYPNHSCEKISCJANXAL",
    "CKXIZTMVOVZJMBXNKPNCGYKDPVICAMWJLNDMMWZMCEJJSMLNGGWTRXOFZPSI",
    "AFBSWFLTLNVAKEGWXNVKLTOCEPCAOEGEGJTHQWVONERINRARRQKIMESFNOGM",
  ],
  "sig": "d1a036931ba8066e5817492ccf0af146f091857e83f2f2c388d11fb7b6c47b0891ec7aa46dbcd0d047d48a155c7f0e962a8e3f8d438425ee27ad88deafcb0d00"
}
```

## Send raw tx:
```bash
$ curl "localhost:8080/v1/send-raw-tx" --json '{"hex_raw_tx": "C872E68E1C0ECCCE3BC6A87BC32E187C59BBA99AB81D7CC37E7D22F7423672A70E4EAF16A2218457BA8B46991B5CCA63E65AE65FF65C575A06743E40E8DA982A0100000000000000C60DAE0000000000A1C0B21A5C15D72275F7968D30A4F0520075F85A0232E180A5FC6C0137CC414F402404CF40773F444A25BCF30B6455B18A18FF7DD105F3223EECA8C566781A00"}'
```



