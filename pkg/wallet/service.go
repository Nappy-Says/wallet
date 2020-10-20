package wallet

import (
	"sync"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"bufio"
	"github.com/Nappy-Says/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")
var ErrFileNotFound = errors.New("file not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment

	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

//Deposit method
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}

	account.Balance += amount
	return nil
}

func (s *Service) Reject(paymentID string) error {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	acc, err := s.FindAccountByID(pay.AccountID)
	if err != nil {
		return err
	}

	pay.Status = types.PaymentStatusFail
	acc.Balance += pay.Amount

	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}


func ReadFile(file *os.File) ([]byte, error) {
	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		content = append(content, buf[:read]...)
	}
	return content, nil
}

func (s *Service) ExportToFile(path string) error {
	r := ""
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	for _, acc := range s.accounts {
		ID := strconv.Itoa(int(acc.ID)) + ";"
		Phone := string(acc.Phone) + ";"
		Balance := strconv.Itoa(int(acc.Balance))

		r += ID
		r += Phone
		r += Balance + "|"
	}
	_, err = file.Write([]byte(r))
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}

	return nil
}

func AddAccountToFile(file *os.File, account *types.Account) error {
	content, err := ReadFile(file)
	if err != nil {
		return err
	}
	r := string(content)

	ID := strconv.Itoa(int(account.ID))
	Phone := string(account.Phone)
	Balance := strconv.Itoa(int(account.Balance))

	r += (ID + ";")
	r += (Phone + ";")
	r += (Balance + "|")
	file.Write([]byte(r))
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	content, err := ReadFile(file)

	var accounts []string = strings.Split(string(content), "|")

	for _, account := range accounts[:len(accounts)-1] {
		var vals []string = strings.Split(account, ";")
		id, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		balance, err := strconv.Atoi(vals[2])
		if err != nil {
			return err
		}
		newAccount := &types.Account{
			ID:      int64(id),
			Phone:   types.Phone(vals[1]),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, newAccount)
	}
	return nil
}

func (s *Service) Export(dir string) error {
	acc:= 0
	for _, account := range s.accounts {
		ID := strconv.FormatInt(account.ID, 10) + ";"
		phone := string(account.Phone) + ";"
		balance := strconv.FormatInt(int64(account.Balance), 10)
		err := WriteToFile(dir+"/accounts.dump", []byte(ID+phone+balance+"\n"))
		if err != nil {
			return err
		}
		acc++
	}
	log.Print("acc: ", acc)

	
	pay := 0
	for _, payment := range s.payments {
		ID := payment.ID + ";"
		AccountID := strconv.FormatInt(payment.AccountID, 10) + ";"
		Amount := strconv.FormatInt(int64(payment.Amount), 10) + ";"
		Category := string(payment.Category) + ";"
		Status := string(payment.Status) + "\n"
		err := WriteToFile(dir+"/payments.dump", []byte(ID+AccountID+Amount+Category+Status))
		if err != nil {
			return err
		}
		pay++
	}
	log.Print("pay: ", pay)


	fav := 0
	for _, favorite := range s.favorites {
		ID := favorite.ID + ";"
		AccountID := strconv.FormatInt(favorite.AccountID, 10) + ";"
		Name := favorite.Name + ";"
		Amount := strconv.FormatInt(int64(favorite.Amount), 10) + ";"
		Category := string(favorite.Category) + "\n"
		err := WriteToFile(dir+"/favorites.dump", []byte(ID+AccountID+Name+Amount+Category))
		fav++
		if err != nil {
			return err
		}
	}
	log.Print("fav: ", fav)
	return nil
}

func WriteToFile(fileName string, data []byte) error {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Print("error: ", closeErr)
		}
	}()
	_, err = file.Write(data)

	if err != nil {
		log.Print("error: ", err)
	}
	return nil
}

func (s *Service) Import(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return err
	}
	for _, file := range files {
		read, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			log.Print(err)
			return err
		}
		defer func() {
			if closeErr := read.Close(); closeErr != nil {
				log.Print(closeErr)
			}
		}()

		reader := bufio.NewReader(read)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Print("EOF: ", line)
				break
			}
			if err != nil {
				log.Print(err)
				return err
			}

			item := strings.Split(line, ";")
			switch file.Name() {
			case "accounts.dump":
				acc := s.convertToAccount(item)
				if acc != nil {
					s.accounts = append(s.accounts, acc)
				}
			case "favorites.dump":
				favorite := s.convertToFavorites(item)
				if favorite != nil {
					s.favorites = append(s.favorites, favorite)
				}
			case "payments.dump":
				payment := s.convertToPayments(item)
				if payment != nil {
					s.payments = append(s.payments, payment)
				}
			default:
				break
			}
		}

	}
	return nil
}

func (s *Service) convertToAccount(item []string) *types.Account {
	ID, _ := strconv.ParseInt(item[0], 10, 64)
	balance, _ := strconv.ParseInt(removeEndLine(item[2]), 10, 64)
	account, err := s.FindAccountByID(ID)
	if err != nil {
		s.nextAccountID++
		return &types.Account{
			ID:      ID,
			Phone:   types.Phone(item[1]),
			Balance: types.Money(balance),
		}
	}
	account.ID = ID
	account.Phone = types.Phone(item[1])
	account.Balance = types.Money(balance)
	return nil
}

func (s *Service) convertToFavorites(item []string) *types.Favorite {
	AccountID, _ := strconv.ParseInt(item[1], 10, 64)
	Amount, _ := strconv.ParseInt(item[3], 10, 64)

	favorite, err := s.FindFavoriteByID(item[0])
	if err != nil {
		return &types.Favorite{
			ID:        item[0],
			AccountID: AccountID,
			Name:      item[2],
			Amount:    types.Money(Amount),
			Category:  types.PaymentCategory(item[4]),
		}
	}
	favorite.ID = item[0]
	favorite.AccountID = AccountID
	favorite.Name = item[2]
	favorite.Amount = types.Money(Amount)
	favorite.Category = types.PaymentCategory(removeEndLine(item[4]))
	return nil
}

func (s *Service) convertToPayments(item []string) *types.Payment {
	AccountID, _ := strconv.ParseInt(item[1], 10, 64)
	Amount, _ := strconv.ParseInt(item[2], 10, 64)

	payment, err := s.FindPaymentByID(item[0])
	if err != nil {
		return &types.Payment{
			ID:        item[0],
			AccountID: AccountID,
			Amount:    types.Money(Amount),
			Category:  types.PaymentCategory(item[3]),
			Status:    types.PaymentStatus(removeEndLine(item[4])),
		}
	}
	payment.ID = item[0]
	payment.AccountID = AccountID
	payment.Amount = types.Money(Amount)
	payment.Category = types.PaymentCategory(item[3])
	payment.Status = types.PaymentStatus(item[4])
	return nil
}

func removeEndLine(balance string) string {
	return strings.TrimRightFunc(balance, func(c rune) bool {
		return c == '\r' || c == '\n'
	})
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	var payments []types.Payment
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments = append(payments, *payment)
		}
	}
	if len(payments) <= 0 {
		return nil, ErrAccountNotFound
	}
	return payments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) > 0 {
		if len(payments) <= records {
			file, _ := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			defer file.Close()

			var str string
			for _, pay := range payments {
				str += fmt.Sprint(pay.ID) + ";" + fmt.Sprint(pay.AccountID) + ";" + fmt.Sprint(pay.Amount) + ";" + fmt.Sprint(pay.Category) + ";" + fmt.Sprint(pay.Status) + "\n"
			}
			file.WriteString(str)
		} else {

			var str string
			num := 0
			num1 := 1
			var file *os.File
			for _, pay := range payments {
				if num == 0 {
					file, _ = os.OpenFile(dir+"/payments"+fmt.Sprint(num1)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
				}
				num++
				str = fmt.Sprint(pay.ID) + ";" + fmt.Sprint(pay.AccountID) + ";" + fmt.Sprint(pay.Amount) + ";" + fmt.Sprint(pay.Category) + ";" + fmt.Sprint(pay.Status) + "\n"
				_, _ = file.WriteString(str)
				if num == records {
					str = ""
					num1++
					num = 0
					file.Close()
				}
			}

		}
	}
	return nil
}

func (s *Service) SumPayments(goroutines int) types.Money {

	if goroutines < 1 {
		goroutines = 1
	}
	pays := (len(s.payments) / goroutines) + 1
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	sum := types.Money(0)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		Sum := types.Money(0)
		go func(iteration int) {
			defer wg.Done()
			pay := iteration * pays
			pay1 := (iteration * pays) + pays
			for i := pay; i < pay1; i++ {
				if i > len(s.payments)-1 {
					break
				} 
				sum += s.payments[i].Amount
			}
			mu.Lock()
			defer mu.Unlock()
			Sum += sum
		}(i)
	}
	wg.Wait()
	return sum
}