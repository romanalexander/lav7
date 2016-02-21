package try

// Try runs tryFunc, catches panic, and executes panicFunc with recovered panic.
func Try(tryFunc func(), panicFunc func(interface{})) {
	defer func() {
		if r := recover(); r != nil {
			panicFunc(r)
		}
	}()
	return
}
