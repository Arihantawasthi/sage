package utils

import "fmt"

func Bold(w int, s string) string {
    padded := fmt.Sprintf("| %-*s ", w, s)
    return "\033[1m" + padded + "\033[0m"
}

func CyanBold(w int, s string) string {
    padded := fmt.Sprintf("%-*s ", w, s)
    return "\033[1;96m" + padded + "\033[0m"
}
