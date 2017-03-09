package main

import (
	"fmt"

	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// type order struct {
// 	Make            string `json:"make"`
// 	Model           string `json:"model"`
// 	Reg             string `json:"reg"`
// 	VIN             int    `json:"VIN"`
// 	Owner           string `json:"owner"`
// 	Scrapped        bool   `json:"scrapped"`
// 	Status          int    `json:"status"`
// 	Colour          string `json:"colour"`
// 	V5cID           string `json:"v5cID"`
// 	LeaseContractID string `json:"leaseContractID"`
// }

// Chaincode is
type Chaincode struct {
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode %s", err)
	}
}

// Init is
func (c *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	err := stub.CreateTable("personfile", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "name", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "age", Type: shim.ColumnDefinition_INT32, Key: false},
	})
	if err != nil {
		return nil, err
	}

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
	case "ping":
		return c.ping(stub)
	case "settable":
		return c.setTable(stub, args)
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
	case "ping":
		return c.ping(stub)
	case "gettable":
		return c.getTable(stub, args)
	}

	fmt.Println("Query did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

func (c *Chaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}

func (c *Chaincode) setTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// col1Val := args[0]
	// col2Int, err := strconv.ParseInt(args[1], 10, 32)
	// if err != nil {
	// 	return nil, errors.New("insertTableOne failed. arg[1] must be convertable to int32")
	// }
	// col2Val := int32(col2Int)
	// col3Int, err := strconv.ParseInt(args[2], 10, 32)
	// if err != nil {
	// 	return nil, errors.New("insertTableOne failed. arg[2] must be convertable to int32")
	// }
	// col3Val := int32(col3Int)

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "liang"}}
	col2 := shim.Column{Value: &shim.Column_Int32{Int32: 18}}
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("personfile", row)
	if err != nil {
		return nil, fmt.Errorf("insert Table operation failed. %s", err)
	}
	if !ok {
		return nil, errors.New("insert Table operation failed. Row with given key already exists")
	}

	return nil, err
}

func (c *Chaincode) getTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// 	if len(args) < 1 {
	// 	return nil, errors.New("getRowTableOne failed. Must include 1 key value")
	// }

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: "liang"}}
	columns = append(columns, col1)
	row, err := stub.GetRow("personfile", columns)
	if err != nil {
		return nil, fmt.Errorf("getRow Table operation failed. %s", err)
	}

	rowString := fmt.Sprintf("%s", row)
	return []byte(rowString), nil
}

//=================================================================================================================================//
//	Read & Write data to ledger
//=================================================================================================================================//

func (c *Chaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting key to query")
	}

	key := args[0]
	valBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valBytes, nil
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
