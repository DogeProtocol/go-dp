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
    function newDeposit(address validatorAddress) external payable;
    function increaseDeposit() external payable;

    //Rotate
    function changeValidator(address newValidatorAddress) external payable;

    //Withdraw
    function initiateWithdrawal() external;
    function completeWithdrawal() external;

    //get data
    function getDepositorCount() external view returns (uint256);
    function getTotalDepositedBalance() external view returns (uint256);
    function listValidators() external view returns (address[] memory);
    function getDepositorOfValidator(address validatorAddress) external view returns (address);
    function getValidatorOfDepositor(address depositorAddress) external view returns (address);
    function getBalanceOfDepositor(address depositorAddress)  external view returns (uint256);

    event OnNewDeposit(
        address indexed depositorAddress,
        address indexed validatorAddress,
        uint256 amount,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnIncreaseDeposit(
        address indexed depositorAddress,
        uint256 amount,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnChangeValidator(
        address indexed depositorAddress,
        address indexed oldValidatorAddress,
        address indexed newValidatorAddress,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnInitiateWithdrawal(
        address depositorAddress,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnCompleteWithdrawal(
        address depositorAddress,
        uint256 blockNumber,
        uint256 blockTime
    );
}

contract StakingContract is IStakingContract {
    using SafeMath for uint256;

    uint256 constant MINIMUM_DEPOSIT = 250000000000000000000000000; //250000000

    address[] private _validatorList;

    //Balance
    mapping (address => uint256) private _depositorBalances;

    //total amount
    uint256 private _totalDepositedBalance = 0;

    //depositor count
    uint256 private _depositorCount = 0;

    //depositor and validator exists
    mapping (address => bool) private _validatorExists;
    mapping (address => bool) private _depositorExists;

    mapping (address => bool) private _validatorEverExisted;
    mapping (address => bool) private _depositorEverExisted;

    //depositor and validator reverse mapping
    mapping (address => address) private _validatorToDepositorMapping;
    mapping (address => address) private _depositorToValidatorMapping;

    function newDeposit(address validatorAddress) override external payable {
        address depositorAddress = msg.sender;
        uint256 depositAmount = msg.value;
        require(depositAmount >= MINIMUM_DEPOSIT, "Deposit amount below minimum deposit amount");

        require(depositorAddress != validatorAddress, "Depositor address cannot be same as Validator address");

        require(_validatorExists[validatorAddress] == false, "Validator already exists");
        require(_validatorEverExisted[validatorAddress] == false, "Validator existed once");
        uint256 validatorBalance = validatorAddress.balance;
        require(validatorBalance == 0, "validator balance should be zero"); //Since we don't check validator credentials, atleast verify if zero balance

        require(_depositorExists[depositorAddress] == false, "Depositor already exists");
        require(_depositorEverExisted[depositorAddress] == false, "Depositor existed once");

        _validatorList.push(validatorAddress);
        _totalDepositedBalance = _totalDepositedBalance.add(depositAmount);
        _depositorCount = _depositorCount.add(1);
        _depositorBalances[depositorAddress] = depositAmount;

        _validatorExists[validatorAddress] = true;
        _depositorExists[depositorAddress] = true;
        _validatorEverExisted[validatorAddress] = true;
        _depositorEverExisted[depositorAddress] = true;

        _validatorToDepositorMapping[validatorAddress] = depositorAddress;
        _depositorToValidatorMapping[depositorAddress] = validatorAddress;

        emit OnNewDeposit(
            depositorAddress,
            validatorAddress,
            depositAmount,
            block.number,
            block.timestamp
        );
    }

    function increaseDeposit() override external payable {

    }

    function changeValidator(address newValidatorAddress) override external payable {

    }

    function initiateWithdrawal() override external {

    }

    function completeWithdrawal() override external {

    }

    function getDepositorCount() override external view returns (uint256) {
        return _depositorCount;
    }

    function getTotalDepositedBalance() override external view returns (uint256) {
        return _totalDepositedBalance;
    }

    function getDepositorOfValidator(address validatorAddress) override external view returns (address) {
        address depositoAddress = _validatorToDepositorMapping[validatorAddress];
        return depositoAddress;
    }

    function getValidatorOfDepositor(address depositorAddress) override external view returns (address) {
        address validatorAddress = _validatorToDepositorMapping[depositorAddress];
        return validatorAddress;
    }

    function listValidators() override external view returns (address[] memory) {
        return _validatorList;
    }

    function getBalanceOfDepositor(address depositorAddress) override external view returns (uint256) {
        return _depositorBalances[depositorAddress];
    }
}