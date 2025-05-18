package utils

import (
	"fmt"
	"strconv"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

func PrintTable(data []models.PListData) {
    headers := []string{"SNo.", "PID", "P_NAME", "NAME", "CMD", "UP TIME", "CPU%", "MEM%"}
    widths := make([]int, len(headers))

    for i, h := range headers {
        widths[i] = len(h)
    }

    padding := 6
    for i, d := range data {
        widths[0] = max(widths[0], len(fmt.Sprintf("%d", i+1)) + padding)
        widths[1] = max(widths[1], len(fmt.Sprintf("%d", d.Pid)) + padding)
        widths[2] = max(widths[2], len(d.PName) + padding)
        widths[3] = max(widths[3], len(d.Name) + padding)
        widths[4] = max(widths[4], len(d.Cmd) + padding)
        widths[5] = max(widths[5], len(d.UpTime) + padding)
        widths[6] = max(widths[6], len(fmt.Sprintf("%0.02f", d.CPUPercent)) + padding)
        widths[7] = max(widths[7], len(fmt.Sprintf("%0.02f", d.MemPrecent)) + padding)
    }
    printBorders(widths, headers)

    for i, h := range headers {
        fmt.Printf("| %s", CyanBold(widths[i], h))
    }
    fmt.Println()

    printBorders(widths, headers)
    for i, d := range data {
        fmt.Printf("| %s", CyanBold(widths[0], strconv.Itoa(i + 1)))
        fmt.Printf("| %-*d ", widths[1], d.Pid)
        fmt.Printf("| %-*s ", widths[2], d.PName)
        fmt.Printf("| %-*s ", widths[3], d.Name)
        fmt.Printf("| %-*s ", widths[4], d.Cmd)
        fmt.Printf("| %-*s ", widths[5], d.UpTime)
        fmt.Printf("| %-*.2f ", widths[6], d.CPUPercent)
        fmt.Printf("| %-*.2f ", widths[7], d.MemPrecent)
        fmt.Println()
    }
    printBorders(widths, headers)
}

func printBorders(widths []int, headers []string) {
    for i, _ := range headers {
        for w := 0; w < widths[i]+3; w++ {
            fmt.Print("-")
        }
    }
    fmt.Println()
}
