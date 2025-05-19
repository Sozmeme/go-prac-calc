package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	pb "prac/proto"

	"google.golang.org/grpc"
)

func TestHTTP(t *testing.T) {
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

func TestGRPC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"localhost:9090",
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)

	req := &pb.CalculationRequest{
		Instructions: []*pb.Instruction{
			{
				Type:  "calc",
				Op:    "+",
				Var:   "x",
				Left:  &pb.Instruction_LeftInt{LeftInt: 2},
				Right: &pb.Instruction_RightInt{RightInt: 3},
			},
			{
				Type: "print",
				Var:  "x",
			},
		},
	}

	resp, err := client.Calculate(ctx, req)
	if err != nil {
		t.Fatalf("Calculate RPC failed: %v", err)
	}

	expected := map[string]int64{
		"x": 5,
	}

	if len(resp.Items) != 1 {
		t.Errorf("Expected 1 result, got %d", len(resp.Items))
	}

	for _, item := range resp.Items {
		if val, ok := expected[item.Var]; !ok {
			t.Errorf("Unexpected variable in response: %s", item.Var)
		} else if item.Value != val {
			t.Errorf("For variable %s, expected %d, got %d", item.Var, val, item.Value)
		}
	}
}
