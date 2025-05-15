package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"prac/calc"
	"prac/grpcserver"
	"sync"

	pb "prac/proto"

	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		startHTTPServer()
	}()

	go func() {
		defer wg.Done()
		startGRPCServer()
	}()

	wg.Wait()
}

func startHTTPServer() {
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		var instructions []calc.Instruction
		calculator := calc.NewCalculator()
		if err := json.NewDecoder(r.Body).Decode(&instructions); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		results, err := calculator.Calculate(instructions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"items": results,
		})
	})

	fmt.Println("HTTP server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startGRPCServer() {
	calculator := calc.NewCalculator()
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, grpcserver.NewCalculatorServer(calculator))

	fmt.Println("gRPC server started at :9090")
	log.Fatal(grpcServer.Serve(lis))
}
