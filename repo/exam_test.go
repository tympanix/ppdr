package repo

import (
	"testing"
)

func TestExamIntroduction(t *testing.T) {

	usrAlice := NewIdentity("Alice")
	usrBob := NewIdentity("Bob")
	usrCharlie := NewIdentity("Charlie")
	usrMallory := NewIdentity("Mallory")

	dr := NewRepo()

	r1 := NewState()
	r2 := NewState()
	r3 := NewState()
	r4 := NewState()
	r5 := NewState()

	r3.addDependency(r1)
	r4.addDependency(r2)
	r5.addDependency(r3)
	r5.addDependency(r4)

	// Setup of data repository
	dr.SetCurrentUser(usrMallory)
	dr.Put(r1)
	dr.SetCurrentUser(usrBob)
	dr.Put(r2)
	dr.SetCurrentUser(usrAlice)
	dr.Put(r3)
	dr.Put(r4)
	dr.SetCurrentUser(usrCharlie)
	dr.Put(r5)

	p1 := "author != user('Mallory') U author = user('Alice')"
	p2 := "author = user('Alice') & [] author != user('Mallory')"
	p3 := "[] author != user('Mallory') | author != user('Mallory') U author = user('Alice')"

	// Query all data resources with all integrity policies
	for r, name := range map[*State]string{
		r1: "r1",
		r2: "r2",
		r3: "r3",
		r4: "r4",
		r5: "r5",
	} {
		for phi, phiName := range map[string]string{
			p1: "p1",
			p2: "p2",
			p3: "p3",
		} {
			_, err := dr.QueryString(r, phi)
			if err == nil {
				t.Logf("DR_{c, %v} |= %s: %v\n", name, phiName, err == nil)
			}
		}
	}
}