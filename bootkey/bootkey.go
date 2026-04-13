// Package bootkey provides methods to extract the bootkey from an offline system hive.
package bootkey

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/unicode"
	"www.velocidex.com/golang/regparser"
)

// ControlSet fallback
var ControlSet = "ControlSet001"

// SBox for key transformation
var SBox = []int{8, 5, 4, 2, 11, 9, 13, 3, 0, 6, 1, 12, 14, 10, 15, 7}

// ReadFile and return the extracted bootkey and any potential errors.
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	return ReadData(f)
}

// ReadData and return the extracted bootkey and any potential errors.
func ReadData(r io.ReaderAt) ([]byte, error) {
	reg, err := regparser.NewRegistry(r)

	if err != nil {
		return nil, err
	}

	// get the current control set
	if k := reg.OpenKey("\\Select"); k != nil {
		for _, v := range k.Values() {
			if v.ValueName() == "Current" {
				ControlSet = fmt.Sprintf("ControlSet%03d", v.ValueData().Uint64)
			}
		}
	}

	buf := bytes.NewBuffer(nil)

	// extract key parts hidden in key class names
	for _, v := range []string{
		"JD", "Skew1", "GBG", "Data",
	} {
		k := reg.OpenKey(fmt.Sprintf("\\%s\\Control\\Lsa\\%s", ControlSet, v))
		b := make([]byte, k.ClassLength())

		_, err = reg.BaseBlock.HiveBin().Reader.ReadAt(b, int64(k.Class()+4096+4))

		if err != nil {
			return nil, err
		}

		buf.Write(b)
	}

	tmp := buf.String()

	// decode unicode step
	if buf.Len() > 32 {
		tmp, err = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder().String(buf.String())

		if err != nil {
			return nil, err
		}
	}

	// decode hex step
	sub, err := hex.DecodeString(tmp)

	if err != nil {
		return nil, err
	}

	key := make([]byte, len(sub))

	// transform key bytes
	for i := 0; i < len(sub); i++ {
		key[i] = sub[SBox[i]]
	}

	return key, nil
}
