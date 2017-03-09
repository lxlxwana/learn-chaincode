package main

import (
	"fmt"

	"errors"

	"strconv"

	"bytes"

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

//=================================================================================================================================//
//	Init Invoke & Query functions
//=================================================================================================================================//

// Init is
func (c *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	err := stub.CreateTable("orders", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "id", Type: shim.ColumnDefinition_UINT64, Key: true},
		&shim.ColumnDefinition{Name: "passenger", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "driver", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "startX", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "startY", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "destinationX", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "destinationY", Type: shim.ColumnDefinition_INT64, Key: false},
		&shim.ColumnDefinition{Name: "preFee", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "actFeeTime", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "actFeeDis", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "startTime", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "pickTime", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "endTime", Type: shim.ColumnDefinition_UINT64, Key: false},
		&shim.ColumnDefinition{Name: "state", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "passInfo", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "driverInfo", Type: shim.ColumnDefinition_STRING, Key: false},
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
	case "setOrder":
		return c.setOrder(stub, args)
	case "enroll":
		return nil, c.enroll(stub, args)
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
	case "getOrder":
		return c.getOrder(stub, args)
	case "test0":
		return c.test0(stub, args)
	case "test1":
		return c.test1(stub, args)
	}

	fmt.Println("Query did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

//=================================================================================================================================//
//=================================================================================================================================//
//=================================================================================================================================//
//=================================================================================================================================//
//=================================================================================================================================//

func (c *Chaincode) test0(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	re, err := stub.GetCallerCertificate()
	if err != nil {
		return nil, err
	}
	return re, nil
}
func (c *Chaincode) test1(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	ok, err := c.isPassengerOne(stub)
	if err != nil {
		return []byte("error"), err
	}
	if ok {
		return []byte("yes"), nil
	}
	return []byte("no"), nil
}

//=================================================================================================================================//
//	set & get state of passenger/driver to/from ledger
//=================================================================================================================================//

func (c *Chaincode) setState(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Expecting 2 arguments")
	}

	return nil, nil
}

//=================================================================================================================================//
//	identity check
//=================================================================================================================================//

// we have 4 user
// user_type1_0 user_type1_1 user_type2_0 user_type2_1
func (c *Chaincode) enroll(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 1 {
		return errors.New("Expecting 1 arguments")
	}

	ca, err := stub.GetCallerCertificate()
	if err != nil {
		return err
	}
	err = stub.PutState(args[0], ca)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chaincode) isPassengerOne(stub shim.ChaincodeStubInterface) (bool, error) {
	re, err := stub.GetCallerCertificate()
	if err != nil {
		return false, err
	}
	ca, err := stub.GetState("user_type1_0")
	if err != nil {
		return false, err
	}
	return bytes.Equal(re, ca), nil
}

func (c *Chaincode) isPassengerTwo(stub shim.ChaincodeStubInterface) (bool, error) {
	re, err := stub.GetCallerCertificate()
	if err != nil {
		return false, err
	}
	ca, err := stub.GetState("user_type1_1")
	if err != nil {
		return false, err
	}
	return bytes.Equal(re, ca), nil
}

func (c *Chaincode) isDriverOne(stub shim.ChaincodeStubInterface) (bool, error) {
	re, err := stub.GetCallerCertificate()
	if err != nil {
		return false, err
	}
	ca, err := stub.GetState("user_type2_0")
	if err != nil {
		return false, err
	}
	return bytes.Equal(re, ca), nil
}

func (c *Chaincode) isDriverTwo(stub shim.ChaincodeStubInterface) (bool, error) {
	re, err := stub.GetCallerCertificate()
	if err != nil {
		return false, err
	}
	ca, err := stub.GetState("user_type2_1")
	if err != nil {
		return false, err
	}
	return bytes.Equal(re, ca), nil
}

//=================================================================================================================================//
//	setOrder & getOrder to/from ledger
//=================================================================================================================================//

func (c *Chaincode) setOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 16 {
		return nil, errors.New("Expecting 16 arguments")
	}

	var id, startTime, pickTime, endTime uint64
	var passenger, driver, passInfo, driverInfo string
	var startX, startY, destinationX, destinationY, temp int64
	var preFee, actFeeTime, actFeeDis, state int32

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("args[0] must be convertable to uint64")
	}
	passenger = args[1]
	driver = args[2]
	startX, err = strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return nil, errors.New("args[3] must be convertable to int64")
	}
	startY, err = strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return nil, errors.New("args[4] must be convertable to int64")
	}
	destinationX, err = strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return nil, errors.New("args[5] must be convertable to int64")
	}
	destinationY, err = strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return nil, errors.New("args[6] must be convertable to int64")
	}
	temp, err = strconv.ParseInt(args[7], 10, 32)
	if err != nil {
		return nil, errors.New("args[7] must be convertable to int32")
	}
	preFee = int32(temp)
	temp, err = strconv.ParseInt(args[8], 10, 32)
	if err != nil {
		return nil, errors.New("args[8] must be convertable to int32")
	}
	actFeeTime = int32(temp)
	temp, err = strconv.ParseInt(args[9], 10, 32)
	if err != nil {
		return nil, errors.New("args[9] must be convertable to int32")
	}
	actFeeDis = int32(temp)
	startTime, err = strconv.ParseUint(args[10], 10, 64)
	if err != nil {
		return nil, errors.New("args[10] must be convertable to uint64")
	}
	pickTime, err = strconv.ParseUint(args[11], 10, 64)
	if err != nil {
		return nil, errors.New("args[11] must be convertable to uint64")
	}
	endTime, err = strconv.ParseUint(args[12], 10, 64)
	if err != nil {
		return nil, errors.New("args[12] must be convertable to uint64")
	}
	temp, err = strconv.ParseInt(args[13], 10, 32)
	if err != nil {
		return nil, errors.New("args[13] must be convertable to int32")
	}
	state = int32(temp)
	passInfo = args[14]
	driverInfo = args[15]

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_Uint64{Uint64: id}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: passenger}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: driver}}
	col4 := shim.Column{Value: &shim.Column_Int64{Int64: startX}}
	col5 := shim.Column{Value: &shim.Column_Int64{Int64: startY}}
	col6 := shim.Column{Value: &shim.Column_Int64{Int64: destinationX}}
	col7 := shim.Column{Value: &shim.Column_Int64{Int64: destinationY}}
	col8 := shim.Column{Value: &shim.Column_Int32{Int32: preFee}}
	col9 := shim.Column{Value: &shim.Column_Int32{Int32: actFeeTime}}
	col10 := shim.Column{Value: &shim.Column_Int32{Int32: actFeeDis}}
	col11 := shim.Column{Value: &shim.Column_Uint64{Uint64: startTime}}
	col12 := shim.Column{Value: &shim.Column_Uint64{Uint64: pickTime}}
	col13 := shim.Column{Value: &shim.Column_Uint64{Uint64: endTime}}
	col14 := shim.Column{Value: &shim.Column_Int32{Int32: state}}
	col15 := shim.Column{Value: &shim.Column_String_{String_: passInfo}}
	col16 := shim.Column{Value: &shim.Column_String_{String_: driverInfo}}

	columns = append(columns, &col1, &col2, &col3, &col4, &col5, &col6, &col7, &col8)
	columns = append(columns, &col9, &col10, &col11, &col12, &col13, &col14, &col15, &col16)
	row := shim.Row{Columns: columns}

	ok, err := stub.InsertRow("orders", row)
	if err != nil {
		return nil, fmt.Errorf("insert Table operation failed. %s", err)
	}
	if !ok {
		ok, err := stub.ReplaceRow("orders", row)
		if err != nil {
			return nil, fmt.Errorf("replace Row operation failed. %s", err)
		}
		if !ok {
			return nil, errors.New("replace Row operation failed. Row with given key does not exist")
		}
	}

	return nil, err
}

func (c *Chaincode) getOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("getOrder failed. Must include at least 1 key value")
	}

	var id uint64
	id, err := strconv.ParseUint(args[0], 10, 64)
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_Uint64{Uint64: id}}
	columns = append(columns, col1)
	row, err := stub.GetRow("orders", columns)
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

func (c *Chaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}
