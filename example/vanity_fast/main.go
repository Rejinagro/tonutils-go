package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {
	threads := flag.Uint64("threads", 8, "parallel threads")
	suffix := flag.String("suffix", "", "desired contract suffix, required")
	caseSensitive := flag.Bool("case", false, "is case sensitive")
	flag.Parse()

	if *suffix == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var counter uint64

	for x := uint64(0); x < *threads; x++ {
		go generateWallets(*suffix, *caseSensitive, &counter)
	}

	log.Println("searching...")
	for {
		time.Sleep(1 * time.Second)
		log.Println("checked", atomic.LoadUint64(&counter), "per second")
		atomic.StoreUint64(&counter, 0)
	}
}

func generateWallets(suffix string, caseSensitive bool, counter *uint64) {
	for {
		_, pk, _ := ed25519.GenerateKey(nil)
		w, _ := wallet.FromPrivateKey(nil, pk, wallet.V4R2)

		s := w.WalletAddress().String()
		if !caseSensitive {
			s = strings.ToLower(s)
			suffix = strings.ToLower(suffix)
		}

		if strings.HasSuffix(s, suffix) {
			atomic.AddUint64(counter, 1)
			log.Println(
				"========== FOUND ==========\n",
				"Address:", w.WalletAddress().String(), "\n", "Private key:", hex.EncodeToString(pk.Seed()),
				"\n========== FOUND ==========",
			)
		}
	}
}
