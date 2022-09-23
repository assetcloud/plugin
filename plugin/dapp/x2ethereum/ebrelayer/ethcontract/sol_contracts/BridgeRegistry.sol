pragma solidity ^0.5.0;

contract BridgeRegistry {

    address public chainBridge;
    address public bridgeBank;
    address public oracle;
    address public valset;
    uint256 public deployHeight;

    event LogContractsRegistered(
        address _chainBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    );
    
    constructor(
        address _chainBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    )
        public
    {
        chainBridge = _chainBridge;
        bridgeBank = _bridgeBank;
        oracle = _oracle;
        valset = _valset;
        deployHeight = block.number;

        emit LogContractsRegistered(
            chainBridge,
            bridgeBank,
            oracle,
            valset
        );
    }
}