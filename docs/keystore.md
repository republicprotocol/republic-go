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
    "id": "0b6e5b62-404f-8d43-a45f-68cbcd984607",
    "version": 3,
    "ecdsa": {
        "address": "8MHXRr47PxLLKNJQPg5EuaSp6oSBqB",
        "crypto": {
            "cipher": "aes-128-ctr",
            "cipherText": "13bd547552d003339f408b58c7ac8a5c9f6fcf6ecbb0f08b15d19314a4a721e8",
            "cipherParams": {
                "iv": "c1248cb6c28399d4b98f035146c60e38"
            },
            "kdf": "scrypt",
            "kdfParams": {
                "dkLength": 32,
                "n": 262144,
                "p": 1,
                "r": 8,
                "salt": "ed4e9fc3a617211650d0ce4047a7359048a9b7d5a1cf7a69cb83a205959cc9a1"
            },
            "mac": "29733978a65792ff2bec5305c6bf1c063787c8ee6ea0aec7c120ff00f28d1111"
        }
    },
    "rsa": {
        "publicKey": "XIYFjElJbj4p0w5Yy4bg7vtu+l46TwNcMvulLxFK8hMs7AnRgC0D6zRC7RQgMWQQ==",
        "crypto": {
            "cipher": "aes-128-ctr",
            "cipherText": "fcf6ecbb0f08b1ac8a5c9f3bd5475525c9f65d19316fcf6ecbb0f08b114a4a721e8",
            "cipherParams": {
                "iv": "d4b98e384c1246c603998cb6c28f0351"
            },
            "kdf": "scrypt",
            "kdfParams": {
                "dkLength": 32,
                "n": 262144,
                "p": 1,
                "r": 8,
                "salt": "5d4e947a7359048172116f7a69fc3a63a200d0ce40cb8a9b7d5a1ce5959cc9a1"
            },
            "mac": "2287c8ee6e5792ff111bec5309733978a5c6bf1c0a0aed16637c7c120ff00f28"
        }
    }
}
```

**Plain text**

```json
{
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
        "d": "DlWyCkIe5TRI1vJWyi3htW8mSZRt9dfngwvO4VxFRec=",
        "x": "3UqFqdND9Obcp4gliv8Gcoih3yhg/Lj5bVU1RNU4uMI=",
        "y": "yxHMmNcjiKRzRuuPk/Kx1avFpPe9dPV4SsIWM5SH4io="
    },
    "rsa": {
        "d": "XIYFjElJbj4p0w5Yy4bg7vtu+l46TwNcMvulLxFK8hy6hR0PnICzLbkgNJkxbC2nGNBurHuQgF1sfDavouHPhi8jJ/JfpPoEoU0NHoMlnnCNwh4suwU8ZcdA6B+wT3p8E8j8jdo0pCBPciJ0JmZttKR7uHxsb9mKGzAFYxg/ejpkcU25hpxIbfp1yjCZK4S59S2jlAZ/PBuGFr5+YoUaBJONrz0QdkZQm6/To9Gp+sQUw5n142WHZTKFE26dARPM7BM8/UpImK3FNtnGPI3d6FFh6Z6JMtEhfLp6ODmqcyVfTgGI6mWCHKWFgabPOrBrMs7AnRgC0D6zRC7RQgMWQQ==",
        "e": 65537,
        "n": "3CV+WTa0AY7ChVFg+Yc3crGvwHnVvc2xPnSpoErkmFkLZKvNFQuFmVuCnk/pckXsSDFwxBd3vv7a/K79fI5V4OheMpyo2fJCbcMD5kasqgeyjFeT+Z+02zqUG14iyjGgg4En18NZur+RiHTxSTAlIY3LxoagRaaSEyDxrb4lf1F5q1WLTW9c5fey8+1W7dBzP27uzmcqFrGqwEEFZrRmVg1vaoWd8lj3N/OkRL5da4e4LKZBP2WevNK9NLK2AWLudWWb2kQA1ZYMSl27iGszgqEdm0Pm2l8x2zib71DK7vYKwW8qwSp7DnTCVI2xachDi2EcPStsj6PSbbCNKzsZlw==",
        "primes": [
            "757Po8tivBRZMyIElNQR1xWXtInSrsUlYOZSr3HNQM+TCUIIgUl8hixCW9Zeytc6/8us8SVl7Gceftl5HTxhkTKlGs5GWR4tTuW8dIdEz/yHLYiAykgAyFox/VvDPpNfzbTH1D66kxlIl3vAkZvFIxEbuNUEFQeguC4WPr1rYS0=",
            "6zHnIfvvfY3bIw6+3aXFsnOtghQB+XkBE+x6Ers/EJ3F/u3fVgeZYocwqD1AJDJ32Ysk8dh2HcApLlCh6Tkjb9mwSPqsVO7kaOJfISyIz22SOLBjUVmwiHR/YTWbpo9Wt7Ok9/gnXogsXtPtm9AY7qFv9hIx5qjvz/E2607D+FM="
        ]
    }
}
```

### Ecdsa

An Ecdsa key can be encoded as JSON. By itself, an Ecdsa key does not provide support for passphrase encryption. To use passphrase encryption, store the Ecdsa key in a Keystore.

**JSON**

```json
{
    "curveParams": {
        "b": "Bw==",
        "bits": 256,
        "n": "/////////////////////rqu3OavSKA7v9JejNA2QUE=",
        "name": "s256",
        "p": "/////////////////////////////////////v///C8=",
        "x": "eb5mfvncu6xVoGKVzocLBwKb/NstzijZWfKBWxb4F5g=",
        "y": "SDradyajxGVdpPv8DhEIqP0XtEimhVQZnEfQj/sQ1Lg="
    },
    "d": "DlWyCkIe5TRI1vJWyi3htW8mSZRt9dfngwvO4VxFRec=",
    "x": "3UqFqdND9Obcp4gliv8Gcoih3yhg/Lj5bVU1RNU4uMI=",
    "y": "yxHMmNcjiKRzRuuPk/Kx1avFpPe9dPV4SsIWM5SH4io="
}
```

### Rsa

An Rsa key can be encoded as JSON. By itself, an Rsa key does not provide support for passphrase encryption. To use passphrase encryption, store the Rsa key in a Keystore.

**JSON**

```json
{
    "d": "XIYFjElJbj4p0w5Yy4bg7vtu+l46TwNcMvulLxFK8hy6hR0PnICzLbkgNJkxbC2nGNBurHuQgF1sfDavouHPhi8jJ/JfpPoEoU0NHoMlnnCNwh4suwU8ZcdA6B+wT3p8E8j8jdo0pCBPciJ0JmZttKR7uHxsb9mKGzAFYxg/ejpkcU25hpxIbfp1yjCZK4S59S2jlAZ/PBuGFr5+YoUaBJONrz0QdkZQm6/To9Gp+sQUw5n142WHZTKFE26dARPM7BM8/UpImK3FNtnGPI3d6FFh6Z6JMtEhfLp6ODmqcyVfTgGI6mWCHKWFgabPOrBrMs7AnRgC0D6zRC7RQgMWQQ==",
    "e": 65537,
    "n": "3CV+WTa0AY7ChVFg+Yc3crGvwHnVvc2xPnSpoErkmFkLZKvNFQuFmVuCnk/pckXsSDFwxBd3vv7a/K79fI5V4OheMpyo2fJCbcMD5kasqgeyjFeT+Z+02zqUG14iyjGgg4En18NZur+RiHTxSTAlIY3LxoagRaaSEyDxrb4lf1F5q1WLTW9c5fey8+1W7dBzP27uzmcqFrGqwEEFZrRmVg1vaoWd8lj3N/OkRL5da4e4LKZBP2WevNK9NLK2AWLudWWb2kQA1ZYMSl27iGszgqEdm0Pm2l8x2zib71DK7vYKwW8qwSp7DnTCVI2xachDi2EcPStsj6PSbbCNKzsZlw==",
    "primes": [
        "757Po8tivBRZMyIElNQR1xWXtInSrsUlYOZSr3HNQM+TCUIIgUl8hixCW9Zeytc6/8us8SVl7Gceftl5HTxhkTKlGs5GWR4tTuW8dIdEz/yHLYiAykgAyFox/VvDPpNfzbTH1D66kxlIl3vAkZvFIxEbuNUEFQeguC4WPr1rYS0=",
        "6zHnIfvvfY3bIw6+3aXFsnOtghQB+XkBE+x6Ers/EJ3F/u3fVgeZYocwqD1AJDJ32Ysk8dh2HcApLlCh6Tkjb9mwSPqsVO7kaOJfISyIz22SOLBjUVmwiHR/YTWbpo9Wt7Ok9/gnXogsXtPtm9AY7qFv9hIx5qjvz/E2607D+FM="
    ]
}
```