package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

func rgb(i int) (int, int, int) {
	var f = 0.1
	return int(math.Sin(f*float64(i)+0)*127 + 128),
		int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func printRainbowChar(c rune, rainbowSeed int) {
	r, g, b := rgb(rainbowSeed)
	fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, c)
}

var rainbowSeed int = 0

func displayContents(content string) {
	for i := 0; i < len(content); i++ {
		printRainbowChar(rune(content[i]), rainbowSeed)
		rainbowSeed += 1
	}
}

func displayChunksFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	const chunkSz = 16
	buf := make([]byte, chunkSz)

	for {
		readTotal, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err)
			}
			break
		}

		chunk := string(buf[:readTotal])

		displayContents(chunk)
	}
}

func expandPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func runOneFile(path string) {
	expandedPath, err := expandPath(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	displayChunksFromFile(expandedPath)
}

func runWithFiles() {
	nArgs := len(os.Args)
	if nArgs <= 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [filename1] [filename2] [...]\n", os.Args[0])
		os.Exit(1)
	}

	for _, path := range os.Args[1:] {
		runOneFile(path)
	}
}

func runWithPipe() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		displayContents(string(input))
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
