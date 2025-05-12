package grpc

import (
	"context"
	"prac/calc"
)

type Server struct {
	UnimplementedCalculatorServiceServer                 
	calculator                           *calc.Calculator
}

func NewServer(calculator *calc.Calculator) *Server {
	return &Server{calculator: calculator}
}

func (s *Server) Calculate(ctx context.Context, req *CalculationRequest) (*CalculationResponse, error) {
	instructions := make([]calc.Instruction, len(req.Instructions))
	for i, instr := range req.Instructions {
		instructions[i] = convertProtoInstruction(instr)
	}

	results, err := s.calculator.Calculate(instructions)
	if err != nil {
		return nil, err
	}

	return &CalculationResponse{
		Items: convertToProtoResults(results),
	}, nil
}

func convertProtoInstruction(instr *Instruction) calc.Instruction {
	res := calc.Instruction{
		Type: instr.Type,
		Op:   instr.Op,
		Var:  instr.Var,
	}

	switch x := instr.Left.(type) {
	case *Instruction_LeftInt:
		res.Left = x.LeftInt
	case *Instruction_LeftVar:
		res.Left = x.LeftVar
	}

	switch x := instr.Right.(type) {
	case *Instruction_RightInt:
		res.Right = x.RightInt
	case *Instruction_RightVar:
		res.Right = x.RightVar
	}

	return res
}

func convertToProtoResults(results []calc.Result) []*Result {
	protoResults := make([]*Result, len(results))
	for i, res := range results {
		protoResults[i] = &Result{
			Var:   res.Var,
			Value: res.Value,
		}
	}
	return protoResults
}
