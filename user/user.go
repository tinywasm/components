package user

// Get is a global function that returns the current user name and area.
// It can be overridden by other packages (e.g., site) to provide real data.
var Get = func() (string, string) {
	return "", ""
}
