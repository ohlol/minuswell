package main

type Output interface {
	Emit(event Event)
}
