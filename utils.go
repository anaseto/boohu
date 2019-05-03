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

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Indefinite(s string, upper bool) (text string) {
	if len(s) > 0 {
		switch s[0] {
		case 'a', 'e', 'i', 'o', 'u':
			if upper {
				text = "An " + s
			} else {
				text = "an " + s
			}
		default:
			if upper {
				text = "A " + s
			} else {
				text = "a " + s
			}
		}
	}
	return text
}

func formatText(text string, width int) string {
	pbuf := bytes.Buffer{}
	wordbuf := bytes.Buffer{}
	col := 0
	wantspace := false
	wlen := 0
	for _, c := range text {
		if c == ' ' || c == '\n' {
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
			if c == '\n' {
				pbuf.WriteRune('\n')
				pbuf.WriteRune('\n')
				col = 0
				wantspace = false
			} else {
				wantspace = true
			}
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
