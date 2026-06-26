package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"coremanager/logic"

	"github.com/spf13/cobra"
)

var verbose bool

func main() {
	rootCmd := &cobra.Command{
		Use:   "coremanager",
		Short: "A simple CLI tool to manage CPU cores on Linux",
		Long:  "CoreManager - Manage CPU cores dynamically to save battery or reduce heat.",
	}

	disableCmd := &cobra.Command{
		Use:     "disable-cores [N|all]",
		Aliases: []string{"dc"},
		Short:   "Disable a number of CPU cores",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDisableCores(args[0])
		},
	}
	disableCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	enableCmd := &cobra.Command{
		Use:     "enable-cores [N|all]",
		Aliases: []string{"ec"},
		Short:   "Enable a number of CPU cores",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnableCores(args[0])
		},
	}
	enableCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	coreCountCmd := &cobra.Command{
		Use:     "core-count",
		Aliases: []string{"cc"},
		Short:   "Display the total and active CPU core counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCoreCount()
		},
	}

	cpuModelCmd := &cobra.Command{
		Use:     "cpu-model",
		Aliases: []string{"cm"},
		Short:   "Display the CPU model name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCPUModel()
		},
	}

	debugCmd := &cobra.Command{
		Use:     "debug-info",
		Aliases: []string{"debug"},
		Short:   "Display detailed CPU information for debugging",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebugInfo()
		},
	}

	rootCmd.AddCommand(disableCmd, enableCmd, coreCountCmd, cpuModelCmd, debugCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDisableCores(arg string) error {
	if strings.ToLower(arg) == "all" || strings.ToLower(arg) == "a" {
		if verbose {
			fmt.Println("Disabling all secondary cores...")
		}
		if err := logic.DisableAll(verbose); err != nil {
			return err
		}
		fmt.Println("All secondary cores disabled.")
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
	fmt.Printf("Disabled %d core(s).\n", target)
	return nil
}

func runEnableCores(arg string) error {
	if strings.ToLower(arg) == "all" || strings.ToLower(arg) == "a" {
		if verbose {
			fmt.Println("Enabling all secondary cores...")
		}
		if err := logic.EnableAll(verbose); err != nil {
			return err
		}
		fmt.Println("All secondary cores enabled.")
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
	fmt.Printf("Enabled %d core(s).\n", target)
	return nil
}

func runCoreCount() error {
	info, err := logic.GetInfo()
	if err != nil {
		return err
	}
	fmt.Printf("Total CPU cores: %d\n", info.TotalCores)
	fmt.Printf("Active CPU cores: %d\n", info.ActiveCores)
	return nil
}

func runCPUModel() error {
	model, err := logic.GetCPUModel()
	if err != nil {
		return err
	}
	fmt.Printf("CPU Model: %s\n", model)
	return nil
}

func runDebugInfo() error {
	info, err := logic.GetInfo()
	if err != nil {
		return err
	}

	fmt.Printf("CPU Model: %s\n", info.Model)
	fmt.Printf("Total CPU cores: %d\n", info.TotalCores)
	fmt.Printf("Active CPU cores: %d\n", info.ActiveCores)
	fmt.Printf("Core states (0=offline, 1=online): %v\n", info.CoreStates)

	fmt.Print("Type 'a' for a lot of info. Press Enter to continue... ")
	var input string
	fmt.Scanln(&input)
	if strings.ToLower(input) == "a" {
		fmt.Println("Detailed info:")
		raw, err := logic.AllCPU()
		if err != nil {
			return err
		}
		fmt.Println(raw)
	}
	return nil
}
