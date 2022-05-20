//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]Session
	mu       sync.Mutex
}

// Session stores the session's data
type Session struct {
	Data      map[string]interface{}
	CreatedAt int64 // Unix Epoch time in seconds
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
	}

	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	m.mu.Lock()
	m.sessions[sessionID] = Session{
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now().Unix(),
	}
	m.mu.Unlock()
	fmt.Printf("New session is created  at %d with ID %s\n", m.sessions[sessionID].CreatedAt, sessionID)
	return sessionID, nil
}

// SessionRoutine checks all sessions in the manager for every seconds and kicks out expired ones
func (m *SessionManager) SessionRoutine() {
	// iterate over the sessions indefinitely and check for expired sessions
	for range time.Tick(5 * time.Second) {
		fmt.Printf("Current EPOCH is %d\n", time.Now().Unix())
		for sessionID, session := range m.sessions {
			fmt.Printf("Time difference between current session %s wiht current time is %d \n", sessionID, time.Now().Unix()-session.CreatedAt)
			if (time.Now().Unix() - session.CreatedAt) > 4 {
				fmt.Printf("Session %s is expired it will be deleted \n", sessionID)
				m.mu.Lock()
				delete(m.sessions, sessionID)
				m.mu.Unlock()
				_, ok := m.sessions[sessionID]
				if !ok {
					fmt.Printf("Session %s has removed\n", sessionID)
				}
			}
		}
	}
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data:      data,
		CreatedAt: time.Now().Unix(),
	}

	return nil
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()

	go func(m *SessionManager) {
		m.SessionRoutine()
	}(m)

	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
	time.Sleep(100 * time.Second)
}
