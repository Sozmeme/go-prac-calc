package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHTTPRawJSONRequest(t *testing.T) {

	rawJSON := `[
		{ "type": "calc", "op": "+", "var": "x", "left": 10, "right": 2 },
		{ "type": "calc", "op": "*", "var": "y", "left": "x", "right": 5 },
		{ "type": "calc", "op": "-", "var": "q", "left": "y", "right": 20 },
		{ "type": "calc", "op": "+", "var": "unusedA", "left": "y", "right": 100 },
		{ "type": "calc", "op": "*", "var": "unusedB", "left": "unusedA", "right": 2 },
		{ "type": "print", "var": "q" },
		{ "type": "calc", "op": "-", "var": "z", "left": "x", "right": 15 },
		{ "type": "print", "var": "z" },
		{ "type": "calc", "op": "+", "var": "ignoreC", "left": "z", "right": "y" },
		{ "type": "print", "var": "x" }
	]`

	resp, err := http.Post("http://localhost:8080/calculate", "application/json", strings.NewReader(rawJSON))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response body: %s", body)
		return
	}

	var response struct {
		Items []struct {
			Var   string `json:"var"`
			Value int64  `json:"value"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expected := map[string]int64{
		"x": 12,
		"y": 60,
		"q": 40,
		"z": -3,
	}

	for _, item := range response.Items {
		val, ok := expected[item.Var]
		if !ok {
			t.Errorf("Unexpected variable in response: %s", item.Var)
			continue
		}
		if item.Value != val {
			t.Errorf("For variable %s, expected %d, got %d", item.Var, val, item.Value)
		}
	}

	if len(response.Items) != 3 {
		t.Errorf("Expected 3 results, got %d", len(response.Items))
	}
}
