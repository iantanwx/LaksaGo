package LaksaGo

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"
)

func Pack(a int, b int) int {
	return a<<16 + b
}

func EncodeHex(src []byte) string {
	return hex.EncodeToString(src)
}

func DecodeHex(src string) []byte {
	ret, _ := hex.DecodeString(src)
	return ret
}

func Sha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func Compress(curve elliptic.Curve, x, y *big.Int, compress bool) []byte {
	return Marshal(curve, x, y, compress)
}

func Marshal(curve elliptic.Curve, x, y *big.Int, compress bool) []byte {
	byteLen := (curve.Params().BitSize + 7) >> 3

	if compress {
		ret := make([]byte, 1+byteLen)
		if y.Bit(0) == 0 {
			ret[0] = 2
		} else {
			ret[0] = 3
		}
		xBytes := x.Bytes()
		copy(ret[1+byteLen-len(xBytes):], xBytes)
		return ret
	}

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point
	xBytes := x.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes := y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)
	return ret
}

func bigIntToBytes(bi *big.Int) []byte {
	b1, b2 := [32]byte{}, bi.Bytes()
	copy(b1[32-len(b2):], b2)
	return b1[:]
}

func Hash(Q []byte, pubKey []byte, msg []byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(Q)
	buffer.Write(pubKey[:33])
	buffer.Write(msg)
	return Sha256(buffer.Bytes())
}

func GenerateMac(derivedKey, cipherText []byte) []byte {
	result := make([]byte, 16+len(cipherText))
	copy(result[0:16], derivedKey[16:])
	copy(result[16:], cipherText[:])
	return Sha256(result)
}

func ToCheckSumAddress(address string) string {
	lowerAddress := strings.ToLower(address)
	ar := strings.ReplaceAll(lowerAddress, "0x", "")
	hash := Sha256(DecodeHex(ar))
	v := new(big.Int).SetBytes(hash)
	sb := strings.Builder{}
	sb.WriteString("0x")

	for i := 0; i < len(ar); i++ {
		if strings.IndexByte("1234567890", ar[i]) != -1 {
			sb.WriteByte(ar[i])
		} else {
			checker := new(big.Int).And(v, new(big.Int).Exp(new(big.Int).SetInt64(2), new(big.Int).SetInt64(int64(255-6*i)), nil))
			r := checker.Cmp(new(big.Int).SetInt64(1))
			if r < 0 {
				sb.WriteString(strings.ToLower(string(ar[i])))
			} else {
				sb.WriteString(strings.ToUpper(string(ar[i])))
			}
		}
	}

	return sb.String()
}

func IntToHex(value, size int) string {
	hexVal := strconv.FormatInt(int64(value), 16)
	hexRep := make([]byte, len(hexVal))

	for i := 0; i < len(hexVal); i++ {
		hexRep[i] = hexVal[i]
	}

	hex := make([]byte, size, size)

	for i := 0; i < size-len(hexVal); i++ {
		hex = append(hex, '0')
	}

	for i := 0; i < len(hexVal); i++ {
		hex = append(hex, hexVal[i])
	}

	var hexFixed [16]byte
	copy(hexFixed[:],hex[len(hex)-16:])
	sb := strings.Builder{}


	for _, v := range hexFixed {
		sb.WriteByte(v)
	}

	return sb.String()

}
