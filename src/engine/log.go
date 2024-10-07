package engine

// Special error that is not fatal when paniced on and called from GUI
// E.g.: try to set bad grid size: panic on UserErr, recover, print error, carry on.
type UserErr string

func (e UserErr) Error() string { return string(e) }

func checkRecoverable(err error) {
	if err != nil {
		panic(UserErr(err.Error()))
	}
}
