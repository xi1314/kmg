package kmgConsole

// @deprecated
func AddAction(action Command) {
	DefaultCommandGroup.AddCommand(action)
}
