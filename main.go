package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/jessevdk/go-flags"
)

type option struct {
	IsReverse []bool `short:"r" long:"reverse" description:"Unicode code points to string. E.g. \\u3042\\u3044\\u3046\\u3048\\u304A -> あいうえお"`
}

func main() {
	var opt option
	args, err := flags.ParseArgs(&opt, os.Args)
	if err != nil {
		panic(err)
	}

	isReverse := 0 < len(opt.IsReverse)
	for _, arg := range args[1:] {
		if isReverse {
			asciiToUtf8(arg)
		} else {
			utf8ToAscii(arg)
		}
	}
}

func utf8ToAscii(s string) {
	if !utf8.ValidString(s) {
		fmt.Println("'%s' includes non-UTF8 value(s).", s)
		return
	}

	for _, r := range s {
		// TODO: ignore ascii characters.
		codepoint := fmt.Sprintf("%U", r)[2:]
		fmt.Printf("\\u%s", codepoint)
	}
	fmt.Println()
}

func asciiToUtf8(s string) {
	// TODO: allow \\U
	// TODO: allow non-codepoint value.
	codepoints := strings.Split(s, "\\u")[1:]
	rs := make([]rune, len(codepoints))
	for i, codepoint := range codepoints {
		n, err := strconv.ParseInt(codepoint, 16, 32)
		if err != nil {
			fmt.Println("'%s' includes non-encoded value(s).", s)
			return
		}
		rs[i] = int32(n)
	}

	for _, r := range rs {
		fmt.Printf("%c", r)
	}
	fmt.Println()
}
