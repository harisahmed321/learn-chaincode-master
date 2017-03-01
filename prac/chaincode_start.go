package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	A := args[0]
	val := args[1]
	err := stub.PutState(A, []byte(val))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" { //read a variable
		fmt.Println("hi there " + function) //error
		return nil, nil
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

// // Check user
// func (t *SimpleChaincode) AddUser(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

// 	var A string // Entities
// 	var err error

// 	if len(args) != 1 {
// 		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
// 	}

// 	user = args[0]

// 	// Get the state from the ledger
// 	checkExistingUser, err := stub.GetState(user)
// 	if err != nil {
// 		jsonResp := "{\"Error\":\"Failed to get state for user " + user + "\"}"
// 		return nil, errors.New(jsonResp)
// 	}

// 	if checkExistingUser == nil {
// 		createUser, err := stub.PutState(user)
// 		if err != null {
// 			jsonResp := "{\"Error\":\"Failed to create user \"}"
// 			return nil, errors.New(jsonResp)
// 		}
// 		jsonResp := "{\"Name\":\"" + user + "\"}"
// 		fmt.Printf("Query Response:%s\n", jsonResp)
// 		return user, nil
// 	}

// }
