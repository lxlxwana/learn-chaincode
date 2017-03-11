package main

import (
	"encoding/json"
	"fmt"
)

type RetOrder struct {
	SName string `json:"sname"`
	DName string `json:"dname"`
	ID    uint64 `json:"id"`
}

func main() {
	// var op = [4]uint64{1, 2, 3, 4}
	// var OrderID uint64
	// OrderID = 19``
	// str := fmt.Sprintf("%d", OrderID)
	// fmt.Println(str)
	// a := 4.21312
	// fmt.Println(int(a))
	var a [2]RetOrder
	a[0].SName = "hello"
	a[0].DName = "hwwwwo"
	a[0].ID = 100
	re, _ := json.Marshal(a)
	fmt.Printf("%s\n", re)
}

// const earthRadius = 6378.137

// func driverSelect(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
// 	var radLng1, radLat1, radLng2, radLat2 float64
// 	radLng1 = x1 * math.Pi / 180.0
// 	radLat1 = y1 * math.Pi / 180.0
// 	radLng2 = x2 * math.Pi / 180.0
// 	radLat2 = y2 * math.Pi / 180.0
// 	a := radLat1 - radLat2
// 	b := radLng1 - radLng2
// 	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
// 	s = s * earthRadius
// 	return s
// }
