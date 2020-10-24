package wallet

import (
	"errors"
	"github.com/Nappy-Says/wallet/pkg/types"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var defaultFavorite = types.Favorite{
	ID:        uuid.New().String(),
	AccountID: 1,
	Name:      types.CategoryIt,
	Amount:    10,
	Category:  types.CategoryIt,
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{
		Service: &Service{},
	}
}

func TestService_FindAccountByID_success(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")

	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("account => %v", account)
	}

}

func TestService_FindAccountByID_notFound(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")

	account, err := service.FindAccountByID(2)

	if err == nil {
		t.Errorf("method returned nil error, account => %v", account)
	}

}

func TestService_Reject_success_user(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")
	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("error => %v", err)
	}

	err = service.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("error => %v", err)
	}

	payment, err := service.Pay(account.ID, 10_00, "Food")

	if err != nil {
		t.Errorf("error => %v", err)
	}

	pay, err := service.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("error => %v", err)
	}

	err = service.Reject(pay.ID)

	if err != nil {
		t.Errorf("error => %v", err)
	}

}

func TestService_Reject_fail_user(t *testing.T) {
	var service Service
	service.RegisterAccount("9127660305")
	account, err := service.FindAccountByID(1)

	if err != nil {
		t.Errorf("account => %v", account)
	}

	err = service.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("error => %v", err)
	}

	payment, err := service.Pay(account.ID, 10_00, "Food")

	if err != nil {
		t.Errorf("account => %v", account)
	}

	pay, err := service.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("payment => %v", payment)
	}

	err = service.Reject(pay.ID + "uu")

	if err == nil {
		t.Errorf("pay => %v", pay)
	}

}

func TestService_RegisterAccount(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		phone types.Phone
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Account
		wantErr bool
	}{
		{
			name:   "user successfully registered.",
			fields: fields{},
			args:   args{phone: "9127660305"},
			want: &types.Account{
				ID:      1,
				Phone:   "9127660305",
				Balance: 0,
			},
			wantErr: false,
		},
		{
			name: "phone already registered",
			fields: fields{
				nextAccountID: 10,
				accounts:      Accounts(),
				payments:      nil,
			},
			args: args{
				phone: "9127660305",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.RegisterAccount(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Deposit(t *testing.T) {
	var accounts []*types.Account
	accounts = append(accounts,
		&types.Account{ID: 1, Phone: "9127660305", Balance: 0},
		&types.Account{ID: 3, Phone: "9127660307", Balance: 2},
		&types.Account{ID: 2, Phone: "9127660306", Balance: 1})

	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		accountID int64
		amount    types.Money
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "amount must be greater than zero",
			fields: fields{},
			args: args{
				accountID: 1,
				amount:    0,
			},
			wantErr: true,
		},
		{
			name: "account not found",
			fields: fields{
				nextAccountID: 0,
				accounts:      accounts,
				payments:      nil,
			},
			args: args{
				accountID: 4,
				amount:    10,
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				nextAccountID: 0,
				accounts:      accounts,
				payments:      nil,
			},
			args: args{
				accountID: 1,
				amount:    10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			if err := s.Deposit(tt.args.accountID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_FindAccountByID(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		accountID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.FindAccountByID(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAccountByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_FindPaymentByID(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		paymentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.FindPaymentByID(tt.args.paymentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPaymentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPaymentByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Pay(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		accountID int64
		amount    types.Money
		category  types.PaymentCategory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Payment
		wantErr bool
	}{
		{
			name:   "amount must be greater than zero",
			fields: fields{},
			args: args{
				accountID: 0,
				amount:    10,
				category:  types.CategoryFood,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "account not found",
			fields: fields{
				nextAccountID: 0,
				accounts:      Accounts(),
				payments:      nil,
			},
			args: args{
				accountID: 10,
				amount:    10,
				category:  types.CategoryIt,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not enough balance in account",
			fields: fields{
				nextAccountID: 0,
				accounts:      Accounts(),
				payments:      nil,
			},
			args: args{
				accountID: 4,
				amount:    10,
				category:  types.CategoryFood,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			got, err := s.Pay(tt.args.accountID, tt.args.amount, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pay() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Reject(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
	}
	type args struct {
		paymentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
			}
			if err := s.Reject(tt.args.paymentID); (err != nil) != tt.wantErr {
				t.Errorf("Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Accounts() []*types.Account {
	var accounts []*types.Account
	accounts = append(
		accounts,
		&types.Account{ID: 1, Phone: "9127660305", Balance: 0},
		&types.Account{ID: 2, Phone: "9127660306", Balance: 1},
		&types.Account{ID: 3, Phone: "9127660307", Balance: 2},
		&types.Account{ID: 4, Phone: "9127660307", Balance: 3})
	return accounts
}

func Payments() []*types.Payment {
	var payments []*types.Payment
	payments = append(
		payments,
		&types.Payment{
			ID:        "e1dceb29-6cc4-48c5-acd4-455530f9d50a",
			AccountID: 3,
			Amount:    1,
			Category:  types.CategoryFood,
			Status:    types.PaymentStatusInProgress,
		},
		&types.Payment{
			ID:        "5f7834e5-1c1e-42ff-bd78-9fd7f1b5da28",
			AccountID: 2,
			Amount:    1,
			Category:  types.CategoryIt,
			Status:    types.PaymentStatusInProgress,
		})
	return payments
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()

	account, err := s.AddAccountWithBalance("9127660305", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, types.CategoryIt)
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}

	newPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("newPayment => %v", newPayment)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	account, err := s.AddAccountWithBalance("9127660305", 100)
	if err != nil {
		t.Errorf("account => %v", account)
		return
	}

	payment, err := s.Pay(account.ID, 10, types.CategoryIt)
	if err != nil {
		t.Errorf("payment => %v", payment)
		return
	}

	favorite, err := s.FavoritePayment(payment.ID, types.CategoryIt)
	if err != nil {
		t.Errorf("favorite => %v", favorite)
		return
	}

	payFromFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("payFromFavorite => %v", payFromFavorite)
		return
	}
}

func TestService_FindFavoriteByID(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		favoriteID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.Favorite
		wantErr bool
	}{
		{
			name: "successFull find",
			fields: fields{
				nextAccountID: 0,
				accounts:      Accounts(),
				payments:      Payments(),
				favorites:     Favorites(),
			},
			args: args{
				favoriteID: defaultFavorite.ID,
			},
			want:    &defaultFavorite,
			wantErr: false,
		},
		{
			name: "successFull find",
			fields: fields{
				nextAccountID: 0,
				accounts:      Accounts(),
				payments:      Payments(),
				favorites:     Favorites(),
			},
			args: args{
				favoriteID: "nonExistingFavoriteID",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			got, err := s.FindFavoriteByID(tt.args.favoriteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFavoriteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFavoriteByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Favorites() []*types.Favorite {
	var favorites []*types.Favorite
	favorites = append(
		favorites,
		&types.Favorite{
			ID:        defaultFavorite.ID,
			AccountID: Accounts()[0].ID,
			Name:      types.CategoryIt,
			Amount:    10,
			Category:  types.CategoryIt,
		}, &types.Favorite{
			ID:        uuid.New().String(),
			AccountID: Accounts()[1].ID,
			Name:      types.CategoryIt,
			Amount:    10,
			Category:  types.CategoryIt,
		}, &types.Favorite{
			ID:        uuid.New().String(),
			AccountID: Accounts()[2].ID,
			Name:      types.CategoryIt,
			Amount:    10,
			Category:  types.CategoryIt,
		}, &types.Favorite{
			ID:        uuid.New().String(),
			AccountID: Accounts()[3].ID,
			Name:      types.CategoryIt,
			Amount:    10,
			Category:  types.CategoryIt,
		})
	return favorites
}

func TestService_ExportToFile(t *testing.T) {
	s := newTestService()
	_, _ = s.AddAccountWithBalance("9127660305", 10)
	_, _ = s.AddAccountWithBalance("9127660306", 11)
	_ = s.ExportToFile("../../data/accounts.txt")
}

func TestService_ImportFromFile(t *testing.T) {
	s := newTestService()
	_ = s.ImportFromFile("../../data/accounts.txt")
}

func TestService_Export(t *testing.T) {
	s := newTestService()
	account1, _ := s.AddAccountWithBalance("9127660305", 10)
	payment, _ := s.Pay(account1.ID, 10, types.CategoryIt)
	_, _ = s.FavoritePayment(payment.ID, types.CategoryIt)

	account2, _ := s.AddAccountWithBalance("9127660306", 11)
	payment2, _ := s.Pay(account2.ID, 10, types.CategoryIt)
	_, _ = s.FavoritePayment(payment2.ID, types.CategoryIt)

	_ = s.Export("../../data/")
}

func TestService_Import(t *testing.T) {
	s := newTestService()
	_ = s.Import("../../data/")
}
func TestService_Import2(t *testing.T) {
	dirname := "../../data/" + uuid.New().String()

	// создаём сервис
	s := newTestService()
	account1, _ := s.AddAccountWithBalance("9127660305", 10)
	payment, _ := s.Pay(account1.ID, 10, types.CategoryIt)
	_, _ = s.FavoritePayment(payment.ID, types.CategoryIt)

	account2, _ := s.AddAccountWithBalance("9127660306", 11)
	payment2, _ := s.Pay(account2.ID, 10, types.CategoryIt)
	_, _ = s.FavoritePayment(payment2.ID, types.CategoryIt)

	err := os.Mkdir(dirname, 0777)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Export(dirname)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = os.Stat(filepath.Join(dirname, "accounts.dump"))
	if err != nil {
		t.Error(err)
		return
	}

	_, err = os.Stat(filepath.Join(dirname, "payments.dump"))
	if err != nil {
		t.Error(err)
		return
	}

	i := newTestService()
	err = i.Import(dirname)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(s.accounts, i.accounts) {
		t.Error(errors.New("imported and exported accounts doesn't match"))
		return
	}

	if !reflect.DeepEqual(s.payments, i.payments) {
		t.Error(errors.New("imported and exported payments doesn't match"))
		return
	}
}

func TestService_HistoryToFiles(t *testing.T) {
	s := newTestService()
	account1, _ := s.AddAccountWithBalance("9127660305", 10000)
	_, _ = s.Pay(account1.ID, 10, types.CategoryIt)
	_, _ = s.Pay(account1.ID, 100, types.CategoryIt)
}

func BenchmarkService_SumPayments(b *testing.B) {
	s := newTestService()
	account1, _ := s.AddAccountWithBalance("9127660305", 10)
	_, _ = s.Pay(account1.ID, 10, types.CategoryIt)
	_, _ = s.Pay(account1.ID, 10, types.CategoryIt)
	_, _ = s.Pay(account1.ID, 10, types.CategoryIt)
	_, _ = s.Pay(account1.ID, 10, types.CategoryIt)

	want := types.Money(10)

	got := s.SumPayments(2)
	if want != got {
		b.Errorf(" error, want => %v got => %v", want, got)
	}
}

func TestService_SumPaymentsWithProgress(t *testing.T) {
	s := newTestService()
	for i := 0; i < 200_000; i++ {
		payment := &types.Payment{
			ID:     uuid.New().String(),
			Amount: types.Money(100),
		}
		s.payments = append(s.payments, payment)
	}

	s.SumPaymentsWithProgress()
	//want := make(chan types.Progress)
	//got := s.SumPaymentsWithProgress()
	//
	//if want != got {
	//	t.Errorf(" error, want => %v got => %v", want, got)
	//}
	//
}
