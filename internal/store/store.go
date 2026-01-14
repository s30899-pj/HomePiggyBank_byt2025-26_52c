package store

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
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

type Expense struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Amount   float32   `json:"amount"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
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
	GetUserByUsername(username string) (*User, error)
	GetAllUsers() ([]User, error)
	EmailExists(email string) (bool, error)
	UsernameExists(username string) (bool, error)
}

type HouseholdStore interface {
	CreateHousehold(name string, description string, createdByID uint) (uint, error)
	GetHouseholdsByUserID(userID uint) ([]Household, error)
	NameExists(name string) (bool, error)
}

type MembershipStore interface {
	CreateMembership(userID uint, householdID uint, role string) error
	//GetByUserAndHousehold(userID, householdID uint) (*Membership, error)
}

type ExpensesStore interface {
	// CreateExpense TODO: check
	CreateExpense(amount float32, category string, date time.Time) error
	GetExpense(amount float32, category string, date time.Time) (*Expense, error)
}

type ReportStore interface {
	CreateReport(periodOfDates time.Time, totalExpenses float32, generationDate time.Time) error
	GetReport(periodOfDates time.Time, totalExpenses float32, generationDate time.Time) (*Report, error)
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
	DeleteSession(userID uint) error
}
