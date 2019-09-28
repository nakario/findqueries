package b

import "github.com/nakario/findqueries/test/c"

func F00() {}

type G func()

func (g G) F00() {}

func C() c.C { return c.C(F00) }
