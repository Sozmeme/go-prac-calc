package calc

import (
	"testing"
)

func TestGroupOperationsSimple(t *testing.T) {
	instructions := []Instruction{
		{Type: "calc", Op: "+", Var: "x", Left: int64(1), Right: int64(2)},
		{Type: "calc", Op: "*", Var: "y", Left: "x", Right: int64(5)},
		{Type: "calc", Op: "-", Var: "z", Left: "y", Right: int64(3)},
		{Type: "calc", Op: "+", Var: "a", Left: int64(10), Right: int64(2)},
		{Type: "calc", Op: "*", Var: "b", Left: "a", Right: int64(2)},
	}

	groups := groupOperations(instructions)

	var (
		group1Vars []string
		group2Vars []string
	)

	for _, group := range groups {
		vars := extractVars(group)
		if containsAll(vars, []string{"x", "y", "z"}) {
			group1Vars = vars
		} else if containsAll(vars, []string{"a", "b"}) {
			group2Vars = vars
		} else {
			t.Errorf("Unexpected group of vars: %v", vars)
		}
	}

	if len(group1Vars) == 0 {
		t.Error("Group with vars x, y, z not found")
	}
	if len(group2Vars) == 0 {
		t.Error("Group with vars a, b not found")
	}
}

func extractVars(instrs []Instruction) []string {
	var vars []string
	for _, instr := range instrs {
		vars = append(vars, instr.Var)
	}
	return vars
}

func containsAll(haystack, needles []string) bool {
	set := make(map[string]bool)
	for _, v := range haystack {
		set[v] = true
	}
	for _, v := range needles {
		if !set[v] {
			return false
		}
	}
	return true
}

func TestProcessCalc(t *testing.T) {
	calc := NewCalculator()
	err := calc.processCalc(Instruction{
		Type: "calc", Op: "+", Var: "a", Left: int64(2), Right: int64(3),
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if calc.vars["a"] != 5 {
		t.Errorf("expected 5, got %d", calc.vars["a"])
	}

	err = calc.processCalc(Instruction{
		Type: "calc", Op: "*", Var: "a", Left: int64(1), Right: int64(1),
	})
	if err == nil {
		t.Errorf("expected error on duplicate variable")
	}
}

func TestGetValueLocked(t *testing.T) {
	calc := NewCalculator()
	calc.vars["x"] = 42

	val, err := calc.getValueLocked("x")
	if err != nil || val != 42 {
		t.Errorf("expected 42, got %d, err: %v", val, err)
	}

	val, err = calc.getValueLocked(int64(10))
	if err != nil || val != 10 {
		t.Errorf("expected 10, got %d, err: %v", val, err)
	}

	val, err = calc.getValueLocked(float64(3.0))
	if err != nil || val != 3 {
		t.Errorf("expected 3, got %d, err: %v", val, err)
	}

	_, err = calc.getValueLocked("undefined")
	if err == nil {
		t.Error("expected error for undefined variable")
	}
}
