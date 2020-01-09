package main

import (
	"fmt"
	"github.com/willf/bitset"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	encode()
	decode()
}
func encode() {
	plaintext, err0 := ioutil.ReadFile("plaintext")
	length := uint(len(plaintext) << 3)
	gen_code(length, "encode_stream")
	code, err1 := ioutil.ReadFile("encode_stream")
	if err1 != nil || err0 != nil {
		println(err1)
		os.Exit(1)
	}
	ciphertext := make([]byte, len(plaintext))
	for r := 0; r < len(plaintext); r++ {
		ciphertext[r] = plaintext[r] ^ code[r]
	}
	err2 := ioutil.WriteFile("ciphertext", ciphertext, 0644)

	if err2 != nil {
		println(err2)
		os.Exit(1)
	}
}
func decode() {
	ciphertext, err0 := ioutil.ReadFile("ciphertext")
	length := uint(len(ciphertext) << 3)
	gen_code(length, "decode_stream")
	code, err1 := ioutil.ReadFile("decode_stream")
	if err1 != nil || err0 != nil {
		println(err1)
		os.Exit(1)
	}
	decodedtext := make([]byte, len(ciphertext))
	for r := 0; r < len(ciphertext); r++ {
		decodedtext[r] = ciphertext[r] ^ code[r]
	}
	err2 := ioutil.WriteFile("decodedtext", decodedtext, 0644)

	if err2 != nil {
		println(err2)
		os.Exit(1)
	}
}

func gen_code(length uint, file_name string) {
	bitset.LittleEndian()
	b := bitset.New(length)
	b.SetTo(1, true)
	b.SetTo(2, false)
	b.SetTo(3, true)
	b.SetTo(4, false)
	b.SetTo(5, true)
	LFSR(b, length)
	_, err := b.WriteTo(create_file_stream(file_name))
	if err != nil {
		println(err)
	}
}
func LFSR(set *bitset.BitSet, length uint) {
	//指向当前位
	i := uint(1)
	for ; i+5 < length; i++ {
		bit1, err := set.NextSet(i - 1)
		if !err {
			fmt.Print("LFSR over!")
			os.Exit(0)
		} else if bit1 == i {
			a1 := 1
			a4 := 0
			//指向当前位的前一位
			r := i - 1
			for j := uint(1); j <= 3; j++ {
				bit4, err0 := set.NextSet(r + j)
				if !err0 {
					break
				} else if bit4 == i+3 {
					a4 = 1
				}
			}
			r++
			a6 := a1 ^ a4
			if a6 != 0 {
				set.SetTo(i+5, true)
			}
		}
	}
}

func create_file_stream(file_name string) io.Writer {
	stream, err := os.Create(file_name)
	if err == nil {
		return stream
	} else {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}
}
