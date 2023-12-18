## Who are Genesis Validators?
Genesis Validators are validators who run blockchain nodes a part of the mainnet Genesis block initialization. New validators can subsequently join after mainnet to create new blocks.

## What are the requirements for becoming a Genesis Validator?
1) Ensure that your address is part of mainnet snapshot at https://snapshot.dpscan.app
2) You should have had atleast 5000000 (5 million) DogeP at the time of snapshot.
3) You should run the DP Blockchain Node on the mainnet date.
4) You should complete the steps below to indicate you want to become a Genesis validator, before 27th Dec 2023, 12:59:59 PM UTC

## What are the node requirements?
1) Atleast 8 cores CPU
2) Atleast 32 GB RAM
3) Atleast 2 TB SSD disk (SSD disk is important)

## How to become a Genesis Validator for Doge Protocol mainnet?

1) Run the following command from Terminal command prompt:
```
     curl https://github.com/DogeProtocol/go-dp/releases/download/v2.0.12/dp-release-windows.zip
```
2) Extract the content of the above zip file to c:\dp
3) Change to c:\dp folder in the Terminal command prompt.
4) Run the following command to create Depositor Account Wallet. Use a strong password on being prompted.
```
     dp --datadir data account new
```
5) Run the following command to create Validator Account Wallet. Use a strong password on being prompted.
```
     dp --datadir data account new
```
6) Backup the data folder containing these wallets safely in different devices and also save copies offline (like in disconnected Pen Drives).

   > [WARNING] 
   > If you loose these wallets or forget the passwords, you will not only be ineligible to become a genesis validator, but will also be not able to get mainnet coins!
   
7) Run the following command to indicate how many coins you would like to deposit for Genesis Validator.  
```
     dputil genesis-sign ETH_ADDRESS DEPOSITOR_QUANTUM_ADDRESS VALIDATOR_QUANTUM_ADDRESS AMOUNT
```
