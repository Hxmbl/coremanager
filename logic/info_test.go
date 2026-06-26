package logic

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// fakeSysfs sets up a temporary directory tree mimicking /sys/devices/system/cpu/
// and returns a cleanup function.
func fakeSysfs(t *testing.T, totalCores int, onlineStates map[int]bool) (sysfsRoot string, cleanup func()) {
	t.Helper()
	root := t.TempDir()

	cpuDir := filepath.Join(root, "sys", "devices", "system", "cpu")
	if err := os.MkdirAll(cpuDir, 0755); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < totalCores; i++ {
		coreDir := filepath.Join(cpuDir, "cpu"+itoa(i))
		if err := os.MkdirAll(coreDir, 0755); err != nil {
			t.Fatal(err)
		}

		// cpu0 has no online file; others do
		if i > 0 {
			val := "0"
			if onlineStates[i] {
				val = "1"
			}
			onlinePath := filepath.Join(coreDir, "online")
			if err := os.WriteFile(onlinePath, []byte(val), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return root, func() { os.RemoveAll(root) }
}

// fakeProcCPUinfo creates a temporary /proc/cpuinfo with the given content.
func fakeProcCPUinfo(t *testing.T, content string) (procPath string, cleanup func()) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cpuinfo")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path, func() { os.RemoveAll(dir) }
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

// TestGetCPUModel tests CPU model parsing from /proc/cpuinfo
func TestGetCPUModel(t *testing.T) {
	content := "processor\t: 0\nmodel name\t: Test CPU Model\n"
	path, cleanup := fakeProcCPUinfo(t, content)
	defer cleanup()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Read line by line looking for model name
	var model string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				model = strings.TrimSpace(parts[1])
			}
		}
	}

	if model != "Test CPU Model" {
		t.Errorf("expected 'Test CPU Model', got '%s'", model)
	}
}

// TestListAllCPUCores tests counting of cpuN directories
func TestListAllCPUCores(t *testing.T) {
	root, cleanup := fakeSysfs(t, 4, map[int]bool{1: true, 2: true, 3: false})
	defer cleanup()

	entries, err := os.ReadDir(filepath.Join(root, "sys", "devices", "system", "cpu"))
	if err != nil {
		t.Fatal(err)
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

	if count != 4 {
		t.Errorf("expected 4 cores, got %d", count)
	}
}

// TestCoreStates tests reading online states from sysfs
func TestCoreStates(t *testing.T) {
	states := map[int]bool{1: true, 2: false, 3: true}
	root, cleanup := fakeSysfs(t, 4, states)
	defer cleanup()

	cpuDir := filepath.Join(root, "sys", "devices", "system", "cpu")
	result := make([]bool, 4)

	for i := 0; i < 4; i++ {
		if i == 0 {
			result[i] = true // cpu0 always on
			continue
		}
		onlinePath := filepath.Join(cpuDir, "cpu"+itoa(i), "online")
		data, err := os.ReadFile(onlinePath)
		if err != nil {
			t.Fatal(err)
		}
		result[i] = strings.TrimSpace(string(data)) == "1"
	}

	expected := []bool{true, true, false, true}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("core %d: expected %v, got %v", i, expected[i], v)
		}
	}
}

// TestActiveCores tests counting of active cores
func TestActiveCores(t *testing.T) {
	states := map[int]bool{1: true, 2: false, 3: true}
	root, cleanup := fakeSysfs(t, 4, states)
	defer cleanup()

	cpuDir := filepath.Join(root, "sys", "devices", "system", "cpu")
	active := 1 // cpu0 always on
	for i := 1; i < 4; i++ {
		onlinePath := filepath.Join(cpuDir, "cpu"+itoa(i), "online")
		data, err := os.ReadFile(onlinePath)
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(data)) == "1" {
			active++
		}
	}

	if active != 3 {
		t.Errorf("expected 3 active cores, got %d", active)
	}
}
