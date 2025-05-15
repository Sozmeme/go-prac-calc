package grpcserver

import (
	"context"
	"prac/calc"
	pb "prac/proto"
)

type calculatorServer struct {
	pb.UnimplementedCalculatorServiceServer
	calcService *calc.Calculator
}

func NewCalculatorServer(calcService *calc.Calculator) *calculatorServer {
	return &calculatorServer{calcService: calcService}
}

func (s *calculatorServer) Calculate(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	instructions := make([]calc.Instruction, len(req.Instructions))
	for i, instr := range req.Instructions {
		instructions[i] = convertProtoInstruction(instr)
	}

	results, err := s.calcService.Calculate(instructions)
	if err != nil {
		return nil, err
	}

	return &pb.CalculationResponse{
		Items: convertToProtoResults(results),
	}, nil
}

func convertProtoInstruction(instr *pb.Instruction) calc.Instruction {
	res := calc.Instruction{
		Type: instr.Type,
		Op:   instr.Op,
		Var:  instr.Var,
	}

	switch x := instr.Left.(type) {
	case *pb.Instruction_LeftInt:
		res.Left = x.LeftInt
	case *pb.Instruction_LeftVar:
		res.Left = x.LeftVar
	}

	switch x := instr.Right.(type) {
	case *pb.Instruction_RightInt:
		res.Right = x.RightInt
	case *pb.Instruction_RightVar:
		res.Right = x.RightVar
	}

	return res
}

func convertToProtoResults(results []calc.Result) []*pb.Result {
	protoResults := make([]*pb.Result, len(results))
	for i, res := range results {
		protoResults[i] = &pb.Result{
			Var:   res.Var,
			Value: res.Value,
		}
	}
	return protoResults
}
