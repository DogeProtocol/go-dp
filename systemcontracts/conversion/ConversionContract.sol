// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <0.8.0;
pragma abicoder v2;

interface IConversionContract {
    function requestConversion(string calldata ethAddress, string calldata ethSignature) external returns (uint8);

    function getConversionStatus(address ethAddress) external view returns (bool);

    function getAmount(address ethAddress) external view returns (uint256);

    function getQuantumAddress(address ethAddress) external view returns (address);

    function setConverted(address ethAddress, address quantumAddress) external returns (uint256);

    event OnRequestConversion(
        address indexed quantumAddress,
        string ethAddress,
        string ethereumSignature
    );

    event OnConversion(
        address indexed quantumAddress,
        address ethAddress,
        uint256 amount
    );
}

contract ConversionContract is IConversionContract {
    mapping (address => uint256) private _snapshotAmountMap; //key is Ethereum Address, value is snapshot amount in wei
    mapping (address => address) private _quantumAddressMap; //key is Ethereum Address, value is quantum address
    mapping (address => bool) private _conversionStatusMap; //key is Ethereum address, value is true/false if conversion is done or not

    function requestConversion(string calldata ethAddress, string calldata ethSignature) override external returns (uint8) {
        emit OnRequestConversion(msg.sender, ethAddress, ethSignature); //do nothing else, request is processed outside the VM in consensus layer
        return 0;
    }

    function getConversionStatus(address ethAddress) override external view returns (bool) {
        return _conversionStatusMap[ethAddress];
    }

    function getAmount(address ethAddress) override external view returns (uint256) {
        return _snapshotAmountMap[ethAddress];
    }

    function getQuantumAddress(address ethAddress) override  external view returns (address) {
        return _quantumAddressMap[ethAddress];
    }

    function setConverted(address ethAddress, address quantumAddress) override external returns (uint256) {
        require(msg.sender == address(0), "Only VM calls are allowed");

        require(_snapshotAmountMap[ethAddress] > 0, "ethAddress Doesn't exist in snapshot");
        require(_conversionStatusMap[ethAddress] == false, "Already converted");

        _conversionStatusMap[ethAddress] = true;
        _quantumAddressMap[ethAddress] = quantumAddress;

        (bool success, ) = quantumAddress.call{value:_snapshotAmountMap[ethAddress]}("");
        require(success,"Transfer Balance failed");

        emit OnConversion(quantumAddress, ethAddress, _snapshotAmountMap[ethAddress]);

        return _snapshotAmountMap[ethAddress];
    }
}