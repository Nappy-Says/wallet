package types

// Money представляет собой денежную сумму в минималных единиц (центы, копейки, дирамы и т.д).
type Money int64

// PaymentCategory представляет собой категорию, в которию был совершен платёж (авто, аптека, рестараны и т.д).
type PaymentCategory string

// PaymentStatus представляет статус плятежей.
type PaymentStatus string

// Предопределеные статусы платежа.
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment представляет игформацию о платиже.
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

// Phone представляет информацию о номере телефона
type Phone string

// Account представляет информацию о счёте пользлователя.
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

// Favorite представляет информацию о Избранное
type Favorite struct {
	ID        string
	AccountID int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}
