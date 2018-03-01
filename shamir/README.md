# Shamir Secret Sharing

The Shamir Secret Sharing library is an implementation of the Shamir Secret Sharing scheme. We recommend the Wikipedia page as the starting point for readers that are not familiar with secret sharing. 

* [Wikipedia](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
* [Shamir, Adi (1979), "How to share a secret", Communications of the ACM, 22 (11): 612–613](https://doi.org/10.1145%2F359168.359176)

## Finite Fields

This library requires the use of a prime number to define a finite field from which secrets can be selected. Informally, you cannot use a secret that is greater than, or equal to, this prime number. The test suite uses the first prime number larger than 1024 bits, which we believe is sufficiently large for the majority of secrets.

```
179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859
```

There is no reason why larger prime numbers would not work, however they have not been included in the test suite.

## How it works

To split a secret into many shares, we need to define the following (these will be application dependent):

1. The number of shares to break the secret into, `N`.
2. The number of shares required to reconstruct the secret, `K`.
3. The prime number defining the finite field, `P`.

The SSS library makes use of `math/big` in the Go standard library. This allows us to use numbers outside the range of standard integers. This is useful for encoding complex data structures into a single integer, and for using large prime numbers to define the finite field.

```go
// 100 shares will be created.
N := int64(100)
// 50 shares will be required to reconstruct the secret.
K := int64(50)
// The secret cannot be larger than this number (the first prime greater than 1024 bits).
P, ok := big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
```

### Splitting a secret

We can split a secret into `N` shares using the `Split` function. We will only need `K` shares, or more, to reconstruct the secret later.

```go
secret := big.NewInt(1234)
shares, err := sss.Split(N, K, prime, secret)
```

That's all there is to it. Remember, you should check the value of `err` to make sure it is `nil`.

### Reconstructing a secret

To reconstruct the secret, we can use any `K` or more shares, and the `Join` function. In the example below, we use the first `K` shares.

```go
secret := sss.Join(prime, shares[:K])
```

## Tests

To run the test suite, install Ginkgo.

```sh
go get github.com/onsi/ginkgo/ginkgo
```

Now we can run the tests.

```sh
ginkgo -v
```