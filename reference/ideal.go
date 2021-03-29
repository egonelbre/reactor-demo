func (cashier *Cashier) PurchaseItem(item ItemID, conn ConnectionID) error {
	user := cashier.users.GetByConn(conn)

	data := cashier.someData
	err := REENTRY cashier.db.PurchaseItem(user, item)
	assert(data == cashier.someData)
	if !err {
		return wrap(err)
	}

	world = cashier.users.FindWorld(user)
	err := REENTRY world.AddItem(user, item)
	if err != nil {
		// TODO: rollback purchase?
		return wrap(err)
	}

	return nil
}