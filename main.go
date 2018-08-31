package main

import (
	"fmt"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "server" {
		RunServer()
	} else if argsWithoutProg[0] == "sync" {
		syncModels()
	} else if argsWithoutProg[0] == "user" {
		username, action, role := argsWithoutProg[1], argsWithoutProg[2], argsWithoutProg[3]
		ManageRole(username, action, role)
	} else if argsWithoutProg[0] == "remove" {
		RemoveUser(argsWithoutProg[1])
	} else if argsWithoutProg[0] == "index" {
		IndexItems()
	} else if argsWithoutProg[0] == "search" {
		SearchItems(argsWithoutProg[1])
	} else if argsWithoutProg[0] == "update-stalled-transactions" {
		updateStalledTransactions()
	} else if argsWithoutProg[0] == "update-old-pending" {
		updateOldAndPending()
	} else if argsWithoutProg[0] == "resend-released" {
		resendReleasedTransactions()
	} else if argsWithoutProg[0] == "staff-stats" {
		staffStats()
	} else if argsWithoutProg[0] == "fix-images" {
		fixImages()
	} else {
		fmt.Println("wrong command")
	}
}
