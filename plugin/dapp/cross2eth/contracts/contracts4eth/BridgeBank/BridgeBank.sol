pragma solidity ^0.5.0;

import "./ChainBank.sol";
import "./EthereumBank.sol";
import "./TransferHelper.sol";
import "../Oracle.sol";
import "../ChainBridge.sol";


/**
 * @title BridgeBank
 * @dev Bank contract which coordinates asset-related functionality.
 *      ChainBank manages the minting and burning of tokens which
 *      represent Chain based assets, while EthereumBank manages
 *      the locking and unlocking of Ethereum and ERC20 token assets
 *      based on Ethereum.
 **/

contract BridgeBank is ChainBank, EthereumBank {

    using SafeMath for uint256;

    address public operator;
    Oracle public oracle;
    ChainBridge public chainBridge;
    string public platformTokenSymbol;
    bool public hasSetPlatformTokenSymbol;

    /*
    * @dev: Constructor, sets operator
    */
    constructor (
        address _operatorAddress,
        address _oracleAddress,
        address _chainBridgeAddress
    )
        public
    {
        operator = _operatorAddress;
        oracle = Oracle(_oracleAddress);
        chainBridge = ChainBridge(_chainBridgeAddress);
    }

    /*
    * @dev: Modifier to restrict access to operator
    */
    modifier onlyOperator() {
        require(
            msg.sender == operator,
            'Must be BridgeBank operator.'
        );
        _;
    }

    /*
     * @dev: Modifier to restrict access to Offline
     */
    modifier onlyOffline() {
        require(
            msg.sender == offlineSave,
            'Must be onlyOffline.'
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the oracle
    */
    modifier onlyOracle()
    {
        require(
            msg.sender == address(oracle),
            "Access restricted to the oracle"
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the chain bridge
    */
    modifier onlyChainBridge()
    {
        require(
            msg.sender == address(chainBridge),
            "Access restricted to the chain bridge"
        );
        _;
    }

   /*
    * @dev: Fallback function allows operator to send funds to the bank directly
    *       This feature is used for testing and is available at the operator's own risk.
    */
    function() external payable onlyOffline {}

    /*
    * @dev: Creates a new BridgeToken
    *
    * @param _symbol: The new BridgeToken's symbol
    * @return: The new BridgeToken contract's address
    */
    function createNewBridgeToken(
        string memory _symbol
    )
        public
        onlyOperator
        returns(address)
    {
        return deployNewBridgeToken(_symbol);
    }

    /*
     * @dev: Mints new BankTokens
     *
     * @param _chainSender: The sender's Chain address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _chainTokenAddress: The currency type
     * @param _symbol: chain token symbol
     * @param _amount: number of chain tokens to be minted
     */
     function mintBridgeTokens(
        bytes memory _chainSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyChainBridge
    {
        return mintNewBridgeTokens(
            _chainSender,
            _intendedRecipient,
            _bridgeTokenAddress,
            _symbol,
            _amount
        );
    }

    /*
     * @dev: Burns bank tokens
     *
     * @param _chainReceiver: The _chain receiver address in bytes.
     * @param _chainTokenAddress: The currency type
     * @param _amount: number of chain tokens to be burned
     */
    function burnBridgeTokens(
        bytes memory _chainReceiver,
        address _chainTokenAddress,
        uint256 _amount
    )
        public
    {
        return burnChainTokens(
            msg.sender,
            _chainReceiver,
            _chainTokenAddress,
             _amount
        );
    }

    /*
     * @dev: addToken2LockList used to add token with the specified address to be
     *       allowed locked from Ethereum
     *
     * @param _token: token contract address
     * @param _symbol: token symbol
     */
     function addToken2LockList(
        address _token,
        string memory _symbol
     )
        public
        onlyOperator
     {
         addToken2AllowLock(_token, _symbol);
     }

    /*
    * @dev: configTokenOfflineSave used to config threshold to trigger tranfer token to offline account
    *       when the balance of locked token reaches
    *
    * @param _token: token contract address
    * @param _symbol:token symbol,just used for double check that token address and symbol is consistent
    * @param _threshold: _threshold to trigger transfer
    * @param _percents: amount to transfer per percents of threshold
    */
    function configLockedTokenOfflineSave(
        address _token,
        string memory _symbol,
        uint256 _threshold,
        uint8 _percents
    )
    public
    onlyOperator
    {
        if (address(0) != _token) {
            require(keccak256(bytes(BridgeToken(_token).symbol())) == keccak256(bytes(_symbol)), "token address and symbol is not consistent");
        } else {
            require(true == hasSetPlatformTokenSymbol, "The platform Token Symbol has not been configured");
            require(keccak256(bytes(platformTokenSymbol)) == keccak256(bytes(_symbol)), "token address and symbol is not consistent");
        }
        configOfflineSave4Lock(_token, _symbol, _threshold, _percents);
    }

    /*
    * @dev: configplatformTokenSymbol used to config platform token symbol,and just could be configured once
    *
    * @param _symbol:token symbol,just used for double check that token address and symbol is consistent
    */
    function configplatformTokenSymbol(string memory _symbol) public onlyOperator
    {
        require(false == hasSetPlatformTokenSymbol, "The platform Token Symbol has been configured");
        platformTokenSymbol = _symbol;
        hasSetPlatformTokenSymbol = true;
    }

   /*
    * @dev: configOfflineSaveAccount used to config offline account to receive token
    *       when the balance of locked token reaches threshold
    *
    * @param _offlineSave: receiver address
    */
    function configOfflineSaveAccount(address payable _offlineSave) public onlyOperator
    {
        offlineSave = _offlineSave;
    }

    /*
    * @dev: Locks received Ethereum funds.
    *
    * @param _recipient: bytes representation of destination address.
    * @param _token: token address in origin chain (0x0 if ethereum)
    * @param _amount: value of deposit
    */
    function lock(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        availableNonce()
        payable
    {
        string memory symbol;

        // Ethereum deposit
        if (msg.value > 0) {
          require(
              _token == address(0),
              "Ethereum deposits require the 'token' address to be the null address"
            );
          require(
              msg.value == _amount,
              "The transactions value must be equal the specified amount (in wei)"
            );
          require(true == hasSetPlatformTokenSymbol, "The platform Token Symbol has not been configured");
          // Set the the symbol to ETH
          symbol = platformTokenSymbol;
          // ERC20 deposit
        } else {

            TransferHelper.safeTransferFrom(_token, msg.sender, address(this), _amount);
            symbol = tokenAddrAllow2symbol[_token];

            require(
                tokenAllow2Lock[keccak256(abi.encodePacked(symbol))] == _token,
                'The token is not allowed to be locked from Ethereum.'
            );
        }

        lockFunds(
            msg.sender,
            _recipient,
            _token,
            symbol,
            _amount
        );
    }

   /*
    * @dev: Unlocks Ethereum and ERC20 tokens held on the contract.
    *
    * @param _recipient: recipient's Ethereum address
    * @param _token: token contract address
    * @param _symbol: token symbol
    * @param _amount: wei amount or ERC20 token count
\   */
     function unlock(
        address payable _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyChainBridge
        hasLockedFunds(
            _token,
            _amount
        )
        canDeliver(
            _token,
            _amount
        )
    {
        unlockFunds(
            _recipient,
            _token,
            _symbol,
            _amount
        );
    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getChainDepositStatus(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLockedChainDeposit(_id);
    }

    /*
    * @dev: Allows access to a Chain deposit's information via its unique identifier.
    *
    * @param _id: The deposit to be viewed.
    * @return: Original sender's Ethereum address.
    * @return: Intended Chain recipient's address in bytes.
    * @return: The lock deposit's currency, denoted by a token address.
    * @return: The amount locked in the deposit.
    * @return: The deposit's unique nonce.
    */
    function viewChainDeposit(
        bytes32 _id
    )
        public
        view
        returns(bytes memory, address payable, address, uint256)
    {
        return getChainDeposit(_id);
    }

}