package gsa

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Client struct {
	PhoneNumbers []string
	ScrapeURLs   []string
	Filter       func([]byte, string) (string, error)
	TextbeltURL  string
}

func New(phoneNumbers, scrapeURLs []string, filter func([]byte, string) (string, error)) *Client {
	c := new(Client)
	c.PhoneNumbers = phoneNumbers
	c.ScrapeURLs = scrapeURLs
	c.Filter = filter
	c.TextbeltURL = "http://textbelt.com/text"

	return c
}

func (c *Client) Run() {
	msgs := c.scrapeAll()

	msgBuffer := make([]string, 0, 0)

	for msg := range msgs {
		msgBuffer = append(msgBuffer, msg)
	}

	if len(msgBuffer) == 0 {
		fmt.Printf("%v | No changes\n", time.Now().Format(time.Stamp))
	} else {
		fmt.Println(strings.Join(msgBuffer, " | "))
		c.alertAll(strings.Join(msgBuffer, " | "))
	}
}

func (c *Client) scrapeAll() <-chan string {
	var wg sync.WaitGroup
	wg.Add(len(c.ScrapeURLs))
	out := make(chan string)

	for _, u := range c.ScrapeURLs {
		go func(scrapeURL string) {
			msg, err := scrape(scrapeURL, c.Filter)
			if err != nil {
				fmt.Println(err)
			} else {
				out <- msg
			}
			wg.Done()
		}(u)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (c *Client) alertAll(msg string) {
	var wg sync.WaitGroup
	wg.Add(len(c.PhoneNumbers))

	for _, p := range c.PhoneNumbers {
		go func(phoneNumber string) {
			err := alert(c.TextbeltURL, phoneNumber, msg)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(p)
	}

	wg.Wait()
}
