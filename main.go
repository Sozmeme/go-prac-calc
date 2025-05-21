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

	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"

	_ "prac/docs"
)

type ResponseWrapper struct {
	Items []calc.Result `json:"items"`
}

// @title Calculator API
// @version 1.0
// @description This is a simple calculator API with both HTTP and gRPC interfaces.
// @host localhost:8080
// @BasePath /
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

// Calculate godoc
// @Summary Calculate operations
// @Description Perform a batch of calculations with 50ms delay per operation
// @Tags Calculator
// @Accept json
// @Produce json
// @Param instructions body []calc.Instruction true "Array of calculation instructions"
// @Success 200 {object} ResponseWrapper
// @Failure 400 {string} string "Invalid request format"
// @Failure 500 {string} string "Internal calculation error"
// @Router /calculate [post]
// @Example request
// [
//
//	{ "type": "calc", "op": "+", "var": "x", "left": 1, "right": 2 },
//	{ "type": "print", "var": "x" }
//
// ]
// @Example response
//
//	{
//	  "items": [
//	    { "var": "x", "value": 3 }
//	  ]
//	}
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseWrapper{Items: results})
	})

	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

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
