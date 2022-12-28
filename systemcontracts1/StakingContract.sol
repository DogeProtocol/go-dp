// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <0.8.0;
pragma abicoder v2;

import "./SafeMath.sol";
import "./IStakingContract.sol";

contract StakingContract is IStakingContract {

    using SafeMath for uint256; 

    //deposit count
    uint256 private _depositCount;

    //deposit balance
    uint256 private _totalDepositBalance;
    mapping (address => uint256) private _balances;

    mapping (bytes32 => bytes) private _validatorKey;
    mapping (bytes32 => address) private _validatorIdSenderMapping;
    mapping (address => bytes32) private _senderValidatorIdMapping;

    //list of validator id 
    address[] private _validatorList; 
   
    constructor() {
        _depositCount = 0;
        _totalDepositBalance = 0;
    }


    function bytes32toaddress(bytes32 data) internal pure returns (address) {
        return address(uint160(uint256(data)));
    }    

    function addresstobytes32(address data) internal pure returns (bytes32) {
        return bytes32(uint256(uint160(data)) << 96);
    }

    function newDeposit(bytes calldata pubkey) override external payable {
        require(pubkey.length > 0, "Public key is invalid");
        require(_validatorIdSenderMapping[_senderValidatorIdMapping[msg.sender]] != msg.sender, "Sender already exists");

        _depositCount = _depositCount.add(1);

        _totalDepositBalance = _totalDepositBalance.add(msg.value);
        _balances[msg.sender] = _balances[msg.sender].add(msg.value);

         bytes32 keyhash = keccak256(pubkey[1:]);  
         address validatorAddress = bytes32toaddress(keyhash);
         bytes32 validatorId = addresstobytes32(validatorAddress);

        _validatorKey[validatorId] = pubkey;
        _validatorIdSenderMapping[validatorId] = msg.sender;
        _senderValidatorIdMapping[msg.sender] = validatorId;
       
        _validatorList.push(validatorAddress);

        emit OnNewDeposit(
            msg.sender,
            validatorId,
            validatorAddress,
            pubkey,
            msg.value,
            block.number,
            block.timestamp
        );
    }

    function withdraw(uint256 value)  override external {
        require(_balances[msg.sender] >= value, "Insufficient funds");
        _totalDepositBalance = _totalDepositBalance.sub(value);
        _balances[msg.sender] = _balances[msg.sender].sub(value);
         msg.sender.transfer(value);

        emit OnWithdrawKey(
            msg.sender,
            value,
            block.number,
            block.timestamp
        );
    }

    function depositCount() override external view returns (uint256) {
        return _depositCount;
    }

    function totalDepositBalance() override external view returns (uint256) {
        return _totalDepositBalance;
    }

    function depositBalanceOf(address depositor) override external view returns (uint256) {
        return _balances[depositor];
    }   

    function listValidator() override external view returns (address[] memory) {
        return _validatorList;
    }

    function getDepositor(address validator) override external view returns (address) {
        bytes32 validatorId = addresstobytes32(validator);
        address depositor = _validatorIdSenderMapping[validatorId];
        return depositor;
    }
}