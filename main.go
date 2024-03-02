//Problem1

// package main

// import (
// 	"fmt"
// )

// func bossBaby(s string) {
// 	length := len(s)
// 	fmt.Println(length)

// 	if length > 1000000 || length < 1 {
// 		fmt.Println("Invalid string length")
// 		return
// 	}
// 	R := 0
// 	S := 0
// 	for i, char := range s {
// 		// Boss Baby shoots first
// 		if i == 0 && char == 'R' {
// 			fmt.Println("Bad boy")
// 			return
// 		}
// 		// Specify only 'R' or 'S'.
// 		if char != 'R' && char != 'S' {
// 			fmt.Println("Invalid character")
// 			return
// 		}
// 		if char == 'R' {
// 			R++
// 		}
// 		if char == 'S' {
// 			S++
// 		}
// 	}
// 	if R > S {
// 		fmt.Println("Good boy")
// 	} else {
// 		fmt.Println("Bad boy")
// 	}
// }

// func main() {
// 	var str string = "SRSSRRR"
// 	bossBaby(str)
// 	var str1 string = "RSSRR"
// 	bossBaby(str1)
// 	var str2 string = "SSSRRRRS"
// 	bossBaby(str2)
// }

//Problem2

// package main

// import "fmt"

// func RescueChickensWithRoof(chickens []int, roof int) {
// 	maxChicken := 0
// 	for i := 0; i < len(chickens); i++ {
// 		maxRoof := chickens[i] + (roof - 1)
// 		max := 0
// 		for j := 0; j < len(chickens); j++ {
// 			if j >= i {
// 				if chickens[j] <= maxRoof {
// 					max++
// 				}
// 			}

// 			if maxChicken < max {
// 				maxChicken = max
// 			}
// 		}
// 	}
// 	fmt.Println(maxChicken)
// }

// func main() {
// 	chickens := []int{2, 5, 10, 12, 15}
// 	roof := 5
// 	RescueChickensWithRoof(chickens, roof)

// }

// Problem3
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Transaction struct {
	Symbol    string `json:"symbol"`
	Price     uint64 `json:"price"`
	Timestamp uint64 `json:"timestamp"`
}

func main() {
	tx := Transaction{
		Symbol:    "ETH",
		Price:     4500,
		Timestamp: 1678912345,
	}

	payload, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post("https://mock-node-wgqbnxruha-as.a.run.app/broadcast", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	fmt.Println("Transaction Hash:", response["tx_hash"])
	for {
		url := "https://mock-node-wgqbnxruha-as.a.run.app/check/" + response["tx_hash"]
		resp, err = http.Get(url)
		if err != nil {
			fmt.Println("Error sending GET request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var statusResponse map[string]string
		if err := json.Unmarshal(body, &statusResponse); err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return
		}
		txStatus := statusResponse["tx_status"]
		fmt.Println("Transaction Status:", txStatus)

		switch txStatus {
		case "CONFIRMED":
			fmt.Println("ธุรกรรมได้รับการประมวลผลและได้รับการยืนยัน")
		case "FAILED":
			fmt.Println("ธุรกรรมล้มเหลวในการประมวลผล")
		case "PENDING":
			fmt.Println("ธุรกรรมกำลังรอการประมวลผล")
		case "DNE":
			fmt.Println("ธุรกรรมไม่มีอยู่")
		default:
			fmt.Println("สถานะธุรกรรมไม่ระบุ")
		}
		if txStatus != "PENDING" {
			break
		}
		time.Sleep(3 * time.Second)
	}
}
