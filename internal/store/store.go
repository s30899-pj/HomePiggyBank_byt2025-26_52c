package store

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type Household struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedByID uint         `json:"created_by_id"`
	CreatedBy   User         `gorm:"foreignKey:CreatedByID" json:"created_by"`
	Memberships []Membership `gorm:"foreignKey:HouseholdID" json:"memberships"`
}

type Membership struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	HouseholdID uint      `json:"household_id"`
	Household   Household `gorm:"foreignKey:HouseholdID" json:"household"`
	Role        string    `json:"role"`
}

type ExpenseCategory string

const (
	CategoryFood          ExpenseCategory = "food"
	CategoryRent          ExpenseCategory = "rent"
	CategoryUtilities     ExpenseCategory = "utilities"
	CategoryTransport     ExpenseCategory = "transport"
	CategoryEntertainment ExpenseCategory = "entertainment"
	CategoryHealth        ExpenseCategory = "health"
	CategoryShopping      ExpenseCategory = "shopping"
	CategoryOther         ExpenseCategory = "other"
)

func (c ExpenseCategory) IsValid() bool {
	switch c {
	case CategoryFood,
		CategoryRent,
		CategoryUtilities,
		CategoryTransport,
		CategoryEntertainment,
		CategoryHealth,
		CategoryShopping,
		CategoryOther:
		return true
	}
	return false
}

type Expense struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `json:"name"`
	Amount      float64         `json:"amount"`
	Category    ExpenseCategory `json:"category"`
	CreatedOn   time.Time       `json:"created_on"`
	CreatedByID uint            `json:"created_by_id"`
	CreatedBy   User            `gorm:"foreignKey:CreatedByID" json:"created_by"`
	HouseholdID uint            `json:"household_id"`
	Household   Household       `gorm:"foreignKey:HouseholdID" json:"household"`
}

type ExpenseShare struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ExpenseID uint    `json:"expense_id"`
	Expense   Expense `gorm:"foreignKey:ExpenseID" json:"expense"`
	UserID    uint    `json:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	Amount    float64 `json:"amount"`
	Paid      bool    `json:"paid"`
}

type Report struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
	TotalExpenses  float64   `json:"total_expenses"`
	PaymentStatus  string    `json:"payment_status"`
	GenerationDate time.Time `json:"generation_date"`
	FileName       string    `json:"file_name"`
}

type UserStore interface {
	CreateUser(username string, email string, password string) error
	GetUser(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetAllUsers() ([]User, error)
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
	DeleteSession(userID uint) error
}

type HouseholdStore interface {
	CreateHousehold(name string, description string, createdByID uint) (uint, error)
	GetHouseholdsByUserID(userID uint) ([]Household, error)
	GetOwnedHouseholdsByUserID(userID uint) ([]Household, error)
	NameExists(name string) (bool, error)
}

type MembershipStore interface {
	CreateMembership(userID uint, householdID uint, role string) error
	GetMembersByHouseholdID(householdID uint) ([]Membership, error)
}

type ExpenseStore interface {
	CreateExpense(name string, amount float64, category ExpenseCategory, createdOn time.Time, householdID, createdByID uint) (uint, error)
	NameExists(name string) (bool, error)
	GetExpensesByHouseholdID(householdID uint) ([]Expense, error)
}

type ExpenseShareStore interface {
	CreateExpenseShare(expenseID uint, userID uint, amount float64) error
	GetExpenseShare(expenseID uint, userID uint) (ExpenseShare, error)
	GetExpensesByUserID(userID uint) ([]ExpenseShare, error)
	UpdateExpenseShare(share ExpenseShare) error
}

type ReportStore interface {
	CreateReport(userID uint, from, to time.Time, paymentStatus string) (Report, error)
	GetReportsByUser(userID uint) ([]Report, error)
	GetReportByFileName(fileName string) (Report, error)
}
