package main

import (
	"fmt"

	"errors"

	"strconv"

	"encoding/json"

	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// State defination
const (
	DRIVER_STATE_HAND           = 0
	DRIVER_STATE_COMPECT        = 1
	DRIVER_STATE_PCIKUP         = 2
	DRIVER_STATE_ONGOING        = 3
	PASSENGER_STATE_HAND        = 0
	PASSENGER_STATE_WAITCOMPECT = 1
	PASSENGER_STATE_WAITPICKUP  = 2
	PASSENGER_STATE_ONGOING     = 3
)

// Chaincode is
type Chaincode struct {
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode %s", err)
	}
}

// Order is
type Order struct {
	ID         uint64  `json:"id"`
	Passenger  string  `json:"passenger"`
	Driver     string  `json:"driver"`
	StartX     float64 `json:"startX"`
	StartY     float64 `json:"startY"`
	DestX      float64 `json:"destX"`
	DestY      float64 `json:"destY"`
	ActFeeTime int32   `json:"actFeeTime"`
	ActFeeDis  int32   `json:"actFeeDis"`
	StartTime  uint64  `json:"startTime"`
	PickTime   uint64  `json:"pickTime"`
	EndTime    uint64  `json:"endTime"`
	State      int32   `json:"state"`
	PassInfo   string  `json:"passInfo"`
	DriverInfo string  `json:"driverInfo"`
}

// User is
type User struct {
	Name            string    `json:"name"`
	X               float64   `json:"x"`
	Y               float64   `json:"y"`
	DriverInfo      string    `json:"dInfo"`
	DriverState     int32     `json:"dState"`
	DriverOrderPool [8]uint64 `json:"orderpool"`
	PassengerInfo   string    `json:"pInfo"`
	PassengerState  int       `json:"pState"`
	Balance         int32     `json:"balance"`
	Role            int       `json:"role"`
	PwdHash         []byte    `jons:"pwdHash"`
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
		&shim.ColumnDefinition{Name: "startX", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "startY", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "destinationX", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "destinationY", Type: shim.ColumnDefinition_STRING, Key: false},
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

	// var newUser User
	// err = c.setUser(stub, "user_type1_0", newUser)
	// if err != nil {
	// 	return nil, err
	// }
	// err = c.setUser(stub, "user_type2_0", newUser)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// Invoke is
func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke is running: " + function)
	switch function {
	case "init":
		return c.Init(stub, "init", args)
	case "write":
		return c.write(stub, args)
	case "ping":
		return c.ping(stub)
	case "enroll":
		return c.enroll(stub, args)
	}
	fmt.Println("Invoke did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

// Query is
func (c *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query is running: " + function)
	switch function {
	case "read":
		return c.read(stub, args)
	case "ping":
		return c.ping(stub)
	case "queryOrder":
		return c.queryOrder(stub, args)
	case "test0":
		return c.test0(stub, args)
	case "test1":
		return c.test1(stub, args)
	case "isEnroll":
		return c.isEnroll(stub, args)
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
	return nil, nil
}
func (c *Chaincode) test1(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================//
// 主流程
//=================================================================================================================================//

// 用户名 密码 起点经度 起点纬度 终点经度 终点纬度 当前时间
func (c *Chaincode) passengerSubmitOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// if(passengerStates[msg.sender] != 0) //乘客必须处于挂起状态才能抢单
	// {
	// 	return 0;
	// }

	// //创建新的订单
	// passengerToOrder[msg.sender] = counterOrderIndex;
	// orders[counterOrderIndex].id = counterOrderIndex;
	// orders[counterOrderIndex].passenger = msg.sender;
	// orders[counterOrderIndex].driver = 0x0;
	// orders[counterOrderIndex].s_x = s_x;
	// orders[counterOrderIndex].s_y = s_y;
	// orders[counterOrderIndex].d_x = d_x;
	// orders[counterOrderIndex].d_y = d_y;
	// orders[counterOrderIndex].distance = 0;//calculateDistance(s_x, d_x, s_y, d_y);
	// orders[counterOrderIndex].preFee = penaltyPrice + calculatePreFee(s_x, s_y, d_x, d_y);
	// orders[counterOrderIndex].actFee = 0;
	// orders[counterOrderIndex].actFeeTime = 0;
	// orders[counterOrderIndex].startTime = time;
	// orders[counterOrderIndex].state = 1;
	// orders[counterOrderIndex].passInfo = passInfo;
	// orders[counterOrderIndex].sName = sName;
	// orders[counterOrderIndex].dName = dName;
	// counterOrderIndex++;
	// passengerStates[msg.sender] = 1; //乘客订单分配中

	// if(!driverSelction(s_x, s_y, counterOrderIndex-1))
	// {
	// 	orders[counterOrderIndex-1].state = 4;
	// 	passengerStates[msg.sender] = 0;
	// 	return 0;
	// }

	// return counterOrderIndex-1;
	return nil, nil
}

//=================================================================================================================================//
//=================================================================================================================================//

// 用户名  密码哈希值 经度 纬度 是否接单
func (c *Chaincode) driverUpdatePosition(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// userName, userPWD := args[0], args[1]
	// if flag, err := c.isDriverOne(stub, userName, userPWD); !flag {
	// 	return nil, err
	// }
	userName := args[0]
	old, err := c.getUser(stub, userName)
	if err != nil {
		return nil, err
	}
	old.X, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, errors.New("args[2] must be convertable to float")
	}
	old.Y, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
		return nil, errors.New("args[3] must be convertable to float")
	}
	err = c.setUser(stub, userName, old)
	if err != nil {
		return nil, err
	}
	return []byte("success update driver position"), nil
}

//=================================================================================================================================//
//	set & get state of passenger/driver to/from ledger
//=================================================================================================================================//

func (c *Chaincode) getDriverState(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	userName := fmt.Sprintf("%s", args[0])
	user, err := c.getUser(stub, userName)
	if err != nil {
		return nil, err
	}
	return []byte(strconv.Itoa(int(user.DriverState))), nil
}

func (c *Chaincode) getPassengerState(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	userName := fmt.Sprintf("%s", args[0])
	user, err := c.getUser(stub, userName)
	if err != nil {
		return nil, err
	}
	return []byte(strconv.Itoa(int(user.PassengerState))), nil
}

//=================================================================================================================================//
//	identity check
//=================================================================================================================================//

//args: 用户名 密码SHA1哈希值 信息 角色类型(1司机 2乘客 0表示未注册)
//return:
func (c *Chaincode) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var newUser User
	newUser.Name = args[0]
	newUser.PwdHash = []byte(args[1])
	role, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("arguments must be convertable to int")
	}
	switch role {
	case 0:
		newUser.DriverInfo = args[2]
	case 1:
		newUser.PassengerInfo = args[2]
	}
	newUser.Balance = 100000000
	err = c.setUser(stub, args[0], newUser)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//args: 用户名 密码
//return: 0未注册 1司机 2乘客
func (c *Chaincode) isEnroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	user, err := c.getUser(stub, args[0])
	if err != nil {
		return nil, err
	}
	if bytes.Equal(user.PwdHash, []byte(args[1])) {
		return []byte(strconv.Itoa(user.Role)), nil
	}
	return []byte("0"), nil
}

// func (c *Chaincode) isPassengerOne(stub shim.ChaincodeStubInterface, userName string, pwdHash string) (bool, error) {
// 	if userName != "user_type1_0" {
// 		return false, nil
// 	}
// 	re, err := stub.GetState("user_type1_0_pwd")
// 	if err != nil {
// 		return false, err
// 	}
// 	return bytes.Equal([]byte(pwdHash), re), nil
// }

// func (c *Chaincode) isDriverOne(stub shim.ChaincodeStubInterface, userName string, pwdHash string) (bool, error) {
// 	if userName != "user_type2_0" {
// 		return false, nil
// 	}
// 	re, err := stub.GetState("user_type2_0_pwd")
// 	if err != nil {
// 		return false, err
// 	}
// 	return bytes.Equal([]byte(pwdHash), re), nil
// }

//=================================================================================================================================//
//	setOrder & getOrder to/from ledger
//=================================================================================================================================//

func (c *Chaincode) setOrder(stub shim.ChaincodeStubInterface, key string, order Order) error {
	jsonResult, err := json.Marshal(order)
	if err != nil {
		return err
	}
	err = stub.PutState(key, jsonResult)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chaincode) getOrder(stub shim.ChaincodeStubInterface, key string) (Order, error) {
	var re Order
	jsonResult, err := stub.GetState(key)
	if err != nil {
		return re, err
	}
	err = json.Unmarshal(jsonResult, re)
	if err != nil {
		return re, err
	}
	return re, nil
}

func (c *Chaincode) writeOrder2Table(stub shim.ChaincodeStubInterface, order Order) ([]byte, error) {
	startX := fmt.Sprintf("%f", order.StartX)
	startY := fmt.Sprintf("%f", order.StartY)
	destinationX := fmt.Sprintf("%f", order.DestX)
	destinationY := fmt.Sprintf("%f", order.DestY)

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_Uint64{Uint64: order.ID}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: order.Passenger}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: order.Driver}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: startX}}
	col5 := shim.Column{Value: &shim.Column_String_{String_: startY}}
	col6 := shim.Column{Value: &shim.Column_String_{String_: destinationX}}
	col7 := shim.Column{Value: &shim.Column_String_{String_: destinationY}}
	col9 := shim.Column{Value: &shim.Column_Int32{Int32: order.ActFeeTime}}
	col10 := shim.Column{Value: &shim.Column_Int32{Int32: order.ActFeeDis}}
	col11 := shim.Column{Value: &shim.Column_Uint64{Uint64: order.StartTime}}
	col12 := shim.Column{Value: &shim.Column_Uint64{Uint64: order.PickTime}}
	col13 := shim.Column{Value: &shim.Column_Uint64{Uint64: order.EndTime}}
	col14 := shim.Column{Value: &shim.Column_Int32{Int32: order.State}}
	col15 := shim.Column{Value: &shim.Column_String_{String_: order.PassInfo}}
	col16 := shim.Column{Value: &shim.Column_String_{String_: order.DriverInfo}}

	columns = append(columns, &col1, &col2, &col3, &col4, &col5, &col6, &col7)
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

func (c *Chaincode) queryOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

//=================================================================================================================================//
//	set & get Order
//=================================================================================================================================//

func (c *Chaincode) setUser(stub shim.ChaincodeStubInterface, key string, user User) error {
	jsonResult, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = stub.PutState(key, jsonResult)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chaincode) getUser(stub shim.ChaincodeStubInterface, key string) (User, error) {
	var re User
	jsonResult, err := stub.GetState(key)
	if err != nil {
		return re, err
	}
	err = json.Unmarshal(jsonResult, re)
	if err != nil {
		return re, err
	}
	return re, nil
}

func (c *Chaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}
