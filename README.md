## Go DP

go-dp is a Golang implementation of the Doge Protocol blockchain node client. This is a fork of [go-ethereum](https://github.com/ethereum/go-ethereum) with additions such as post quantum cryptography.

[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://discord.gg/bbbMPyzJTM)

## Prerequisites

### Post Quantum Cryptography (liboqs)

1) Follow the steps in https://github.com/DogeProtocol/liboqs to build the liboqs binaries (Post Quantum Cryptography Libraries)
2) Install Package Manager:
#### Linux
```
apt-get install -y pkg-config
```
#### Windows
- 1) Download and extract https://download.gnome.org/binaries/win64/dependencies/pkg-config_0.23-2_win64.zip to c:\pkg-config
- 2) Download and extract https://download.gnome.org/binaries/win64/glib/2.26/glib_2.26.1-1_win64.zip to c:\pkg-config
- 3) Download and extract https://download.gnome.org/binaries/win64/dependencies/gettext-runtime-dev_0.18.1.1-2_win64.zip to c:\pkg-config
- 4) Add c:\pkg-config to the PATH environment variable
4) Edit go-dp/.config/liboqs.pc to point to the liboqs build folder from Step 1)
5) Set the PKG_CONFIG_PATH environment variable to point to the absolute path of go-dp/.config/liboqc.pc file in your computer.
6) Add the build/bin folder where liboqs.dll (example C:\liboqs\build\bin) is located to the PATH environment variable.

## Building geth

Building `geth` requires both a Go (version 1.14 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies and pre-requisites detailed above are installed, run:

```
cd cmd/geth
go build -o c:\build 
```

### Running geth
Check the [documentation](https://dpdocs.org/testnet-setup.html) portal for information on running the blockchain node client.

## Major changes from [go-ethereum](https://github.com/ethereum/go-ethereum)

go-dp is a fork of the Go Ethereum Client (go-ethereum) with the following changes:

1) Falcon, which is a NIST standardized post-quantum digital signature scheme, is used to secure accounts. This is a change from using ECDSA which is vulnerable to quantum computers (Shor's algorithm).

2) NTRU HRSS, which is a post-quantum KEM scheme, is used to secure inter-node communication.

3) These cryptographic schemes have been added in the following package:
   (https://github.com/DogeProtocol/go-dp/tree/dogep/cryptopq)

4) UDP Discovery Protocol has been updated to use session based error correction (Reed Solomon Codes) using the [KCP GO](https://github.com/xtaci/kcp-go) library.
The discovery protocol in Ethereum relies on a single datagram packet size being smaller than commonly accepted MTU size. 
Since post-quantum cryptography scheme public-keys such as that of Falcon are large, they will not fit into a single UDP packet. 
Hence this change to session based UDP has been made in Doge Protocol. These changes are in (https://github.com/DogeProtocol/go-dp/tree/dogep/p2p/discover)

5) RLPX Protocol has been completely rewritten and modularized, to use post-quantum cryptography model; the final client and server encryption keys 
are derived similar to TLS as detailed in RFC 8446. NTRU HRSS is used for key exchange and the key material thus derived 
is used as input to HMAC HKDF functions (RFC 5869). However, unlike TLS, instead of trusting the certificate, 
the key of the other node is instead trusted. The private key corresponds to the Falcon key-pair used to secure the account 
using digital signatures. Since Falcon is susceptible to timing attacks, 
future iterations may switch to Dilithium for securing inter-node communication, while account security is maintained with Falcon.
These changes are at (https://github.com/DogeProtocol/go-dp/tree/dogep/p2p/rlpx)

6) A new consensus engine (Proof-of-Stake) has been added, along with Withdraw and Deposit tools. Only part of the proof-of-stake implementation is complete, the remaining is work-in-progress. 

7) Testnet command line tools for periodically creating test transactions, smart contracts etc. have been added.  

## Known Issues

1) Commits to fix tests are pending sanitization, before merge.
2) The transaction metadata contains values names 'v', 'r', 's'; these are specific to ECDSA, a carry-over from the Ethereum codebase.
These values can be used for public key recovery from the transaction metadata in Ethereum. 
These variables will be removed in forthcoming testnets because the cryptography scheme used in Doge Protocol is Falcon (for quantum resistance).

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
