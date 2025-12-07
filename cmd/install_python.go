package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const pythonVersion = "3.14"

var (
	forcePythonInstall bool
	skipVerify         bool
)

// pythonCmd represents the python subcommand
var pythonCmd = &cobra.Command{
	Use:   "python",
	Short: "Install Python 3.14 for netcheck PY check type",
	Long: `Install Python 3.14 to enable netcheck's Python script functionality.

This command will attempt to install Python 3.14 using the appropriate
package manager for your operating system:
  - Windows: winget, chocolatey, or direct download
  - macOS: Homebrew
  - Linux: System package manager (apt, dnf, yum, or zypper)

The command will first check if Python 3.14 is already installed and
skip installation unless --force is specified.`,
	RunE: installPython,
}

func init() {
	installCmd.AddCommand(pythonCmd)
	pythonCmd.Flags().BoolVar(&forcePythonInstall, "force", false, "force installation even if Python is already installed")
	pythonCmd.Flags().BoolVar(&skipVerify, "skip-verify", false, "skip verification after installation")
}

func installPython(cmd *cobra.Command, args []string) error {
	fmt.Println("Python 3.14 Installation for netcheck")
	fmt.Println("======================================")
	fmt.Println()

	// Check if Python is already installed
	if !forcePythonInstall {
		if version, installed := checkPythonInstalled(); installed {
			fmt.Printf("✓ Python is already installed: %s\n", version)
			fmt.Println()
			fmt.Println("Use --force to reinstall")
			return nil
		}
	}

	fmt.Printf("Installing Python %s for %s/%s...\n\n", pythonVersion, runtime.GOOS, runtime.GOARCH)

	var err error
	switch runtime.GOOS {
	case "windows":
		err = installPythonWindows()
	case "darwin":
		err = installPythonMacOS()
	case "linux":
		err = installPythonLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Verify installation unless skipped
	if !skipVerify {
		fmt.Println()
		fmt.Println("Verifying installation...")
		if version, installed := checkPythonInstalled(); installed {
			fmt.Printf("✓ Python successfully installed: %s\n", version)
		} else {
			fmt.Println("⚠ Warning: Python installation completed but verification failed")
			fmt.Println("  You may need to restart your terminal or add Python to your PATH")
		}
	}

	return nil
}

func checkPythonInstalled() (string, bool) {
	// Try python3 first
	cmd := exec.Command("python3", "--version")
	output, err := cmd.CombinedOutput()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return version, true
	}

	// Fall back to python
	cmd = exec.Command("python", "--version")
	output, err = cmd.CombinedOutput()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return version, true
	}

	return "", false
}

func installPythonWindows() error {
	fmt.Println("Attempting Windows installation methods...")
	fmt.Println()

	// Try winget first (Windows 10/11)
	if _, err := exec.LookPath("winget"); err == nil {
		fmt.Println("→ Using winget (Windows Package Manager)")
		cmd := exec.Command("winget", "install", "-e", "--id", "Python.Python.3.14")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via winget")
			return nil
		}
		fmt.Println("⚠ winget installation failed, trying next method...")
		fmt.Println()
	}

	// Try chocolatey
	if _, err := exec.LookPath("choco"); err == nil {
		fmt.Println("→ Using Chocolatey")
		cmd := exec.Command("choco", "install", "python", "--version=3.14.0", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via Chocolatey")
			return nil
		}
		fmt.Println("⚠ Chocolatey installation failed, trying next method...")
		fmt.Println()
	}

	// Manual installation instructions
	fmt.Println("→ Automated installation not available")
	fmt.Println()
	fmt.Println("Please install Python 3.14 manually:")
	fmt.Println("1. Visit: https://www.python.org/downloads/")
	fmt.Println("2. Download Python 3.14 installer for Windows")
	fmt.Println("3. Run the installer")
	fmt.Println("4. IMPORTANT: Check 'Add Python to PATH' during installation")
	fmt.Println()

	return fmt.Errorf("automatic installation not available - please install manually")
}

func installPythonMacOS() error {
	fmt.Println("Attempting macOS installation methods...")
	fmt.Println()

	// Check for Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("→ Using Homebrew")
		fmt.Println("Running: brew install python@3.14")
		cmd := exec.Command("brew", "install", "python@3.14")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via Homebrew")

			// Link python3.14 to python3
			fmt.Println()
			fmt.Println("Setting up Python 3.14 as default python3...")
			linkCmd := exec.Command("brew", "link", "python@3.14")
			linkCmd.Stdout = os.Stdout
			linkCmd.Stderr = os.Stderr
			_ = linkCmd.Run() // Don't fail if linking fails

			return nil
		}
		fmt.Println("⚠ Homebrew installation failed")
		fmt.Println()
	} else {
		fmt.Println("⚠ Homebrew not found")
		fmt.Println()
		fmt.Println("To install Homebrew:")
		fmt.Println("/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
		fmt.Println()
	}

	// Manual installation instructions
	fmt.Println("Alternative manual installation:")
	fmt.Println("1. Visit: https://www.python.org/downloads/")
	fmt.Println("2. Download Python 3.14 installer for macOS")
	fmt.Println("3. Run the installer package")
	fmt.Println()

	return fmt.Errorf("automatic installation failed - please install Homebrew or install manually")
}

func installPythonLinux() error {
	fmt.Println("Attempting Linux installation methods...")
	fmt.Println()

	// Detect package manager
	var pkgManager string
	var installCmd *exec.Cmd

	if _, err := exec.LookPath("apt"); err == nil {
		pkgManager = "apt (Debian/Ubuntu)"
		fmt.Printf("→ Using %s\n", pkgManager)
		fmt.Println()
		fmt.Println("Note: Python 3.14 may not be available in default repositories")
		fmt.Println("Attempting to install python3...")

		// Update package list first
		updateCmd := exec.Command("sudo", "apt", "update")
		updateCmd.Stdout = os.Stdout
		updateCmd.Stderr = os.Stderr
		_ = updateCmd.Run()

		installCmd = exec.Command("sudo", "apt", "install", "-y", "python3", "python3-pip")
	} else if _, err := exec.LookPath("dnf"); err == nil {
		pkgManager = "dnf (Fedora/RHEL)"
		fmt.Printf("→ Using %s\n", pkgManager)
		installCmd = exec.Command("sudo", "dnf", "install", "-y", "python3", "python3-pip")
	} else if _, err := exec.LookPath("yum"); err == nil {
		pkgManager = "yum (CentOS/RHEL)"
		fmt.Printf("→ Using %s\n", pkgManager)
		installCmd = exec.Command("sudo", "yum", "install", "-y", "python3", "python3-pip")
	} else if _, err := exec.LookPath("zypper"); err == nil {
		pkgManager = "zypper (openSUSE)"
		fmt.Printf("→ Using %s\n", pkgManager)
		installCmd = exec.Command("sudo", "zypper", "install", "-y", "python3", "python3-pip")
	} else if _, err := exec.LookPath("pacman"); err == nil {
		pkgManager = "pacman (Arch Linux)"
		fmt.Printf("→ Using %s\n", pkgManager)
		installCmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "python", "python-pip")
	} else {
		return fmt.Errorf("no supported package manager found (apt, dnf, yum, zypper, pacman)")
	}

	fmt.Println()
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Println()
		fmt.Println("⚠ Package manager installation failed")
		fmt.Println()
		fmt.Println("Alternative methods:")
		fmt.Println("1. Build from source: https://www.python.org/downloads/source/")
		fmt.Println("2. Use pyenv: https://github.com/pyenv/pyenv")
		fmt.Println("3. Use deadsnakes PPA (Ubuntu): sudo add-apt-repository ppa:deadsnakes/ppa")
		fmt.Println()
		return err
	}

	fmt.Println("✓ Installation completed via package manager")
	return nil
}
