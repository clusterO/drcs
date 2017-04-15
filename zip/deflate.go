package zip

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

func compress(input []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decompress(input []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(input))
	defer r.Close()
	return ioutil.ReadAll(r)
}