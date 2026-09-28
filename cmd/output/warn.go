package output

import "charm.land/lipgloss/v2"

var warnStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Yellow)

func PrintWarn(message string) error {
	_, err := lipgloss.Println(warnStyle.Render(message))
	return err
}
