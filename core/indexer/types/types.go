package types

import (
	"math/big"
	"time"

	"github.com/ElrondNetwork/elrond-go/data/state"
)

// Transaction is a structure containing all the fields that need
//  to be saved for a transaction. It has all the default fields
//  plus some extra information for ease of search and filter
type Transaction struct {
	Hash                 string        `json:"-"`
	MBHash               string        `json:"miniBlockHash"`
	BlockHash            string        `json:"-"`
	Nonce                uint64        `json:"nonce"`
	Round                uint64        `json:"round"`
	Value                string        `json:"value"`
	Receiver             string        `json:"receiver"`
	Sender               string        `json:"sender"`
	ReceiverShard        uint32        `json:"receiverShard"`
	SenderShard          uint32        `json:"senderShard"`
	GasPrice             uint64        `json:"gasPrice"`
	GasLimit             uint64        `json:"gasLimit"`
	GasUsed              uint64        `json:"gasUsed"`
	Fee                  string        `json:"fee"`
	Data                 []byte        `json:"data"`
	Signature            string        `json:"signature"`
	Timestamp            time.Duration `json:"timestamp"`
	Status               string        `json:"status"`
	SearchOrder          uint32        `json:"searchOrder"`
	EsdtTokenIdentifier  string        `json:"token,omitempty"`
	EsdtValue            string        `json:"esdtValue,omitempty"`
	SenderUserName       []byte        `json:"senderUserName,omitempty"`
	ReceiverUserName     []byte        `json:"receiverUserName,omitempty"`
	Logs                 *TxLog        `json:"logs,omitempty"`
	SmartContractResults []*ScResult   `json:"-"`
	RcvAddrBytes         []byte        `json:"-"`
}

// GetGasLimit will return transaction gas limit
func (t *Transaction) GetGasLimit() uint64 {
	return t.GasLimit
}

// GetGasPrice will return transaction gas price
func (t *Transaction) GetGasPrice() uint64 {
	return t.GasPrice
}

// GetData will return transaction data field
func (t *Transaction) GetData() []byte {
	return t.Data
}

// GetRcvAddr will return transaction receiver address
func (t *Transaction) GetRcvAddr() []byte {
	return t.RcvAddrBytes
}

// GetValue wil return transaction value
func (t *Transaction) GetValue() *big.Int {
	bigIntValue, ok := big.NewInt(0).SetString(t.Value, 10)
	if !ok {
		return big.NewInt(0)
	}

	return bigIntValue
}

// Receipt is a structure containing all the fields that need to be save for a Receipt
type Receipt struct {
	Hash      string        `json:"-"`
	Value     string        `json:"value"`
	Sender    string        `json:"sender"`
	Data      string        `json:"data,omitempty"`
	TxHash    string        `json:"txHash"`
	Timestamp time.Duration `json:"timestamp"`
}

// TxLog holds all the data needed for a log structure
type TxLog struct {
	Address string  `json:"scAddress"`
	Events  []Event `json:"events"`
}

// Event holds all the data needed for an event structure
type Event struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     []string `json:"topics"`
	Data       string   `json:"data"`
}

// ScResult is a structure containing all the fields that need to be saved for a smart contract result
type ScResult struct {
	Hash                string        `json:"-"`
	Nonce               uint64        `json:"nonce"`
	GasLimit            uint64        `json:"gasLimit"`
	GasPrice            uint64        `json:"gasPrice"`
	Value               string        `json:"value"`
	Sender              string        `json:"sender"`
	Receiver            string        `json:"receiver"`
	RelayerAddr         string        `json:"relayerAddr,omitempty"`
	RelayedValue        string        `json:"relayedValue,omitempty"`
	Code                string        `json:"code,omitempty"`
	Data                []byte        `json:"data,omitempty"`
	PreTxHash           string        `json:"prevTxHash"`
	OriginalTxHash      string        `json:"originalTxHash"`
	CallType            string        `json:"callType"`
	CodeMetadata        []byte        `json:"codeMetaData,omitempty"`
	ReturnMessage       string        `json:"returnMessage,omitempty"`
	Timestamp           time.Duration `json:"timestamp"`
	EsdtTokenIdentifier string        `json:"token,omitempty"`
	EsdtValue           string        `json:"esdtValue,omitempty"`
}

// Block is a structure containing all the fields that need
//  to be saved for a block. It has all the default fields
//  plus some extra information for ease of search and filter
type Block struct {
	Nonce                 uint64        `json:"nonce"`
	Round                 uint64        `json:"round"`
	Epoch                 uint32        `json:"epoch"`
	Hash                  string        `json:"-"`
	MiniBlocksHashes      []string      `json:"miniBlocksHashes"`
	NotarizedBlocksHashes []string      `json:"notarizedBlocksHashes"`
	Proposer              uint64        `json:"proposer"`
	Validators            []uint64      `json:"validators"`
	PubKeyBitmap          string        `json:"pubKeyBitmap"`
	Size                  int64         `json:"size"`
	SizeTxs               int64         `json:"sizeTxs"`
	Timestamp             time.Duration `json:"timestamp"`
	StateRootHash         string        `json:"stateRootHash"`
	PrevHash              string        `json:"prevHash"`
	ShardID               uint32        `json:"shardId"`
	TxCount               uint32        `json:"txCount"`
	AccumulatedFees       string        `json:"accumulatedFees"`
	DeveloperFees         string        `json:"developerFees"`
	EpochStartBlock       bool          `json:"epochStartBlock"`
	SearchOrder           uint64        `json:"searchOrder"`
}

//ValidatorsPublicKeys is a structure containing fields for validators public keys
type ValidatorsPublicKeys struct {
	PublicKeys []string `json:"publicKeys"`
}

// AccountInfo holds (serializable) data about an account
type AccountInfo struct {
	Address         string  `json:"address,omitempty"`
	Nonce           uint64  `json:"nonce,omitempty"`
	Balance         string  `json:"balance"`
	BalanceNum      float64 `json:"balanceNum"`
	TokenIdentifier string  `json:"token,omitempty"`
	Properties      string  `json:"properties,omitempty"`
	IsSender        bool    `json:"-"`
}

// AccountBalanceHistory represents an entry in the user accounts balances history
type AccountBalanceHistory struct {
	Address         string `json:"address"`
	Timestamp       int64  `json:"timestamp"`
	Balance         string `json:"balance"`
	TokenIdentifier string `json:"token,omitempty"`
	IsSender        bool   `json:"isSender,omitempty"`
}

// Miniblock is a structure containing miniblock information
type Miniblock struct {
	Hash              string        `json:"-"`
	SenderShardID     uint32        `json:"senderShard"`
	ReceiverShardID   uint32        `json:"receiverShard"`
	SenderBlockHash   string        `json:"senderBlockHash"`
	ReceiverBlockHash string        `json:"receiverBlockHash"`
	Type              string        `json:"type"`
	Timestamp         time.Duration `json:"timestamp"`
}

// TPS is a structure containing all the fields that need to
//  be saved for a shard statistic in the database
type TPS struct {
	LiveTPS               float64  `json:"liveTPS"`
	PeakTPS               float64  `json:"peakTPS"`
	BlockNumber           uint64   `json:"blockNumber"`
	RoundNumber           uint64   `json:"roundNumber"`
	RoundTime             uint64   `json:"roundTime"`
	AverageBlockTxCount   *big.Int `json:"averageBlockTxCount"`
	TotalProcessedTxCount *big.Int `json:"totalProcessedTxCount"`
	AverageTPS            *big.Int `json:"averageTPS"`
	CurrentBlockNonce     uint64   `json:"currentBlockNonce"`
	NrOfShards            uint32   `json:"nrOfShards"`
	NrOfNodes             uint32   `json:"nrOfNodes"`
	LastBlockTxCount      uint32   `json:"lastBlockTxCount"`
	ShardID               uint32   `json:"shardID"`
}

// KibanaResponse -
type KibanaResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// Options structure holds the indexer's configuration options
type Options struct {
	IndexerCacheSize int
	UseKibana        bool
}

// ValidatorRatingInfo is a structure containing validator rating information
type ValidatorRatingInfo struct {
	PublicKey string  `json:"-"`
	Rating    float32 `json:"rating"`
}

// RoundInfo is a structure containing block signers and shard id
type RoundInfo struct {
	Index            uint64        `json:"round"`
	SignersIndexes   []uint64      `json:"signersIndexes"`
	BlockWasProposed bool          `json:"blockWasProposed"`
	ShardId          uint32        `json:"shardId"`
	Timestamp        time.Duration `json:"timestamp"`
}

// EpochInfo holds the information about epoch
type EpochInfo struct {
	AccumulatedFees string `json:"accumulatedFees"`
	DeveloperFees   string `json:"developerFees"`
}

// AccountEGLD is a structure that is needed for EGLD accounts
type AccountEGLD struct {
	Account  state.UserAccountHandler
	IsSender bool
}