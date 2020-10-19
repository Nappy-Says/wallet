package main

import "github.com/Nappy-Says/wallet/pkg/wallet"


func main() {
	svc := &wallet.Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.ExportToFile("../data/export.txt")
	svc.ImportFromFile("../data/import.txt")
}
