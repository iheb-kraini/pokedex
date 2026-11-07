package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type map_res struct {
	Next    string `json:"next"`
	Prev    string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

// fetches a page, checking the cache first
func getPage(c *config, url string) (map_res, error) {
	var page map_res

	data, ok := c.cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return page, err
		}
		defer res.Body.Close()

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return page, err
		}
		if res.StatusCode > 299 {
			return page, fmt.Errorf("response failed with status %d: %s", res.StatusCode, data)
		}
		c.cache.Add(url, data)
	}

	err := json.Unmarshal(data, &page)
	return page, err
}

func showPage(c *config, url string) error {
	page, err := getPage(c, url)
	if err != nil {
		return err
	}
	c.next, c.prev = page.Next, page.Prev
	for _, item := range page.Results {
		fmt.Println(item.Name)
	}
	return nil
}

func command_map(c *config) error {
	if c.next == "" {
		fmt.Println("last page")
		return nil
	}
	return showPage(c, c.next)
}

func command_bmap(c *config) error {
	if c.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return showPage(c, c.prev)
}
