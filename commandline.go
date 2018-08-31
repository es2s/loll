package main

import (
	"fmt"

	"io"
	"os"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/marketplace"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func ManageRole(username, action, role string) {
	user, _ := marketplace.FindUserByUsername(username)
	if user == nil {
		fmt.Println("No such user")
		return
	}
	if action == "grant" && role == "seller" {
		user.IsSeller = !user.IsSeller
	} else if action == "grant" && role == "admin" {
		user.IsAdmin = !user.IsAdmin
	} else {
		fmt.Println("Wrong action")
		return
	}
	user.Save()
}

func RemoveUser(username string) {
	user, _ := marketplace.FindUserByUsername(username)
	if user == nil {
		panic("user not found")
	}
	user.Remove()
}

func IndexItems() {
	util.Log.Debug("[Index] Indexing items...")
	for _, item := range marketplace.GetAllItems() {
		util.Log.Debug("[Index] %s", item.Name)
		err := item.Index()
		if err != nil {
			util.Log.Error("%s", err)
		}
	}
}

func SearchItems(text string) {
	util.Log.Debug("[Index] Searching items...")
	marketplace.SearchItems(text)
}

func syncModels() {
	marketplace.SyncModels()

}
func updateStalledTransactions() {
	util.Log.Debug("[Transactions] UpdatingStalledTransactions")
	marketplace.TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions()
	marketplace.TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount()
}

func updateOldAndPending() {
	marketplace.UpdateOldFailedAndPendingTransactions()
}

func resendReleasedTransactions() {
	marketplace.ResendReleasedTransactions()
}

func fixImages() {

	copyFileContents := func(src, dst string) (err error) {
		in, err := os.Open(src)
		if err != nil {
			return
		}
		defer in.Close()
		out, err := os.Create(dst)
		if err != nil {
			return
		}
		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()
		if _, err = io.Copy(out, in); err != nil {
			return
		}
		err = out.Sync()
		return
	}

	items := marketplace.FindActiveItems()
	for i, _ := range items {
		filaname := "./data/images/" + items[i].Uuid + ".jpeg"
		placeholder := "./etc/images/green-owl.jpg"
		if _, err := os.Stat(filaname); os.IsNotExist(err) {
			copyFileContents(placeholder, filaname)
		}
	}
}

func staffStats() {
	interval := "2018-08-11 22:03"
	sTickets, err := marketplace.StaffSupportTicketsResolutionStats(interval)
	if err != nil {
		return
	}
	sDisputes, err := marketplace.StaffSupportDisputesResolutionStats(interval)
	if err != nil {
		return
	}
	sItems, err := marketplace.StaffItemApprovalStats(interval)
	if err != nil {
		return
	}

	var (
		text = fmt.Sprintf(
			`
Support Agent | Ticket Status | Number Of Tickets
--- | --- | ---
`)
	)
	for _, si := range sTickets {
		text += fmt.Sprintf("%s | TICKET %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}
	for _, si := range sDisputes {
		text += fmt.Sprintf("%s | DISPUTE %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}
	for _, si := range sItems {
		text += fmt.Sprintf("%s | ITEM %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}

	println(text)
}
