package configs

func TranslateCommand(cmd string) string {
	switch cmd {
	case "lift_on":
		return "ATO=1"
	case "lift_off":
		return "ATO=2"
	default:
		return cmd
	}
}
