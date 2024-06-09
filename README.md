## Go DP

go-dp is a Golang implementation of the Doge Protocol blockchain node client. This is a fork of [go-ethereum](https://github.com/ethereum/go-ethereum). DP quantum coin is a quantum resistant blockchain.

[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://discord.gg/bbbMPyzJTM)

### Prerequisites

Requires GO version 1.21 or later.

#### Linux (Ubuntu)

Has only been tested on Ubuntu version 22. Lower Ubuntu versions might need openssl library installation (libcrypto).

##### Setup
```
- 1) Open a new bash window and navigate to the go-dp folder.
- 2) Run ./install-ubuntu.sh
- 3) Make sure you add the environment variables to your bash profile as described at the end of the output of the previous command.
```

##### Building
- 1) Open a new terminal and navigate to the go-dp folder.
- 2) Run go build -o YOUR_BUILD_FOLDER ./...

#### Windows
Ensure you have allowed Powershell local script execution. You can enable this by running the following command in Powershell window that is opened as an administrator:

```
Set-ExecutionPolicy RemoteSigned
```

##### Setup 
- 1) Open a new terminal and navigate to the go-dp folder.
- 2) Run ./install.ps1
 
#### Building     

- 1) Open a new command prompt and navigate to the go-dp folder. Note that this method doesn't work in Powershell or Terminal, hence use command prompt.
- 2) Run templibs/setenv.cmd 
- 3) Run go build -o YOUR_BUILD_FOLDER ./...

#### Mac

Has only been tested on Apple M1.

##### Setup
```
- 1) Ensure brew is installed. To install brew, follow the instructions at https://brew.sh
- 2) Open a new Terminal window and navigate to the go-dp folder.
- 3) Run ./install-mac.sh
- 4) Make sure you add the environment variables to your shell profile as described at the end of the output of the previous command.
```

##### Building
- 1) Open a new terminal and navigate to the go-dp folder. Ensure that appropriate environment variables from the prerequisites section have been set.
- 2) Run go build -o YOUR_BUILD_FOLDER ./...

### Running geth
Check the [documentation](https://dpdocs.org) portal for information on running the blockchain node client.

## Major changes from [go-ethereum](https://github.com/ethereum/go-ethereum)

go-dp is a fork of the Go Ethereum Client (go-ethereum) with the following changes:

1) [Hybrid-PQC](https://github.com/DogeProtocol/hybrid-pqc) that uses a combiner of Dilithium, ed25519 and SPHINCS+ in breakglass mode, is used to secure accounts. This is a change from Ethereum which is vulnerable to quantum computers (Shor's algorithm).

2) Kyber, which is a post-quantum KEM scheme, is used to secure inter-node communication.

3) These cryptographic schemes have been added in the following package:
   (https://github.com/DogeProtocol/go-dp/tree/dogep/crypto)

4) Addresses are 32 bytes instead of 20 bytes in Ethereum, for increased security.

5) RLPX Protocol has been completely rewritten and modularized, to use post-quantum cryptography model; the final client and server encryption keys 
are derived similar to TLS as detailed in RFC 8446. Kyber is used for key exchange and the key material thus derived 
is used as input to HMAC HKDF functions (RFC 5869). However, unlike TLS, instead of trusting the certificate, 
the key of the other node is instead trusted. The private key corresponds to the hybrid pqc key-pair used to secure the account 
using digital signatures. These changes are at (https://github.com/DogeProtocol/go-dp/tree/dogep/p2p/rlpx)

6) A new consensus engine (Proof-of-Stake) has been added.  It uses 3 phase BFT consensus, for deterministic finality. The timeout values are adjusted to improve liveness, within the bounds of FLP theorm.

## Known Issues

1) Commits to fix tests are pending sanitization, before merge.
2) The transaction metadata contains values names 'v', 'r', 's'; these are specific to Ethereum.
These values can be used for public key recovery from the transaction metadata in Ethereum. 

## Contributing

Thank you for considering to help out with the source code! We welcome contributions
from anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to go-dp, please fork, fix, commit and send a pull request
 to review and merge into the main code base. If you wish to submit
more complex changes though, please check up first in [our Discord Server](https://discord.gg/bbbMPyzJTM)
to ensure those changes are in line with the general philosophy of the project and/or get
some early feedback which can make both your efforts much lighter as well as our review
and merge procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting)
   guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary)
   guidelines.
 * Pull requests need to be based on and opened against the `dogep` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "eth, rpc: make trace configs optional"

## License
The go-dp library maintains the same licensing model of go-ethereum. The library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
