package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Structure
type User struct {
	UserID     string  `json:"userid"`
	UserAmount float64 `json:"useramount"`
}

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
	// if len(args) != 1 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	// }
	// var key string
	// for i := 0; i < 3; i++ {
	// 	key = args[i]
	// 	err := stub.PutState(key, []byte(args[i+1]))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "addUser" {
		return t.addUser(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for fun
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running addUser()")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting UserId and UserAmount only")
	}
	var user User
	err = json.Unmarshal([]byte(args[0]), &user)
	if err != nil {
		fmt.Println("error invalid")
		return nil, error.New("error invalid user")
	}
	key = user.UserID
	value = strconv.FormatFloat(user.UserAmount, 'f', -1, 64)
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

type Transaction struct {
	UserA             string  `json:"usera"`
	UserAval          float64 `json:"useraval"`
	UserB             string  `json:"userb"`
	UserBval          float64 `json:"userbval"`
	UserAdmin         string  `json:"useradmin"`
	UserAdminval      float64 `json:"useradminval"`
	TransactionAmount float64 `json:"transactionamount"`
	Tax               float64 `json:"tax"`
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var A, B, Admin string
	//var Aval, Bval, AdminVal, X float64
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting Transactions")
	}

	var trans Transaction

	A = args[0]
	B = args[1]
	X, err := strconv.ParseFloat(args[2], 64)

	Admin = "admin"
	Adminbytes, err := stub.GetState(Admin)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Adminbytes == nil {
		return nil, errors.New("Entity not found")
	}
	AdminVal, _ := strconv.ParseFloat(string(Adminbytes), 64)

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ := strconv.ParseFloat(string(Avalbytes), 64)

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ := strconv.ParseFloat(string(Bvalbytes), 64)

	// Perform the execution
	Aval = Aval - X
	AdminVal = AdminVal + (X * 0.20)
	Bval = Bval + (X - (X * 0.20))

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.FormatFloat(Aval, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(Admin, []byte(strconv.FormatFloat(AdminVal, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.FormatFloat(Bval, 'f', -1, 64)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
