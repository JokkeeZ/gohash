package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

const DefaultAlg = "sha256"

type Algorithm struct {
	name string
	hash hash.Hash
}

type Color int

const (
	ResetColor   Color = 0
	Cyan         Color = 36
	DefaultColor Color = 39
	DarkGray     Color = 90
	LightRed     Color = 91
	LightGreen   Color = 92
)

var algorithms = []Algorithm{
	{name: "md5", hash: md5.New()},
	{name: "sha1", hash: sha1.New()},
	{name: "sha224", hash: sha256.New224()},
	{name: "sha384", hash: sha512.New384()},
	{name: "sha256", hash: sha256.New()},
	{name: "sha512", hash: sha512.New()},
}

func getAlgorithm(algorithm string) *hash.Hash {
	for _, a := range algorithms {
		if strings.ToLower(algorithm) == a.name {
			return &a.hash
		}
	}

	return nil
}

func getUserInput() *string {
	var input string
	_, err := fmt.Scanln(&input)

	if err != nil {
		print("Invalid input", 31)
		os.Exit(1)
	}

	return &input
}

func main() {
	argc := len(os.Args)

	if argc == 1 {
		print("/**\n", DarkGray)
		print(" * You can specify algorithm with cli argument\n", DarkGray)
		print(" * Otherwise by default SHA256 will be used\n", DarkGray)

		for _, a := range algorithms {
			print(" * ./gohash -a %s\n", DarkGray, a.name)
		}

		print(" */\n", DarkGray)
	}

	alg := DefaultAlg

	if argc == 2 {
		alg = os.Args[1]
	}

	algorithm := getAlgorithm(alg)

	if algorithm == nil {
		print("Invalid algorithm\n", LightRed)
		return
	}

	fmt.Print("Drag and drop a file to make a checksum for it: ")
	filePath := getUserInput()

	printReport(*filePath, *algorithm)
}

func printReport(filePath string, alg hash.Hash) {
	checksum := fileChecksum(filePath, alg)
	printChecksums(checksum)

	print("Compare to checksum: ", DefaultColor)
	inputChecksum := getUserInput()

	printChecksums(checksum)
	printChecksums(*inputChecksum)

	if checksum == *inputChecksum {
		print("✓ Checksums match\n", LightGreen)
	} else {
		print("✕ Checksums DO NOT match\n", LightRed)
	}
}

func fileChecksum(filePath string, algorithm hash.Hash) string {
	if fi, err := os.Stat(filePath); err != nil || !fi.Mode().IsRegular() {
		print("Input contains invalid filepath. (or file do not exist)\n", LightRed)
		os.Exit(1)
	}

	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = io.Copy(algorithm, file)

	if err != nil {
		panic(err)
	}

	hash := algorithm.Sum(nil)
	return hex.EncodeToString(hash[:])
}

func print(text string, color Color, f ...any) {
	format := fmt.Sprintf(text, f...)

	if len(f) == 0 {
		format = text
	}

	fmt.Printf("\033[0;%dm%s\033[0m", color, format)
}

func printChecksums(checksum string) {
	print("Checksum\t: ", DefaultColor)

	for _, e := range checksum {
		if e >= '0' && e <= '9' {
			print("%s", Cyan, string(e))
		} else {
			print("%s", DefaultColor, string(e))
		}
	}

	print("\n", ResetColor)
}
