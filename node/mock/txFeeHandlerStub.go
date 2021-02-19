package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/process"
)

// FeeHandlerStub -
type FeeHandlerStub struct {
	MaxGasLimitPerBlockCalled     func() uint64
	SetMinGasPriceCalled          func(minasPrice uint64)
	SetMinGasLimitCalled          func(minGasLimit uint64)
	ComputeGasLimitCalled         func(tx process.TransactionWithFeeHandler) uint64
	ComputeMoveBalanceFeeCalled   func(tx process.TransactionWithFeeHandler) *big.Int
	ComputeTxFeeCalled            func(tx process.TransactionWithFeeHandler) *big.Int
	CheckValidityTxValuesCalled   func(tx process.TransactionWithFeeHandler) error
	DeveloperPercentageCalled     func() float64
	MinGasPriceCalled             func() uint64
	GasPriceModifierCalled        func() float64
	ComputeFeeForProcessingCalled func(tx process.TransactionWithFeeHandler, gasToUse uint64) *big.Int
	GenesisTotalSupplyCalled      func() *big.Int
}

// ComputeFeeForProcessing -
func (fhs *FeeHandlerStub) ComputeFeeForProcessing(tx process.TransactionWithFeeHandler, gasToUse uint64) *big.Int {
	if fhs.ComputeFeeForProcessingCalled != nil {
		return fhs.ComputeFeeForProcessingCalled(tx, gasToUse)
	}
	return big.NewInt(0)
}

// GasPriceModifier -
func (fhs *FeeHandlerStub) GasPriceModifier() float64 {
	if fhs.GasPriceModifierCalled != nil {
		return fhs.GasPriceModifierCalled()
	}
	return 1.0
}

// MinGasPrice -
func (fhs *FeeHandlerStub) MinGasPrice() uint64 {
	if fhs.MinGasPriceCalled != nil {
		return fhs.MinGasPriceCalled()
	}
	return 0
}

// DeveloperPercentage -
func (fhs *FeeHandlerStub) DeveloperPercentage() float64 {
	return fhs.DeveloperPercentageCalled()
}

// MaxGasLimitPerBlock -
func (fhs *FeeHandlerStub) MaxGasLimitPerBlock(uint32) uint64 {
	return fhs.MaxGasLimitPerBlockCalled()
}

// ComputeGasLimit -
func (fhs *FeeHandlerStub) ComputeGasLimit(tx process.TransactionWithFeeHandler) uint64 {
	return fhs.ComputeGasLimitCalled(tx)
}

// ComputeMoveBalanceFee -
func (fhs *FeeHandlerStub) ComputeMoveBalanceFee(tx process.TransactionWithFeeHandler) *big.Int {
	return fhs.ComputeMoveBalanceFeeCalled(tx)
}

// ComputeTxFee -
func (fhs *FeeHandlerStub) ComputeTxFee(tx process.TransactionWithFeeHandler) *big.Int {
	return fhs.ComputeTxFeeCalled(tx)
}

// CheckValidityTxValues -
func (fhs *FeeHandlerStub) CheckValidityTxValues(tx process.TransactionWithFeeHandler) error {
	if fhs.CheckValidityTxValuesCalled != nil {
		return fhs.CheckValidityTxValuesCalled(tx)
	}

	return nil
}

// GenesisTotalSupply -
func (fhs *FeeHandlerStub) GenesisTotalSupply() *big.Int {
	if fhs.GenesisTotalSupplyCalled != nil {
		return fhs.GenesisTotalSupplyCalled()
	}

	return big.NewInt(0)
}

// IsInterfaceNil returns true if there is no value under the interface
func (fhs *FeeHandlerStub) IsInterfaceNil() bool {
	return fhs == nil
}
