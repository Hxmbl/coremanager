package logic

import (
	"fmt"
	"os/exec"
)

// EnableAll enables all secondary CPU cores (cpu1..cpuN).
func EnableAll(verbose bool) error {
	total, err := ListAllCPUCores()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Enabling all %d secondary cores...\n", total-1)
	}

	for i := 1; i < total; i++ {
		if verbose {
			fmt.Printf("  cpu%d -> enabling\n", i)
		}
		if err := setCoreOnline(i, true); err != nil {
			return fmt.Errorf("failed to enable cpu%d: %w", i, err)
		}
	}

	if verbose {
		fmt.Println("Done.")
	}
	return nil
}

// DisableAll disables all secondary CPU cores (cpu1..cpuN).
func DisableAll(verbose bool) error {
	total, err := ListAllCPUCores()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Disabling all %d secondary cores...\n", total-1)
	}

	for i := 1; i < total; i++ {
		if verbose {
			fmt.Printf("  cpu%d -> disabling\n", i)
		}
		if err := setCoreOnline(i, false); err != nil {
			return fmt.Errorf("failed to disable cpu%d: %w", i, err)
		}
	}

	if verbose {
		fmt.Println("Done.")
	}
	return nil
}

// Enable enables exactly `count` offline cores, starting from the lowest-numbered offline core.
func Enable(count int, verbose bool) error {
	if count < 1 {
		return fmt.Errorf("must enable at least 1 core")
	}

	total, err := ListAllCPUCores()
	if err != nil {
		return err
	}

	active, err := GetActiveCores()
	if err != nil {
		return err
	}

	states, err := GetCoreStates()
	if err != nil {
		return err
	}

	disabled := total - active

	if verbose {
		fmt.Printf("Requested: enable %d core(s)\n", count)
		fmt.Printf("Currently: %d active, %d disabled, %d total\n", active, disabled, total)
	}

	if count > disabled {
		return fmt.Errorf("cannot enable %d core(s). Only %d disabled core(s) available", count, disabled)
	}

	enabledCount := 0
	for i := 1; i < total; i++ {
		if !states[i] {
			if verbose {
				fmt.Printf("  cpu%d -> enabling\n", i)
			}
			if err := setCoreOnline(i, true); err != nil {
				return fmt.Errorf("failed to enable cpu%d: %w", i, err)
			}
			enabledCount++
			if enabledCount >= count {
				break
			}
		}
	}

	if verbose {
		fmt.Printf("Done. Enabled %d core(s). Active cores: %d\n", enabledCount, active+enabledCount)
	}
	return nil
}

// Disable disables exactly `count` online cores (excluding cpu0), starting from the lowest-numbered online core.
func Disable(count int, verbose bool) error {
	if count < 1 {
		return fmt.Errorf("must disable at least 1 core")
	}

	total, err := ListAllCPUCores()
	if err != nil {
		return err
	}

	active, err := GetActiveCores()
	if err != nil {
		return err
	}

	states, err := GetCoreStates()
	if err != nil {
		return err
	}

	canDisable := active - 1 // cpu0 always on

	if verbose {
		fmt.Printf("Requested: disable %d core(s)\n", count)
		fmt.Printf("Currently: %d active, %d total\n", active, total)
	}

	if count > canDisable {
		return fmt.Errorf("cannot disable %d core(s). Only %d core(s) can be disabled (cpu0 must stay on)", count, canDisable)
	}

	disabledCount := 0
	for i := 1; i < total; i++ {
		if states[i] {
			if verbose {
				fmt.Printf("  cpu%d -> disabling\n", i)
			}
			if err := setCoreOnline(i, false); err != nil {
				return fmt.Errorf("failed to disable cpu%d: %w", i, err)
			}
			disabledCount++
			if disabledCount >= count {
				break
			}
		}
	}

	if verbose {
		fmt.Printf("Done. Disabled %d core(s). Active cores: %d\n", disabledCount, active-disabledCount)
	}
	return nil
}

// setCoreOnline writes 1 or 0 to the sysfs online file for the given core.
func setCoreOnline(coreID int, online bool) error {
	value := "0"
	if online {
		value = "1"
	}
	path := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", coreID)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %s | sudo tee %s", value, path))
	return cmd.Run()
}
