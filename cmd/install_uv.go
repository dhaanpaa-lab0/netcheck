package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	forceUVInstall bool
	skipUVVerify   bool
)

// uvCmd represents the uv subcommand
var uvCmd = &cobra.Command{
	Use:   "uv",
	Short: "Install UV (Python package installer) for managing Python dependencies",
	Long: `Install UV, an extremely fast Python package installer and resolver.

UV is a modern package manager for Python that can replace pip, pip-tools,
poetry, and more. It's useful for managing dependencies in Python scripts
used with netcheck's PY check type.

This command will attempt to install UV using the appropriate method for
your operating system:
  - Windows: Official installer script, pip, or cargo
  - macOS: Homebrew or official installer script
  - Linux: Official installer script or cargo

The command will first check if UV is already installed and skip
installation unless --force is specified.

Learn more: https://github.com/astral-sh/uv`,
	RunE: installUV,
}

func init() {
	installCmd.AddCommand(uvCmd)
	uvCmd.Flags().BoolVar(&forceUVInstall, "force", false, "force installation even if UV is already installed")
	uvCmd.Flags().BoolVar(&skipUVVerify, "skip-verify", false, "skip verification after installation")
}

func installUV(cmd *cobra.Command, args []string) error {
	fmt.Println("UV Python Package Installer - Installation for netcheck")
	fmt.Println("========================================================")
	fmt.Println()

	// Check if UV is already installed
	if !forceUVInstall {
		if version, installed := checkUVInstalled(); installed {
			fmt.Printf("✓ UV is already installed: %s\n", version)
			fmt.Println()
			fmt.Println("Use --force to reinstall")
			return nil
		}
	}

	fmt.Printf("Installing UV for %s/%s...\n\n", runtime.GOOS, runtime.GOARCH)

	var err error
	switch runtime.GOOS {
	case "windows":
		err = installUVWindows()
	case "darwin":
		err = installUVMacOS()
	case "linux":
		err = installUVLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Verify installation unless skipped
	if !skipUVVerify {
		fmt.Println()
		fmt.Println("Verifying installation...")
		if version, installed := checkUVInstalled(); installed {
			fmt.Printf("✓ UV successfully installed: %s\n", version)
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  uv pip install <package>  # Install Python packages")
			fmt.Println("  uv venv                   # Create virtual environment")
			fmt.Println("  uv help                   # Show help")
		} else {
			fmt.Println("⚠ Warning: UV installation completed but verification failed")
			fmt.Println("  You may need to restart your terminal or add UV to your PATH")
			fmt.Println("  Default UV location:")
			fmt.Println("    - Windows: %USERPROFILE%\\.cargo\\bin\\uv.exe")
			fmt.Println("    - macOS/Linux: ~/.cargo/bin/uv")
		}
	}

	return nil
}

func checkUVInstalled() (string, bool) {
	// Check for uv command
	cmd := exec.Command("uv", "--version")
	output, err := cmd.CombinedOutput()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return version, true
	}

	return "", false
}

func installUVWindows() error {
	fmt.Println("Attempting Windows installation methods...")
	fmt.Println()

	// Method 1: Official installer script (PowerShell)
	if _, err := exec.LookPath("powershell"); err == nil {
		fmt.Println("→ Using official installer script (PowerShell)")
		fmt.Println()
		fmt.Println("Running: powershell -c \"irm https://astral.sh/uv/install.ps1 | iex\"")
		fmt.Println()

		cmd := exec.Command("powershell", "-Command", "irm https://astral.sh/uv/install.ps1 | iex")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println()
			fmt.Println("✓ Installation completed via PowerShell installer")
			fmt.Println()
			fmt.Println("Note: You may need to restart your terminal for PATH changes to take effect")
			return nil
		}
		fmt.Println("⚠ PowerShell installer failed, trying next method...")
		fmt.Println()
	}

	// Method 2: pip (if Python is available)
	if _, err := exec.LookPath("pip"); err == nil {
		fmt.Println("→ Using pip")
		fmt.Println("Running: pip install uv")
		fmt.Println()

		cmd := exec.Command("pip", "install", "uv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via pip")
			return nil
		}
		fmt.Println("⚠ pip installation failed, trying next method...")
		fmt.Println()
	}

	// Method 3: cargo (if Rust is available)
	if _, err := exec.LookPath("cargo"); err == nil {
		fmt.Println("→ Using cargo (Rust package manager)")
		fmt.Println("Running: cargo install uv")
		fmt.Println()
		fmt.Println("Note: This may take several minutes to compile...")

		cmd := exec.Command("cargo", "install", "uv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via cargo")
			return nil
		}
		fmt.Println("⚠ cargo installation failed")
		fmt.Println()
	}

	// Manual installation instructions
	fmt.Println("→ Automated installation not available")
	fmt.Println()
	fmt.Println("Please install UV manually:")
	fmt.Println()
	fmt.Println("Method 1 - PowerShell (recommended):")
	fmt.Println("  powershell -c \"irm https://astral.sh/uv/install.ps1 | iex\"")
	fmt.Println()
	fmt.Println("Method 2 - pip (if Python is installed):")
	fmt.Println("  pip install uv")
	fmt.Println()
	fmt.Println("Method 3 - Direct download:")
	fmt.Println("  Visit: https://github.com/astral-sh/uv/releases")
	fmt.Println()

	return fmt.Errorf("automatic installation not available - please install manually")
}

func installUVMacOS() error {
	fmt.Println("Attempting macOS installation methods...")
	fmt.Println()

	// Method 1: Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("→ Using Homebrew")
		fmt.Println("Running: brew install uv")
		fmt.Println()

		cmd := exec.Command("brew", "install", "uv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via Homebrew")
			return nil
		}
		fmt.Println("⚠ Homebrew installation failed, trying next method...")
		fmt.Println()
	}

	// Method 2: Official installer script (curl)
	if _, err := exec.LookPath("curl"); err == nil {
		fmt.Println("→ Using official installer script (curl)")
		fmt.Println()
		fmt.Println("Running: curl -LsSf https://astral.sh/uv/install.sh | sh")
		fmt.Println()

		cmd := exec.Command("sh", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println()
			fmt.Println("✓ Installation completed via curl installer")
			fmt.Println()
			fmt.Println("Note: You may need to restart your terminal or run:")
			fmt.Println("  source $HOME/.cargo/env")
			return nil
		}
		fmt.Println("⚠ curl installer failed, trying next method...")
		fmt.Println()
	}

	// Method 3: pip (if Python is available)
	if _, err := exec.LookPath("pip3"); err == nil {
		fmt.Println("→ Using pip3")
		fmt.Println("Running: pip3 install uv")
		fmt.Println()

		cmd := exec.Command("pip3", "install", "uv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via pip3")
			return nil
		}
		fmt.Println("⚠ pip3 installation failed")
		fmt.Println()
	}

	// Manual installation instructions
	fmt.Println("→ Automated installation not available")
	fmt.Println()
	fmt.Println("Please install UV manually:")
	fmt.Println()
	fmt.Println("Method 1 - Homebrew (recommended):")
	fmt.Println("  brew install uv")
	fmt.Println()
	fmt.Println("Method 2 - Official installer:")
	fmt.Println("  curl -LsSf https://astral.sh/uv/install.sh | sh")
	fmt.Println()
	fmt.Println("Method 3 - pip:")
	fmt.Println("  pip3 install uv")
	fmt.Println()

	return fmt.Errorf("automatic installation failed - please install manually")
}

func installUVLinux() error {
	fmt.Println("Attempting Linux installation methods...")
	fmt.Println()

	// Method 1: Official installer script (curl) - most reliable
	if _, err := exec.LookPath("curl"); err == nil {
		fmt.Println("→ Using official installer script (curl)")
		fmt.Println()
		fmt.Println("Running: curl -LsSf https://astral.sh/uv/install.sh | sh")
		fmt.Println()

		cmd := exec.Command("sh", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println()
			fmt.Println("✓ Installation completed via curl installer")
			fmt.Println()
			fmt.Println("Note: You may need to restart your terminal or run:")
			fmt.Println("  source $HOME/.cargo/env")
			return nil
		}
		fmt.Println("⚠ curl installer failed, trying next method...")
		fmt.Println()
	}

	// Method 2: pip (if Python is available)
	pythonCmd := "pip3"
	if _, err := exec.LookPath("pip3"); err != nil {
		if _, err := exec.LookPath("pip"); err == nil {
			pythonCmd = "pip"
		} else {
			pythonCmd = ""
		}
	}

	if pythonCmd != "" {
		fmt.Printf("→ Using %s\n", pythonCmd)
		fmt.Printf("Running: %s install uv\n", pythonCmd)
		fmt.Println()

		cmd := exec.Command(pythonCmd, "install", "uv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Printf("✓ Installation completed via %s\n", pythonCmd)
			return nil
		}
		fmt.Printf("⚠ %s installation failed, trying next method...\n", pythonCmd)
		fmt.Println()
	}

	// Method 3: cargo (if Rust is available)
	if _, err := exec.LookPath("cargo"); err == nil {
		fmt.Println("→ Using cargo (Rust package manager)")
		fmt.Println("Running: cargo install uv")
		fmt.Println()
		fmt.Println("Note: This may take several minutes to compile...")

		cmd := exec.Command("cargo", "install", "uv")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via cargo")
			return nil
		}
		fmt.Println("⚠ cargo installation failed")
		fmt.Println()
	}

	// Manual installation instructions
	fmt.Println("→ Automated installation not available")
	fmt.Println()
	fmt.Println("Please install UV manually:")
	fmt.Println()
	fmt.Println("Method 1 - Official installer (recommended):")
	fmt.Println("  curl -LsSf https://astral.sh/uv/install.sh | sh")
	fmt.Println()
	fmt.Println("Method 2 - pip:")
	fmt.Println("  pip3 install uv")
	fmt.Println()
	fmt.Println("Method 3 - cargo (requires Rust):")
	fmt.Println("  cargo install uv")
	fmt.Println()

	return fmt.Errorf("automatic installation failed - please install manually")
}
