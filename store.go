package main

import "sync"

type StoredContent struct {
	sync.RWMutex
	content []byte
}
