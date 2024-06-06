package main

import "fmt"

func main() {
	// info, err := GetUserInfo("asdf")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("%v\n", info)

	err := SendNtfy("", "topic")

	if err != nil {
		fmt.Println(err)
		return
	}
}
