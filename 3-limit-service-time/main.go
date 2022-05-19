//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool { // make the rate-limiting based on the user (accumulated)
	if !u.IsPremium {
		if u.TimeUsed >= 10 {
			fmt.Printf("user %d has rate-limited\n", u.ID)
			return false
		} else {
			// start a concurrent timer if user got any time left
			// when process done, time automatically stops
			go func(u *User) {
				for range time.Tick(1 * time.Second) {
					if u.TimeUsed >= 10 {
						return
					}
					u.TimeUsed++
				}
			}(u)
		}
	}
	process()
	return true
}

func main() {
	RunMockServer()
}
