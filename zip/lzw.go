package zip

import (
	"fmt"
	"strconv"
)

type LZW struct {
	dictionary map[string]int
	nextCode   int
}

func NewLZW() *LZW {
	return &LZW{
		dictionary: make(map[string]int),
		nextCode:   256,
	}
}

func (lzw *LZW) Compress(input string) []int {
	output := []int{}
	prefix := ""
	for _, c := range input {
		current := prefix + string(c)
		if _, ok := lzw.dictionary[current]; ok {
			prefix = current
		} else {
			output = append(output, lzw.dictionary[prefix])
			lzw.dictionary[current] = lzw.nextCode
			lzw.nextCode++
			prefix = string(c)
		}
	}
	if prefix != "" {
		output = append(output, lzw.dictionary[prefix])
	}
	return output
}

func (lzw *LZW) Decompress(input []int) string {
	output := ""
	previous := ""
	for _, code := range input {
		current := ""
		if val, ok := lzw.dictionary[fmt.Sprint(code)]; ok {
			current = fmt.Sprint(val)
		} else if code == lzw.nextCode {
			current = previous + string(previous[0])
		} else {
			panic("Invalid compressed data")
		}
		output += current
		if previous != "" {
			lzw.dictionary[fmt.Sprint(lzw.nextCode)], _ = strconv.Atoi(previous + string(current[0]))
			lzw.nextCode++
		}
		previous = current
	}
	return output
}
