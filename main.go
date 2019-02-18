package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jessevdk/go-flags"
)

type option struct {
	IsReverse []bool `short:"r" long:"reverse" description:"Unicode code points to string. E.g. \\u3042\\u3044\\u3046\\u3048\\u304A -> あいうえお"`
}

var re = regexp.MustCompile(`\\[uU][0-9a-fA-F]{4}`)

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
		fmt.Printf("'%s' includes non-UTF8 value(s).\n", s)
		return
	}

	for _, r := range s {
		if r <= unicode.MaxASCII && !unicode.IsControl(r) {
			fmt.Print(string(r))
		} else {
			codepoint := fmt.Sprintf("%U", r)[2:]
			fmt.Printf("\\u%s", codepoint)
		}
	}
	fmt.Println()
}

func asciiToUtf8(s string) {
	// TODO: handle surrogate pairs.

	match := re.FindString(s)
	if match == "" {
		fmt.Println(s)
		return
	}

	codepoint := match[2:]
	n, err := strconv.ParseInt(codepoint, 16, 32)
	if err != nil {
		fmt.Printf("'%s' can't be converted to number.\n", codepoint)
		return
	}

	r := int32(n)
	s = strings.Replace(s, match, string(r), -1)

	asciiToUtf8(s)
}
