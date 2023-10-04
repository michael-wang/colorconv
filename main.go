package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	flag.Usage = func() {
		prog := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "%s #FFBB80\n", prog)
		fmt.Fprintf(os.Stderr, "or\n")
		fmt.Fprintf(os.Stderr, "%s 255,187,128\n", prog)
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	arg := flag.Arg(0)
	var color [3]uint8
	var err error
	if strings.HasPrefix(arg, "#") {
		color, err = parseHex(arg)
	} else {
		color, err = parseRGB(arg)
	}

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	out(color)
}

// Expected s: "FFBB80"
func parseHex(s string) (c [3]uint8, err error) {
	if len(s) < 6 || 7 < len(s) {
		return c, fmt.Errorf("invalid hex value, expected: #FFBB80, but got: %s", s)
	}

	// Strip leading #
	if s[0] == '#' {
		s = s[1:]
	}

	s = strings.ToLower(s)

	for i := 0; i < 3; i++ {
		v, err := strconv.ParseInt(s[i*2:2+i*2], 16, 16)
		if err != nil {
			return c, err
		}
		c[i] = uint8(v)
	}
	return c, nil
}

func parseRGB(s string) (c [3]uint8, err error) {
	s = stripSpaces(s)
	ss := strings.Split(s, ",")
	if len(ss) != 3 {
		return c, fmt.Errorf("invalid rgb value, expected: 12,34,56 but got: %s", s)
	}

	for i := 0; i < 3; i++ {
		v, err := strconv.Atoi(ss[i])
		if err != nil {
			return c, err
		}

		if v < 0 || 0xff < v {
			return c, fmt.Errorf("invalid color, expect [0, 255] but got %s", ss[i])
		}
		c[i] = uint8(v)
	}

	return c, nil
}

func out(c [3]uint8) {
	fmt.Printf("RGB: %d,%d,%d\n", c[0], c[1], c[2])
	fmt.Printf("Hex: #%2X%2X%2X\n", c[0], c[1], c[2])
}

func stripSpaces(s string) (out string) {
	var b strings.Builder
	b.Grow(len(s))
	for _, c := range s {
		if !unicode.IsSpace(c) {
			b.WriteRune(c)
		}
	}
	return b.String()
}
