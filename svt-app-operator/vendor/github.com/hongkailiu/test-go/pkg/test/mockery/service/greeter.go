package service

type db struct{}

// DB is fake Database interface.
type DB interface {
	FetchMessage(lang string) (string, error)
	FetchDefaultMessage() (string, error)
}

// Greeter is the targeting struct for test
type Greeter struct {
	Database DB
	Lang     string
}

// GreeterService is service to greet your friends.
type GreeterService interface {
	Greet() string
	GreetInDefaultMsg() string
}

func (d *db) FetchMessage(lang string) (string, error) {
	// in real life, this code will call an external db
	// but for this sample we will just return the hardcoded example value
	if lang == "en" {
		return "hello", nil
	}
	if lang == "es" {
		return "holla", nil
	}
	return "bzzzz", nil
}

func (d *db) FetchDefaultMessage() (string, error) {
	return "default message", nil
}

// Greet returns a greeting string
func (g Greeter) Greet() string {
	msg, _ := g.Database.FetchMessage(g.Lang) // call Database to get the message based on the Lang
	return "Message is: " + msg
}

// GreetInDefaultMsg returns a default greeting string
func (g Greeter) GreetInDefaultMsg() string {
	msg, _ := g.Database.FetchDefaultMessage() // call Database to get the default message
	return "Message is: " + msg
}

// NewDB generates a new db
func NewDB() DB {
	return new(db)
}

// NewGreeter generates a new greeter
func NewGreeter(db DB, lang string) GreeterService {
	return Greeter{db, lang}
}
