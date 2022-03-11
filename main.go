package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

const (
	defaultGoroutinesNumber = 10
	chanBufSize             = 10
	parallelFlagName        = "parallel"
	httpScheme              = "http://"
)

func main() {
	var goroutinesNumber int64
	flag.Int64Var(&goroutinesNumber, parallelFlagName, defaultGoroutinesNumber, "number of goroutines")
	flag.Parse()

	links := flag.Args()
	if len(links) == 0 {
		return
	}

	logger := log.New(os.Stderr, "ERROR: ", log.Ltime)

	sh := New(&Hasher{}, &Getter{}, logger)
	result := sh.GetResponsesHashes(links, int(goroutinesNumber))

	for link, hash := range result {
		fmt.Println(link, hash)
	}
}

type hasher interface {
	Sum([]byte) string
}

type getter interface {
	Get(link string) ([]byte, error)
}

type SiteHasher struct {
	hasher hasher
	getter getter
	logger logger
}

func New(h hasher, g getter, l logger) *SiteHasher {
	return &SiteHasher{
		hasher: h,
		getter: g,
		logger: l,
	}
}

func (s *SiteHasher) GetResponsesHashes(rawLinks []string, goroutinesNumber int) map[string]string {
	wg := sync.WaitGroup{}

	linksChan := make(chan string, chanBufSize)
	result := make(map[string]string, len(rawLinks))
	mu := sync.Mutex{}

	go func() {
		defer close(linksChan)
		for _, rawLink := range rawLinks {
			link, err := validateLink(rawLink)
			if err != nil {
				if s.logger != nil {
					s.logger.Println(err)
				}
				continue
			}
			linksChan <- link
		}
	}()

	for i := 0; i < goroutinesNumber; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for link := range linksChan {
				body, err := s.getter.Get(link)
				if err != nil {
					if s.logger != nil {
						s.logger.Println(err)
					}
					continue
				}

				hash := s.hasher.Sum(body)
				// Could use sync.map instead
				mu.Lock()
				result[link] = hash
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	return result
}

func validateLink(rawLink string) (string, error) {
	urlStruct, err := url.Parse(rawLink)
	if err != nil {
		return "", err
	}

	link := urlStruct.String()
	if urlStruct.Scheme == "" {
		link = httpScheme + link
	}

	_, err = url.ParseRequestURI(link)
	if err != nil {
		return "", err
	}

	return link, nil
}

type Hasher struct{}

func (h *Hasher) Sum(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

type Getter struct{}

func (g *Getter) Get(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type logger interface {
	Println(v ...interface{})
}
