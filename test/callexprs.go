package a

import "github.com/nakario/findqueries/test/b"

func F00() {}

func F10(i int) {}

func F02() (int, int) { return 1, 2 }

func F20(i, j int) {}

func F01F() func() { return func(){} }

func F1F0(f func()) { f() }

type G func()

type H = G

func (g G) F00() {}

type I interface {
	F00()
}

func calls() {
	F00()
	(F00)()
	((F00))()
	f := F00
	f()
	(f)()
	F10(1)
	F20(F02())
	F01F()()
	f = F01F()
	f()
	F1F0(F01F())
	g := G(F00)
	g()
	(g)()
	g.F00()
	(g.F00)()
	G(g.F00)()
	G(F00)()
	(G(F00))()
	H(F00)()
	func(){}()
	(func(){})()
	G(func(){})()
	f = func(){
		F00()
	}
	I(G(F00)).F00()
	interface{F00()}(G(F00)).F00()
	I(G(F00)).(G).F00()
	b.F00()
	(b.F00)()
	b.G(F00)()
	f = b.G(F00)
	f()
	bg := I(b.G(F00))
	bg.F00()
	b.C().F00()
	I(b.C()).F00()
}
