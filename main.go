package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// THREADS is the maximum number of goroutines
var THREADS = runtime.GOMAXPROCS(-1)

func main() {
	fmt.Printf("Current time: %v\n", time.Now().Unix())

	// args := os.Args[1:]
	fmt.Printf("Goroutines: %v\n", THREADS)

	header := os.Args[1]
	fmt.Printf("Commit Header: %v\n", header)

	target := os.Args[2]
	fmt.Printf("Target prefix: %s\n", target)

	message := resolve(header, target)
	fmt.Printf("Message for target prefix: %s\n", message)
}

func resolve(header, target string) string {
	fmt.Printf("Mining for target String: %s\n", target)
	start := time.Now()
	ch := make(chan string)
	halter := make(chan bool)

	for i := 0; i < THREADS; i++ {
		go unhash(header, target, i, ch, halter)
	}

	msg := <-ch
	close(halter)

	fmt.Printf("Message text %s found in %v\n", msg, time.Since(start))
	return msg
}

func unhash(header, target string, start int, ch chan string, quit chan bool) int {
	hexHash := make([]byte, 64)

	for i := start; true; i += THREADS {
		select {
		case <-quit:
			return -1
		default:
			hashVal := fmt.Sprintf("%x", i)
			headerSize := len(header) + len(hashVal) + 1
			seedStr := fmt.Sprintf("commit %v\000", headerSize) + header + hashVal
			hash := sha1.Sum([]byte(seedStr))
			hex.Encode(hexHash, hash[:])
			if strings.HasPrefix(string(hexHash), target) {
				fmt.Printf("String hashed: %v\n", seedStr)
				fmt.Printf("Total size: %v\n", headerSize)
				ch <- hashVal
				return 0
			}
		}
	}

	return -1
}

func pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}
