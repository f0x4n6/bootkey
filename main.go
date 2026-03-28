// Extract the BootKey from an offline system hive.
//
// Usage:
//
//	bootkey system
//
// The arguments are:
//
//	system
//	    The system registry hive (required).
package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/text/encoding/unicode"
	"www.velocidex.com/golang/regparser"
)

const offset = 4096 + 4

var keys = []string{"JD", "Skew1", "GBG", "Data"}
var sbox = []int{8, 5, 4, 2, 11, 9, 13, 3, 0, 6, 1, 12, 14, 10, 15, 7}

func extract(path string) (key []byte, err error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() { _ = f.Close() }()

	reg, err := regparser.NewRegistry(f)

	if err != nil {
		return nil, err
	}

	cs := "ControlSet001"

	// get current control set
	if k := reg.OpenKey("\\Select"); k != nil {
		for _, v := range k.Values() {
			if v.ValueName() == "Current" {
				cs = fmt.Sprintf("ControlSet%03d", v.ValueData().Uint64)
			}
		}
	}

	buf := bytes.NewBuffer(nil)

	// extract key parts hidden in key class names
	for _, v := range keys {
		k := reg.OpenKey(fmt.Sprintf("\\%s\\Control\\Lsa\\%s", cs, v))
		b := make([]byte, k.ClassLength())

		_, err = reg.BaseBlock.HiveBin().Reader.ReadAt(b, int64(k.Class()+offset))

		if err != nil {
			return nil, err
		}

		buf.Write(b)
	}

	tmp := buf.String()

	// decode unicode step
	if buf.Len() > 32 {
		dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		tmp, err = dec.String(buf.String())

		if err != nil {
			return nil, err
		}
	}

	// decode hex step
	sub, err := hex.DecodeString(tmp)

	if err != nil {
		return nil, err
	}

	// transform key bytes
	for i := 0; i < len(sub); i++ {
		key = append(key, sub[sbox[i]])
	}

	return key, nil
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" {
		_, _ = fmt.Fprintln(os.Stderr, "usage: bootkey system")
		os.Exit(2)
	}

	key, err := extract(os.Args[1])

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%x\n", key)
}
