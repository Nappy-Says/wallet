package wallet

import (
	"testing"
  "github.com/Nappy-Says/wallet/pkg/types"
)

func TestService_RegisterAccount_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_FindAccoundByIdmethod_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(2)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestDeposit(t *testing.T) {
	svc := Service{}

	svc.RegisterAccount("+992000000001")

	err := svc.Deposit(1, 100_00)
	if err != nil {
		t.Error("something wrong while paying")
	}

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_Reject_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Reject(pay.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}

func TestService_Reject_fail(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	editPayID := pay.ID + "mr.virus :)"
	err = svc.Reject(editPayID)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}

func TestService_Repeat_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err = svc.Repeat(pay.ID)
	if err != nil {
		t.Errorf("Repeat(): Error(): can't pay for an account(%v): %v", pay.ID, err)
	}
}

func TestService_Favorite_success_user(t *testing.T) {
	svc := Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Pay() Error() can't pay for an account(%v): %v", account, err)
	}

	favorite, err := svc.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Errorf("FavoritePayment() Error() can't for an favorite(%v): %v", favorite, err)
	}

	paymentFavorite, err := svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite() Error() can't for an favorite(%v): %v", paymentFavorite, err)
	}
}

func TestService_Export_success_user(t *testing.T) {
	var svc Service

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")

	err := svc.ExportToFile("export.txt")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

}

func TestService_Import_success_user(t *testing.T) {
	var svc Service

	err := svc.ImportFromFile("export.txt")

	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

}

func TestService_Export_success(t *testing.T) {
	svc := Service{}

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.RegisterAccount("+992000000004")

	err := svc.Export("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

	err = svc.Import("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
}

func TestService_ExportHistory_success_user(t *testing.T) {
	svc := Service{}

	acc, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", acc)
	}

	err = svc.Deposit(acc.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = svc.Pay(acc.ID, 1, "Cafe")
	_, err = svc.Pay(acc.ID, 2, "Auto")
	_, err = svc.Pay(acc.ID, 3, "MarketShop")
	if err != nil {
		t.Errorf("method Pay returned not nil error, err => %v", err)
	}

	payments, err := svc.ExportAccountHistory(acc.ID)
	if err != nil {
		t.Errorf("method ExportAccountHistory returned not nil error, err => %v", err)
	}

	err = svc.HistoryToFiles(payments, "../../data", 2)
	if err != nil {
		t.Errorf("method HistoryToFiles returned not nil error, err => %v", err)
	}
}

func BenchmarkSumPayment_user(b *testing.B) {
	var svc Service

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		b.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		b.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = svc.Pay(account.ID, 1, "Cafe")
	_, err = svc.Pay(account.ID, 2, "Cafe")
	_, err = svc.Pay(account.ID, 3, "Cafe")
	_, err = svc.Pay(account.ID, 4, "Cafe")
	_, err = svc.Pay(account.ID, 5, "Cafe")
	_, err = svc.Pay(account.ID, 6, "Cafe")
	_, err = svc.Pay(account.ID, 7, "Cafe")


	_, err = svc.Pay(account.ID, 8, "Cafe")
	_, err = svc.Pay(account.ID, 9, "Cafe")
	_, err = svc.Pay(account.ID, 10, "Cafe")
	_, err = svc.Pay(account.ID, 11, "Cafe")
	if err != nil {
		b.Errorf("method Pay returned not nil error, err => %v", err)
	}

	want := types.Money(66)

	got := svc.SumPayments(2)
	if want != got {
		b.Errorf(" error, want => %v got => %v", want, got)
	}

}
