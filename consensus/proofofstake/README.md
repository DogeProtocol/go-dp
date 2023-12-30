## Who are Genesis Validators?
Genesis Validators are validators who run blockchain nodes a part of the mainnet Genesis block initialization. Validators can run nodes by depositing coins and become eligible to mine coins. You can also join as a validator subsequently anytime after mainnet. 

## What are the requirements for becoming a Genesis Validator?
1) Ensure that your address is part of mainnet snapshot at https://snapshot.dpscan.app
2) You should have had atleast 5000000 (5 million) DogeP coins (check the coins for your address in the snapshot portal).
3) You should run the DP Blockchain Node on the mainnet date.
4) You should complete the steps below to indicate that you want to become a Genesis validator, before 28th Dec 2023, 12:59:59 PM UTC.

   > [WARNING] 
   > If register to become a Genesis Validator, but do not run the validator node, you will loose coins due to slashings from the proof-of-stake consensus mechanism!

## How to become a Genesis Validator for Doge Protocol mainnet?

In order to associate your Ethereum account that has the DogeP tokens and get corresponding coins in the quantum mainnet and become a validator, a cryptographic technique named cross-signing is followed. As part of this, you will create new quantum wallets and then sign a payload with both your Ethereum keys as well as Quantum keys. This payload will comprise of address of your Ethereum address and Quantum address in addition to other meta-data.

You may follow these steps on any device or computer, but it is highly recommended that you encrypt the disks with Bitlocker or Veracrypt; this way, if your computer or disk is stolen, or if you thrash your computer (or disk) in the future and the disk gets on to a bad actor's hands, your keys, wallets and coins will still be protected.

1) Run the following command from Terminal command prompt (Windows):
```
     curl https://github.com/DogeProtocol/go-dp/releases/download/v2.0.13/dp-release-windows.zip
```
2) Extract the content of the above zip file to c:\dp.
3) Change to c:\dp folder in the command prompt or Terminal.
4) Run the following command to create Depositor Account Wallet. Use a strong password on being prompted.
```
     dp --datadir data account new
```
5) Run the following command to create Validator Account Wallet. Use a strong password on being prompted.
```
     dp --datadir data account new
```
6) Backup the data folder containing these wallets safely in different devices and also save copies offline (like in new Pen Drives stored offline safely and dedicatedly used for this purpose).

   > [WARNING] 
   > If you loose these wallets or forget the passwords after registering for becoming a genesis validator, you will not only be ineligible to become a genesis validator, but will also be not able to get mainnet coins!
      
8) Set the following environment variables in the command prompt;
```
     set DP_KEY_FILE_DIR=c:\dp\data\keystore
     set DP_DEPOSITOR_ACC_PWD=Enter_Password_For_Depositor_Account_From_Step_4
     set DP_VALIDATOR_ACC_PWD=Enter_Password_For_Validator_Account_From_Step_5
```

9) Run the following command to complete the quantum signing part of the cross-sign operation.

    Replace ETH_ADDRESS with the address that had the DogeP tokens as of snapshot date.
   
    Replace DEPOSITOR_QUANTUM_ADDRESS from Step 4
   
    Replace VALIDATOR_QUANTUM_ADDRESS from Step 5
   
    Replace AMOUNT with the number of coins used for running the Validator node. This value should be less than or equal to the number of coins 
         specified for your Ethereum address in the snapshot portal. The remaining will be added to your depositor address.
   
```
     dputil genesis-sign ETH_ADDRESS DEPOSITOR_QUANTUM_ADDRESS VALIDATOR_QUANTUM_ADDRESS AMOUNT
```

10) The above command will create a json file. This is the part of the cross-signing in which the quantum part of the signing is complete. Backup this json file. Open this json file in a text editor. You will notice that the field "ethereumSignature" is empty while other fields have values; this is because the Ethereum signature will need to be generated in the subsequent steps.
    
12) Open this json file in a text editor. Copy the Message field, without the quotes. The message sentence starts with "I agree" and ends with a full-stop. Ensure no other extra space of character is copied. An example message is shown below.

```
I AGREE TO BECOME A GENESIS VALIDATOR FOR MAINNET. MY ETH ADDRESS IS 0xbC22f18344B750Cc46F37375ed7fb607b9C649cc. MY CORRESPONDING DEPOSITOR QUANTUM ADDRESS IS 0x313858E1599c63FE32d6fdbC2AaD8c70ad19591b01648bC36a836C9Fc99558B8 AND VALIDATOR QUANTUM ADDRESS IS 0x610a5e81242D27B086D81e950503d202052819b36fa3ccB2Dd3EFd4A35F9b127. VALIDATOR AMOUNT IS 50000000 DOGEP.
```

12) Next you will be navigating to a 3rd party website. Since you will have to connect your Ethereum wallet that held the DogeP tokens as of snapshot date, it is recommended that you do not have any coins or Ethereum or tokens in the wallet (you may move them to a different wallet, for safety).

13) Open the following website:

    https://app.mycrypto.com/sign-message

14) At this point, once the above page is loaded, it is recommended that you disconnect your device from the internet. In case the 3rd party website is compromised, disconnecting from the internet will help in protecting the wallet. Doge Protocol Community or community developers are not responsible for any problems that you may face.

15) Connect your Ethereum wallet using one of the options provided.

16) Paste the message copied from Step 11, in the Message field.

17) Next click Sign message and complete the steps to create the signature. You should see a Textbox for the Signature.

18) In the Signature textbox, identify the field that says "Sig" and copy the value without quotes. This value is the Ethereum part of the cross-signing and starts with the value 0x.

19) Paste this copied sig in the json file from step 10, in the value of the field that says "ethereumSignature" and save the file. Close the browser tab from Step 13 and then close the browser window completely before connecting back to the internet.

20) Now that the cross-signing is completed, the json file should be submitted for indicating that you want to become a genesis validator. To do this, follow either Step 21 or Step 22.

21) You may simply submit a Git Pull Request in this project by adding the json file, under the consensus/proofofstake folder. You can simply re-submit if you made a mistake, before the 28th Dec cut-off date.

22) The other option is to submit the json file via the Genesis Staking Contract that has been created in the Base chain. The base chain is used since the Ethereum chain has a high gas fee. The genesis contract code and details are available at https://basescan.org/address/0x07370dc2139b1ffc486dfd798dc00b69d7ad2bf7#code

Simply call the addGenesisValidator function from the Ethereum address, passing the individual json fields.

```
  addGenesisValidator(string memory depositorAddress, string memory validatorAddress, 
        string memory ethSign, string memory quantumSign, uint256 depositAmount)
```

23) Next, once this is done, watchout for instructions on running mainnet node (Post Dec 28th, 2023). You should be prepared to run this node on the mainnet date, else you might loose coins due to slashings by being an inactive validator.

## What are the node requirements?

It is recommended that you use a dedicated computer for running a node. Though the Doge Protocol blockchain dynamically adjusts according to the capability of nodes, for optimal functionality, the below configuration is recommended. You should get this device ready soon.

1) Atleast 8 cores CPU.
2) Atleast 32 GB RAM.
3) Atleast 2 TB SSD disk (SSD disk is important).
4) Atleast 100 Mbps internet download speed and 50 Mbps upload speed.
5) Unlimited internet data usage from your internet plan.
6) Stable internet connectivity with no downtime. It is recommended to have two internet providers, just in case one of them goes down.
7) Stable power supply; it is recommended to have a backup power mechanism, in case you have a power-cut. This backup power should be both for the computer running the blockchain node as well as the internet modem or devices used for internet.

## Important, please read

As part of mainnet launch, only conversion transactions will be enabled, to allow getting coins from DogeP tokens. Regular send/receive transactions will be enabled on a block that will be created roughly on April 14th 2024, World Quantum Day. This logic is part of the blockchain itself. This is to make it fair and give enough time for those who have DogeP tokens as of the mainnet snapshot, to get their coins. Additionally, ability to mine coins will also be enabled only on April 14th, 2024. This is to tive a fair chance and enough time for new block validators to join, than just genesis block validators.
