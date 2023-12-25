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
        address erc20,
        string  depositor,
        string  validator,
        string  ethSign,
        string  quantumSign,
        uint256 amount,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnCancelGenesisValidator(
        address  erc20,
        string  depositor,
        string  validator,
        uint256 blockNumber,
        uint256 blockTime
    );

    //Date and time (GMT): Thursday, December 28, 2023 11:59:59 PM
    uint256 _genesisValidatorAllowedCutOffDate = 1703807999;

    address[] private _erc20AddressList;

    //Balance
    mapping (string => uint256) private _depositorBalances;

    //erc20Address, depositor and validator exists
    mapping (address => bool) private _erc20AddressExists;
    mapping (string => bool) private _validatorExists;
    mapping (string => bool) private _depositorExists;

    //erc20Address, depositor and validator reverse mapping
    mapping (address => string) private _erc20AddressToDepositorMapping;
    mapping (address => string) private _erc20AddressToValidatorMapping;
   
    //sign
    mapping (address => string) private _erc20AddressToEthSignMapping;
    mapping (string => string) private _depositorAddressToQuantumSignMapping;

    constructor(){       

    }

    function addGenesisValidator(string memory depositor, string memory validator, 
        string memory ethSign, string memory quantumSign, uint256 depositAmount) 
        external 
        returns (bool)
    { 
         //If cutoff date exceeds
        if(block.timestamp > _genesisValidatorAllowedCutOffDate) {
            return false;
        }
       
        require(_erc20AddressExists[msg.sender] == false, "Caller already exists");
        require(keccak256(abi.encodePacked(depositor)) != keccak256(abi.encodePacked(validator)), "Depositor address cannot be same as Validator address");
        require(_depositorExists[depositor] == false, "Depositor already exists");
        require(_validatorExists[validator] == false, "Validator already exists");

        _erc20AddressList.push(msg.sender);
        _depositorBalances[depositor] = depositAmount;
       
        _erc20AddressExists[msg.sender] = true;
        _validatorExists[validator] = true;
        _depositorExists[depositor] = true;

        _erc20AddressToDepositorMapping[msg.sender] = depositor;
        _erc20AddressToValidatorMapping[msg.sender] = validator;

        _erc20AddressToEthSignMapping[msg.sender] = ethSign;
        _depositorAddressToQuantumSignMapping[depositor] = quantumSign;

        emit OnAddGenesisValidator(
            msg.sender,  depositor,  validator, 
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
     
        require(_erc20AddressExists[msg.sender] == true, "Caller is not a genesis validator");

        string memory depositor = _erc20AddressToDepositorMapping[msg.sender];
        string memory validator = _erc20AddressToValidatorMapping[msg.sender];

        _depositorBalances[depositor] = 0;

        _erc20AddressExists[msg.sender] = false;
        _validatorExists[validator] = false;
        _depositorExists[depositor] = false;

        emit OnCancelGenesisValidator(
            msg.sender, 
            depositor, 
            validator,
            block.number, 
            block.timestamp
        );

        return true;
    }

   function getGenesisValidator(address erc20)
         external view returns  (address, string memory, string memory,
         string memory, string memory, bool, uint256)
    {
   
       string memory depositor = _erc20AddressToDepositorMapping[erc20];
       string memory validator = _erc20AddressToValidatorMapping[erc20];

       string memory ethSign =  _erc20AddressToEthSignMapping[erc20];
       string memory quantumSign =  _depositorAddressToQuantumSignMapping[depositor];

       bool status = _erc20AddressExists[erc20];
       uint256 depositorBalances = _depositorBalances[depositor];
       return(erc20, depositor, validator, ethSign, quantumSign, 
            status, depositorBalances);
    }

    function listGenesisValidators()  external view returns (address[] memory) {
        return _erc20AddressList;
    }
}