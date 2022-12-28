// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <0.8.0;
pragma abicoder v2;

interface IStakingContract {

    //Deposit
    function newDeposit(bytes calldata pubkey)  external payable;

    //Reward
    function rewardDeposit()  external payable;

    //Withdraw
    function withdraw()  external payable;

    //get data
    function depositCount() external view returns (uint256);
    function totalDepositBalance() external view returns (uint256);
    function depositBalanceOf(address owner)  external view returns (uint256);
    function listValidator() external view returns (address[] memory);

    event OnNewDeposit(
        address indexed sender, 
        bytes32 indexed validatorId, 
        address indexed validatorAddress,
        bytes   pubkey,
        uint256 value,
        uint256 blockNumber,
        uint256 blockTime
    );

    event OnRewardDepositKey(
        address sender,
        bytes32 validatorId,
        address reward,
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
