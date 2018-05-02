# Keystore

## Keystore Specification

The JSON encoding for Keystores, Ecdsa keys, and Rsa keys.

### Keystore

A Keystore can be encoded as JSON, in an encrypted format, or a plain-text format.

**Encrypted**

The encrypted JSON encoding uses scrypt to encrypt the Ecdsa and Rsa keys using a passphrase.

Both keys are marshaled into JSON, and the resulting bytes are encrypted using a passphrase. This is done because it allows us to decrypt and unmarshal the keys directly back into JSON objects, which are simpler to work than binary blobs. Both keys are encrypted individually, using the same passphrase. 

```json
{
    "ecdsa": {
        "address": "8MKaREsC8WsWt6quzDPrifQuz3bi8v",
        "crypto": {
            "cipher": "aes-128-ctr",
            "cipherText": "e943912f330ce579fa6ef24ec67e5ae6229f76402d945dd655edb211c41fc8bb3086bcc586c60e8e56aeeddbe5b6cda0188f2f657b46a9afa669585c98b46edf12db68738b1b100744103cf2a10cabe16ddf7098f876e9cd3e63b1d2e8aacc4fc454b627048ef24d4a05d66b187a44a20fa29b8e89e3212e82fc92b9e0d481983d42754eb7878165489202556b790ff5946fb933b4ea3551544db5dc988af3eaeb02ca25e0237a3b4e21eebaa1d89adea01c79a5337c3c18593d6cfc959b233f6bde59b14e72472a9d5c05468261fb3be207d2d914862d0cb866702261b3d73fd53b95e4a02e388cd4e1b3baa3f3ead77043310df2be515a254a3ee35440d956fd7526b276a9659cd71be9294194b7b52a2f0dad3f6c3693702af121f82634506ab328cf92bf0f7091cec6610682c60d3972961ac87de319a43582e55e848071e1b9f7981962442aad31d54aad801c47111ce9fb4fc964428a265ee6fb8346843d9ca7c4b0f3288743585fb0c7db861695fa34c902b8a41250e43ede5ff3b89b637adbc0aa646aec538b597d18caf84feb9abc522a60113a9f40",
            "cipherParams": {
                "iv": "742ff38d9b9bbaf7873c4e29d7e9e94a"
            },
            "kdf": "scrypt",
            "kdfParams": {
                "dkLength": 32,
                "n": 262144,
                "p": 1,
                "r": 8,
                "salt": "f7aebd20d936138a13c5038b3252ceef6cd92e18a37c72e8ba9b6e25d1b2bf91"
            },
            "mac": "2630fbae36c5e57270bc6700724e165930fe937f71d17e674b3824c308e44fcf"
        }
    },
    "id": "7480e750-3256-43be-b8f4-17037258534b",
    "rsa": {
        "publicKey": "AAAAAAABAAHiFaLEpLghepkqeZZWNfP0Dk3OWUPrmbjNFdZlkEr8vQ3qrNBTy9tAK967OYMia8CWBBeZK/fSrx6aI/UCrX3uaF2UdHhE6wE1S/RDckKbT/3nuT6ZxdGSt+UP75VH2MqpJ6QgcGc0oGXrkGrap10e5ANIpeI3A6c/EKBtzBSXd3jla4I68a8CICygmaylDuM8Pfx4SziNvYC7kcZSPDgkjkfIg/OPGvNSJh9hJDO4QUSoMDhsRvolAPU8DmJWS3nPdCHDj2VNLV+LZacDX4XSiAxraEPs/OM82FmOBirxSjxkhRSth2jsgO6J4wKyNRSDTR1EB3kihIktDhbucpTX",
        "crypto": {
            "cipher": "aes-128-ctr",
            "cipherText": "d3f55f59ac98ba895f7ead06174b8541f5a8554061560933d031fdabaff9c7a87404c6c54246212e01416b6278e8a6e6434aa9c25f92d6e9dbbd29e2a7506989134a891818acc38bde3d4cb202a778a911a625055b0e7bce834b76a526dfdee03b2fe23bae5347c50633983a7604f4e1ab9898f5a89ec654f0f487384b4b0c95d98f26885dbf1fd86b676b7b8d1d4dd77679eaaea4cb5997de54556bccc1163fdf43e77113d76689fe89dbe70040eaf836751d14bea528ff384138838ac8fb1c07d8415c9cf3b4d4d1221d4763d993e9ecb5120571b991d4a508c22da46c90a665f25a86784953324dcf78f569fd4b563d2331083b4952ce0d6f4f5bd93443eb1222bc7545c3d0d61a5ee5fe63d6f5aab8ca64a00c7fa50bfd212771b1ad9d0519a42f5ec52d75098cb7d09e75d20d4d53250d70a02fac68bd7dc24d88d838db7a2d91faf941b13bff6c0e31ff4b893e4ec3dbb5b9f64ee9084eab33e3e2c35a77c54baacd46714e496963c59584cbf27bd733a06a38ba56c2470bead392f20554bb40b724faf3c98012d1e5f0a2287093c7dd865dcd97bd92f65a926b93d16680439cf26ba1bb6a3db636224fae48b3a3adc74e05fcdc5d91ff23571136f9a55e00f0799a4949d891ff5fa7eaa320582c92ba28cfa2a66cf14f899acce84428c7f46cb86afc002e9f5584921c3d23f96e2f90d36c5ef39ad6b57df7249d551dc6980ad9c8a59bcb44a80d5ecb52b3220feaded3c171d8499586f4c97662a68e6a2cabbfcba60352cc39116d1218db788a1656054d836307ee265a9ed2df01ec8331eb35fa04e5bbbf58e9a8c0e44e3ba54628c2702d74b44ba3ca5b29f0ab6585683abe828bcc8c0ff887cfcd25e2289f1df8dfe2a3d4cf179635a6cc20bf01b9ce46708f66a599ab7889d87a5ea169958450580340920f21603ea5ed339034616daee3217e1e663d21e4977906a95d22b112a2c2cd766104e3be3ec7fc2453e51af7aeef80eb4460297363782c0744fc4e0eeae2c4f326211bddca23a74cabde5c1ad77e7aac62822504d7d948f8d7d044369de18c6def7691e3e7247e41c6991b7189ec88b88891d91d8fc7b59ab08d6a016a1479f3b3a50f1c757f0b94987cedc6be5d9dad996f28231a88918aa6dd35342f2aff491916886a6dd470e7f7f2c0f338ee62f7e89db9b68f93670a6cfe2494a31c72dc394cf4c25a95cdff758b2146d26c067babbe277e24763d515f99a6702b3bf1bd15e1faca00b97d4d768a89aa9aa2f4854adc60ad1bd3d361a2b828b5e3b4692e8a8a63a1f5e7bdc67fd7952c182c4a9671f68cc17d171a03c5703d0b1d9f05ec8d9278c8efbdc4cd2c32941edd69c6e414766626d549dba77eec965d5e627c440808ad2a29ab9ba1651d478651ac64dc727ee8c3df721702ae0ca0c408e8562d8c1adfb6e84db63bcbec913070dbae8fb05a8def60e6bd17a22da75ece0a26bf94fc9cf69b1f64dcf6d4cad37ac5aa753480e2bf36757daf124bd9",
            "cipherParams": {
                "iv": "058a75a02d50e4421725456e0ae0d3e2"
            },
            "kdf": "scrypt",
            "kdfParams": {
                "dkLength": 32,
                "n": 262144,
                "p": 1,
                "r": 8,
                "salt": "063a0dc0c3971362feaa0d92121cadaae34c789d81c171381b5f6e7308bf67ee"
            },
            "mac": "8ee9cceb2880d52992bf4951699a1a90ba98da7a20a5f5b3593189746c48354e"
        }
    },
    "version": "3"
}
```

**Plain text**

```json
{
    "id": "eb1c2636-35ba-4384-b947-aad3adc10156",
    "version": "3",
    "ecdsa": {
        "curveParams": {
            "b": "Bw==",
            "bits": 256,
            "n": "/////////////////////rqu3OavSKA7v9JejNA2QUE=",
            "name": "s256",
            "p": "/////////////////////////////////////v///C8=",
            "x": "eb5mfvncu6xVoGKVzocLBwKb/NstzijZWfKBWxb4F5g=",
            "y": "SDradyajxGVdpPv8DhEIqP0XtEimhVQZnEfQj/sQ1Lg="
        },
        "d": "pZd8XNd7aEBIizV9q88AE29yHLecTskHdocwjzekoP4=",
        "x": "SWVfXtZWaOY2ziVLQIFoJK81YS3RBuCuoYpElWYv7A4=",
        "y": "hddhN4yU+r6flKfxZ0hPAOQWjS6TA5eJOn6W0cgFFR0="
    },
    "rsa": {
        "d": "bJoLAC8nWIvOiBSTsLJZIscFs4YPt+cC9IuKhCShdBnwQXmTQJn2rY8V1TrwkbbGbmSEPLpy1KZwYRuD9S8qfyTj0YZ9/sOPKlYCXCqcbbwjjWR6oT9EcMEFMTzDeffRtHdRpIgZTII5i99K4NpKEhNcLh94RwGhj/oaVy6ZXENwRGJso9Do85/NsCvH13M8ClnwRm4NIqK4ANRyx84kM6gPR3ZjKlmUI3zDtdvf7Yu27GKC1p1S8Ku/PZYO11+V7TN+2EdfPe9KTb6pMnsIt5qrl+yupgYNmkkYHi+icSFuHyT3yVTlS4kSNNVlJdw2EWwz+4C6exUF6iV1luBiUQ==",
        "e": 65537,
        "n": "uE2wi8EZwiiko06UO86xMDGxxIycKqEoUJlc80DbFMV3zytfOHli1fGfzTNBHVy9uw0Tne9D8Q2mPWkhWULAk1HPA6KAPYcZZtnBTcPWvkHLiOPGEt2UeitkrvdmdpNSS3H3kR7/CUOknOAAAD8TIkFMNDU72Z4G00dAPmzaMZjdYyOAoKrkYIktC4mDOpMeJrfm3Cs9eoVYvm6oTiQF8J+5CS/F1roq6t225sxykpM8nUkehO7vAQOxehIeAKCr+BRNBNpg3nK856yv97+esEKq4rzIGSmcjEc77R0TCbJ6wAGo2/pMe2RHMrtTHhGJOglMehYCSes8/i2b1sXk/w==",
        "primes": [
            "8OGrxN7jIZHxxld3sRPQRSZOxVbk2JWPVxzkgkc5iykE6Q0s5uL5qvkU5PCoX/893ZKhd04P6QOqIZcz6G5PMjZQuRJ36zOVhdPNP4xGGbm0ZLvN6A++brJf8IsrUuxp8UddfaKDSMMJcO08+B3bUlsbFnEkhwFmpTeflHouhgU=",
            "w971gYM+uHLrjgAfFDHnbhdfr+IeehpMk8p6QTRRlmkt1pOsJU7NrNDinMd9D+a3QKCAVbxujuR5IU4amM/cACKrx4WuJzKR9fXdCH7bP/onGuZ7240G3v1uNLPrrJEiF/OuovTBz/iWddHVNY09+IVDsJTI30pFc4v8HLtBCjM="
        ]
    }
}
```

### Ecdsa

An Ecdsa key can be encoded as JSON. By itself, an Ecdsa key does not provide support for passphrase encryption. To use passphrase encryption, store the Ecdsa key in a Keystore.

**JSON**

```json
{
    "address": "8MKaREsC8WsWt6quzDPrifQuz3bi8v",
    "curveParams": {
        "b": "Bw==",
        "bits": 256,
        "n": "/////////////////////rqu3OavSKA7v9JejNA2QUE=",
        "name": "s256",
        "p": "/////////////////////////////////////v///C8=",
        "x": "eb5mfvncu6xVoGKVzocLBwKb/NstzijZWfKBWxb4F5g=",
        "y": "SDradyajxGVdpPv8DhEIqP0XtEimhVQZnEfQj/sQ1Lg="
    },
    "d": "pZd8XNd7aEBIizV9q88AE29yHLecTskHdocwjzekoP4=",
    "x": "SWVfXtZWaOY2ziVLQIFoJK81YS3RBuCuoYpElWYv7A4=",
    "y": "hddhN4yU+r6flKfxZ0hPAOQWjS6TA5eJOn6W0cgFFR0="
}
```

- `curveParams` — An object defining the parameters of the elliptic curve. Currently, an s256 secp256k1 is the only curve supported.
- `address` — Republic Protocol address generated from the public key.
- `d` — Big integer encoded as big-endian bytes, used for the private key.
- `x` — Big integer encoded as big-endian bytes, used for the public key.
- `y` — Big integer encoded as big-endian bytes, used for the public key.

### Rsa

An Rsa key can be encoded as JSON. By itself, an Rsa key does not provide support for passphrase encryption. To use passphrase encryption, store the Rsa key in a Keystore.

**JSON**

```json
{
    "d": "bJoLAC8nWIvOiBSTsLJZIscFs4YPt+cC9IuKhCShdBnwQXmTQJn2rY8V1TrwkbbGbmSEPLpy1KZwYRuD9S8qfyTj0YZ9/sOPKlYCXCqcbbwjjWR6oT9EcMEFMTzDeffRtHdRpIgZTII5i99K4NpKEhNcLh94RwGhj/oaVy6ZXENwRGJso9Do85/NsCvH13M8ClnwRm4NIqK4ANRyx84kM6gPR3ZjKlmUI3zDtdvf7Yu27GKC1p1S8Ku/PZYO11+V7TN+2EdfPe9KTb6pMnsIt5qrl+yupgYNmkkYHi+icSFuHyT3yVTlS4kSNNVlJdw2EWwz+4C6exUF6iV1luBiUQ==",
    "e": 65537,
    "n": "uE2wi8EZwiiko06UO86xMDGxxIycKqEoUJlc80DbFMV3zytfOHli1fGfzTNBHVy9uw0Tne9D8Q2mPWkhWULAk1HPA6KAPYcZZtnBTcPWvkHLiOPGEt2UeitkrvdmdpNSS3H3kR7/CUOknOAAAD8TIkFMNDU72Z4G00dAPmzaMZjdYyOAoKrkYIktC4mDOpMeJrfm3Cs9eoVYvm6oTiQF8J+5CS/F1roq6t225sxykpM8nUkehO7vAQOxehIeAKCr+BRNBNpg3nK856yv97+esEKq4rzIGSmcjEc77R0TCbJ6wAGo2/pMe2RHMrtTHhGJOglMehYCSes8/i2b1sXk/w==",
    "primes": [
        "8OGrxN7jIZHxxld3sRPQRSZOxVbk2JWPVxzkgkc5iykE6Q0s5uL5qvkU5PCoX/893ZKhd04P6QOqIZcz6G5PMjZQuRJ36zOVhdPNP4xGGbm0ZLvN6A++brJf8IsrUuxp8UddfaKDSMMJcO08+B3bUlsbFnEkhwFmpTeflHouhgU=",
        "w971gYM+uHLrjgAfFDHnbhdfr+IeehpMk8p6QTRRlmkt1pOsJU7NrNDinMd9D+a3QKCAVbxujuR5IU4amM/cACKrx4WuJzKR9fXdCH7bP/onGuZ7240G3v1uNLPrrJEiF/OuovTBz/iWddHVNY09+IVDsJTI30pFc4v8HLtBCjM="
    ]
}
```

- `e` — An integer, used for the private key.
- `n` — Big integer encoded as big-endian bytes, used for the private key.
- `d` — Big integer encoded as big-endian bytes, used for the public key.
- `primes` — An array of big integer encoded as big-endian bytes, used for the private key.