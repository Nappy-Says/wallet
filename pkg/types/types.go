package types

type Money int64
type PaymentCategory string
type PaymentStatus string
type Phone string

const (
	PaymentStatusOK         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

const (
	CategoryFood = "food"
	CategoryIt   = "it"
	CategoryShop = "shop"
)

type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}

type Progress struct {
	Part   int
	Result Money
}
