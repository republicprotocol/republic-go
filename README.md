# Shamir Secret Sharing

The Shamir Secret Sharing (SSS) library is a Go implementation of the Shamir Secret Sharing Scheme. The original paper describing the scheme was published by Adi Shamir in 1979, but the Wikipedia page is a good starting point for readers that are not familiar with secret sharing. 

* [Shamir, Adi (1979), "How to share a secret", Communications of the ACM, 22 (11): 612–613](https://doi.org/10.1145%2F359168.359176)
* [Wikipedia](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)

## Finite Fields

This library requires the use of a prime number to define a finite field from which secrets can be selected. Its test suite uses the first prime number larger than 1024 bits, which we believe is sufficiently large for the majority of secrets. There is no reason why large prime numbers would not work, however larger primes are not included in the test suite.

```
179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859
```

## How it works

To split a secret into many shares, we first need three things (these will be application dependent):

1. The number of shares to break the secret into, `N`.
2. The number of shares required to reconstruct the secret, `K`.
3. The prime number defining the finite field, `P`.

The SSS library makes use of the `math/big` library defined by the Go standard library. This allows us to use numbers outside the standard range of integers. This is useful for encoding complex data in a single integer, and for using large prime numbers to define the finite field.

```go
// 100 shares will be created.
N := int64(100)
// 50 shares will be required to reconstruct the secret.
K := int64(50)
// The secret cannot be larger than this number (the first prime greater than 1024 bits).
P, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
```

Now that we have these values, we can split a secret using the `Split` function.

```go
secret := big.NewInt(1234)
shares, err := Split(N, K, prime, secret)
```

That's all there is to it. Remember, you should also check the value of `err` to make sure it is `nil`. To reconstruct the secret, we can use any `K` or more shares, and the `Join` function. In the example below, we take the first `K` shares and use these to call the `Join` function.

```go
secret, err := Join(prime, shares[:K])
```

That's it! We can now use the Shamir Secret Sharing scheme.

## Republic

The SSS library was developed by the Republic Protocol team. For more information, see our website https://republicprotocol.com.

## Contributors

* Loong loong@republicprotocol.com
* Noah noah@republicprotocol.com