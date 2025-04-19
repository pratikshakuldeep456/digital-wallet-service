package main

import (
	"fmt"
	"pratikshakuldeep456/digital-wallet-service/pkg/dws"
	"pratikshakuldeep456/digital-wallet-service/pkg/dws/payment_methods"
	"time"
)

func main() {
	fmt.Println("=== DIGITAL WALLET SERVICE TEST SUITE ===")
	fmt.Println("Date: Saturday, April 19, 2025")
	fmt.Println("----------------------------------------")

	// Get the singleton instance of the wallet service
	service := dws.GetDigitalWalletService()

	// Initialize the Transactions map if not already done in GetDigitalWalletService
	if service.Transactions == nil {
		service.Transactions = make(map[int][]*dws.Transaction)
	}

	fmt.Println("\n=== USER MANAGEMENT SCENARIOS ===")

	// Create users
	user1 := service.CreateUser("Pratiksha", "12345", dws.Profile{
		Email:   "pratiksha@example.com",
		Address: "Initial Address",
	})
	fmt.Printf("User created: %s (ID: %d)\n", user1.Name, user1.Id)

	user2 := service.CreateUser("John", "67890", dws.Profile{
		Email:   "john@example.com",
		Address: "123 Main St",
	})
	fmt.Printf("User created: %s (ID: %d)\n", user2.Name, user2.Id)

	// Update profile information
	fmt.Println("\n=== PROFILE UPDATE SCENARIOS ===")
	email := "pratiksha.updated@example.com"
	address := "New Address 123"

	updatedUser, _ := service.UpdateProfile(user1.Id, &email, &address)
	if updatedUser != nil {
		fmt.Printf("Profile updated: Email: %s, Address: %s\n",
			updatedUser.Profile.Email, updatedUser.Profile.Address)
	} else {
		fmt.Println("Failed to update profile: User not found")
	}

	// Try updating a non-existent user
	nonExistentID := 9999
	invalidUser, _ := service.UpdateProfile(nonExistentID, &email, &address)
	if invalidUser == nil {
		fmt.Println("Correctly handled non-existent user update")
	}

	// Account management
	fmt.Println("\n=== ACCOUNT MANAGEMENT SCENARIOS ===")

	// Create accounts for users
	account1 := &dws.Account{
		Id:        dws.GenerateId(),
		AccountNo: "ACC-001",
		Balance:   1.0,
		Currency:  dws.USD,
		User:      user1,
	}

	account2 := &dws.Account{
		Id:        dws.GenerateId(),
		AccountNo: "ACC-002",
		Balance:   2.0,
		Currency:  dws.INR,
		User:      user2,
	}

	// Add accounts to users
	err := user1.AddAccount(account1)
	if err != nil {
		fmt.Printf("Error adding account to user1: %v\n", err)
	} else {
		fmt.Printf("Account %s added to %s successfully\n", account1.AccountNo, user1.Name)
	}

	err = user2.AddAccount(account2)
	if err != nil {
		fmt.Printf("Error adding account to user2: %v\n", err)
	} else {
		fmt.Printf("Account %s added to %s successfully\n", account2.AccountNo, user2.Name)
	}

	// Payment methods
	fmt.Println("\n=== PAYMENT METHOD SCENARIOS ===")

	// Add multiple payment methods
	bankPayment := payment_methods.NewBank(
		"Example Bank",
		"Pratiksha Kuldeep",
		"1234567890",
		"EXBK0001234",
		user1,
	)
	service.AddPaymentMethod(user1.Id, bankPayment)

	creditCardPayment := payment_methods.NewCreditCard(
		"1234-5678-9012-3456",
		"12/25",
		"123",
		"Pratiksha Kuldeep",
		user1,
	)
	service.AddPaymentMethod(user1.Id, creditCardPayment)

	upiPayment := payment_methods.NewUpi(
		"pratiksha@upi",
		user1,
	)
	service.AddPaymentMethod(user1.Id, upiPayment)

	fmt.Printf("Added 3 payment methods for %s\n", user1.Name)

	// List payment methods
	paymentMethods := service.GetPaymentMethods(user1.Id)
	fmt.Printf("%s has %d payment methods\n", user1.Name, len(paymentMethods))

	// Test GetId() method on payment methods
	fmt.Println("\n=== PAYMENT METHOD ID TESTING ===")
	for i, method := range paymentMethods {
		fmt.Printf("Payment method #%d: Type=%s, ID=%d\n",
			i+1, method.PaymentMethod(), method.GetId())
	}

	// Transaction scenarios
	fmt.Println("\n=== TRANSACTION SCENARIOS ===")

	// Deposit funds
	depositAmount := 500.0
	err = account1.Deposit(depositAmount)
	if err != nil {
		fmt.Printf("Error depositing funds: %v\n", err)
	} else {
		fmt.Printf("Deposited %.2f %s to account %s. New balance: %.2f\n",
			depositAmount, account1.Currency, account1.AccountNo, account1.GetBalance())
	}

	// Withdraw funds - successful case
	withdrawAmount := 200.0
	err = account1.Withdraw(withdrawAmount)
	if err != nil {
		fmt.Printf("Error withdrawing funds: %v\n", err)
	} else {
		fmt.Printf("Withdrew %.2f %s from account %s. New balance: %.2f\n",
			withdrawAmount, account1.Currency, account1.AccountNo, account1.GetBalance())
	}

	// Withdraw funds - insufficient balance
	largeWithdrawal := 1000.0
	err = account1.Withdraw(largeWithdrawal)
	if err != nil {
		fmt.Printf("Expected error for large withdrawal: %v\n", err)
	} else {
		fmt.Printf("Withdrew %.2f %s from account %s. New balance: %.2f\n",
			largeWithdrawal, account1.Currency, account1.AccountNo, account1.GetBalance())
	}

	// Fund account from payment method
	fmt.Println("\n=== FUND ACCOUNT FROM PAYMENT METHOD ===")

	// Assuming you've implemented the FundAccountFromPaymentMethod method
	fundAmount := 100.0
	_, err = service.FundAccountFromPaymentMethod(user1.Id, account1.Id, bankPayment.GetId(), fundAmount, dws.USD)
	if err != nil {
		fmt.Printf("Error funding account from bank: %v\n", err)
	} else {
		fmt.Printf("Successfully funded account %s with %.2f %s using bank payment method\n",
			account1.AccountNo, fundAmount, dws.USD)
		fmt.Printf("New balance: %.2f %s\n", account1.GetBalance(), account1.Currency)
	}

	// Try funding with credit card
	fundAmount = 50.0
	_, err = service.FundAccountFromPaymentMethod(user1.Id, account1.Id, creditCardPayment.GetId(), fundAmount, dws.USD)
	if err != nil {
		fmt.Printf("Error funding account from credit card: %v\n", err)
	} else {
		fmt.Printf("Successfully funded account %s with %.2f %s using credit card payment method\n",
			account1.AccountNo, fundAmount, dws.USD)
		fmt.Printf("New balance: %.2f %s\n", account1.GetBalance(), account1.Currency)
	}

	// Try funding with UPI with currency conversion
	fundAmount = 1000.0
	_, err = service.FundAccountFromPaymentMethod(user1.Id, account1.Id, upiPayment.GetId(), fundAmount, dws.INR)
	if err != nil {
		fmt.Printf("Error funding account from UPI: %v\n", err)
	} else {
		fmt.Printf("Successfully funded account %s with %.2f %s using UPI payment method\n",
			account1.AccountNo, fundAmount, dws.INR)
		fmt.Printf("New balance: %.2f %s\n", account1.GetBalance(), account1.Currency)
	}

	// Transfer funds between accounts with currency conversion
	fmt.Println("\n=== FUND TRANSFER SCENARIOS ===")
	transferAmount := 50.0
	err, success := service.TransferFunds(account1, account2, transferAmount, dws.USD)
	if err != nil {
		fmt.Printf("Error transferring funds: %v\n", err)
	} else if success {
		fmt.Printf("Transferred %.2f %s from %s to %s\n",
			transferAmount, dws.USD, account1.AccountNo, account2.AccountNo)
		fmt.Printf("%s balance: %.2f %s\n", account1.AccountNo, account1.GetBalance(), account1.Currency)
		fmt.Printf("%s balance: %.2f %s\n", account2.AccountNo, account2.GetBalance(), account2.Currency)
	}

	// Try transferring with insufficient funds
	largeTransfer := 10000.0
	err, success = service.TransferFunds(account1, account2, largeTransfer, dws.USD)
	if err != nil {
		fmt.Printf("Expected error for large transfer: %v\n", err)
	} else if !success {
		fmt.Printf("Transfer of %.2f %s failed as expected\n", largeTransfer, dws.USD)
	}

	// User deletion scenarios
	fmt.Println("\n=== USER DELETION SCENARIOS ===")

	// Create a temporary user for deletion testing
	tempUser := service.CreateUser("TempUser", "99999", dws.Profile{
		Email:   "temp@example.com",
		Address: "Temp Address",
	})
	fmt.Printf("Temporary user created: %s (ID: %d)\n", tempUser.Name, tempUser.Id)

	// Delete the user
	err, isDeleted := tempUser.Delete()
	if err != nil {
		fmt.Printf("Error deleting user: %v\n", err)
	} else {
		fmt.Printf("User %s deleted successfully: %v\n", tempUser.Name, isDeleted)
	}

	// Try to add an account to a deleted user
	newAccount := &dws.Account{
		Id:        dws.GenerateId(),
		AccountNo: "ACC-TEMP",
		Balance:   0.0,
		Currency:  dws.USD,
	}

	err = tempUser.AddAccount(newAccount)
	if err != nil {
		fmt.Printf("Expected error when adding account to deleted user: %v\n", err)
	} else {
		fmt.Printf("Unexpectedly succeeded in adding account to deleted user\n")
	}

	// Try to delete an already deleted user
	err, _ = tempUser.Delete()
	if err != nil {
		fmt.Printf("Expected error when deleting already deleted user: %v\n", err)
	} else {
		fmt.Printf("Unexpectedly succeeded in deleting already deleted user\n")
	}

	// Concurrency scenario
	fmt.Println("\n=== CONCURRENCY SCENARIOS ===")

	// Create a new account for concurrent operations
	concurrentAccount := &dws.Account{
		Id:        dws.GenerateId(),
		AccountNo: "CONC-001",
		Balance:   100.0,
		Currency:  dws.USD,
		User:      user2,
	}
	user2.AddAccount(concurrentAccount)
	fmt.Printf("Created concurrent testing account with initial balance: %.2f %s\n",
		concurrentAccount.GetBalance(), concurrentAccount.Currency)

	// Simulate concurrent deposits and withdrawals
	go func() {
		for i := 0; i < 5; i++ {
			concurrentAccount.Deposit(10.0)
			fmt.Printf("Concurrent deposit #%d: +10.00\n", i+1)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(50 * time.Millisecond) // Ensure deposits start first
	go func() {
		for i := 0; i < 5; i++ {
			concurrentAccount.Withdraw(5.0)
			fmt.Printf("Concurrent withdrawal #%d: -5.00\n", i+1)
			time.Sleep(150 * time.Millisecond)
		}
	}()

	// Wait for concurrent operations to complete
	time.Sleep(1 * time.Second)
	fmt.Printf("Final balance after concurrent operations: %.2f %s\n",
		concurrentAccount.GetBalance(), concurrentAccount.Currency)

	// Transaction history
	fmt.Println("\n=== TRANSACTION HISTORY ===")
	transactions := account1.GetTransactions()
	fmt.Printf("Account %s has %d transactions\n", account1.AccountNo, len(transactions))

	// Display transaction details if available
	if len(transactions) > 0 {
		fmt.Println("\nTransaction details for account " + account1.AccountNo + ":")
		for i, txn := range transactions {
			srcAcct := "External"
			if txn.SrcAccount != nil {
				srcAcct = txn.SrcAccount.AccountNo
			}

			destAcct := "External"
			if txn.DesAccount != nil {
				destAcct = txn.DesAccount.AccountNo
			}

			fmt.Printf("%d. Type: %s, Amount: %.2f, From: %s, To: %s, Date: %s\n",
				i+1, txn.TransactionType, txn.Amount, srcAcct, destAcct,
				txn.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	userTransactions := service.GetTransactionHistory(user2.Id)
	if userTransactions != nil {
		fmt.Printf("\nUser %s has %d transactions\n", user2.Name, len(userTransactions))

		// Display user transaction details if available
		if len(userTransactions) > 0 {
			fmt.Println("\nTransaction details for user " + user2.Name + ":")
			for i, txn := range userTransactions {
				srcAcct := "External"
				if txn.SrcAccount != nil {
					srcAcct = txn.SrcAccount.AccountNo
				}

				destAcct := "External"
				if txn.DesAccount != nil {
					destAcct = txn.DesAccount.AccountNo
				}

				fmt.Printf("%d. Type: %s, Amount: %.2f, From: %s, To: %s, Date: %s\n",
					i+1, txn.TransactionType, txn.Amount, srcAcct, destAcct,
					txn.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
	}

	// Payment method removal
	fmt.Println("\n=== PAYMENT METHOD REMOVAL ===")
	service.RemovePaymentMethod(user1.Id, dws.UpiTransfer)
	paymentMethods = service.GetPaymentMethods(user1.Id)
	fmt.Printf("After removal, %s has %d payment methods\n", user1.Name, len(paymentMethods))

	// Edge case: Try to fund account with removed payment method
	fmt.Println("\n=== EDGE CASE: USING REMOVED PAYMENT METHOD ===")
	_, err = service.FundAccountFromPaymentMethod(user1.Id, account1.Id, upiPayment.GetId(), 100.0, dws.USD)
	if err != nil {
		fmt.Printf("Expected error when using removed payment method: %v\n", err)
	} else {
		fmt.Printf("Unexpectedly succeeded in using removed payment method\n")
	}

	fmt.Println("\n=== PROGRAM COMPLETED ===")
}
