package main

import (
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func Create() *Store {
	s := &Store{
		data: make(map[string]string),
	}
	return s
}

type Transaction struct {
	s        *Store
	writable bool
}

func (tx *Transaction) Set(key, value string) {
	tx.s.data[key] = value
}

func (tx *Transaction) Get(key string) string {
	return tx.s.data[key]
}

func (tx *Transaction) lock() {
	if tx.writable {
		tx.s.mu.Lock()
	} else {
		tx.s.mu.RLock()
	}
}

func (tx *Transaction) unlock() {
	if tx.writable {
		tx.s.mu.Unlock()
	} else {
		tx.s.mu.RUnlock()
	}
}

func (s *Store) Begin(writable bool) (*Transaction, error) {
	tx := &Transaction{
		s:        s,
		writable: writable,
	}
	tx.lock()

	return tx, nil
}

func (s *Store) managed(writable bool, fn func(tx *Transaction) error) (err error) {
	var tx *Transaction
	tx, err = s.Begin(writable)
	if err != nil {
		return
	}

	defer func() {
		tx.unlock()
	}()

	err = fn(tx)
	return
}

func (s *Store) View(fn func(tx *Transaction) error) error {
	return s.managed(false, fn)
}

func (s *Store) Update(fn func(tx *Transaction) error) error {
	return s.managed(true, fn)
}
