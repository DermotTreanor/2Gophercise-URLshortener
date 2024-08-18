package main

import (
	"fmt"
	"os"
	//"net/http"
)

func main() {
	//fmt.Println("hello")
	//result, err := http.Get("https://youtube.com")
	//if err != nil {
	//	fmt.Println("main: there was a poblem.", err)
	//}
	//fmt.Printf("Here is what we got back: \n%v\n\n", result)
	fd, err := os.Open("../main/data.yaml")
	if err != nil {
		fmt.Printf("main: there was an error opening the file: %v", err)
	} else {
		info, err := fd.Stat()
		if err != nil {
			fmt.Printf("main: unable to get file inforamtion: %v", err)

		} else {
			fmt.Println(info.Size())
		}
	}
}
