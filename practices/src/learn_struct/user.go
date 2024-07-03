package learnstruct

import (
	"fmt"
	"time"
)

type user struct {
	firstName string
	lastName  string
	age       int
	createdAt time.Time
}

func GetUsers() {
	users := []user{
        {firstName: "Jay", lastName: "Kumar", age: 25, createdAt: time.Now()},
        {firstName: "Gowtham", lastName: "Sankar", age: 25, createdAt: time.Now().AddDate(-2, 0, 0)},
    }

    for _, user := range users {
        fmt.Printf("User: %s %s, Age: %d, Created At: %s\n", user.firstName, user.lastName, user.age, user.createdAt.Format("2006-01-02"))
    }
}