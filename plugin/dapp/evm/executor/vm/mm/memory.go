// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm

import (
	"fmt"

	"github.com/assetcloud/chain/common/log/log15"
	"github.com/holiman/uint256"
)

// Memory 内存操作封装，在EVM中使用此对象模拟物理内存
type Memory struct {
	// Store 内存中存储的数据
	Store []byte
	// LastGasCost 上次开辟内存消耗的Gas
	LastGasCost uint64
}

// NewMemory 创建内存对象结构
func NewMemory() *Memory {
	return &Memory{}
}

// Set 设置内存中的值， value => offset:offset + size
func (m *Memory) Set(offset, size uint64, value []byte) (err error) {
	if size > 0 {
		// 偏移量+大小一定不会大于内存长度
		if offset+size > uint64(len(m.Store)) {
			err = fmt.Errorf("INVALID memory access, memory size:%v, offset:%v, size:%v", len(m.Store), offset, size)
			log15.Crit(err.Error())
			//panic("invalid memory: store empty")
			return err
		}
		copy(m.Store[offset:offset+size], value)
	}
	return nil
}

// Set32 从offset开始设置32个字节的内存值，如果值长度不足32个字节，左零值填充
func (m *Memory) Set32(offset uint64, val *uint256.Int) (err error) {

	// 确保长度足够设置值
	if offset+32 > uint64(len(m.Store)) {
		err = fmt.Errorf("INVALID memory access, memory size:%v, offset:%v, size:%v", len(m.Store), offset, 32)
		log15.Crit(err.Error())
		//panic("invalid memory: store empty")
		return err
	}
	// 先填充零值
	copy(m.Store[offset:offset+32], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	// Fill in relevant bits
	val.WriteToSlice(m.Store[offset:])

	return nil
}

// Resize 扩充内存到指定大小
func (m *Memory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.Store = append(m.Store, make([]byte, size-uint64(m.Len()))...)
	}
}

// Get 获取内存中制定偏移量开始的指定长度的数据，返回数据的拷贝而非引用
func (m *Memory) Get(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.Store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.Store[offset:offset+size])

		return
	}

	return
}

// Get returns offset + size as a new slice
func (m *Memory) GetCopy(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.Store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.Store[offset:offset+size])

		return
	}

	return
}

// GetPtr 同Get操作，不过这里返回的是数据引用
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(m.Store) > int(offset) {
		return m.Store[offset : offset+size]
	}

	return nil
}

// Len 返回内存中已开辟空间的大小（以字节计算）
func (m *Memory) Len() int {
	return len(m.Store)
}

// Data 返回内存中的原始数据引用
func (m *Memory) Data() []byte {
	return m.Store
}

// Print 打印内存中的数据（调试用）
func (m *Memory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.Store))
	if len(m.Store) > 0 {
		addr := 0
		for i := 0; i+32 <= len(m.Store); i += 32 {
			fmt.Printf("%03d: % x\n", addr, m.Store[i:i+32])
			addr++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("####################")
}

// 计算所需的内存偏移量和数据大小，计算所需内存大小
func calcMemSize64(off, l *uint256.Int) (uint64, bool) {
	if !l.IsUint64() {
		return 0, true
	}
	return calcMemSize64WithUint(off, l.Uint64())
}

// calcMemSize64WithUint calculates the required memory size, and returns
// the size and whether the result overflowed uint64
// Identical to calcMemSize64, but length is a uint64
func calcMemSize64WithUint(off *uint256.Int, length64 uint64) (uint64, bool) {
	// if length is zero, memsize is always zero, regardless of offset
	if length64 == 0 {
		return 0, false
	}
	// Check that offset doesn't overflow
	offset64, overflow := off.Uint64WithOverflow()
	if overflow {
		return 0, true
	}
	val := offset64 + length64
	// if value < either of it's parts, then it overflowed
	return val, val < offset64
}
