package dws

import (
	"errors"
	"fmt"
)

type User struct {
	Id        int
	Name      string
	Phone     string
	Profile   *Profile
	Accounts  []*Account
	IsDeleted bool
}

func CreateUser(name, phone string, profile Profile) *User {
	return &User{
		Id:        GenerateId(),
		Name:      name,
		Phone:     phone,
		Profile:   &profile,
		IsDeleted: false}
}

func (u *User) UpdateProfile(email, address *string) {

	if email != nil {
		u.Profile.Email = *email
	}
	if address != nil {
		u.Profile.Address = *address
	}

	fmt.Printf("User %d profile updated: %+v\n", u.Id, u.Profile)
}

func (u *User) Delete() (error, bool) {
	if u.IsDeleted {
		return errors.New("User already deleted"), false
	}
	u.IsDeleted = true
	return nil, true
}
func (u *User) GetAccountStatus() bool {
	return u.IsDeleted
}

func (u *User) AddAccount(a *Account) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	data := u.GetAccountStatus()
	if data {
		return errors.New("User is deleted")
	}
	u.Accounts = append(u.Accounts, a)
	return nil
}

func (u *User) Removeccount(id int) error {
	data := u.GetAccountStatus()
	if data {
		return errors.New("User is deleted")
	}
	for i, account := range u.Accounts {
		if account.Id == id {
			u.Accounts = append(u.Accounts[:i], u.Accounts[i+1:]...)
			return nil
		}

	}
	return errors.New("Account not found")
}
