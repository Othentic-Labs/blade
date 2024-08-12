package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygon/polygon-edge/helper/common"
	"github.com/0xPolygon/polygon-edge/helper/hex"
	"github.com/0xPolygon/polygon-edge/types"
	"github.com/valyala/fastjson"
)

var (
	defaultArena fastjson.ArenaPool
	defaultPool  fastjson.ParserPool
)

const jsonRPCMetric = "json_rpc"

// For union type of transaction and types.Hash
type transactionOrHash interface {
	getHash() types.Hash
}

type transaction struct {
	Nonce       argUint64          `json:"nonce"`
	GasPrice    *argBig            `json:"gasPrice,omitempty"`
	GasTipCap   *argBig            `json:"maxPriorityFeePerGas,omitempty"`
	GasFeeCap   *argBig            `json:"maxFeePerGas,omitempty"`
	Gas         argUint64          `json:"gas"`
	To          *types.Address     `json:"to"`
	Value       argBig             `json:"value"`
	Input       argBytes           `json:"input"`
	V           argBig             `json:"v"`
	R           argBig             `json:"r"`
	S           argBig             `json:"s"`
	Hash        types.Hash         `json:"hash"`
	From        types.Address      `json:"from"`
	BlockHash   *types.Hash        `json:"blockHash"`
	BlockNumber *argUint64         `json:"blockNumber"`
	TxIndex     *argUint64         `json:"transactionIndex"`
	ChainID     *argBig            `json:"chainId,omitempty"`
	Type        argUint64          `json:"type"`
	AccessList  types.TxAccessList `json:"accessList,omitempty"`
}

func (t transaction) getHash() types.Hash { return t.Hash }

// Redefine to implement getHash() of transactionOrHash
type transactionHash types.Hash

func (h transactionHash) getHash() types.Hash { return types.Hash(h) }

func (h transactionHash) MarshalText() ([]byte, error) {
	return []byte(types.Hash(h).String()), nil
}

func toPendingTransaction(t *types.Transaction) *transaction {
	return toTransaction(t, nil, nil)
}

// toTransaction converts types.Transaction struct to JSON RPC transaction format
func toTransaction(
	t *types.Transaction,
	header *types.Header,
	txIndex *int,
) *transaction {
	v, r, s := t.RawSignatureValues()
	res := &transaction{
		Nonce: argUint64(t.Nonce()),
		Gas:   argUint64(t.Gas()),
		To:    t.To(),
		Value: argBig(*t.Value()),
		Input: t.Input(),
		V:     argBig(*v),
		R:     argBig(*r),
		S:     argBig(*s),
		Hash:  t.Hash(),
		From:  t.From(),
		Type:  argUint64(t.Type()),
	}

	if header != nil {
		// transaction is already mined
		res.BlockNumber = argUintPtr(header.Number)
		res.BlockHash = &header.Hash
		res.GasPrice = argBigPtr(t.GetGasPrice(header.BaseFee))
	} else if t.GasPrice() != nil {
		// transaction is pending (within the tx pool)
		res.GasPrice = argBigPtr(t.GasPrice())
	}

	if t.Type() == types.DynamicFeeTxType {
		if t.GasTipCap() != nil {
			res.GasTipCap = argBigPtr(t.GasTipCap())
		}

		if t.GasFeeCap() != nil {
			res.GasFeeCap = argBigPtr(t.GasFeeCap())
		}

		if res.GasPrice == nil {
			res.GasPrice = res.GasFeeCap
		}
	}

	if t.ChainID() != nil {
		chainID := argBig(*(t.ChainID()))
		res.ChainID = &chainID
	}

	if txIndex != nil {
		res.TxIndex = argUintPtr(uint64(*txIndex))
	}

	if t.AccessList() != nil {
		res.AccessList = t.AccessList()
	}

	return res
}

type header struct {
	ParentHash      types.Hash  `json:"parentHash"`
	Sha3Uncles      types.Hash  `json:"sha3Uncles"`
	Miner           argBytes    `json:"miner"`
	StateRoot       types.Hash  `json:"stateRoot"`
	TxRoot          types.Hash  `json:"transactionsRoot"`
	ReceiptsRoot    types.Hash  `json:"receiptsRoot"`
	LogsBloom       types.Bloom `json:"logsBloom"`
	Difficulty      argUint64   `json:"difficulty"`
	TotalDifficulty argUint64   `json:"totalDifficulty"`
	Number          argUint64   `json:"number"`
	GasLimit        argUint64   `json:"gasLimit"`
	GasUsed         argUint64   `json:"gasUsed"`
	Timestamp       argUint64   `json:"timestamp"`
	ExtraData       argBytes    `json:"extraData"`
	MixHash         types.Hash  `json:"mixHash"`
	Nonce           types.Nonce `json:"nonce"`
	Hash            types.Hash  `json:"hash"`
	BaseFee         argUint64   `json:"baseFeePerGas,omitempty"`
}

type accessListResult struct {
	Accesslist types.TxAccessList `json:"accessList"`
	Error      error              `json:"error,omitempty"`
	GasUsed    argUint64          `json:"gasUsed"`
}

type block struct {
	header
	Size         argUint64           `json:"size"`
	Transactions []transactionOrHash `json:"transactions"`
	Uncles       []types.Hash        `json:"uncles"`
}

func (b *block) Copy() *block {
	bb := new(block)
	*bb = *b

	bb.Miner = make([]byte, len(b.Miner))
	copy(bb.Miner[:], b.Miner[:])

	bb.ExtraData = make([]byte, len(b.ExtraData))
	copy(bb.ExtraData[:], b.ExtraData[:])

	return bb
}

func toBlock(b *types.Block, fullTx bool) *block {
	h := b.Header
	resHeader := header{
		ParentHash:      h.ParentHash,
		Sha3Uncles:      h.Sha3Uncles,
		Miner:           argBytes(h.Miner),
		StateRoot:       h.StateRoot,
		TxRoot:          h.TxRoot,
		ReceiptsRoot:    h.ReceiptsRoot,
		LogsBloom:       h.LogsBloom,
		Difficulty:      argUint64(h.Difficulty),
		TotalDifficulty: argUint64(h.Difficulty), // not needed for POS
		Number:          argUint64(h.Number),
		GasLimit:        argUint64(h.GasLimit),
		GasUsed:         argUint64(h.GasUsed),
		Timestamp:       argUint64(h.Timestamp),
		ExtraData:       argBytes(h.ExtraData),
		MixHash:         h.MixHash,
		Nonce:           h.Nonce,
		Hash:            h.Hash,
		BaseFee:         argUint64(h.BaseFee),
	}

	res := &block{
		header:       resHeader,
		Size:         argUint64(b.Size()),
		Transactions: []transactionOrHash{},
		Uncles:       []types.Hash{},
	}

	for idx, txn := range b.Transactions {
		if fullTx {
			res.Transactions = append(
				res.Transactions,
				toTransaction(
					txn,
					b.Header,
					&idx,
				),
			)
		} else {
			res.Transactions = append(
				res.Transactions,
				transactionHash(txn.Hash()),
			)
		}
	}

	for _, uncle := range b.Uncles {
		res.Uncles = append(res.Uncles, uncle.Hash)
	}

	return res
}

func toHeader(h *types.Header) *header {
	res := &header{
		ParentHash:      h.ParentHash,
		Sha3Uncles:      h.Sha3Uncles,
		Miner:           argBytes(h.Miner),
		StateRoot:       h.StateRoot,
		TxRoot:          h.TxRoot,
		ReceiptsRoot:    h.ReceiptsRoot,
		LogsBloom:       h.LogsBloom,
		Difficulty:      argUint64(h.Difficulty),
		TotalDifficulty: argUint64(h.Difficulty), // not needed for POS
		Number:          argUint64(h.Number),
		GasLimit:        argUint64(h.GasLimit),
		GasUsed:         argUint64(h.GasUsed),
		Timestamp:       argUint64(h.Timestamp),
		ExtraData:       argBytes(h.ExtraData),
		MixHash:         h.MixHash,
		Nonce:           h.Nonce,
		Hash:            h.Hash,
		BaseFee:         argUint64(h.BaseFee),
	}

	return res
}

type receipt struct {
	Root              types.Hash     `json:"root"`
	CumulativeGasUsed argUint64      `json:"cumulativeGasUsed"`
	LogsBloom         types.Bloom    `json:"logsBloom"`
	Logs              []*Log         `json:"logs"`
	Status            argUint64      `json:"status"`
	TxHash            types.Hash     `json:"transactionHash"`
	TxIndex           argUint64      `json:"transactionIndex"`
	BlockHash         types.Hash     `json:"blockHash"`
	BlockNumber       argUint64      `json:"blockNumber"`
	GasUsed           argUint64      `json:"gasUsed"`
	ContractAddress   *types.Address `json:"contractAddress"`
	FromAddr          types.Address  `json:"from"`
	ToAddr            *types.Address `json:"to"`
	Type              argUint64      `json:"type"`
}

func toReceipt(src *types.Receipt, tx *types.Transaction,
	txIndex uint64, header *types.Header, logs []*Log) *receipt {
	return &receipt{
		Root:              src.Root,
		CumulativeGasUsed: argUint64(src.CumulativeGasUsed),
		LogsBloom:         src.LogsBloom,
		Status:            argUint64(*src.Status),
		TxHash:            tx.Hash(),
		TxIndex:           argUint64(txIndex),
		BlockHash:         header.Hash,
		BlockNumber:       argUint64(header.Number),
		GasUsed:           argUint64(src.GasUsed),
		ContractAddress:   src.ContractAddress,
		FromAddr:          tx.From(),
		ToAddr:            tx.To(),
		Logs:              logs,
		Type:              argUint64(tx.Type()),
	}
}

type Log struct {
	Address     types.Address `json:"address"`
	Topics      []types.Hash  `json:"topics"`
	Data        argBytes      `json:"data"`
	BlockNumber argUint64     `json:"blockNumber"`
	TxHash      types.Hash    `json:"transactionHash"`
	TxIndex     argUint64     `json:"transactionIndex"`
	BlockHash   types.Hash    `json:"blockHash"`
	LogIndex    argUint64     `json:"logIndex"`
	Removed     bool          `json:"removed"`
}

func toLogs(srcLogs []*types.Log, baseIdx, txIdx uint64, header *types.Header, txHash types.Hash) []*Log {
	logs := make([]*Log, len(srcLogs))
	for i, srcLog := range srcLogs {
		logs[i] = toLog(srcLog, baseIdx+uint64(i), txIdx, header, txHash)
	}

	return logs
}

func toLog(src *types.Log, logIdx, txIdx uint64, header *types.Header, txHash types.Hash) *Log {
	return &Log{
		Address:     src.Address,
		Topics:      src.Topics,
		Data:        argBytes(src.Data),
		BlockNumber: argUint64(header.Number),
		BlockHash:   header.Hash,
		TxHash:      txHash,
		TxIndex:     argUint64(txIdx),
		LogIndex:    argUint64(logIdx),
	}
}

type argBig big.Int

func argBigPtr(b *big.Int) *argBig {
	v := argBig(*b)

	return &v
}

func (a *argBig) UnmarshalText(input []byte) error {
	buf, err := decodeToHex(input)
	if err != nil {
		return err
	}

	b := new(big.Int)
	b.SetBytes(buf)
	*a = argBig(*b)

	return nil
}

func (a argBig) MarshalText() ([]byte, error) {
	b := (*big.Int)(&a)

	return []byte("0x" + b.Text(16)), nil
}

func argAddrPtr(a types.Address) *types.Address {
	return &a
}

type argUint64 uint64

func argUintPtr(n uint64) *argUint64 {
	v := argUint64(n)

	return &v
}

func (u argUint64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 10)
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, uint64(u), 16)

	return buf, nil
}

func (u *argUint64) UnmarshalText(input []byte) error {
	str := strings.Trim(string(input), "\"")

	num, err := common.ParseUint64orHex(&str)
	if err != nil {
		return err
	}

	*u = argUint64(num)

	return nil
}

func (u *argUint64) UnmarshalJSON(buffer []byte) error {
	return u.UnmarshalText(buffer)
}

type argBytes []byte

func argBytesPtr(b []byte) *argBytes {
	bb := argBytes(b)

	return &bb
}

func (b argBytes) MarshalText() ([]byte, error) {
	return encodeToHex(b), nil
}

func (b *argBytes) UnmarshalText(input []byte) error {
	hh, err := decodeToHex(input)
	if err != nil {
		return nil
	}

	aux := make([]byte, len(hh))
	copy(aux[:], hh[:])
	*b = aux

	return nil
}

func decodeToHex(b []byte) ([]byte, error) {
	str := string(b)
	str = strings.TrimPrefix(str, "0x")

	if len(str)%2 != 0 {
		str = "0" + str
	}

	return hex.DecodeString(str)
}

func encodeToHex(b []byte) []byte {
	str := hex.EncodeToString(b)
	if len(str)%2 != 0 {
		str = "0" + str
	}

	return []byte("0x" + str)
}

// txnArgs is the transaction argument for the rpc endpoints
type txnArgs struct {
	From       *types.Address      `json:"from"`
	To         *types.Address      `json:"to"`
	Gas        *argUint64          `json:"gas"`
	GasPrice   *argBytes           `json:"gasPrice,omitempty"`
	GasTipCap  *argBytes           `json:"maxPriorityFeePerGas,omitempty"`
	GasFeeCap  *argBytes           `json:"maxFeePerGas,omitempty"`
	Value      *argBytes           `json:"value"`
	Data       *argBytes           `json:"data"`
	Input      *argBytes           `json:"input"`
	Nonce      *argUint64          `json:"nonce"`
	Type       *argUint64          `json:"type"`
	AccessList *types.TxAccessList `json:"accessList,omitempty"`
	ChainID    *argUint64          `json:"chainId,omitempty"`
}

// data retrieves the transaction calldata. Input field is preferred.
func (args *txnArgs) data() []byte {
	if args.Input != nil {
		return *args.Input
	}

	if args.Data != nil {
		return *args.Data
	}

	return nil
}

func (args *txnArgs) setDefaults(priceLimit uint64, eth *Eth) error {
	if err := args.setFeeDefaults(priceLimit, eth.store); err != nil {
		return err
	}

	if args.Nonce == nil {
		args.Nonce = argUintPtr(eth.store.GetNonce(*args.From))
	}

	if args.Gas == nil {
		// These fields are immutable during the estimation, safe to
		// pass the pointer directly.
		data := args.data()
		callArgs := txnArgs{
			From:      args.From,
			To:        args.To,
			GasPrice:  args.GasPrice,
			GasTipCap: args.GasTipCap,
			GasFeeCap: args.GasFeeCap,
			Value:     args.Value,
			Data:      argBytesPtr(data),
		}

		estimatedGas, err := eth.EstimateGas(&callArgs, nil)
		if err != nil {
			return err
		}

		estimatedGasUint64, ok := estimatedGas.(argUint64)
		if !ok {
			return errors.New("estimated gas not a uint64")
		}

		args.Gas = &estimatedGasUint64
	}

	// If chain id is provided, ensure it matches the local chain id. Otherwise, set the local
	// chain id as the default.
	want := eth.chainID

	if args.ChainID != nil {
		have := (uint64)(*args.ChainID)
		if have != want {
			return fmt.Errorf("chainId does not match node's (have=%v, want=%v)", have, want)
		}
	} else {
		args.ChainID = argUintPtr(want)
	}

	return nil
}

// setFeeDefaults fills in default fee values for unspecified tx fields.
func (args *txnArgs) setFeeDefaults(priceLimit uint64, store ethStore) error {
	// If both gasPrice and at least one of the EIP-1559 fee parameters are specified, error.
	if args.GasPrice != nil && (args.GasFeeCap != nil || args.GasTipCap != nil) {
		return errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}

	// If the tx has completely specified a fee mechanism, no default is needed.
	// This allows users who are not yet synced past London to get defaults for
	// other tx values. See https://github.com/ethereum/go-ethereum/pull/23274
	// for more information.
	eip1559ParamsSet := args.GasFeeCap != nil && args.GasTipCap != nil

	// Sanity check the EIP-1559 fee parameters if present.
	if args.GasPrice == nil && eip1559ParamsSet {
		maxFeePerGas := new(big.Int).SetBytes(*args.GasFeeCap)
		maxPriorityFeePerGas := new(big.Int).SetBytes(*args.GasTipCap)

		if maxFeePerGas.Sign() == 0 {
			return errors.New("maxFeePerGas must be non-zero")
		}

		if maxFeePerGas.Cmp(maxPriorityFeePerGas) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.GasFeeCap, args.GasTipCap)
		}

		args.Type = argUintPtr(uint64(types.DynamicFeeTxType))

		return nil // No need to set anything, user already set MaxFeePerGas and MaxPriorityFeePerGas
	}

	// Sanity check the non-EIP-1559 fee parameters.
	head := store.Header()
	isLondon := store.GetForksInTime(head.Number).London

	if args.GasPrice != nil && !eip1559ParamsSet {
		// Zero gas-price is not allowed after London fork
		if new(big.Int).SetBytes(*args.GasPrice).Sign() == 0 && isLondon {
			return errors.New("gasPrice must be non-zero after london fork")
		}

		return nil // No need to set anything, user already set GasPrice
	}

	// Now attempt to fill in default value depending on whether London is active or not.
	if isLondon {
		// London is active, set maxPriorityFeePerGas and maxFeePerGas.
		if err := args.setLondonFeeDefaults(head, store); err != nil {
			return err
		}
	} else {
		if args.GasFeeCap != nil || args.GasTipCap != nil {
			return errors.New("maxFeePerGas and maxPriorityFeePerGas are not valid before London is active")
		}

		// London not active, set gas price.
		avgGasPrice := store.GetAvgGasPrice()

		args.GasPrice = argBytesPtr(common.BigMax(new(big.Int).SetUint64(priceLimit), avgGasPrice).Bytes())
	}

	return nil
}

// setLondonFeeDefaults fills in reasonable default fee values for unspecified fields.
func (args *txnArgs) setLondonFeeDefaults(head *types.Header, store ethStore) error {
	// Set maxPriorityFeePerGas if it is missing.
	if args.GasTipCap == nil {
		tip, err := store.MaxPriorityFeePerGas()
		if err != nil {
			return err
		}

		args.GasTipCap = argBytesPtr(tip.Bytes())
	}

	// Set maxFeePerGas if it is missing.
	if args.GasFeeCap == nil {
		// Set the max fee to be 2 times larger than the previous block's base fee.
		// The additional slack allows the tx to not become invalidated if the base
		// fee is rising.
		val := new(big.Int).Add(
			new(big.Int).SetBytes(*args.GasTipCap),
			new(big.Int).Mul(new(big.Int).SetUint64(head.BaseFee), big.NewInt(2)),
		)
		args.GasFeeCap = argBytesPtr(val.Bytes())
	}

	// Both EIP-1559 fee parameters are now set; sanity check them.
	if new(big.Int).SetBytes(*args.GasFeeCap).Cmp(new(big.Int).SetBytes(*args.GasTipCap)) < 0 {
		return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.GasFeeCap, args.GasTipCap)
	}

	args.Type = argUintPtr(uint64(types.DynamicFeeTxType))

	return nil
}

type progression struct {
	Type          string    `json:"type"`
	StartingBlock argUint64 `json:"startingBlock"`
	CurrentBlock  argUint64 `json:"currentBlock"`
	HighestBlock  argUint64 `json:"highestBlock"`
}

type feeHistoryResult struct {
	OldestBlock   argUint64     `json:"oldestBlock"`
	BaseFeePerGas []argUint64   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio  []float64     `json:"gasUsedRatio"`
	Reward        [][]argUint64 `json:"reward,omitempty"`
}

func convertToArgUint64Slice(slice []uint64) []argUint64 {
	argSlice := make([]argUint64, len(slice))
	for i, value := range slice {
		argSlice[i] = argUint64(value)
	}

	return argSlice
}

func convertToArgUint64SliceSlice(slice [][]uint64) [][]argUint64 {
	argSlice := make([][]argUint64, len(slice))
	for i, value := range slice {
		argSlice[i] = convertToArgUint64Slice(value)
	}

	return argSlice
}

type OverrideAccount struct {
	Nonce     *argUint64                 `json:"nonce"`
	Code      *argBytes                  `json:"code"`
	Balance   *argUint64                 `json:"balance"`
	State     *map[types.Hash]types.Hash `json:"state"`
	StateDiff *map[types.Hash]types.Hash `json:"stateDiff"`
}

func (o *OverrideAccount) ToType() types.OverrideAccount {
	res := types.OverrideAccount{}

	if o.Nonce != nil {
		res.Nonce = (*uint64)(o.Nonce)
	}

	if o.Code != nil {
		res.Code = *o.Code
	}

	if o.Balance != nil {
		res.Balance = new(big.Int).SetUint64(*(*uint64)(o.Balance))
	}

	if o.State != nil {
		res.State = *o.State
	}

	if o.StateDiff != nil {
		res.StateDiff = *o.StateDiff
	}

	return res
}

// StateOverride is the collection of overridden accounts
type StateOverride map[types.Address]OverrideAccount

// MarshalJSON marshals the StateOverride to JSON
func (s StateOverride) MarshalJSON() ([]byte, error) {
	a := defaultArena.Get()
	defer a.Reset()

	o := a.NewObject()

	for addr, obj := range s {
		oo := a.NewObject()
		if obj.Nonce != nil {
			oo.Set("nonce", a.NewString(fmt.Sprintf("0x%x", *obj.Nonce)))
		}

		if obj.Balance != nil {
			oo.Set("balance", a.NewString(fmt.Sprintf("0x%x", obj.Balance)))
		}

		if obj.Code != nil {
			oo.Set("code", a.NewString("0x"+hex.EncodeToString(*obj.Code)))
		}

		if obj.State != nil {
			ooo := a.NewObject()
			for k, v := range *obj.State {
				ooo.Set(k.String(), a.NewString(v.String()))
			}

			oo.Set("state", ooo)
		}

		if obj.StateDiff != nil {
			ooo := a.NewObject()
			for k, v := range *obj.StateDiff {
				ooo.Set(k.String(), a.NewString(v.String()))
			}

			oo.Set("stateDiff", ooo)
		}

		o.Set(addr.String(), oo)
	}

	res := o.MarshalTo(nil)

	defaultArena.Put(a)

	return res, nil
}

// CallMsg contains parameters for contract calls
type CallMsg struct {
	From       types.Address  // the sender of the 'transaction'
	To         *types.Address // the destination contract (nil for contract creation)
	Gas        uint64         // if 0, the call executes with near-infinite gas
	GasPrice   *big.Int       // wei <-> gas exchange ratio
	GasFeeCap  *big.Int       // EIP-1559 fee cap per gas
	GasTipCap  *big.Int       // EIP-1559 tip per gas
	Value      *big.Int       // amount of wei sent along with the call
	Data       []byte         // input data, usually an ABI-encoded contract method invocation
	Type       uint64
	AccessList types.TxAccessList // EIP-2930 access list
	ChainID    *big.Int           // ChainID
	Nonce      uint64             // nonce
}

// MarshalJSON implements the Marshal interface.
func (c *CallMsg) MarshalJSON() ([]byte, error) {
	a := defaultArena.Get()
	defer a.Reset()

	o := a.NewObject()
	o.Set("from", a.NewString(c.From.String()))

	if c.Gas != 0 {
		o.Set("gas", a.NewString(fmt.Sprintf("0x%x", c.Gas)))
	}

	if c.To != nil {
		o.Set("to", a.NewString(c.To.String()))
	}

	if len(c.Data) != 0 {
		o.Set("data", a.NewString("0x"+hex.EncodeToString(c.Data)))
	}

	if c.GasPrice != nil {
		o.Set("gasPrice", a.NewString(fmt.Sprintf("0x%x", c.GasPrice)))
	}

	if c.Value != nil {
		o.Set("value", a.NewString(fmt.Sprintf("0x%x", c.Value)))
	}

	if c.GasFeeCap != nil {
		o.Set("maxFeePerGas", a.NewString(fmt.Sprintf("0x%x", c.GasFeeCap)))
	}

	if c.GasTipCap != nil {
		o.Set("maxPriorityFeePerGas", a.NewString(fmt.Sprintf("0x%x", c.GasTipCap)))
	}

	if c.AccessList != nil {
		o.Set("accessList", c.AccessList.MarshalJSONWith(a))
	}

	if c.ChainID != nil {
		o.Set("chainID", a.NewString(fmt.Sprintf("0x%x", c.ChainID)))
	}

	o.Set("nonce", a.NewString(fmt.Sprintf("0x%x", c.Nonce)))
	o.Set("type", a.NewString(fmt.Sprintf("0x%x", c.Type)))

	res := o.MarshalTo(nil)

	defaultArena.Put(a)

	return res, nil
}

// FeeHistory represents the fee history data returned by an rpc node
type FeeHistory struct {
	OldestBlock  uint64     `json:"oldestBlock"`
	Reward       [][]uint64 `json:"reward,omitempty"`
	BaseFee      []uint64   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64  `json:"gasUsedRatio"`
}

// UnmarshalJSON unmarshals the FeeHistory object from JSON
func (f *FeeHistory) UnmarshalJSON(data []byte) error {
	var raw feeHistoryResult

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	f.OldestBlock = uint64(raw.OldestBlock)

	if raw.Reward != nil {
		f.Reward = make([][]uint64, 0, len(raw.Reward))

		for _, r := range raw.Reward {
			elem := make([]uint64, 0, len(r))
			for _, i := range r {
				elem = append(elem, uint64(i))
			}

			f.Reward = append(f.Reward, elem)
		}
	}

	f.BaseFee = make([]uint64, 0, len(raw.BaseFeePerGas))
	for _, i := range raw.BaseFeePerGas {
		f.BaseFee = append(f.BaseFee, uint64(i))
	}

	f.GasUsedRatio = raw.GasUsedRatio

	return nil
}

// Transaction is the json rpc transaction object
// (types.Transaction object, expanded with block number, hash and index)
type Transaction struct {
	*types.Transaction

	// BlockNumber is the number of the block in which the transaction was included.
	BlockNumber uint64 `json:"blockNumber"`

	// BlockHash is the hash of the block in which the transaction was included.
	BlockHash types.Hash `json:"blockHash"`

	// TxnIndex is the index of the transaction within the block.
	TxnIndex uint64 `json:"transactionIndex"`
}

// UnmarshalJSON unmarshals the transaction object from JSON
func (t *Transaction) UnmarshalJSON(data []byte) error {
	p := defaultPool.Get()
	defer defaultPool.Put(p)

	v, err := p.Parse(string(data))
	if err != nil {
		return err
	}

	t.Transaction = new(types.Transaction)
	if err := t.Transaction.UnmarshalJSONWith(v); err != nil {
		return err
	}

	if types.HasJSONKey(v, "blockNumber") {
		t.BlockNumber, err = types.UnmarshalJSONUint64(v, "blockNumber")
		if err != nil {
			return err
		}
	}

	if types.HasJSONKey(v, "blockHash") {
		t.BlockHash, err = types.UnmarshalJSONHash(v, "blockHash")
		if err != nil {
			return err
		}
	}

	if types.HasJSONKey(v, "transactionIndex") {
		t.TxnIndex, err = types.UnmarshalJSONUint64(v, "transactionIndex")
		if err != nil {
			return err
		}
	}

	return nil
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw argBytes           `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}
