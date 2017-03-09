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
	}

	fmt.Println("Query did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

func (c *Chaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}

func (c *Chaincode) table(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var table shim.Table
	table.Name = "person file"
	var def0, def1 shim.ColumnDefinition
	def0.Name, def0.Type, def0.Key = "name", shim.ColumnDefinition_STRING, true
	def1.Name, def1.Type, def1.Key = "age", shim.ColumnDefinition_UINT32, false
	table.ColumnDefinitions = append(table.ColumnDefinitions, &def0, &def1)
	err := stub.CreateTable(table.Name, table.ColumnDefinitions)
	if err != nil {
		return nil, err
	}

	var row shim.Row
	var col0, col1 shim.Column
	var myname shim.Column_String_
	myname.String_ = "liang"
	var myage shim.Column_Int32
	myage.Int32 = 18
	col0.Value = &myname
	col1.Value = &myage
	row.Columns = append(row.Columns, &col0, &col1)
	Ok, err := stub.InsertRow(table.Name, row)
	if !Ok {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Row already exists for the given key")
	}

	var keys []shim.Column
	var key shim.Column
	key.Value = &myname
	keys = append(keys, key)
	newrow, err := stub.GetRow(table.Name, keys)
	if err != nil {
		return nil, err
	}

	var result []byte
	result = append(result, []byte(newrow.Columns[0].GetString_())...)
	result = append(result, []byte(fmt.Sprintf("%d", newrow.Columns[1].GetInt32()))...)
	return result, err
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
