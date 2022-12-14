package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestSplitNFTContent(t *testing.T) {
//	hash := "7b8c47ff0f29187c4fd7b9404d6d8671c3a05d041a2126753722fe940f30e2d3"
//	fmt.Println("len", len(hash))
//	a, b, err := SplitNFTContent(hash)
//	assert.Nil(t, err)
//	t.Log("a", a.Text(16), "b", b.Text(16))
//	t.Log("a", a.BitLen(), "b", b.BitLen())
//}

func TestFindExponent(t *testing.T) {
	s := "12304"
	r, err := ZkFindExponentPart(s)
	assert.Nil(t, err)
	assert.True(t, r == 0)

	s = "123040"
	r, err = ZkFindExponentPart(s)
	assert.Nil(t, err)
	assert.True(t, r == 1)

	s = "0"
	r, err = ZkFindExponentPart(s)
	assert.Nil(t, err)
	assert.True(t, r == 0)

	s = "12"
	for i := 0; i < 33; i++ {
		s += "0"
	}
	r, err = ZkFindExponentPart(s)
	assert.Nil(t, err)
	assert.True(t, r == 31)
	//fmt.Println("s",s)
	//fmt.Println("s.len",len(s),"exp",r,"s",s[0:len(s)-r])

}

func TestFindManExpPart(t *testing.T) {
	s := "12304"
	m, e, err := ZkTransferManExpPart(s)
	assert.Nil(t, err)
	assert.True(t, m == s)
	assert.True(t, e == 0)

	s = "123040"
	m, e, err = ZkTransferManExpPart(s)
	assert.Nil(t, err)
	assert.True(t, m == "12304")
	assert.True(t, e == 1)

	s = "0"
	m, e, err = ZkTransferManExpPart(s)
	assert.Nil(t, err)
	assert.True(t, m == "0")
	assert.True(t, e == 0)

	s = "12"
	for i := 0; i < 31; i++ {
		s += "0"
	}
	m, e, err = ZkTransferManExpPart(s)
	assert.Nil(t, err)
	assert.True(t, m == "12")
	assert.True(t, e == 31)

	s = "12"
	for i := 0; i < 30; i++ {
		s += "0"
	}
	m, e, err = ZkTransferManExpPart(s)
	assert.Nil(t, err)
	assert.True(t, m == "12")
	assert.True(t, e == 30)

	s = "12"
	for i := 0; i < 32; i++ {
		s += "0"
	}
	m, e, err = ZkTransferManExpPart(s)
	assert.Nil(t, err)
	assert.True(t, m == "120")
	assert.True(t, e == 31)
}
