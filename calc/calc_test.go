package calc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestProcessCalc(t *testing.T) {
	calc := NewCalculator()
	err := calc.processCalc(Instruction{
		Type: "calc", Op: "+", Var: "a", Left: int64(2), Right: int64(3),
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	val, ok := calc.vars.Load("a")
	if !ok || val.(int64) != 5 {
		t.Errorf("expected 5, got %v", val)
	}

	err = calc.processCalc(Instruction{
		Type: "calc", Op: "*", Var: "a", Left: int64(1), Right: int64(1),
	})
	if err == nil {
		t.Errorf("expected error on duplicate variable")
	}
}

func TestGetValue(t *testing.T) {
	calc := NewCalculator()
	calc.vars.Store("x", int64(42))

	val, err := calc.getValue("x")
	if err != nil || val != 42 {
		t.Errorf("expected 42, got %d, err: %v", val, err)
	}

	val, err = calc.getValue(int64(10))
	if err != nil || val != 10 {
		t.Errorf("expected 10, got %d, err: %v", val, err)
	}

	val, err = calc.getValue(float64(3.0))
	if err != nil || val != 3 {
		t.Errorf("expected 3, got %d, err: %v", val, err)
	}

	_, err = calc.getValue("undefined")
	if err == nil {
		t.Error("expected error for undefined variable")
	}
}

func TestSimpleScenario(t *testing.T) {
	instructions := []Instruction{
		{Type: "calc", Op: "+", Var: "x", Left: int64(10), Right: int64(2)},
		{Type: "calc", Op: "*", Var: "y", Left: "x", Right: int64(5)},
		{Type: "calc", Op: "-", Var: "q", Left: "y", Right: int64(20)},
		{Type: "print", Var: "q"},
		{Type: "print", Var: "x"},
	}

	calc := NewCalculator()
	results, err := calc.Calculate(instructions)

	if err != nil {
		t.Errorf("Calculate failed: %v", err)
	}

	expectedResults := map[string]int64{
		"x": 12,
		"q": 40,
	}

	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(results))
	}

	for _, res := range results {
		expectedVal, ok := expectedResults[res.Var]
		if !ok {
			t.Errorf("Unexpected result variable: %s", res.Var)
			continue
		}
		if res.Value != expectedVal {
			t.Errorf("For %s expected %d, got %d", res.Var, expectedVal, res.Value)
		}
	}
}

func TestComplexScenario(t *testing.T) {
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

	var instructions []Instruction
	if err := json.Unmarshal([]byte(rawJSON), &instructions); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	calc := NewCalculator()
	results, err := calc.Calculate(instructions)

	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	expectedResults := map[string]int64{
		"x": 12,
		"q": 40,
		"z": -3,
	}

	expectedPrintCount := 3
	if len(results) != expectedPrintCount {
		t.Errorf("Expected %d printed results, got %d", expectedPrintCount, len(results))
	}

	for _, res := range results {
		expectedVal, ok := expectedResults[res.Var]
		if !ok {
			t.Errorf("Unexpected printed variable: %s", res.Var)
			continue
		}
		if res.Value != expectedVal {
			t.Errorf("For %s expected %d, got %d", res.Var, expectedVal, res.Value)
		}
	}

	computedVars := []struct {
		name  string
		value int64
	}{
		{"y", 60},
		{"unusedA", 160},
		{"unusedB", 320},
		{"ignoreC", 57},
	}

	for _, v := range computedVars {
		if val, ok := calc.vars.Load(v.name); !ok {
			t.Errorf("Variable %s was not computed", v.name)
		} else if val.(int64) != v.value {
			t.Errorf("For %s expected %d, got %d", v.name, v.value, val)
		}
	}
}

func Test100IndependentOperations(t *testing.T) {
	calc := NewCalculator()
	n := 100
	instructions := make([]Instruction, n)

	for i := 0; i < n; i++ {
		instructions[i] = Instruction{
			Type:  "calc",
			Op:    "+",
			Var:   fmt.Sprintf("v%d", i),
			Left:  int64(i),
			Right: int64(1),
		}
	}

	for i := 0; i < n; i++ {
		instructions = append(instructions, Instruction{
			Type: "print",
			Var:  fmt.Sprintf("v%d", i),
		})
	}

	results, err := calc.Calculate(instructions)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	if len(results) != n {
		t.Fatalf("Expected 100 results, got %d", len(results))
	}

	for i, res := range results {
		expected := int64(i + 1)
		if res.Value != expected {
			t.Errorf("For v%d expected %d, got %d", i, expected, res.Value)
		}
	}
}