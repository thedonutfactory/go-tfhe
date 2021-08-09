# go-tfhe
Golang implementation of TFHE Homomorphic Encryption Library

TFHE, or Fully Homomorphic Encryption Library over the Torus, is a scheme created by (Ilaria Chillotti)[https://github.com/ilachill] et al, that implements a fast, fully bootstrapped circuit environment for running programs within the homomorphic realm.

## WTF is a "fully bootstrapped circuit environment for running programs within the homomorphic realm"

In simple terms, fully homomorphic encryption allows two parties, Alice and Bob, to execute programs on each others's computer systems, without knowing the inputs and output of the data. For example, let's say that Bob owns a cloud processing company that crunches health datasets for Alice. Bob has access to cutting-edge bare metal machines with a lot of processing power and is happy to sell Alice that processing power. However, due to HIPPA compliance requirements ( and a general, altruistic respect for an individual's privacy ), Alice cannot actually share the data.

How do we solve this problem with today's cryptography? Well, we can encrypt the data over the wire, send it to Bob's cloud processing company, and then securely hand him the private key to decrypt, process and reencrypt the data. However, we are still giving Bob's prying eyes access to very sensitive, private health data of Alice's customers.

Enter fully homomorphic encryption. Using an FHE cryptographic runtime, Alice can build a special homomorphic software program designed to process her customer data. She gives this special program over to Bob, where he installs it onto one of his meaty servers. Now, here's the magic... Alice can fully encrypt all of her customer data, give it to Bob, who executes the homomorphic software program to process it, returning to Alice the fully encrypted results, all without ever seeing any of her customer's data unencrypted! Alice decrypts the resulting data with the same key she used to encrypt it's inputs, knowing full well that her data was always safe, even when being processed on Bob's servers.

Modern day cryptographic miracle, eh?

So, if Alice writes a simple homomorphic program to add two numbers and return the results, classically she would create something like this:

```
function addTwoNumbers(int32 a, int32 b) {
  return a + b
}
```

So if she were to run the program function:

`addTwoNumbers(2, 3)`, resulting in `5`

Simple, but Bob and other prying eyes could see the input and output values of the function.

Now, if she used a fully homomorphic runtime environment, she would be able to first encrypt the values of `a` and `b`, pass these encrypted values to Bob, where he executes the homomorphic version of the same function giving back an encrypted result:

1. Alice encrypts a's value of `2` resulting in `1d8b4cf854c`
2. Alice encrypts a's value of `3` resulting in `32c4feed996`
3. Alice asks Bob to run the program function `addTwoNumbers(1d8b4cf854c, 32c4feed996)`. (To Bob, the numbers are gobbely gloop)
4. Bob gets the result `489f719cad` and returns it to Alice
5. Alice decrypts the result with her key revealing the number to be `5`. Magic.
