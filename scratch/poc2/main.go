package main

import "fmt"

type Event interface{}

type Users struct{}

func (u *Users) GetUserID(connid int) int             { return connid }
func (u *Users) FindWorld(userid int) (*World, error) { return &World{}, nil }

type Cashier struct {
	Users    *Users
	Database *Database
}

func (cashier *Cashier) React(ev Event) error {
	switch ev := p.(type) {
	case PurchaseItem:
		return cashier.PurchaseItem(ev)
	}
}

type PurchaseItem struct {
	ConnID int
	ItemID int
}

func (cashier *Cashier) PurchaseItem(ev PurchaseItem) error {
	userID := cashier.Users.ConnToUser(ev.ConnID)

	err := cashier.Database.PurchaseItem(userID, ev.ItemID)
	if err != nil {
		return fmt.Errorf("failed to PurchaseItem: %w", err)
	}

	world, err := cashier.Users.FindWorld(userID)
	if err != nil {
		return fmt.Errorf("failed to FindWorld: %w", err)
	}

	err = world.AddItem(userID, ev.ItemID)
	if err != nil {
		return fmt.Errorf("failed to AddItem: %w", err)
	}

	return nil
}

type World struct{}

func (world *World) AddItem(userID int, itemID int) error {
	return nil
}

type Database struct{}
