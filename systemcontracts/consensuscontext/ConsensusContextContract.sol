// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <0.8.0;
pragma abicoder v2;

interface IConsensusContextContract {
    function setContext (string calldata contextId, bytes32 Context) external;

    function getContext (string calldata contextId) external view returns (bytes32);

    function deleteContext (string calldata contextId) external;
}

//This contract can be used to set
contract ConsensusContextContract is IConsensusContextContract {
    mapping (string => bytes32) private _contextMap; //key is context id, value is context

    function setContext(string calldata contextId, bytes32 context) override external {
        require(msg.sender == address(0), "Only VM calls are allowed");
        _contextMap[contextId] = context;
    }

    function getContext(string calldata contextId) override external view returns (bytes32) {
        return _contextMap[contextId];
    }

    function deleteContext (string calldata contextId) override external {
        require(msg.sender == address(0), "Only VM calls are allowed");
        delete _contextMap[contextId];
    }
}