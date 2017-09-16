package main

import (
	"bytes"
	"crypto/rand"
	"log"
	"math/big"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func RandInt(n int) int {
	if n <= 0 {
		return 0
	}
	x, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		log.Println(err, ":", n)
		return n / 2
	}
	return int(x.Int64())
}

func formatText(text string, width int) string {
	pbuf := bytes.Buffer{}
	wordbuf := bytes.Buffer{}
	col := 0
	wantspace := false
	wlen := 0
	for _, c := range text {
		// it's enough to test with just spaces for now
		//if unicode.IsSpace(c) && c != 0xa0 {
		if c == ' ' {
			if wlen == 0 {
				continue
			}
			if col+wlen > width {
				if wantspace {
					pbuf.WriteRune('\n')
					col = 0
				}
			} else if wantspace {
				pbuf.WriteRune(' ')
				col++
			}
			pbuf.Write(wordbuf.Bytes())
			col += wlen
			wordbuf.Reset()
			wlen = 0
			wantspace = true
			continue
		}
		wordbuf.WriteRune(c)
		wlen++
	}
	if wordbuf.Len() > 0 {
		if wantspace {
			if wlen+col > width {
				pbuf.WriteRune('\n')
			} else {
				pbuf.WriteRune(' ')
			}
		}
		pbuf.Write(wordbuf.Bytes())
	}
	return pbuf.String()
}
