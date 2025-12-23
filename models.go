package main

import time

type Todo struct {
	ID int
	Title string
	Done bool
	Priority string
	Category string
	CreatedAt time.Time
	DueDate *time.Time
}
