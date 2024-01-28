package agent

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/oktavarium/go-gauger/internal/agent/internal/flags"
	"golang.org/x/sync/errgroup"
)

// Run - запускает агент сбора и отправки метрик на сервер
func Run() error {
	flagsConfig, err := flags.LoadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	publicKeyData, err := os.ReadFile(flagsConfig.CryptoKey)
	if err != nil {
		return fmt.Errorf("error on reading public key file: %w", err)
	}

	pkPEM, _ := pem.Decode(publicKeyData)
	if pkPEM.Type != "RSA PUBLIC KEY" {
		return fmt.Errorf("wrong key type: %w", err)
	}

	publicKey, err := x509.ParsePKCS1PublicKey(pkPEM.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}

	eg, egCtx := errgroup.WithContext(context.Background())

	chMetrics := collector(
		egCtx,
		readMetrics,
		eg, time.Duration(flagsConfig.PollInterval))
	chPsMetrics := collector(
		egCtx,
		readPsMetrics,
		eg, time.Duration(flagsConfig.PollInterval))

	unitedCh := fanIn(chMetrics, chPsMetrics)
	for i := 0; i < flagsConfig.RateLimit; i++ {
		go sender(egCtx, flagsConfig.Address, flagsConfig.HashKey, publicKey,
			flagsConfig.ReportInterval, unitedCh)
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

// fanIn - метод мультиплексирования входящих данных от множества
func fanIn(chs ...<-chan []byte) <-chan []byte {
	chOut := make(chan []byte, len(chs))
	var wg sync.WaitGroup
	wg.Add(len(chs))

	output := func(ch <-chan []byte) {
		defer wg.Done()
		for v := range ch {
			chOut <- v
		}
	}

	for _, ch := range chs {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(chOut)
	}()

	return chOut
}
