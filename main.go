package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
)

// settings defines a structure for config.json file
type settings struct {
	Items []int `json:"items"`
}

// default item values provided in the task
var items []int = []int{
	250,
	500,
	1000,
	2000,
	5000,
}

// calculatePacks uses Greedy algorithm to
// fill the result with items
func calculatePacks(target int) map[int]int {
	result := make(map[int]int)
	l := len(items)

	// if target exceeds the max value in the items range
	// we need to put it into the range
	// for example, if target is 12001, we put 2x5000 into the result
	// and continue with the remaining of 2001
	if target > items[l-1] {
		result[items[l-1]] = target / items[l-1]
		target %= items[l-1]
	}

	// if target is withing items items range, we start
	// our greedy algorithm from the end of the items
	// to the beginning
	for i := len(items) - 2; i >= 0 && target > 0; i-- {
		// if target is within the range of items[i] and items[i+1]
		// we need to check which one is closer to the target
		// and put it into the result
		if items[i] <= target && target <= items[i+1] {
			// if target + minimal item is greater than the next item (i.e. items[i+1]),
			// it is better to put the next item into the result.
			// for example, if target is 1800 and the next item is 2000,
			// it is better to put 2000 into the result than 1000 + 1000 or 1000 + 500 + 500
			if items[i+1]-target-items[0] < 0 {
				result[items[i+1]] = items[i+1] / target
				target = target - items[i+1]
			} else {
				// otherwise, we put the current item (i.e. items[i] into the result
				// for example, if target is 1700 and the current item is 1000
				result[items[i]] = target / items[i]
				target = target % items[i]
			}
		}
	}

	// if we still have some target left, we put the minimal item into the result
	// for example, if target is 1 or 249, we put 250 into the result
	if target > 0 {
		result[items[0]]++
	}

	return result
}

// writeResponse writes a response with the provided status code and message
func writeResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}

// handleCalculatePacks calculate packs for the provided items
// as a 'items' query parameter
func handleCalculatePacks(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("items") {
		writeResponse(w, http.StatusBadRequest, "items must be provided")
		return
	}

	itemNum, err := strconv.Atoi(r.URL.Query().Get("items"))
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "items must be an integer")
		return
	}

	packs := calculatePacks(itemNum)
	data, err := json.Marshal(packs)
	if err != nil {
		log.Printf("failed to convert packs to JSON: %s", err)
		writeResponse(w, http.StatusInternalServerError, "failed to convert packs to JSON")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// runCalculateTests defines unit test for calculatePacks function
// with default items values (i.e. 250, 500, 1000, 2000, 5000)
func runCalculateTests() error {
	testCases := []struct {
		input  int
		output map[int]int
	}{
		{
			input:  1,
			output: map[int]int{250: 1},
		},
		{
			input:  250,
			output: map[int]int{250: 1},
		},
		{
			input:  251,
			output: map[int]int{500: 1},
		},
		{
			input:  501,
			output: map[int]int{500: 1, 250: 1},
		},
		{
			input:  750,
			output: map[int]int{500: 1, 250: 1},
		},
		{
			input:  3000,
			output: map[int]int{2000: 1, 1000: 1},
		},
		{
			input:  3001,
			output: map[int]int{2000: 1, 1000: 1, 250: 1},
		},
		{
			input:  12001,
			output: map[int]int{250: 1, 2000: 1, 5000: 2},
		},
		{
			input:  20000,
			output: map[int]int{5000: 4},
		},
	}

	for _, tc := range testCases {
		packs := calculatePacks(tc.input)
		// if obtained result doesn't match expected
		// we return an error to stop the program
		if !reflect.DeepEqual(tc.output, packs) {
			return fmt.Errorf("expected output %v doesn't match obtained %v for input %d", tc.output, packs, tc.input)
		}
	}

	return nil
}

func setupConfig() {
	// default path expected is ./config.json
	configFilePath := "config.json"

	// obtaining a custom config path from command line
	if len(os.Args) > 1 && os.Args[1] != "" {
		configFilePath = os.Args[1]
	}

	// if config file doesn't exist, we don't need to do anything
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Printf("couldn't get configFilePath from %s. Error: %s", configFilePath, err)
		return
	}

	var stg settings
	if err := json.Unmarshal(data, &stg); err != nil {
		log.Printf("failed to unmarshal config file. Error: %s", err)
		return
	}

	items = stg.Items
}

func main() {
	/*
	   sample curl query is:
	   curl http://localhost:8080/calculate-packs?items=250
	*/

	if err := runCalculateTests(); err != nil {
		log.Printf("unit test failed with error '%s'", err)
		return
	}

	// parse config.json file if exists
	// to setup custom item values
	setupConfig()

	// sorting is needed for greedy algorithm
	sort.Ints(items)

	// setting up HTTP server
	http.HandleFunc("/calculate-packs", handleCalculatePacks)

	log.Println("starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("failed to listen at :8080. Error:", err)
	}
}
