package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hxmbl/coremanager/logic"

	"github.com/spf13/cobra"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

func green(s string) string  { return colorGreen + s + colorReset }
func red(s string) string    { return colorRed + s + colorReset }
func yellow(s string) string { return colorYellow + s + colorReset }
func blue(s string) string   { return colorBlue + s + colorReset }

var verbose bool

func main() {
	rootCmd := &cobra.Command{
		Use:   "coremanager",
		Short: "A simple CLI tool to manage CPU cores on Linux",
		Long:  "CoreManager - Manage CPU cores dynamically to save battery or reduce heat.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	disableCmd := &cobra.Command{
		Use:     "dc [N|all]",
		Aliases: []string{"disable-cores"},
		Short:   "Disable a number of CPU cores",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisableCores(args[0])
		},
	}
	disableCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	enableCmd := &cobra.Command{
		Use:     "ec [N|all]",
		Aliases: []string{"enable-cores"},
		Short:   "Enable a number of CPU cores",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnableCores(args[0])
		},
	}
	enableCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	coreCountCmd := &cobra.Command{
		Use:     "cc",
		Aliases: []string{"core-count"},
		Short:   "Display the total and active CPU core counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCoreCount()
		},
	}

	cpuModelCmd := &cobra.Command{
		Use:     "cm",
		Aliases: []string{"cpu-model"},
		Short:   "Display the CPU model name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCPUModel()
		},
	}

	debugCmd := &cobra.Command{
		Use:     "debug",
		Aliases: []string{"debug-info"},
		Short:   "Display detailed CPU information for debugging",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugInfo()
		},
	}

	rootCmd.AddCommand(disableCmd, enableCmd, coreCountCmd, cpuModelCmd, debugCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, red(fmt.Sprintf("Error: %v", err)))
		os.Exit(1)
	}
}

func runDisableCores(arg string) error {
	if strings.ToLower(arg) == "all" || strings.ToLower(arg) == "a" {
		if verbose {
			fmt.Println(yellow("Disabling all secondary cores..."))
		}
		if err := logic.DisableAll(verbose); err != nil {
			return err
		}
		fmt.Println(green("All secondary cores disabled."))
		return nil
	}

	target, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number or 'all'", arg)
	}
	if target < 1 {
		return fmt.Errorf("must disable at least 1 core")
	}

	if err := logic.Disable(target, verbose); err != nil {
		return err
	}
	fmt.Println(green(fmt.Sprintf("Disabled %d core(s).", target)))
	return nil
}

func runEnableCores(arg string) error {
	if strings.ToLower(arg) == "all" || strings.ToLower(arg) == "a" {
		if verbose {
			fmt.Println(yellow("Enabling all secondary cores..."))
		}
		if err := logic.EnableAll(verbose); err != nil {
			return err
		}
		fmt.Println(green("All secondary cores enabled."))
		return nil
	}

	target, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number or 'all'", arg)
	}
	if target < 1 {
		return fmt.Errorf("must enable at least 1 core")
	}

	if err := logic.Enable(target, verbose); err != nil {
		return err
	}
	fmt.Println(green(fmt.Sprintf("Enabled %d core(s).", target)))
	return nil
}

func runCoreCount() error {
	info, err := logic.GetInfo()
	if err != nil {
		return err
	}
	fmt.Println(blue(fmt.Sprintf("Total CPU cores: %d", info.TotalCores)))
	fmt.Println(blue(fmt.Sprintf("Active CPU cores: %d", info.ActiveCores)))
	return nil
}

func runCPUModel() error {
	model, err := logic.GetCPUModel()
	if err != nil {
		return err
	}
	fmt.Println(blue(fmt.Sprintf("CPU Model: %s", model)))
	return nil
}

func runDebugInfo() error {
	info, err := logic.GetInfo()
	if err != nil {
		return err
	}

	fmt.Println(blue(fmt.Sprintf("CPU Model: %s", info.Model)))
	fmt.Println(blue(fmt.Sprintf("Total CPU cores: %d", info.TotalCores)))
	fmt.Println(blue(fmt.Sprintf("Active CPU cores: %d", info.ActiveCores)))
	fmt.Println(blue(fmt.Sprintf("Core states (0=offline, 1=online): %v", info.CoreStates)))

	fmt.Print("Type 'a' for a lot of info. Press Enter to continue... ")
	var input string
	fmt.Scanln(&input)
	if strings.ToLower(input) == "a" {
		fmt.Println(blue("Detailed info:"))
		raw, err := logic.AllCPU()
		if err != nil {
			return err
		}
		fmt.Println(raw)
	}
	return nil
}
