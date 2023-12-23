// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.0 <0.8.20;

pragma abicoder v2;

library SafeMath {
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        }
        uint256 c = a * b;
        assert(c / a == b);
        return c;
    }

    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a / b;
        return c;
    }

    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        assert(b <= a);
        return a - b;
    }

    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        assert(c >= a);
        return c;
    }

    function ceil(uint256 a, uint256 m) internal pure returns (uint256) {
        uint256 c = add(a,m);
        uint256 d = sub(c,1);
        return mul(div(d,m),m);
    }
}

contract DpGenesisValidator {

    using SafeMath for uint256;

    event OnAddGenesisValidator(
        string  erc20Address,
        string  depositorAddress,
        string  validatorAddress,
        string  ethSign,
        string  quantumSign,
        uint256 amount,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnCancelGenesisValidator(
        string  erc20Address,
        string  depositorAddress,
        string  validatorAddress,
        uint256 blockNumber,
        uint256 blockTime
    );

    //Wednesday, December 27, 2023 11:59:59 PM
    uint256 _genesisValidatorAllowedCutOffDate = 1703721599;  

    string[] private _erc20AddressList;

    //Balance
    mapping (string => uint256) private _depositorBalances;

    //erc20Address, depositor and validator exists
    mapping (string => bool) private _erc20AddressExists;
    mapping (string => bool) private _validatorExists;
    mapping (string => bool) private _depositorExists;

    //erc20Address, depositor and validator reverse mapping
    mapping (string => string) private _erc20AddressToDepositorMapping;
    mapping (string => string) private _erc20AddressToValidatorMapping;
   
    //sign
    mapping (string => string) private _erc20AddressToEthSignMapping;
    mapping (string => string) private _depositorAddressToQuantumSignMapping;

    constructor(){       

    }

    function addGenesisValidator(string memory depositorAddress, string memory validatorAddress, 
        string memory ethSign, string memory quantumSign, uint256 depositAmount) 
        external 
        returns (bool)
    { 
         //If cutoff date exceeds
        if(block.timestamp > _genesisValidatorAllowedCutOffDate) {
            return false;
        }

        string memory erc20Address = toLower(addressToString(msg.sender));
        depositorAddress = toLower(depositorAddress);
        validatorAddress = toLower(validatorAddress);

        require(_erc20AddressExists[erc20Address] == false, "Caller already exists");
        require(keccak256(abi.encodePacked(depositorAddress)) != keccak256(abi.encodePacked(validatorAddress)), "Depositor address cannot be same as Validator address");
        require(_depositorExists[depositorAddress] == false, "Depositor already exists");
        require(_validatorExists[validatorAddress] == false, "Validator already exists");

        _erc20AddressList.push(erc20Address);
        _depositorBalances[depositorAddress] = depositAmount;
       
        _erc20AddressExists[erc20Address] = true;
        _validatorExists[validatorAddress] = true;
        _depositorExists[depositorAddress] = true;

        _erc20AddressToDepositorMapping[erc20Address] = depositorAddress;
        _erc20AddressToValidatorMapping[erc20Address] = validatorAddress;

        _erc20AddressToEthSignMapping[erc20Address] = ethSign;
        _depositorAddressToQuantumSignMapping[depositorAddress] = quantumSign;

        emit OnAddGenesisValidator(
            erc20Address,  depositorAddress,  validatorAddress, 
            ethSign, quantumSign, depositAmount,
            block.number, 
            block.timestamp
        );

        return true;
    }

    function cancelGenesisValidator()
        external
        returns (bool)
    {
        //If cutoff date exceeds
        if(block.timestamp > _genesisValidatorAllowedCutOffDate) {
            return false;
        }
        
        string memory erc20Address = toLower(addressToString(msg.sender));

        require(_erc20AddressExists[erc20Address] == true, "Caller is not a genesis validator");

        string memory depositorAddress = _erc20AddressToDepositorMapping[erc20Address];
        string memory validatorAddress = _erc20AddressToValidatorMapping[erc20Address];

        _depositorBalances[depositorAddress] = 0;

        _erc20AddressExists[erc20Address] = false;
        _validatorExists[validatorAddress] = false;
        _depositorExists[depositorAddress] = false;

        emit OnCancelGenesisValidator(
            erc20Address, 
            depositorAddress, 
            validatorAddress,
            block.number, 
            block.timestamp
        );

        return true;
    }

   function getGenesisValidator(string memory erc20Address)
         external view returns  (string memory, string memory, string memory,
         string memory, string memory, bool, uint256)
    {
       erc20Address = toLower(erc20Address);
       string memory depositorAddress = _erc20AddressToDepositorMapping[erc20Address];
       string memory validatorAddress =  _erc20AddressToValidatorMapping[erc20Address];

       string memory ethSign =  _erc20AddressToEthSignMapping[erc20Address];
       string memory quantumSign =  _depositorAddressToQuantumSignMapping[depositorAddress];

       bool status = _erc20AddressExists[erc20Address];
       uint256 depositorBalances = _depositorBalances[depositorAddress];
       return(erc20Address, depositorAddress, validatorAddress, ethSign, quantumSign, 
            status, depositorBalances);
    }

    function listGenesisValidators()  external view returns (string[] memory) {
        return _erc20AddressList;
    }

    function addressToString(address account) private pure returns(string memory) {
        bytes32 value = bytes32(uint256(uint160(account)));
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(2 + value.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint i = 0; i < value.length; i++) {
            str[2+i*2] = alphabet[uint(uint8(value[i] >> 4))];
            str[3+i*2] = alphabet[uint(uint8(value[i] & 0x0f))];
        }
        return string(str);
    }

    function toLower(string memory str) internal pure returns (string memory) {
        bytes memory bStr = bytes(str);
        bytes memory bLower = new bytes(bStr.length);
        for (uint i = 0; i < bStr.length; i++) {
            // Uppercase character...
            if ((uint8(bStr[i]) >= 65) && (uint8(bStr[i]) <= 90)) {
                // So we add 32 to make it lowercase
                bLower[i] = bytes1(uint8(bStr[i]) + 32);
            } else {
                bLower[i] = bStr[i];
            }
        }
        return string(bLower);
    }
    
}