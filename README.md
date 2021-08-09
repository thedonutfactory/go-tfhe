# 🍩 go-tfhe
Golang implementation of [TFHE Homomorphic Encryption Library](https://tfhe.github.io/tfhe/)

TFHE, or Fully Homomorphic Encryption Library over the Torus, is a scheme created by [Ilaria Chillotti](https://github.com/ilachill) et al, that implements a fast, fully bootstrapped circuit environment for running programs within the homomorphic realm.

## Show me the Code

The following is a simple fully homomorphic 8-bit integer addition circuit program. As you can see, fully homomorphic encryption constructs and evaluates boolean circuits, just as traditional computing environments do. This allows developers to produce FHE programs using [boolean logic gates](https://en.wikipedia.org/wiki/Logic_gate).

```
import (
  "fmt"
  "github.com/TheDonutFactory/go-tfhe"
)

func main() {
  // generate params
  const nbBits = 8
  var minimumLambda int32 = 100
  params := tfhe.NewDefaultGateBootstrappingParameters(minimumLambda)
  inOutParams := params.InOutParams
  // generate the secret keyset
  keyset := tfhe.NewRandomGateBootstrappingSecretKeyset(params)
  // Encrypt the input
  x := tfhe.NewLweSampleArray(nbBits, inOutParams)
  y := tfhe.NewLweSampleArray(nbBits, inOutParams)
  xBits := toBits(22)
  yBits := toBits(33)
  for i := 0; i < nbBits; i++ {
    tfhe.BootsSymEncrypt(x[i], int32(xBits[i]), keyset)
    tfhe.BootsSymEncrypt(y[i], int32(yBits[i]), keyset)
  }
  // output sum
  sum := tfhe.NewLweSampleArray(nbBits+1, inOutParams)
  fullAdder(sum, x, y, nbBits, keyset)
}

func fullAdder(sum []*LweSample, x []*LweSample, y []*LweSample, nbBits int, keyset *TFheGateBootstrappingSecretKeySet) {
  inOutParams := keyset.Params.InOutParams
  // carry bits
  carry := NewLweSampleArray(2, inOutParams)
  tfhe.BootsSymEncrypt(carry[0], 0, keyset) // first carry initialized to 0
	temp := NewLweSampleArray(3, inOutParams)
  for i := 0; i < nbBits; i++ {
    //sumi = xi XOR yi XOR carry(i-1)
    BootsXOR(temp[0], x[i], y[i], keyset.Cloud) // temp = xi XOR yi
    BootsXOR(sum[i], temp[0], carry[0], keyset.Cloud)
    
    // carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
    BootsAND(temp[1], x[i], y[i], keyset.Cloud)        // temp1 = xi AND yi
    BootsAND(temp[2], carry[0], temp[0], keyset.Cloud) // temp2 = carry AND temp
    BootsXOR(carry[1], temp[1], temp[2], keyset.Cloud)
    BootsCOPY(carry[0], carry[1], keyset.Cloud)
  }
  BootsCOPY(sum[nbBits], carry[0], keyset.Cloud)
}
```

### WTF is a "fully bootstrapped circuit environment for running programs within the homomorphic realm"

In simple terms, fully homomorphic encryption allows two parties, Alice and Bob, to execute programs on each others's computer systems, without knowing the inputs and output of the data. For example, let's say that Bob owns a cloud processing company that crunches health datasets for Alice. Bob has access to cutting-edge bare metal machines with a lot of processing power and is happy to sell Alice that processing power. However, due to HIPPA compliance requirements ( and a general, altruistic respect for an individual's privacy ), Alice cannot actually share the data.

How do we solve this problem with today's cryptography? Well, we can encrypt the data over the wire, send it to Bob's cloud processing company, and then securely hand him the private key to decrypt, process and reencrypt the data. However, we are still giving Bob's prying eyes access to very sensitive, private health data of Alice's customers.

Enter fully homomorphic encryption. Using an FHE cryptographic runtime, Alice can build a special homomorphic software program designed to process her customer data. She gives this special program over to Bob, where he installs it onto one of his meaty servers. Now, here's the magic... Alice can fully encrypt all of her customer data, give it to Bob, who executes the homomorphic software program to process it, returning to Alice the fully encrypted results, all without ever seeing any of her customer's data unencrypted! Alice decrypts the resulting data with the same key she used to encrypt it's inputs, knowing full well that her data was always safe, even when being processed on Bob's servers.

Modern day cryptographic miracle.

### FHE Illustrated

So, if Alice writes a simple program to add two numbers and return the results, classically she would create something like this:

```
function addTwoNumbers(int32 a, int32 b) {
  return a + b
}
```

So if she were outsource Bob to run the program function:

`addTwoNumbers(2, 3)` it would result in `5`

Simple, but Bob and his employees could see the input and output values of the function.

Now, if she used a fully homomorphic runtime environment, she would be able to first encrypt the values of `a` and `b`, pass these encrypted values to Bob, where he executes the homomorphic version of the same function giving back an encrypted result. For illustration purposes, the interaction might look something like this:

1. Alice encrypts a's value of `2` resulting in `1d8b4cf854c`
2. Alice encrypts a's value of `3` resulting in `32c4feed996`
3. Alice asks Bob to run the program function `addTwoNumbers(1d8b4cf854c, 32c4feed996)`. (To Bob, the numbers are encrypted nonsense)
4. Bob gets the result `489f719cad` and returns it to Alice
5. Alice decrypts the result with her key revealing the number to be `5`. At no point did Bob ever see Alice's input or output data, but he performed all of the processing for her. Magic.

## References

[CGGI19]: I. Chillotti, N. Gama, M. Georgieva, and M. Izabachène. TFHE: Fast Fully Homomorphic Encryption over the Torus. In Journal of Cryptology, volume 33, pages 34–91 (2020). [PDF](https://eprint.iacr.org/2018/421.pdf)

[CGGI16]: I. Chillotti, N. Gama, M. Georgieva, and M. Izabachène. Faster fully homomorphic encryption: Bootstrapping in less than 0.1 seconds. In Asiacrypt 2016 (Best Paper), pages 3-33. [PDF](https://eprint.iacr.org/2016/870.pdf)


