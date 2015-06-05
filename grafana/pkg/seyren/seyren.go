package seyren

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Represents the Seyren API
type Seyren struct {
	URL string
}

func New(url string) *Seyren {
	seyren := Seyren{}
	seyren.URL = url
	return &seyren
}

func extractIdFromURL(location string) string {
	u, err := url.Parse(location)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	snips := strings.Split(u.Path, "/")
	return snips[len(snips)-1]
}

// FIXME: delete the existing alert.
func (s *Seyren) PostCheck(c *Check) (string, error) {
	checkURL := fmt.Sprintf("%s/api/checks", s.URL)

	checkJSON, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", checkURL, bytes.NewBuffer(checkJSON))
	if err != nil {
		return "", err
	}
	log.Printf("JSON check: %s\n", checkJSON)
	req.Header.Set("Content-Type", "application/json")

	// Set basic auth
	username := os.Getenv("SEYREN_USERNAME")
	password := os.Getenv("SEYREN_PASSWORD")
	if username == "" && password == "" {
		log.Printf("NOTE: No authentication information found for seyren.\n")
	} else {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")

	if location == "" {
		log.Printf("Couldn't find location\n")
		return "", errors.New("Couldn't find a Location header. Can't figure out the ID. FUCK YOU")
	}
	id := extractIdFromURL(location)
	log.Printf("Id is %s\n", id)
	return id, nil
}

func (s *Seyren) DeleteCheck(c *Check) (string, error) {
	checkURL := fmt.Sprintf("%s/api/checks/%s", s.URL, c.Id)
	req, err := http.NewRequest("DELETE", checkURL, nil)
	if err != nil {
		return "", err
	}
	// Set basic auth
	username := os.Getenv("SEYREN_USERNAME")
	password := os.Getenv("SEYREN_PASSWORD")
	if username == "" && password == "" {
		log.Printf("NOTE: No authentication information found for seyren.")
	} else {
		req.SetBasicAuth(username, password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.Status, nil
}
