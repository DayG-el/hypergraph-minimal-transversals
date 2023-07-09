package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type InputData struct {
	Data []string `json:"data"`
}

type OutputData struct {
	Results []string `json:"results"`
}

type mybit int32

var max_size int

type mybits []mybit

func popcnt(v mybit) int {
	count := 0
	for v != 0 {
		count++
		v &= v - 1
	}
	return count
}

func b2s(v mybit) string {
	var builder strings.Builder
	for i := max_size - 1; i >= 0; i-- {
		if v&(1<<i) != 0 {
			builder.WriteString("1")
		} else {
			builder.WriteString("0")
		}
	}
	return builder.String()
}

func a2b(str string) mybit {
	t := mybit(0)
	tokens := strings.Split(str, " ")
	for _, token := range tokens {
		if token != "" {
			v, _ := strconv.Atoi(token)
			if v > max_size {
				max_size = v
			}
			t |= 1 << (v - 1)
		}
	}
	return t
}

func b2a(t mybit) string {
	var builder strings.Builder
	for i := max_size - 1; i >= 0; i-- {
		if t&(1<<i) != 0 {
			builder.WriteString(strconv.Itoa(i + 1))
			builder.WriteString(" ")
		}
	}
	return strings.TrimSpace(builder.String())
}

func s2b(str string) mybit {
	t := mybit(0)
	for i := 0; i < len(str); i++ {
		if str[i] == '1' {
			t |= 1 << (len(str) - i - 1)
		}
	}
	return t
}

func loadDatFile(data []string) mybits {
	var v mybits
	for _, line := range data {
		if line != "" {
			v = append(v, a2b(strings.TrimSpace(line)))
		}
	}
	return v
}

func loadBitsFile(data []string) mybits {
	var v mybits
	for _, line := range data {
		if line != "" {
			v = append(v, s2b(line))
		}
	}
	return v
}

func checkMinimal(t mybit, k int, e mybits) bool {
	v := t
	for v != 0 {
		t2 := v & -v
		t3 := t ^ t2
		flag := true
		for i := 0; i < k; i++ {
			flag = flag && ((e[i] & t3) != 0)
			if !flag {
				break
			}
		}
		if flag {
			return false
		}
		v ^= t2
	}
	return true
}

func checkMinimal2(t mybit, k int, e mybits) bool {
	t2 := mybit(0)
	for i := 0; i < k; i++ {
		if popcnt(t&e[i]) == 1 {
			t2 |= t & e[i]
		}
	}
	return t2 == t
}

func search(k int, t mybit, e mybits, r *mybits) {
	if k == len(e) {
		*r = append(*r, t)
		return
	}
	if t&e[k] != 0 {
		search(k+1, t, e, r)
	} else {
		v := e[k]
		for v != 0 {
			t2 := v & -v
			if checkMinimal2(t|t2, k+1, e) {
				search(k+1, t|t2, e, r)
			}
			v ^= t2
		}
	}
}

func processData(data []string) []string {
	var e mybits
	bits := false // Assuming data is not in bit format
	if bits {
		e = loadBitsFile(data)
	} else {
		e = loadDatFile(data)
	}

	var r mybits
	search(0, 0, e, &r)

	results := make([]string, 0, len(r))
	for _, val := range r {
		results = append(results, b2a(val))
	}

	return results
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the input data
	var inputData InputData
	if err := json.Unmarshal(body, &inputData); err != nil {
		http.Error(w, "Failed to parse input data", http.StatusBadRequest)
		return
	}

	// Process the input data
	results := processData(inputData.Data)

	// Prepare the output data
	outputData := OutputData{Results: results}

	// Marshal the output data
	responseBody, err := json.Marshal(outputData)
	if err != nil {
		http.Error(w, "Failed to marshal output data", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")

	// Write the response
	if _, err := w.Write(responseBody); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/process", processRequest)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
