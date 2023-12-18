package internal

import (
	"fmt"
	"github.com/lowk3v/takeaddr/internal/defiscan"
	"github.com/lowk3v/takeaddr/utils"
	"os"
	"sync"
)

func Explorer(addresses []string) {
	// optimize the above for loop faster with goroutine
	var wg sync.WaitGroup
	defiYield := defiscan.New()

	for _, address := range addresses {
		var err error
		msg := "No smart contract address found"
		if utils.IsBlacklist(address) {
			continue
		}
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			_, msg, err = defiYield.Scan(address)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stdout, "%s\t%s\n", address, err.Error())
				return
			}
			_, _ = fmt.Fprintf(os.Stdout, "%s\t%s\n", address, msg)
		}(address)
	}
	wg.Wait()
}
