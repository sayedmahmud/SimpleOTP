// Entries are the main data structure of the database. They are encrypted and stored
package database

import (
	"encoding/gob"
	"errors"
	"os"
	"strings"
)

type Entry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Secret      string `json:"secret"`
}

type Entries struct {
	Entries map[string]string `json:"entries"` // Key is the hashed name, value is the base64 encoded encrypted entry
}

func (e *Entries) Get(name string) (*Entry, error) {
	// Hash the name
	hashedName := hash(name)
	// Get the entry from the map
	encryptedEntry, ok := e.Entries[string(hashedName[:])]
	if !ok {
		return nil, errors.New("entry not found")
	}
	// Decrypt the entry
	decryptedEntry, err := Decrypt(encryptedEntry)
	if err != nil {
		return nil, err
	}
	return decryptedEntry, nil
}

func (e *Entries) Search(name string) ([]string, error) {
	names, err := e.List()
	if err != nil {
		return nil, err
	}
	var matches []string
	for _, n := range names {
		// Check if name is a substring of n
		if strings.Contains(strings.ToLower(n), strings.ToLower(name)) {
			matches = append(matches, n)
		}
	}
	return matches, nil
}

func (e *Entries) Add(entry Entry) {
	// Hash the name
	hashedName := hash(entry.Name)
	// Add the entry to the map
	e.Entries[string(hashedName[:])] = Encrypt(&entry)

}

func (e *Entries) Remove(name string) {
	// Hash the name
	hashedName := hash(name)
	// Remove the entry from the map
	delete(e.Entries, string(hashedName[:]))
}

func (e *Entries) List() ([]string, error) {
	names := make([]string, len(e.Entries))
	i := 0
	for _, entry := range e.Entries {
		decryptedEntry, err := Decrypt(entry)
		if err != nil {
			return nil, err
		}
		names[i] = decryptedEntry.Name
		i++
	}
	return names, nil

}

func (e *Entries) Save() error {
	// Gob encode the entries
	file, err := os.OpenFile("entries.gob", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entries) Load() error {
	file, err := os.Open("entries.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(e)
	if err != nil {
		return err
	}
	return nil
}
