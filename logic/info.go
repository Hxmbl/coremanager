package logic

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CPUInfo holds information about the system's CPU.
type CPUInfo struct {
	Model        string
	TotalCores   int
	ActiveCores  int
	CoreStates   []bool
	PhysicalCores int
	LogicalCores  int
}

// GetCPUModel reads /proc/cpuinfo to get the CPU model name.
func GetCPUModel() (string, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("failed to open /proc/cpuinfo: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "Unknown", nil
}

// ListAllCPUCores counts the number of CPU cores by checking /sys/devices/system/cpu/ for cpuN directories.
func ListAllCPUCores() (int, error) {
	entries, err := os.ReadDir("/sys/devices/system/cpu/")
	if err != nil {
		return 0, fmt.Errorf("failed to read /sys/devices/system/cpu/: %w", err)
	}

	count := 0
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "cpu") {
			suffix := name[3:]
			if _, err := strconv.Atoi(suffix); err == nil {
				count++
			}
		}
	}
	return count, nil
}

// GetActiveCores counts how many CPU cores are currently online via sysfs.
func GetActiveCores() (int, error) {
	total, err := ListAllCPUCores()
	if err != nil {
		return 0, err
	}

	active := 1 // cpu0 is always on
	for i := 1; i < total; i++ {
		onlinePath := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", i)
		state, err := readOnlineState(onlinePath)
		if err != nil {
			continue
		}
		if state {
			active++
		}
	}
	return active, nil
}

// GetCoreStates returns a slice of booleans indicating whether each core is online.
func GetCoreStates() ([]bool, error) {
	total, err := ListAllCPUCores()
	if err != nil {
		return nil, err
	}

	states := make([]bool, total)
	for i := 0; i < total; i++ {
		onlinePath := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", i)
		if _, err := os.Stat(onlinePath); os.IsNotExist(err) {
			// cpu0 has no online file, always on
			states[i] = true
		} else {
			state, err := readOnlineState(onlinePath)
			if err != nil {
				states[i] = true
				continue
			}
			states[i] = state
		}
	}
	return states, nil
}

// AllCPU returns the raw content of /proc/cpuinfo for debugging.
func AllCPU() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("failed to read /proc/cpuinfo: %w", err)
	}
	return string(data), nil
}

// GetInfo returns a fully populated CPUInfo struct.
func GetInfo() (*CPUInfo, error) {
	model, err := GetCPUModel()
	if err != nil {
		return nil, err
	}

	total, err := ListAllCPUCores()
	if err != nil {
		return nil, err
	}

	active, err := GetActiveCores()
	if err != nil {
		return nil, err
	}

	states, err := GetCoreStates()
	if err != nil {
		return nil, err
	}

	return &CPUInfo{
		Model:        model,
		TotalCores:   total,
		ActiveCores:  active,
		CoreStates:   states,
		PhysicalCores: 0, // psutil equivalent not easily portable; 0 if unknown
		LogicalCores:  active,
	}, nil
}

// readOnlineState reads a sysfs online file and returns true if the core is online (value "1").
func readOnlineState(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(data)) == "1", nil
}
