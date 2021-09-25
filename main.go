package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
)

func rgb(i int) (int, int, int) {
	var f = 0.1
	return int(math.Sin(f*float64(i)+0)*127 + 128),
		int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func readFile(path string) (string, error) {
	// TODO: relative path
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func printRainbowChar(c rune, rainbowSeed int) {
	r, g, b := rgb(rainbowSeed)
	fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, c)
}

func displayContents(content string) {
	for i := 0; i < len(content); i++ {
		printRainbowChar(rune(content[i]), i)
	}
}

func runWithFiles() {
	nArgs := len(os.Args)
	if nArgs <= 1 {
		fmt.Printf("usage: %s [filename1] [filename2] [...]\n", os.Args[0])
		os.Exit(1)
	}

	for _, file := range os.Args[1:] {
		content, err := readFile(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		displayContents(content)
	}
}

func runWithPipe() {
	reader := bufio.NewReader(os.Stdin)
	i := 0
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		printRainbowChar(input, i)
		i += 1
	}
}

func main() {
	info, _ := os.Stdin.Stat()

	if info.Mode()&os.ModeCharDevice != 0 {
		runWithFiles()
	} else {
		runWithPipe()
	}
}
