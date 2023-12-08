TestNet
set DP_URL=http://172.31.34.126:8545
set DP_ALLOC_ACCOUNT=46f8c16c50b122a568c96fb5e97e44ca9cd205ce
set DP_ALLOC_ACCOUNT_PASSWORD=dummy
set TOKENS_INFO=C:\t2build\tokens.json
set DP_DATA_PATH=/data/
set DP_ACCOUNT_PASSWORD=dummy

1) startTestAccountByCoin
        * Primary Account
        * Generate new account (.85 seconds)
        * Transfer coin from account - to account dynamically

2) startTestAccountByContract
        * Account wise dynamically create token contract
        * create dynamic token name, symbol, decimal, total supply

        Primary Conversion
        -----------------
        * remix-backup-token.zip
        * use remix solidity compiler
        * compile and get ABI and Byte code
        * Example copy and past file name token.abi, token.bin ("object": " ") past without string "
        * Use command abigen due to go-ethereum exe
             abigen --abi token.abi --pkg main --type Token --out token.go --bin token.bin
             abigen --abi greeter.abi --pkg main --type Greeter --out greeter.go --bin greeter.bin

3) startTestAccountByToken
         * Get token contract
         * sent token from address - dynamic select to address

4) startTestAccountTokenByToken Pending