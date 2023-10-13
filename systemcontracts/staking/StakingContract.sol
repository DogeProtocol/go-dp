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

    //Rotate
    function changeValidator(address newValidatorAddress) external payable;

    //get data
    function getDepositorCount() external view returns (uint256);
    function getTotalDepositedBalance() external view returns (uint256);
    function listValidators() external view returns (address[] memory);
    function getDepositorOfValidator(address validatorAddress) external view returns (address);
    function getValidatorOfDepositor(address depositorAddress) external view returns (address);
    function getBalanceOfDepositor(address depositorAddress) external view returns (uint256);
    function getNetBalanceOfDepositor(address depositorAddress) external view returns (uint256);
    function getDepositorRewards(address depositorAddress) external view returns (uint256);
    function getDepositorSlashings(address depositorAddress) external view returns (uint256);
    function addDepositorSlashing(address depositorAddress, uint256 slashAmount) external returns (uint256);
    function addDepositorReward(address depositorAddress, uint256 rewardAmount) external returns (uint256);

    //Withdraw
    function initiateWithdrawal() external;
    function completeWithdrawal() external;

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
        address indexed newValidatorAddress
    );

    event OnInitiateWithdrawal(
        address depositorAddress
    );

    event OnCompleteWithdrawal(
        address depositorAddress
    );

    event OnSlashing(address indexed depositorAddress, uint256 slashedAmount);
    event OnReward(address indexed depositorAddress, uint256 rewardAmount);
}

//TODO: THIS IMPLEMENTATION IS INCOMPLETE AND NOT READY

contract StakingContract is IStakingContract {
    using SafeMath for uint256;

    uint256 constant MINIMUM_DEPOSIT = 250000000000000000000000000; //250000000
    uint256 constant WITHDRAWAL_BLOCK_DELAY = 250000;

    address[] private _validatorList;

    //depositor balance
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

    //Slashings
    mapping (address => uint256) private _depositorSlashings;

    //Rewards
    mapping (address => uint256) private _depositorRewards;

    //Withdrawal Request, depositor to withdrawalBlock
    mapping (address => uint256) private _depositorWithdrawalRequests;

    function newDeposit(address validatorAddress) override external payable {
        address depositorAddress = msg.sender;
        uint256 depositAmount = msg.value;
        require(depositAmount >= MINIMUM_DEPOSIT, "Deposit amount below minimum deposit amount");
        require(depositorAddress != validatorAddress, "Depositor address cannot be same as Validator address");
        require(validatorAddress != address(0), "Invalid validator");

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

    function changeValidator(address newValidatorAddress) override external payable {
        require(_validatorExists[newValidatorAddress] == false, "Validator already exists");
        require(_depositorExists[newValidatorAddress] == false, "Validator is a depositor");
        require(_validatorEverExisted[newValidatorAddress] == false, "Validator already existed");
        require(_depositorEverExisted[newValidatorAddress] == false, "Depositor already existed");
        require(newValidatorAddress.balance == 0, "validator balance should be zero"); //Since we don't check validator credentials, atleast verify if zero balance
        require(newValidatorAddress != address(0), "Invalid validator");

        address depositorAddress = msg.sender;
        require(depositorAddress != newValidatorAddress, "Depositor address cannot be same as Validator address");

        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");
        require(_depositorWithdrawalRequests[depositorAddress] == 0, "Withdrawal is pending");

        _validatorExists[newValidatorAddress] = true;
        _validatorEverExisted[newValidatorAddress] = true;
        _validatorToDepositorMapping[newValidatorAddress] = depositorAddress;
        _depositorToValidatorMapping[depositorAddress] = newValidatorAddress;
        _validatorList.push(newValidatorAddress);

        address oldValidatorAddress = _depositorToValidatorMapping[depositorAddress];
        _validatorExists[oldValidatorAddress] = false;
        delete _validatorToDepositorMapping[oldValidatorAddress];

        emit OnChangeValidator(depositorAddress, oldValidatorAddress, newValidatorAddress);
    }

    function getNetBalanceOfDepositor(address depositorAddress) override external view returns (uint256) {
        if (_depositorExists[depositorAddress] == false) {
            return 0;
        }

        if (_depositorWithdrawalRequests[depositorAddress] > 0) {
            return 0;
        }

        uint256 balance = _depositorBalances[depositorAddress].add(_depositorRewards[depositorAddress]);
        if (balance <= _depositorSlashings[depositorAddress]) {
            return 0;
        }

        return balance.sub(_depositorSlashings[depositorAddress]);
    }

    function initiateWithdrawal() override external {
        address depositorAddress = msg.sender;
        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");
        require(_depositorWithdrawalRequests[depositorAddress] == 0, "Depositor withdrawal request exists");
        require(_depositorBalances[depositorAddress] > 0, "Depositor balance is zero");

        _depositorWithdrawalRequests[depositorAddress] = block.number + WITHDRAWAL_BLOCK_DELAY;
        _depositorExists[depositorAddress] == false;

        emit OnInitiateWithdrawal(depositorAddress);
    }

    function completeWithdrawal() override external {
        address depositorAddress = msg.sender;
        require(_depositorWithdrawalRequests[depositorAddress] > 0, "Depositor withdrawal request does not exist");
        require(block.number > _depositorWithdrawalRequests[depositorAddress], "Depositor withdrawal request pending");

        uint256 amount = this.getNetBalanceOfDepositor(depositorAddress);

        delete _depositorBalances[depositorAddress];
        delete _depositorRewards[depositorAddress];
        delete _depositorSlashings[depositorAddress];
        delete _depositorExists[depositorAddress];

        (bool success, ) = depositorAddress.call{value:amount}("");
        // success should be true
        require(success,"Withdraw failed");

        emit OnCompleteWithdrawal(depositorAddress);
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

    function getDepositorSlashings(address depositorAddress) override external view returns (uint256) {
        return _depositorSlashings[depositorAddress];
    }

    function getDepositorRewards(address depositorAddress) override external view returns (uint256) {
        return _depositorRewards[depositorAddress];
    }

    function addDepositorSlashing(address depositorAddress, uint256 slashAmount) override external returns (uint256) {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _depositorSlashings[depositorAddress] = _depositorSlashings[depositorAddress].add(slashAmount);
        emit OnSlashing(depositorAddress, slashAmount);
        return _depositorSlashings[depositorAddress];
    }

    function addDepositorReward(address depositorAddress, uint256 rewardAmount) override external returns (uint256) {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _depositorRewards[depositorAddress] = _depositorRewards[depositorAddress].add(rewardAmount);
        emit OnReward(depositorAddress, rewardAmount);
        return _depositorRewards[depositorAddress];
    }
}