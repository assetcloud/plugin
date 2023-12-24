// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/holiman/uint256"

	"github.com/assetcloud/plugin/plugin/dapp/evm/executor/vm/common"
)

// Tracer 接口用来在合约执行过程中收集跟踪数据。
// CaptureState 会在EVM解释每条指令时调用。
// 需要注意的是，传入的引用参数不允许修改，否则会影响EVM解释执行；如果需要使用其中的数据，请复制后使用。
type Tracer interface {
	// CaptureStart 开始记录
	CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value uint64) error
	// CaptureState 保存状态
	CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, rData []byte, contract *Contract, depth int, err error) error
	// CaptureFault 保存错误
	CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error
	// CaptureEnd 结束记录
	CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error
}

// JSONLogger 使用json格式打印日志
type JSONLogger struct {
	encoder *json.Encoder
}

// Storage represents a contract's storage.
type Storage map[common.Hash]common.Hash

// Copy duplicates the current storage.
func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}
	return cpy
}

// LogConfig are the configuration options for structured logger the EVM
type LogConfig struct {
	DisableMemory     bool // disable memory capture
	DisableStack      bool // disable stack capture
	DisableStorage    bool // disable storage capture
	DisableReturnData bool // disable return data capture
	Debug             bool // print output during capture end
	Limit             int  // maximum length of output, but zero means unlimited
}

// StructLog 指令执行状态信息
type StructLog struct {
	// Pc pc指针
	Pc uint64 `json:"pc"`
	// Op 操作码
	Op OpCode `json:"op"`
	// Gas gas
	Gas uint64 `json:"gas"`
	// GasCost 花费
	GasCost uint64 `json:"gasCost"`
	// Memory 内存对象
	Memory []string `json:"memory"`
	// MemorySize 内存大小
	MemorySize int `json:"memSize"`
	// Stack 栈对象
	Stack []*big.Int `json:"stack"`
	// ReturnStack 返回栈
	ReturnStack []uint32 `json:"returnStack"`
	// ReturnData 返回数据
	ReturnData []byte `json:"returnData"`
	// Storage 存储对象
	Storage map[common.Hash]common.Hash `json:"-"`
	// Depth 调用深度
	Depth int `json:"depth"`
	// RefundCounter 退款统计
	RefundCounter uint64 `json:"refund"`
	// Err 错误信息
	Err error `json:"-"`
}

// NewJSONLogger 创建新的日志记录器
func NewJSONLogger(writer io.Writer) *JSONLogger {
	return &JSONLogger{json.NewEncoder(writer)}
}

// CaptureStart 开始记录
func (logger *JSONLogger) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value uint64) error {
	return nil
}

// CaptureState 输出当前虚拟机状态
func (logger *JSONLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, rData []byte, contract *Contract, depth int, err error) error {
	log := StructLog{
		Pc:         pc,
		Op:         op,
		Gas:        gas,
		GasCost:    cost,
		MemorySize: memory.Len(),
		Storage:    nil,
		Depth:      depth,
		Err:        err,
	}
	log.Memory = formatMemory(memory.Data())
	log.Stack = formatStack(stack.Data())
	log.ReturnData = rData
	return logger.encoder.Encode(log)
}

func formatStack(data []uint256.Int) (res []*big.Int) {
	for _, v := range data {
		res = append(res, v.ToBig())
	}
	return
}

func formatMemory(data []byte) (res []string) {
	for idx := 0; idx < len(data); idx += 32 {
		res = append(res, common.Bytes2HexTrim(data[idx:idx+32]))
	}
	return
}

// CaptureFault 目前实现为空
func (logger *JSONLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {
	return nil
}

// CaptureEnd 结束记录
func (logger *JSONLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	type endLog struct {
		Output  string        `json:"output"`
		GasUsed int64         `json:"gasUsed"`
		Time    time.Duration `json:"time"`
		Err     string        `json:"error,omitempty"`
	}

	if err != nil {
		return logger.encoder.Encode(endLog{common.Bytes2Hex(output), int64(gasUsed), t, err.Error()})
	}
	return logger.encoder.Encode(endLog{common.Bytes2Hex(output), int64(gasUsed), t, ""})
}

type mdLogger struct {
	out io.Writer
	cfg *LogConfig
}

// NewMarkdownLogger creates a logger which outputs information in a format adapted
// for human readability, and is also a valid markdown table
func NewMarkdownLogger(cfg *LogConfig, writer io.Writer) *mdLogger {
	l := &mdLogger{writer, cfg}
	if l.cfg == nil {
		l.cfg = &LogConfig{}
	}
	return l
}

func (t *mdLogger) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value uint64) error {
	if !create {
		fmt.Fprintf(t.out, "From: `%v`\nTo: `%v`\nData: `0x%x`\nGas: `%d`\nValue `%v` wei\n",
			from.String(), to.String(),
			input, gas, value)
	} else {
		fmt.Fprintf(t.out, "From: `%v`\nCreate at: `%v`\nData: `0x%x`\nGas: `%d`\nValue `%v` wei\n",
			from.String(), to.String(),
			input, gas, value)
	}

	fmt.Fprintf(t.out, `
|  Pc   |      Op     | Cost |   Stack   |   RStack  |  Refund |
|-------|-------------|------|-----------|-----------|---------|
`)
	return nil
}

func (t *mdLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, rData []byte, contract *Contract, depth int, err error) error {
	fmt.Fprintf(t.out, "| %4d  | %10v  |  %3d |", pc, op, cost)

	if !t.cfg.DisableStack {
		// format stack
		var a []string
		for _, elem := range stack.data {
			a = append(a, fmt.Sprintf("%v", elem.String()))
		}
		b := fmt.Sprintf("[%v]", strings.Join(a, ","))
		fmt.Fprintf(t.out, "%10v |", b)
	}
	fmt.Fprintf(t.out, "%10v |", env.StateDB.GetRefund())
	fmt.Fprintln(t.out, "")
	if err != nil {
		fmt.Fprintf(t.out, "Error: %v\n", err)
	}
	return nil
}

func (t *mdLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost uint64, memory *Memory, stack *Stack, contract *Contract, depth int, err error) error {

	fmt.Fprintf(t.out, "\nError: at pc=%d, op=%v: %v\n", pc, op, err)

	return nil
}

func (t *mdLogger) CaptureEnd(output []byte, gasUsed uint64, tm time.Duration, err error) error {
	fmt.Fprintf(t.out, "\nOutput: `0x%x`\nConsumed gas: `%d`\nError: `%v`\n",
		output, gasUsed, err)
	return nil
}
