package kmgConsole

// @deprecated
// please use AddCommand instead.
func AddAction(action Command) {
	DefaultCommandGroup.AddCommand(action)
}
