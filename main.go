package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	col_orange string = "\033[93m"
	col_end           = "\033[0m"
)

/*
References:
 - https://www.raspberrypi.com/documentation/computers/os.html#get_throttled
 - https://www.raspberrypi.com/documentation/computers/config_txt.html#monitoring-core-temperature
*/

func main() {
	outputStr, err := runVcgencmd()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}

	value, err := parseThrottledStatus(outputStr)
	if err != nil {
		fmt.Printf("Unexpected output, unable to parse: %s\n", outputStr)
		return
	}
	printHumandReadableStatus(value)
	fmt.Println(wrap_warning(fmt.Sprintf("Throttled status [%s]\n", outputStr)))
}

func wrap_warning(str string) string {
	return fmt.Sprintf("%s%s%s", col_orange, str, col_end)
}

func icon(b bool) string {
	if b {
		return "⚠️"
	} else {
		return "✔️"
	}
}

func printHumandReadableStatus(value int) {
	// See https://www.raspberrypi.com/documentation/computers/os.html#get_throttled for details
	uv_now := value&0x1 != 0
	uv_past := value&0x10000 != 0
	fcap_now := value&0x2 != 0
	fcap_past := value&0x20000 != 0
	throt_now := value&0x4 != 0
	throt_past := value&0x40000 != 0
	tlim_now := value&0x8 != 0
	tlim_past := value&0x80000 != 0

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Category", "CURR", "PAST"})
	t.AppendRow(table.Row{"Under-voltage", icon(uv_now), icon(uv_past)})
	t.AppendRow(table.Row{"Arm freq capped", icon(fcap_now), icon(fcap_past)})
	t.AppendRow(table.Row{"Throttled", icon(throt_now), icon(throt_past)})
	t.AppendRow(table.Row{"Soft temp limit", icon(tlim_now), icon(tlim_past)})
	t.SetStyle(table.StyleLight)
	t.Render()
}

func runVcgencmd() (string, error) {
	cmdName := "vcgencmd"
	cmdArgs := []string{"get_throttled"}
	cmd := exec.Command(cmdName, cmdArgs...)

	// Run the command and capture its output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	// Convert the byte slice to a string and remove leading/trailing whitespace
	return strings.TrimSpace(string(output)), nil
}

func parseThrottledStatus(input string) (int, error) {
	if !strings.HasPrefix(input, "throttled=0x") {
		return 0, errors.New("input does not start with 'throttled=0x'")
	}
	// Remove the prefix
	hexStr := strings.TrimPrefix(input, "throttled=0x")

	// Parse the remaining string as hexadecimal
	throttledInt, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}
	return int(throttledInt), nil
}
