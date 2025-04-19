package main

import (
	"fmt"
	"pratikshakuldeep456/digital-wallet-service/pkg/dws"
)

func main() {

	service := dws.GetDigitalWalletService()
	user1 := service.CreateUser("Pratiksha", "12345", dws.Profile{})
	fmt.Println(user1.Name, user1.Profile, user1.Phone)
	email := "xyz@gmail.com"
	address := "271206"
	service.UpdateProfile(user1.Id, &email, &address)
	user2 := service.UpdateProfile(user1.Id, nil, &address)
	fmt.Println(user2.Name, *user2.Profile, user2.Phone)
	err, isDeleted := user1.Delete()
	if err != nil {
		fmt.Println("Error deleting user:", err)
	} else {
		fmt.Println("User deleted successfully:", isDeleted)
	}
	// err1, data := user2.Delete()
	// if err1 != nil {
	// 	fmt.Println("Error deleting user:", err1)
	// } else {
	// 	fmt.Println("User deleted successfully:", data)
	// }
	err3 := user2.AddAccount(&dws.Account{})
	if err3 != nil {
		fmt.Println("Error adding account:", err3)
	} else {
		fmt.Println("Account added successfully")
	}

}
