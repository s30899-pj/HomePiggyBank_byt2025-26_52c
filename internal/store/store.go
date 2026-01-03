package store

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Expenses struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Amount   float32   `json:"amount"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	//UserID    uint   `json:"user_id"`
	//HouseholdsID    uint   `json:"households_id"`
}

type Households struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	HouseholdsName string    `json:"households_name"`
	Date           time.Time `json:"date"`
}

type Membership struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//UserID    uint   `json:"user_id"`
	//HouseholdsID    uint   `json:"households_id"`
}

type Report struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PeriodOfDates  time.Time `json:"period_of_dates"`
	TotalExpenses  float32   `json:"total_expenses"`
	GenerationDate time.Time `json:"GenerationDate"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type UserStore interface {
	CreateUser(username string, email string, password string) error
	GetUser(email string) (*User, error)
}

type ExpensesStore interface {
	// CreateExpense TODO: check
	CreateExpense(amount float32, category string, date time.Time) error
	GetExpense(amount float32, category string, date time.Time) (*Expenses, error)
}

type HouseholdsStore interface {
	CreateHousehold(householdsName string, date time.Time) error
	GetHousehold(householdsName string, date time.Time) (*Households, error)
}

type MembershipStore interface {
	// CreateMembership TODO: check
	CreateMembership(user User, households Households) error
	GetMembership(user User, households Households) (*Membership, error)
}

type ReportStore interface {
	CreateReport(periodOfDates time.Time, totalExpenses float32, generationDate time.Time) error
	GetReport(periodOfDates time.Time, totalExpenses float32, generationDate time.Time) (*Report, error)
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
}
