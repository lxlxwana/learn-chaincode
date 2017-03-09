package my

import (
	"fmt"

	"errors"

	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Chaincode struct {
}

// func main() {
// 	err := shim.Start(new(Chaincode))
// 	if err != nil {
// 		fmt.Printf("Error starting chaincode %s", err)
// 	}
// }

func (c *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var zero = []byte{'0'}
	var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err = stub.PutState("A", zero)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("B", zero)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("A+B", zero)

	return nil, nil
}

func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "init" {
		return c.Init(stub, "init", args)
	} else if function == "add" {
		return c.add(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (c *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "result" {
		return c.result(stub, args)
	}
	fmt.Println("query did not find func " + function)

	return nil, errors.New("Received unknown function query" + function)
}

func (c *Chaincode) result(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	key = "A+B"
	valBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valBytes, nil
}

func (c *Chaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var A, B, Plus string
	var Aval, Bval, Pval int
	var err error
	fmt.Println("running add()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	A = "A"
	B = "B"
	Plus = "A+B"
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value")
	}
	Bval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value")
	}
	Pval = Aval + Bval
	fmt.Printf("Aval = %d, Bval = %d, Pval = %d", Aval, Bval, Pval)

	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}
	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}
	err = stub.PutState(Plus, []byte(strconv.Itoa(Pval)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
