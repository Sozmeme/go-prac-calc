package calc

import (
	"fmt"
	"sync"
	"time"
)

func NewCalculator() *Calculator {
	return &Calculator{vars: make(map[string]int64)}
}

var operations = map[string]func(int64, int64) int64{
	"+": func(a, b int64) int64 { time.Sleep(50 * time.Millisecond); return a + b },
	"-": func(a, b int64) int64 { time.Sleep(50 * time.Millisecond); return a - b },
	"*": func(a, b int64) int64 { time.Sleep(50 * time.Millisecond); return a * b },
}

func (c *Calculator) Calculate(instructions []Instruction) ([]Result, error) {
	var results []Result

	var calcOps []Instruction
	printQueue := make([]Instruction, 0)

	for _, instr := range instructions {
		if instr.Type == "print" {
			printQueue = append(printQueue, instr)
		} else {
			calcOps = append(calcOps, instr)
		}
	}

	groups := groupOperations(calcOps)

	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		go func(grp []Instruction) {
			defer wg.Done()
			for _, op := range grp {
				if err := c.processCalc(op); err != nil {
					return
				}
			}
		}(group)
	}

	go func() {
		wg.Wait()
		for _, printInstr := range printQueue {
			val, ok := c.vars[printInstr.Var]
			if !ok {
				continue
			}
			results = append(results, Result{Var: printInstr.Var, Value: val})
		}
	}()

	wg.Wait()
	return results, nil
}

func groupOperations(instructions []Instruction) [][]Instruction {
	instMap := make(map[string]Instruction)
	graph := make(map[string][]string)
	reverseGraph := make(map[string][]string)

	for _, instr := range instructions {
		instMap[instr.Var] = instr

		var deps []string
		if v, ok := instr.Left.(string); ok {
			deps = append(deps, v)
		}
		if v, ok := instr.Right.(string); ok {
			deps = append(deps, v)
		}
		graph[instr.Var] = deps

		for _, dep := range deps {
			reverseGraph[dep] = append(reverseGraph[dep], instr.Var)
		}
	}

	visited := make(map[string]bool)
	var groups [][]Instruction

	for varName := range instMap {
		if visited[varName] {
			continue
		}

		// Собираем компоненту зависимых инструкций
		groupVars := make(map[string]struct{})
		collectComponent(varName, graph, reverseGraph, visited, groupVars)

		var group []Instruction
		for v := range groupVars {
			group = append(group, instMap[v])
		}
		groups = append(groups, group)
	}

	return groups
}

func collectComponent(v string, g, rev map[string][]string, visited map[string]bool, group map[string]struct{}) {
	if visited[v] {
		return
	}
	visited[v] = true
	group[v] = struct{}{}

	for _, neighbor := range g[v] {
		collectComponent(neighbor, g, rev, visited, group)
	}
	for _, neighbor := range rev[v] {
		collectComponent(neighbor, g, rev, visited, group)
	}
}


func (c *Calculator) processCalc(instr Instruction) error {
	if _, exists := c.vars[instr.Var]; exists {
		return fmt.Errorf("variable %s already exists", instr.Var)
	}

	left, err := c.getValueLocked(instr.Left)
	if err != nil {
		return err
	}

	right, err := c.getValueLocked(instr.Right)
	if err != nil {
		return err
	}

	op, ok := operations[instr.Op]
	if !ok {
		return fmt.Errorf("unknown operation %s", instr.Op)
	}

	c.vars[instr.Var] = op(left, right)
	return nil
}

func (c *Calculator) getValueLocked(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case float64:
		return int64(val), nil
	case string:
		if stored, ok := c.vars[val]; ok {
			return stored, nil
		}
		return 0, fmt.Errorf("variable %s not defined", val)
	default:
		return 0, fmt.Errorf("invalid value type")
	}
}