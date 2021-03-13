package main

type AnimationRendering struct {
}

func (r *AnimationRendering) React() {
	r.ProcessInput()
	r.Update()
	r.Render()
}

func (db *DB) GetAccountBalance(self ReactorID, user UserID, onComplete func(AccountBalance, error)) {

}

func (r *Reactor) BuyItem(user UserID, item ItemID) {
	item := &ItemBuy{
		Reactor: r,
		User:    user,
		Item:    item,
	}
	item.Run()
}

type BuyItemContext struct {
	*Reactor
	User UserID
	Item ItemID

	Pending Pending

	Errors  Errors
	Balance AccountBalance
	Items   []Item
	Item    Item
}

func (b *BuyItemContext) Run() {
	b.WaitFor(
		b.DB.AccountBalance(b.ID, b.User, b.BalanceReceived),
		b.DB.StoreItems(b.ID, b.User, b.ItemsReceieved),
	).Then(b.Select)
}

func (b *BuyItemContext) BalanceReceived(acc AccountBalance, err error) {
	if b.Errors.Include(err) {
		b.Abort(err)
		b.Client.ShowMessage(b.User, err)
		return
	}

	b.Balance = acc
}

func (b *BuyItemContext) ItemsReceived(items []Item, err error) {
	if b.Errors.Include(err) {
		b.Abort(err)
		b.Client.SendMessage(b.User, err)
		return
	}

	b.Items = items
}

func (b *BuyItemContext) Select() {
	b.Client.SelectItemToBuy(b.User, b.Items, b.Balance, b.ItemSelected)
}

func (b *BuyItemContext) ItemSelected(item Item, err error) {
	b.Item = item
	b.DB.BuyItem(b.User, item, b.ItemBought)
}

func (b *BuyItemContext) ItemBought(err error) {
	if b.Errors.Include(err) {
		b.Client.SendMessage(b.User, err)
		return
	}

	b.Users.AddItem(b.Item)
}
