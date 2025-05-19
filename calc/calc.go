package calc

import (
	"fmt"
	"sync"
	"time"
)

func NewCalculator() *Calculator {
	return &Calculator{
		vars:  sync.Map{},
		ready: make(map[string]*sync.WaitGroup),
	}
}

var operations = map[string]func(int64, int64) int64{
	"+": func(a, b int64) int64 { time.Sleep(50 * time.Millisecond); return a + b },
	"-": func(a, b int64) int64 { time.Sleep(50 * time.Millisecond); return a - b },
	"*": func(a, b int64) int64 { time.Sleep(50 * time.Millisecond); return a * b },
}

func (c *Calculator) Calculate(instructions []Instruction) ([]Result, error) {
	var calcOps []Instruction
	var printOps []Instruction

	for _, instr := range instructions {
		switch instr.Type {
		case "print":
			printOps = append(printOps, instr)
		case "calc":
			calcOps = append(calcOps, instr)
		default:
			return nil, fmt.Errorf("unknown operation: '%s'", instr.Type)
		}
	}

	for _, instr := range calcOps {
		if _, exists := c.ready[instr.Var]; !exists {
			c.ready[instr.Var] = &sync.WaitGroup{}
			c.ready[instr.Var].Add(1)
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(calcOps))

	for _, instr := range calcOps {
		wg.Add(1)
		go func(instr Instruction) {
			defer wg.Done()
			for _, dep := range getDependencies(instr) {
				c.ready[dep].Wait()
			}
			if err := c.processCalc(instr); err != nil {
				errChan <- err
				return
			}
			c.ready[instr.Var].Done()
		}(instr)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		return nil, err
	}

	var results []Result
	for _, printInstr := range printOps {
		if val, ok := c.vars.Load(printInstr.Var); ok {
			results = append(results, Result{Var: printInstr.Var, Value: val.(int64)})
		}
	}

	c.Reset()

	return results, nil
}

func getDependencies(instr Instruction) []string {
	var deps []string
	if v, ok := instr.Left.(string); ok {
		deps = append(deps, v)
	}
	if v, ok := instr.Right.(string); ok {
		deps = append(deps, v)
	}
	return deps
}

func (c *Calculator) processCalc(instr Instruction) error {
	if _, exists := c.vars.Load(instr.Var); exists {
		return fmt.Errorf("variable %s already exists", instr.Var)
	}

	left, err := c.getValue(instr.Left)
	if err != nil {
		return err
	}

	right, err := c.getValue(instr.Right)
	if err != nil {
		return err
	}

	op, ok := operations[instr.Op]
	if !ok {
		return fmt.Errorf("unknown operation %s", instr.Op)
	}

	c.vars.Store(instr.Var, op(left, right))
	return nil
}

func (c *Calculator) getValue(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case float64:
		return int64(val), nil
	case string:
		if stored, ok := c.vars.Load(val); ok {
			return stored.(int64), nil
		}
		return 0, fmt.Errorf("variable %s not defined", val)
	default:
		return 0, fmt.Errorf("invalid value type")
	}
}

func (c *Calculator) Reset() {
	c.vars.Range(func(key, value interface{}) bool {
		c.vars.Delete(key)
		return true
	})

	c.ready = make(map[string]*sync.WaitGroup)
}
