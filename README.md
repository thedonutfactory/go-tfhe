<img src="docs/images/gopher.png" alt="FHE Gopher" width="200"/>

# Fully Homomorphic Encryption for Gophers 🍩

go-tfhe is the Golang implementation of [TFHE Homomorphic Encryption Library](https://tfhe.github.io/tfhe/)

TFHE, or Fully Homomorphic Encryption Library over the Torus, is a scheme developed by [Ilaria Chillotti](https://github.com/ilachill) et al, that implements a fast, fully bootstrapped circuit environment for running programs within the homomorphic realm. ( for the uninitiated, a quick rundown of FHE is [here](#eli5) )

## Show me the Code

The following snippet is a simple fully homomorphic 8-bit integer addition circuit. As you can see, fully homomorphic encryption constructs and evaluates boolean circuits, just as traditional computing environments do. This allows developers to produce FHE programs using [boolean logic gates](https://en.wikipedia.org/wiki/Logic_gate).

```golang
import (
  "fmt"
  "github.com/thedonutfactory/go-tfhe"
)

func fullAdder(sum []*LweSample, x []*LweSample, y []*LweSample, nbBits int, 
               keyset *TFheGateBootstrappingSecretKeySet) {
  inOutParams := keyset.Params.InOutParams
  // carry bit
  carry := NewLweSampleArray(2, inOutParams)
  BootsSymEncrypt(carry[0], 0, keyset)
  temp := NewLweSampleArray(3, inOutParams)
  for i := 0; i < nbBits; i++ {
    //sum[i] = x[i] XOR y[i] XOR carry(i-1)
    BootsXOR(temp[0], x[i], y[i], keyset.Cloud) // temp = xi XOR yi
    BootsXOR(sum[i], temp[0], carry[0], keyset.Cloud)
    
    // carry = (x[i] AND y[i]) XOR (carry(i-1) AND (x[i] XOR y[i]))
    BootsAND(temp[1], x[i], y[i], keyset.Cloud)        // temp1 = xi AND yi
    BootsAND(temp[2], carry[0], temp[0], keyset.Cloud) // temp2 = carry AND temp
    BootsXOR(carry[1], temp[1], temp[2], keyset.Cloud)
    BootsCOPY(carry[0], carry[1], keyset.Cloud)
  }
  BootsCOPY(sum[nbBits], carry[0], keyset.Cloud)
}
```

### FHE Illustrated - like literally illustrated, with cute gophers

In spite of strong advances in confidential computing technologies, critical information is encrypted only temporarily – while not in use – and remains unencrypted during computation in most present-day computing infrastructures. Fully homomorphic encryption addresses this flaw by providing a mechanism for computation on fully encrypted data.

For example, let's say that Bob owns a cloud processing company that crunches health datasets for Allan. Bob has access to cutting-edge bare metal machines with a lot of processing power and is happy to sell Allan that processing power. However, due to HIPAA compliance requirements ( and a general, altruistic respect for an individual's privacy ), Allan cannot actually share the data.

How do we solve this problem with today's cryptography? Well, we can encrypt the data over the wire, send it to Bob's cloud processing company, and then securely hand him the private key to decrypt, process and reencrypt the data. Again, this only provides protection while the data is in transit or at rest. During computation, the data must be decrypted first. This still allows Bob access to very sensitive, private health data of Allan's customers.

<img src="docs/images/enc1.png" alt="FHE Gopher"/>

Enter fully homomorphic encryption. Using an FHE cryptographic runtime, Allan can build a special homomorphic software program designed to process her customer data. He gives this special program over to Bob, where he installs it onto one of his meaty servers. Now, here's the magic: Allan can fully encrypt all of his customer data, give it to Bob, who executes the homomorphic software program to process it (because he's now a clearly a wizard), returning to Allan the fully encrypted results, all without ever seeing any of his customer's data unencrypted! Allan decrypts the resulting data with the same key she used to encrypt it's inputs, knowing full well that his data was always safe, even when being processed on Bob's servers.

<img src="docs/images/enc2.png" alt="FHE Gopher"/>

Homomorphic encryption means that the data is never decrypted, yet it is still able to be processed by a third party. Modern day cryptographic miracle.

### FHE Workflow

So, if Allan writes a simple program to add two numbers and return the results, classically he would create something like this:

```solidity
function add(int8 a, int8 b) {
  return a + b
}
```

So if he were outsource Bob to run the program function:

`add(2, 3)` it would result in `5`

Simple, but Bob and his employees could see the input and output values of the function.

Now, if he used a fully homomorphic runtime environment, he would be able to first encrypt the values of `a` and `b`, pass these encrypted values to Bob, where he executes the homomorphic version of the same function giving back an encrypted result. For illustration purposes, the interaction might look something like this:

1. Allan encrypts a's value of `2` resulting in `1d8b4cf854c`
2. Allan encrypts a's value of `3` resulting in `32c4feed996`
3. Allan asks Bob to run the program function `add(1d8b4cf854c, 32c4feed996)`. (To Bob, the numbers are encrypted nonsense)
4. Bob gets the result `489f719cad` and returns it to Allan
5. Allan decrypts the result with her key revealing the number to be `5`. At no point did Bob ever see Allan's input or output data, but he performed all of the processing for her. Magic.

## Potential FHE Use Cases

* Private Cloud Computing on public Clouds ( FHE on AWS, Azure or GCP is secure )
* Distributed, trustless computing / MFA
* End to end encrypted blockchain smart contracts
* Trustless Voting

## References

[CGGI19]: I. Chillotti, N. Gama, M. Georgieva, and M. Izabachène. TFHE: Fast Fully Homomorphic Encryption over the Torus. In Journal of Cryptology, volume 33, pages 34–91 (2020). [PDF](https://eprint.iacr.org/2018/421.pdf)

[CGGI16]: I. Chillotti, N. Gama, M. Georgieva, and M. Izabachène. Faster fully homomorphic encryption: Bootstrapping in less than 0.1 seconds. In Asiacrypt 2016 (Best Paper), pages 3-33. [PDF](https://eprint.iacr.org/2016/870.pdf)


