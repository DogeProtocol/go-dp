// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <0.8.0;
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

interface IStakingContract {
    //Deposit
    function newDeposit(bytes32 keyhash)  external payable;

    //Withdraw
    function withdraw(uint256 value)  external;

    //get data
    function depositCount() external view returns (uint256);
    function totalDepositBalance() external view returns (uint256);
    function depositBalanceOf(address owner)  external view returns (uint256);
    function listValidator() external view returns (address[] memory);
    function getDepositor(address validator) external view returns (address);

    event OnNewDeposit(
        address indexed sender,
        bytes32 indexed validatorId,
        address indexed validatorAddress,
        uint256 value,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnWithdrawKey(
        address sender,
        uint256 value,
        uint256 blockNumber,
        uint256 blockTime
    );
}

contract StakingContract is IStakingContract {

    using SafeMath for uint256;

    //deposit count
    uint256 private _depositCount = 0;

    //deposit balance
    uint256 private _totalDepositBalance = 0;
    mapping (address => uint256) private _balances;

    mapping (bytes32 => address) private _validatorIdSenderMapping;
    mapping (address => bytes32) private _senderValidatorIdMapping;

    //list of validator id
    address[] private _validatorList;

    function bytes32toaddress(bytes32 data) internal pure returns (address) {
        return address(uint160(uint256(data)));
    }

    function addresstobytes32(address data) internal pure returns (bytes32) {
        return bytes32(uint256(uint160(data)) << 96);
    }

    function newDeposit(bytes32 keyhash) override external payable {
        require(_validatorIdSenderMapping[_senderValidatorIdMapping[msg.sender]] != msg.sender, "Sender already exists");

        _depositCount = _depositCount.add(1);

        _totalDepositBalance = _totalDepositBalance.add(msg.value);
        _balances[msg.sender] = _balances[msg.sender].add(msg.value);

        address validatorAddress = bytes32toaddress(keyhash);
        bytes32 validatorId = addresstobytes32(validatorAddress);

        _validatorIdSenderMapping[validatorId] = msg.sender;
        _senderValidatorIdMapping[msg.sender] = validatorId;

        _validatorList.push(validatorAddress);

        emit OnNewDeposit(
            msg.sender,
            validatorId,
            validatorAddress,
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