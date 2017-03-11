package main

import (
	"fmt"

	"errors"

	"strconv"

	"encoding/json"

	"math"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// State defination
const (
	DRIVER_STATE_HAND          = 0
	DRIVER_STATE_PCIKUP        = 1
	DRIVER_STATE_ONGOING       = 2
	PASSENGER_STATE_HAND       = 0
	PASSENGER_STATE_WAITCOMPET = 1
	PASSENGER_STATE_WAITPICKUP = 2
	PASSENGER_STATE_ONGOING    = 3
	ORDER_STATE_INVALID        = 0
	ORDER_STATE_WAITCOMPET     = 1
	ORDER_STATE_ONGOING        = 2
	ORDER_STATE_FINISH         = 3
	ORDER_STATE_ABORT          = 4
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

// OrderPool is
type OrderPool struct {
	Total uint     `json:"total"`
	IDs   []uint64 `json:"ids"`
	Act   []bool   `json:"act"`
}

// User is
type User struct {
	Name           string  `json:"name"`
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	DriverInfo     string  `json:"dInfo"`
	DriverState    int32   `json:"dState"`
	PassengerInfo  string  `json:"pInfo"`
	PassengerState int32   `json:"pState"`
	Balance        int32   `json:"balance"`
	Role           int     `json:"role"`
	PwdHash        string  `json:"pwdHash"`
	OrderID        uint64  `json:"orderID"`
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

	err = stub.PutState("@counterOrder", []byte(fmt.Sprintf("%d", 1)))
	if err != nil {
		return nil, err
	}

	var newOP OrderPool
	jsonResult, err := json.Marshal(newOP)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("@orderPool", jsonResult)

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
	case "submit":
		return c.passengerSubmitOrder(stub, args)
	case "compet":
		return c.driverCompetOrder(stub, args)
	case "pickup":
		return c.driverPickUp(stub, args)
	case "finish":
		return c.driverFinishOrder(stub, args)
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
	case "isEnroll":
		return c.isEnroll(stub, args)
	case "getdriverstate":
		return c.getDriverState(stub, args)
	case "getpassstate:":
		return c.getPassengerState(stub, args)
	case "queryorderpool":
		return c.driverQueryOrderPool(stub, args)
	case "queryorderentry":
		return c.queryOrderFromEntry(stub, args)
	case "queryordertable":
		return c.queryOrderFromTable(stub, args)
	}

	fmt.Println("Query did not find func: " + function)
	return nil, errors.New("Received unknown function " + function)
}

//=================================================================================================================================//
//=================================================================================================================================//
//=================================================================================================================================//
//=================================================================================================================================//
//=================================================================================================================================//

//=================================================================================================================================//
// 主流程
//=================================================================================================================================//

// 用户名 密码 起点经度 起点纬度 终点经度 终点纬度 当前时间 起点地名 终点地名
func (c *Chaincode) passengerSubmitOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	passenger, err := c.getUser(stub, args[0])
	if err != nil {
		return nil, err
	}
	if passenger.PassengerState != PASSENGER_STATE_HAND {
		return nil, errors.New("Passenger must be at state \"hand\"")
	}

	var newOrder Order
	newOrder.ID, err = c.getNextID(stub)
	if err != nil {
		return nil, err
	}
	newOrder.Passenger = args[0]
	newOrder.StartX, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, err
	}
	newOrder.StartY, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
		return nil, err
	}
	newOrder.DestX, err = strconv.ParseFloat(args[4], 64)
	if err != nil {
		return nil, err
	}
	newOrder.DestY, err = strconv.ParseFloat(args[5], 64)
	if err != nil {
		return nil, err
	}
	newOrder.StartTime, err = strconv.ParseUint(args[6], 10, 64)
	if err != nil {
		return nil, err
	}
	newOrder.State = ORDER_STATE_WAITCOMPET

	err = c.setOrder(stub, fmt.Sprintf("%d", newOrder.ID), newOrder)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(fmt.Sprintf("@SPlace_%d", newOrder.ID), []byte(args[7]))
	if err != nil {
		return nil, err
	}
	err = stub.PutState(fmt.Sprintf("@DPlace_%d", newOrder.ID), []byte(args[8]))
	if err != nil {
		return nil, err
	}

	passenger.PassengerState = PASSENGER_STATE_WAITCOMPET
	passenger.OrderID = newOrder.ID
	err = c.setUser(stub, args[0], passenger)
	if err != nil {
		return nil, err
	}
	c.addOrderPool(stub, newOrder.ID)
	return []byte("success submit order"), nil
}

// 用户名 密码 订单号
func (c *Chaincode) driverCompetOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	driver, err := c.getUser(stub, args[0])
	if err != nil {
		return nil, err
	}
	if driver.DriverState != DRIVER_STATE_HAND {
		return nil, errors.New("Driver must be at state \"hand\"")
	}

	orderState, err := c.getOrderState(stub, args[2])
	if err != nil {
		return nil, err
	}
	if orderState != ORDER_STATE_WAITCOMPET {
		return nil, errors.New("order must be at state \"wait compet\"")
	}

	order, err := c.getOrder(stub, args[2])
	if err != nil {
		return nil, err
	}
	order.Driver = args[0]
	order.DriverInfo = driver.DriverInfo

	driver.DriverState = DRIVER_STATE_PCIKUP
	err = c.setUser(stub, args[0], driver)
	if err != nil {
		return nil, err
	}
	err = c.setPassengerState(stub, order.Passenger, PASSENGER_STATE_WAITPICKUP)
	if err != nil {
		return nil, err
	}
	c.deleteOrderPool(stub, order.ID)
	return []byte("success compet order"), nil
}

// 用户名 密码 当前时间
func (c *Chaincode) driverPickUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	driver, err := c.getUser(stub, args[0])
	if err != nil {
		return nil, err
	}
	if driver.DriverState != DRIVER_STATE_PCIKUP {
		return nil, errors.New("Driver must be at state \"pickup\"")
	}

	order, err := c.getOrder(stub, fmt.Sprintf("%d", driver.OrderID))
	if err != nil {
		return nil, err
	}

	now, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return nil, err
	}
	order.PickTime = now

	driver.DriverState = DRIVER_STATE_ONGOING
	err = c.setUser(stub, args[0], driver)
	if err != nil {
		return nil, err
	}
	err = c.setPassengerState(stub, order.Passenger, PASSENGER_STATE_ONGOING)
	if err != nil {
		return nil, err
	}
	return []byte("success pickup"), nil
}

// 用户名 密码 当前时间

const unitPriceTime = 1
const startingPrice = 1100
const unitPriceDistance = 1500

func (c *Chaincode) driverFinishOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	driver, err := c.getUser(stub, args[0])
	if err != nil {
		return nil, err
	}
	if driver.DriverState != DRIVER_STATE_ONGOING {
		return nil, errors.New("Driver must be at state \"ongoing\"")
	}
	order, err := c.getOrder(stub, fmt.Sprintf("%d", driver.OrderID))
	if err != nil {
		return nil, err
	}
	passenger, err := c.getUser(stub, order.Passenger)
	if err != nil {
		return nil, err
	}

	now, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return nil, err
	}
	order.EndTime = now
	order.State = ORDER_STATE_FINISH
	order.ActFeeTime = int32((order.EndTime - order.StartTime) * unitPriceTime)
	order.ActFeeDis = startingPrice + int32(distance(order.StartX, order.StartY, order.DestX, order.DestY)*unitPriceDistance)

	driver.DriverState = DRIVER_STATE_HAND
	err = c.setUser(stub, args[0], driver)
	if err != nil {
		return nil, err
	}
	passenger.PassengerState = PASSENGER_STATE_HAND
	err = c.setUser(stub, order.Passenger, passenger)
	if err != nil {
		return nil, err
	}
	err = c.writeOrder2Table(stub, order)
	if err != nil {
		return []byte("success finish, but cannot write to table"), nil
	}
	return []byte("success finish "), nil
}

//=================================================================================================================================//
//=================================================================================================================================//

// 用户名  密码哈希值 经度 纬度 是否接单
func (c *Chaincode) driverUpdatePosition(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

// 订单号
func (c *Chaincode) getDriverPosition(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	order, err := c.getOrder(stub, args[0])
	if err != nil {
		return nil, err
	}
	driver, err := c.getUser(stub, order.Driver)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("{\"driverPosition\":[%f,%f]", driver.X, driver.Y)), nil
}

// RetOrder is
type RetOrder struct {
	SName string `json:"sname"`
	DName string `json:"dname"`
	ID    uint64 `json:"id"`
}

// 用户名  密码
func (c *Chaincode) driverQueryOrderPool(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	userName := args[0]
	user, err := c.getUser(stub, userName)
	if err != nil {
		return nil, err
	}
	op, err := c.getOrderPool(stub, user)
	if err != nil {
		return nil, err
	}
	var ret [4]RetOrder
	for i := 0; i < 4; i++ {
		if op[i] != 0 {
			sName, err := stub.GetState(fmt.Sprintf("@SPlace_%d", op[i]))
			if err != nil {
				return nil, err
			}
			ret[i].SName = fmt.Sprintf("%s", sName)
			dName, err := stub.GetState(fmt.Sprintf("@DPlace_%d", op[i]))
			if err != nil {
				return nil, err
			}
			ret[i].DName = fmt.Sprintf("%s", dName)
			ret[i].ID = op[i]
		}
	}
	re, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}
	return re, nil
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

func (c *Chaincode) getOrderState(stub shim.ChaincodeStubInterface, idx string) (int32, error) {
	order, err := c.getOrder(stub, idx)
	if err != nil {
		return 0, err
	}
	return order.State, nil
}

func (c *Chaincode) setOrderState(stub shim.ChaincodeStubInterface, idx string, newState int32) error {
	old, err := c.getOrder(stub, idx)
	if err != nil {
		return err
	}
	old.State = newState
	err = c.setOrder(stub, idx, old)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chaincode) setDriverState(stub shim.ChaincodeStubInterface, userName string, newState int32) error {
	old, err := c.getUser(stub, userName)
	if err != nil {
		return err
	}
	old.DriverState = newState
	err = c.setUser(stub, userName, old)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chaincode) setPassengerState(stub shim.ChaincodeStubInterface, userName string, newState int32) error {
	old, err := c.getUser(stub, userName)
	if err != nil {
		return err
	}
	old.PassengerState = newState
	err = c.setUser(stub, userName, old)
	if err != nil {
		return err
	}
	return nil
}

//=================================================================================================================================//
//	identity check
//=================================================================================================================================//

//args: 用户名 密码SHA1哈希值 信息 角色类型(1司机 2乘客 0表示未注册)
//return:
func (c *Chaincode) enroll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var newUser User
	newUser.Name = args[0]
	newUser.PwdHash = args[1]
	role, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("role must be convertable to int")
	}
	switch role {
	case 1:
		newUser.DriverInfo = args[2]
	case 2:
		newUser.PassengerInfo = args[2]
	}
	newUser.Role = role
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
	if args[1] == user.PwdHash {
		return []byte(strconv.Itoa(user.Role)), nil
	}
	return []byte("0"), nil
}

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

func (c *Chaincode) getOrder(stub shim.ChaincodeStubInterface, idx string) (Order, error) {
	var re Order
	jsonResult, err := stub.GetState(idx)
	if err != nil {
		return re, err
	}
	err = json.Unmarshal(jsonResult, &re)
	if err != nil {
		return re, err
	}
	return re, nil
}

func (c *Chaincode) writeOrder2Table(stub shim.ChaincodeStubInterface, order Order) error {
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
		return fmt.Errorf("insert Table operation failed. %s", err)
	}
	if !ok {
		ok, err := stub.ReplaceRow("orders", row)
		if err != nil {
			return fmt.Errorf("replace Row operation failed. %s", err)
		}
		if !ok {
			return errors.New("replace Row operation failed. Row with given key does not exist")
		}
	}

	return err
}

//订单编号
func (c *Chaincode) queryOrderFromTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

//订单编号
func (c *Chaincode) queryOrderFromEntry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	jsonResult, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
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
	err = json.Unmarshal(jsonResult, &re)
	if err != nil {
		return re, err
	}
	return re, nil
}

func (c *Chaincode) ping(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return []byte("Hello, world!"), nil
}

//=================================================================================================================================//
//	some globe varaible
//=================================================================================================================================//

// 获得下一个空的订单ID，并且ID加1
func (c *Chaincode) getNextID(stub shim.ChaincodeStubInterface) (uint64, error) {
	var id uint64
	idByte, err := stub.GetState("@counterOrder")
	if err != nil {
		return 0, err
	}
	id, err = strconv.ParseUint(string(idByte), 10, 64)
	if err != nil {
		return 0, err
	}
	err = stub.PutState("@counterOrder", []byte(fmt.Sprintf("%d", id+1)))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Chaincode) getOrderPool(stub shim.ChaincodeStubInterface, driver User) ([4]uint64, error) {
	var result = [4]uint64{0, 0, 0, 0}
	var op OrderPool
	idx := 0
	jsonResult, err := stub.GetState("@orderPool")
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(jsonResult, &op)
	if err != nil {
		return result, err
	}

	var i uint
	for i = 0; i < op.Total && idx < 4; i++ {
		if !op.Act[i] {
			continue
		}
		order, err := c.getOrder(stub, fmt.Sprintf("%d", op.IDs[i]))
		if err != nil {
			return result, nil
		}
		if driverSelect(order, driver) {
			result[idx] = op.IDs[i]
			idx++
		}
	}

	return result, nil
}

func (c *Chaincode) addOrderPool(stub shim.ChaincodeStubInterface, id uint64) error {
	var op OrderPool
	jsonResult, err := stub.GetState("@orderPool")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonResult, &op)
	if err != nil {
		return err
	}

	op.Total++
	op.IDs = append(op.IDs, id)
	op.Act = append(op.Act, true)

	jsonResult, err = json.Marshal(op)
	if err != nil {
		return err
	}
	err = stub.PutState("@orderPool", jsonResult)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chaincode) deleteOrderPool(stub shim.ChaincodeStubInterface, id uint64) error {
	var op OrderPool
	jsonResult, err := stub.GetState("@orderPool")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonResult, &op)
	if err != nil {
		return err
	}

	op.Act[id] = false

	jsonResult, err = json.Marshal(op)
	if err != nil {
		return err
	}
	err = stub.PutState("@orderPool", jsonResult)
	if err != nil {
		return err
	}

	return nil
}

const earthRadius = 6378.137
const threadhold = 100.0

func driverSelect(order Order, driver User) bool {
	if distance(order.StartX, order.StartY, driver.X, driver.Y) < threadhold {
		return true
	}
	return false
}

func distance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	var radLng1, radLat1, radLng2, radLat2 float64
	radLng1 = x1 * math.Pi / 180.0
	radLat1 = y1 * math.Pi / 180.0
	radLng2 = x2 * math.Pi / 180.0
	radLat2 = y2 * math.Pi / 180.0
	a := radLat1 - radLat2
	b := radLng1 - radLng2
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * earthRadius
	return s
}
