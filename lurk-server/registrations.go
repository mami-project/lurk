package main

import "fmt"

var lastId int

type LurkRegistration struct {
	Id       int
	CSR      string
	Lifetime uint
	// TODO(tho) other stuff related to its state on the ACME side
}

type LurkRegistrations []LurkRegistration

var regs LurkRegistrations

func init() {
	// TODO(tho)
}

func LookupLurkRegistration(id int) LurkRegistration {
	for _, r := range regs {
		if r.Id == id {
			return r
		}
	}
	// return empty Req if not found
	return LurkRegistration{}
}

// XXX(tho) race
// TODO(tho) uuid + map
func CreateLurkRegistration(r LurkRegistration) error {
	lastId += 1
	r.Id = lastId
	regs = append(regs, r)
	return nil
}

func DeleteLurkRegistration(id int) error {
	for i, r := range regs {
		if r.Id == id {
			regs = append(regs[:i], regs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Could not find LURK registration with id of %d to delete", id)
}
