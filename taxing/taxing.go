package main

import (
	"fmt"

	"errors"

	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Chaincode is
type Chaincode struct {
	data int
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode %s", err)
	}
}

// Init is
func (c *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	c.data = 100
	return nil, nil
}

// Invoke is
func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke is running: " + function)

	// function selle
	switch function {
	case "init":
		return c.Init(stub, "init", args)
	case "write":
		return c.write(stub, args)
	}

	fmt.Println("Invoke did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

// Query is
func (c *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query is running: " + function)

	// function selection
	switch function {
	case "read":
		return c.read(stub, args)
	}

	fmt.Println("Query did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

func (c *Chaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting key to query")
	}

	//key := args[0]
	//valBytes, err := stub.GetState(key)
	// if err != nil {
	// 	jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
	// 	return nil, errors.New(jsonResp)
	// }
	return []byte(strconv.Itoa(c.data)), nil
}

func (c *Chaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, key and value to invoke")
	}

	key := args[0]
	value := args[1]
	err := stub.PutState(key, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
