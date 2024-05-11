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

    //Pause
    function pauseValidation() external;

    //Resume
    function resumeValidation() external;

    //Withdraw
    function completeWithdrawal() external returns (uint256);

    //Rewards and Slashing
    function addDepositorSlashing(address depositorAddress, uint256 slashAmount) external returns (uint256);
    function addDepositorReward(address depositorAddress, uint256 rewardAmount) external returns (uint256);

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
    function getWithdrawalBlock(address depositorAddress) external view returns (uint256);
    function isValidationPaused(address validatorAddress) external view returns (bool);

    function doesValidatorExist(address validatorAddress) external view returns (bool);
    function didValidatorEverExist(address validatorAddress) external view returns (bool);
    function doesDepositorExist(address depositorAddress) external view returns (bool);
    function didDepositorEverExist(address depositorAddress) external view returns (bool);

    //Staking V2 functions

    struct StakingDetails {
        address Depositor;
        address Validator;
        uint256 Balance;
        uint256 NetBalance;
        uint256 BlockRewards;
        uint256 Slashings;
        bool    IsValidationPaused;
        uint256 WithdrawalBlock;
        uint256 WithdrawalAmount;
        uint256 LastNilBlockNumber;
        uint256 NilBlockCount;
    }

    //Rotate
    function changeValidator(address newValidatorAddress) external;

    //Deposit
    function increaseDeposit() external payable;

    //Withdrawal
    function initiatePartialWithdrawal(uint256 amount) external returns (uint256);
    function completePartialWithdrawal() external returns (uint256);

    function getStakingDetails(address validatorAddress) external view returns (StakingDetails calldata);

    //Liveness
    function setNilBlock(address validatorAddress) external;
    function resetNilBlock(address validatorAddress) external;

    event OnNewDeposit(
        address indexed depositorAddress,
        address indexed validatorAddress,
        uint256 amount,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnPauseValidation(
        address depositorAddress,
        address validatorAddress
    );

    event OnResumeValidation(
        address depositorAddress,
        address validatorAddress
    );

    event OnCompleteWithdrawal(
        address depositorAddress,
        uint256 netBalance
    );

    event OnSlashing(address indexed depositorAddress,
        uint256 slashedAmount);

    event OnReward(address indexed depositorAddress,
        uint256 rewardAmount);

    //Staking V2 events
    event OnChangeValidator(address indexed depositorAddress, address indexed oldValidatorAddress, address indexed newValidatorAddress);

    event OnIncreaseDeposit(address indexed depositorAddress, uint256 oldBalance, uint256 newBalance);

    event OnInitiatePartialWithdrawal(address indexed depositorAddress, uint256 withdrawalBlock, uint256 withdrawalQuantity);
    event OnCompletePartialWithdrawal(address indexed depositorAddress, uint256 withdrawalQuantity);
}

contract StakingContract is IStakingContract {
    using SafeMath for uint256;

    uint256 constant MINIMUM_DEPOSIT = 5000000000000000000000000; //5000000
    uint256 constant WITHDRAWAL_BLOCK_DELAY = 32000;

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

    //Whether validation is paused
    mapping (address => bool) private _validationPaused;

    //StakingV2 variables

    //Withdrawal Request, depositor to withdrawalBlock
    mapping (address => uint256) private _depositorPartialWithdrawalBlockMapping;
    mapping (address => uint256) private _depositorPartialWithdrawalAmountMapping;

    //Liveness Mapping
    mapping (address => uint256) private _validatorLastNilBlock;
    mapping (address => uint256) private _validatorNilBlockCount;

    function newDeposit(address validatorAddress) override external payable {
        address depositorAddress = msg.sender;
        uint256 depositAmount = msg.value;
        require(depositAmount >= MINIMUM_DEPOSIT, "Deposit amount below minimum deposit amount");
        require(depositorAddress != validatorAddress, "Depositor address cannot be same as Validator address");
        require(validatorAddress != address(0), "Invalid validator");

        require(_validatorExists[validatorAddress] == false, "Validator already exists");
        require(_validatorEverExisted[validatorAddress] == false, "Validator existed once");

        require(_validatorExists[depositorAddress] == false, "Validator already exists as new depositor");
        require(_validatorEverExisted[depositorAddress] == false, "Validator existed once as new depositor");

        uint256 validatorBalance = validatorAddress.balance;
        require(validatorBalance == 0, "validator balance should be zero"); //Since we don't check validator signature, atleast verify if zero balance

        require(_depositorExists[depositorAddress] == false, "Depositor already exists");
        require(_depositorEverExisted[depositorAddress] == false, "Depositor existed once");

        require(_depositorExists[validatorAddress] == false, "Depositor already exists as new validator once");
        require(_depositorEverExisted[validatorAddress] == false, "Depositor existed once as new validator");

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

    function pauseValidation() override external {
        address depositorAddress = msg.sender;
        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");

        address validatorAddress = _validatorToDepositorMapping[depositorAddress];
        require(_validationPaused[validatorAddress] == false, "Validation is already paused");
        _validationPaused[validatorAddress] = true;

        emit OnPauseValidation(depositorAddress, validatorAddress);
    }

    function resumeValidation() override external {
        address depositorAddress = msg.sender;
        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");

        address validatorAddress = _validatorToDepositorMapping[depositorAddress];
        require(_validationPaused[validatorAddress] == true, "Validation is not paused");
        _validationPaused[validatorAddress] = false;

        emit OnResumeValidation(depositorAddress, validatorAddress);
    }

    //Legacy function, will be removed in staking v3
    function completeWithdrawal() override external returns (uint256) {
        address depositorAddress = msg.sender;
        require(_depositorWithdrawalRequests[depositorAddress] > 0, "Depositor withdrawal request does not exist");

        uint256 balance = _depositorBalances[depositorAddress].add(_depositorRewards[depositorAddress]);
        require(balance > _depositorSlashings[depositorAddress], "balance is negative");
        uint256 netBalance = balance.sub(_depositorSlashings[depositorAddress]);

        delete _depositorBalances[depositorAddress];
        delete _depositorRewards[depositorAddress];
        delete _depositorSlashings[depositorAddress];
        delete _depositorWithdrawalRequests[depositorAddress];

        _totalDepositedBalance = _totalDepositedBalance.sub(netBalance);
	    _depositorExists[depositorAddress] = true;

        (bool success, ) = depositorAddress.call{value:netBalance}("");
        // success should be true
        require(success,"Withdraw failed");

        emit OnCompleteWithdrawal(depositorAddress, netBalance);

        return netBalance;
    }

    function addDepositorSlashing(address depositorAddress, uint256 slashAmount) override external returns (uint256) {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _depositorSlashings[depositorAddress] = _depositorSlashings[depositorAddress].add(slashAmount);

        address zeroAddress = address(0);
        (bool success, ) = zeroAddress.call{value:slashAmount}("");
        // success should be true
        require(success,"transfer to zeroAddress failed");

        emit OnSlashing(depositorAddress, slashAmount);
        return _depositorSlashings[depositorAddress];
    }

    function addDepositorReward(address depositorAddress, uint256 rewardAmount) override external returns (uint256) {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _depositorRewards[depositorAddress] = _depositorRewards[depositorAddress].add(rewardAmount);
        emit OnReward(depositorAddress, rewardAmount);
        return _depositorRewards[depositorAddress];
    }

    function getDepositorCount() override external view returns (uint256) {
        return _depositorCount;
    }

    function getTotalDepositedBalance() override external view returns (uint256) {
        return _totalDepositedBalance;
    }

    function listValidators() override external view returns (address[] memory) {
        return _validatorList;
    }

    function getDepositorOfValidator(address validatorAddress) override external view returns (address) {
        address depositorAddress = _validatorToDepositorMapping[validatorAddress];
        return depositorAddress;
    }

    function getValidatorOfDepositor(address depositorAddress) override external view returns (address) {
        address validatorAddress = _depositorToValidatorMapping[depositorAddress];
        return validatorAddress;
    }

    function getBalanceOfDepositor(address depositorAddress) override external view returns (uint256) {
        return _depositorBalances[depositorAddress];
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

    function getDepositorRewards(address depositorAddress) override external view returns (uint256) {
        return _depositorRewards[depositorAddress];
    }

    function getDepositorSlashings(address depositorAddress) override external view returns (uint256) {
        return _depositorSlashings[depositorAddress];
    }

    function getWithdrawalBlock(address depositorAddress) override external view returns (uint256) {
        uint256 withdrawalBlock = 0;
        if(_depositorPartialWithdrawalBlockMapping[depositorAddress] > 0) {
            withdrawalBlock = _depositorPartialWithdrawalBlockMapping[depositorAddress];
        } else if(_depositorWithdrawalRequests[depositorAddress] > 0) {
            withdrawalBlock = _depositorWithdrawalRequests[depositorAddress];
        }
        return withdrawalBlock;
    }

    function isValidationPaused(address validatorAddress) override external view returns (bool) {
        return _validationPaused[validatorAddress];
    }

    function doesValidatorExist(address validatorAddress) override external view returns (bool) {
        return _validatorExists[validatorAddress];
    }

    function didValidatorEverExist(address validatorAddress) override external view returns (bool) {
        return _validatorEverExisted[validatorAddress];
    }

    function doesDepositorExist(address depositorAddress) override external view returns (bool) {
        return _depositorExists[depositorAddress];
    }

    function didDepositorEverExist(address depositorAddress) override external view returns (bool) {
        return _depositorEverExisted[depositorAddress];
    }

    function changeValidator(address newValidatorAddress) override external {
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

        address oldValidatorAddress = _depositorToValidatorMapping[depositorAddress];

        _validatorExists[newValidatorAddress] = true;
        _validatorEverExisted[newValidatorAddress] = true;

        uint256 lastNilBlock  = _validatorLastNilBlock[oldValidatorAddress];
        uint256 nilBlockCount = _validatorNilBlockCount[oldValidatorAddress];
        if(lastNilBlock > 0) {
            _validatorLastNilBlock[newValidatorAddress] = lastNilBlock;
            _validatorNilBlockCount[newValidatorAddress] = nilBlockCount;

            delete _validatorLastNilBlock[oldValidatorAddress];
            delete _validatorNilBlockCount[oldValidatorAddress];
        }

        _validatorToDepositorMapping[newValidatorAddress] = depositorAddress;
        _depositorToValidatorMapping[depositorAddress] = newValidatorAddress;
        _validatorList.push(newValidatorAddress);

        _validatorExists[oldValidatorAddress] = false;
        delete _validatorToDepositorMapping[oldValidatorAddress];

        bool validationPaused = _validationPaused[oldValidatorAddress];
        if(validationPaused == true) {
            delete(_validationPaused[oldValidatorAddress]);
            _validationPaused[newValidatorAddress] = validationPaused;
        }

        emit OnChangeValidator(depositorAddress, oldValidatorAddress, newValidatorAddress);
    }

    function increaseDeposit() override external payable {
        address depositorAddress = msg.sender;
        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");
        require(_depositorWithdrawalRequests[depositorAddress] == 0, "Depositor withdrawal request exists");

        uint256 depositAmount = msg.value;
        require(depositAmount > 0, "Deposit amount is zero");

        _totalDepositedBalance = _totalDepositedBalance.add(depositAmount);

        uint256 oldBalance = _depositorBalances[depositorAddress];
        uint256 newBalance = oldBalance.add(depositAmount);

        _depositorBalances[depositorAddress] = newBalance;

        emit OnIncreaseDeposit(depositorAddress, oldBalance, newBalance);
    }

    function initiatePartialWithdrawal(uint256 amount) override external returns (uint256) {
        address depositorAddress = msg.sender;
        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");
        require(_depositorWithdrawalRequests[depositorAddress] == 0, "Depositor withdrawal request exists");
        require(_depositorPartialWithdrawalBlockMapping[depositorAddress] == 0, "Depositor partial withdrawal request exists");

        uint256 netBalance = this.getNetBalanceOfDepositor(depositorAddress);
        require(netBalance >= amount, "Depositor net balance is low");

        //First withdraw from rewards and then from balance

        uint256 rewardsAmount = _depositorRewards[depositorAddress];
	    uint256 debitAmount = _depositorSlashings[depositorAddress].add(amount);
        if(rewardsAmount >= debitAmount) {
            _depositorRewards[depositorAddress] = rewardsAmount.sub(debitAmount);
        } else {
            delete _depositorRewards[depositorAddress];
            uint256 remaining = debitAmount.sub(rewardsAmount);
            _totalDepositedBalance = _totalDepositedBalance.sub(remaining);
            _depositorBalances[depositorAddress] = _depositorBalances[depositorAddress].sub(remaining);
        }

        delete _depositorSlashings[depositorAddress];

        _depositorPartialWithdrawalBlockMapping[depositorAddress] = block.number;
        _depositorPartialWithdrawalAmountMapping[depositorAddress] = amount;

        emit OnInitiatePartialWithdrawal(depositorAddress, block.number + WITHDRAWAL_BLOCK_DELAY, amount);

        return amount;
    }

    function completePartialWithdrawal() override external returns (uint256) {
        address depositorAddress = msg.sender;
        require(_depositorExists[depositorAddress] == true, "Depositor does not exist");
        require(_depositorWithdrawalRequests[depositorAddress] == 0, "Depositor withdrawal request exists");
        require(_depositorPartialWithdrawalBlockMapping[depositorAddress] > 0, "Depositor partial withdrawal request does not exist");
        require((_depositorPartialWithdrawalBlockMapping[depositorAddress].add(WITHDRAWAL_BLOCK_DELAY)) >= block.number, "Depositor partial withdrawal request cutoff block not reached");

        uint256 amount = _depositorPartialWithdrawalAmountMapping[depositorAddress];
        delete _depositorPartialWithdrawalBlockMapping[depositorAddress];
        delete _depositorPartialWithdrawalAmountMapping[depositorAddress];

        (bool success, ) = depositorAddress.call{value:amount}("");
        // success should be true
        require(success,"Withdraw failed");

        emit OnCompletePartialWithdrawal(depositorAddress, amount);

        return amount;
    }

    function setNilBlock(address validatorAddress) override external {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _validatorLastNilBlock[validatorAddress] = block.number;
        _validatorNilBlockCount[validatorAddress] = _validatorNilBlockCount[validatorAddress].add(1);
    }

    function resetNilBlock(address validatorAddress) override external {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _validatorLastNilBlock[validatorAddress] = 0;
        delete _validatorNilBlockCount[validatorAddress];
    }

    function getStakingDetails(address validatorAddress) override external view returns (StakingDetails memory) {
        require(_validatorExists[validatorAddress] == true, "Validator does not exist");

        address depositorAddress = _validatorToDepositorMapping[validatorAddress];

        StakingDetails memory stakingDetails;

        uint256 withdrawalBlock = 0;
        uint256 withdrawalAmount = 0;
        uint256 depositorNetBalance = this.getNetBalanceOfDepositor(depositorAddress);

        if(_depositorPartialWithdrawalBlockMapping[depositorAddress] > 0) {
            withdrawalBlock = _depositorPartialWithdrawalBlockMapping[depositorAddress];
            withdrawalAmount = _depositorPartialWithdrawalAmountMapping[depositorAddress];
        } else if(_depositorWithdrawalRequests[depositorAddress] > 0) {
            withdrawalBlock = _depositorWithdrawalRequests[depositorAddress];
            withdrawalAmount = this.getNetBalanceOfDepositor(depositorAddress);
        }

        stakingDetails = StakingDetails(depositorAddress, validatorAddress, _depositorBalances[depositorAddress], depositorNetBalance, _depositorRewards[depositorAddress],
            _depositorSlashings[depositorAddress], _validationPaused[validatorAddress], withdrawalBlock, withdrawalAmount,
            _validatorLastNilBlock[validatorAddress], _validatorNilBlockCount[validatorAddress]);

        return stakingDetails;
    }
}