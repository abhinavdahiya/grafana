package seyren

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type Check struct {
	Id          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Target      string `json:"target"`
	Warn        string `json:"warn"`
	Error       string `json:"error"`
	Enabled     bool   `json:"enabled"`
	Live        bool   `json:"live,omitempty"`
	From        string `json:"from,omitempty"`
	Until       string `json:"until,omitempty"`
}

// For now just one target per graph.
func extractTarget(source []interface{}) string {
	for _, _tmap := range source {
		tmap := _tmap.(map[string]interface{})
		target, ok := tmap["target"]
		if ok {
			return target.(string)
		}
	}
	return ""
}

// A Row represents the row in grafana dashboard
// returns a list of Seyren Checks
func CheckFromPanel(panel map[string]interface{}) (*Check, error) {
	//	panels := _panels.([]interface{})
	if panel["type"].(string) != "graph" {
		log.Printf("Ignoring this panel as it's not a graph - %v\n", panel["title"])
		return nil, errors.New("Not a graph")
	}
	// Gather all the thresholds
	grid := panel["grid"].(map[string]interface{})
	var thresholds = []float64{}
	for k, v := range grid {
		if matched, _ := regexp.MatchString("threshold\\d+$", k); matched {
			if v != nil {
				thresholds = append(thresholds, v.(float64))
			}
		}
	}
	if len(thresholds) < 2 {
		log.Printf("No thresholds found for this graph. Ignoring")
		return nil, errors.New("No thresholds")
	}

	id := panel["alertID"].(string)

	// Create the check here
	check := Check{}
	check.Id = id
	check.Name = panel["title"].(string)
	check.Description = fmt.Sprintf("Check added from graphana for '%s'", check.Name)
	check.Target = extractTarget(panel["targets"].([]interface{}))
	// threshold1 is Warn
	check.Warn = strconv.FormatFloat(thresholds[0], 'f', 6, 64)
	check.Error = strconv.FormatFloat(thresholds[1], 'f', 6, 64)
	check.Enabled = true
	check.Live = false
	// All done. We now have a check
	return &check, nil
}
