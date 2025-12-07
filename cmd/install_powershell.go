package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const powershellVersion = "7"

var (
	forcePowerShellInstall bool
	skipPowerShellVerify   bool
)

// powershellCmd represents the powershell subcommand
var powershellCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Install PowerShell 7 for netcheck PS check type",
	Long: `Install PowerShell 7 (pwsh) to enable netcheck's PowerShell script functionality.

This command will attempt to install PowerShell 7 using the appropriate
package manager for your operating system:
  - Windows: winget, chocolatey, or direct download
  - macOS: Homebrew
  - Linux: Package manager (apt, dnf, yum, zypper) or snap

PowerShell 7 is cross-platform and runs on Windows, macOS, and Linux.
The command will first check if PowerShell 7 is already installed and
skip installation unless --force is specified.`,
	RunE: installPowerShell,
}

func init() {
	installCmd.AddCommand(powershellCmd)
	powershellCmd.Flags().BoolVar(&forcePowerShellInstall, "force", false, "force installation even if PowerShell is already installed")
	powershellCmd.Flags().BoolVar(&skipPowerShellVerify, "skip-verify", false, "skip verification after installation")
}

func installPowerShell(cmd *cobra.Command, args []string) error {
	fmt.Println("PowerShell 7 Installation for netcheck")
	fmt.Println("========================================")
	fmt.Println()

	// Check if PowerShell is already installed
	if !forcePowerShellInstall {
		if version, installed := checkPowerShellInstalled(); installed {
			fmt.Printf("✓ PowerShell is already installed: %s\n", version)
			fmt.Println()
			fmt.Println("Use --force to reinstall")
			return nil
		}
	}

	fmt.Printf("Installing PowerShell %s for %s/%s...\n\n", powershellVersion, runtime.GOOS, runtime.GOARCH)

	var err error
	switch runtime.GOOS {
	case "windows":
		err = installPowerShellWindows()
	case "darwin":
		err = installPowerShellMacOS()
	case "linux":
		err = installPowerShellLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Verify installation unless skipped
	if !skipPowerShellVerify {
		fmt.Println()
		fmt.Println("Verifying installation...")
		if version, installed := checkPowerShellInstalled(); installed {
			fmt.Printf("✓ PowerShell successfully installed: %s\n", version)
		} else {
			fmt.Println("⚠ Warning: PowerShell installation completed but verification failed")
			fmt.Println("  You may need to restart your terminal or add PowerShell to your PATH")
		}
	}

	return nil
}

func checkPowerShellInstalled() (string, bool) {
	// Check for pwsh (PowerShell 7+)
	cmd := exec.Command("pwsh", "--version")
	output, err := cmd.CombinedOutput()
	if err == nil {
		version := strings.TrimSpace(string(output))
		return version, true
	}

	return "", false
}

func installPowerShellWindows() error {
	fmt.Println("Attempting Windows installation methods...")
	fmt.Println()

	// Try winget first (Windows 10/11)
	if _, err := exec.LookPath("winget"); err == nil {
		fmt.Println("→ Using winget (Windows Package Manager)")
		cmd := exec.Command("winget", "install", "-e", "--id", "Microsoft.PowerShell")
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
		cmd := exec.Command("choco", "install", "powershell-core", "-y")
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
	fmt.Println("Please install PowerShell 7 manually:")
	fmt.Println("1. Visit: https://aka.ms/powershell-release?tag=stable")
	fmt.Println("2. Download the MSI installer for Windows")
	fmt.Println("3. Run the installer")
	fmt.Println("4. PowerShell 7 will be installed as 'pwsh.exe'")
	fmt.Println()
	fmt.Println("Alternative - Using Windows PowerShell:")
	fmt.Println("  iex \"& { $(irm https://aka.ms/install-powershell.ps1) } -UseMSI\"")
	fmt.Println()

	return fmt.Errorf("automatic installation not available - please install manually")
}

func installPowerShellMacOS() error {
	fmt.Println("Attempting macOS installation methods...")
	fmt.Println()

	// Check for Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		fmt.Println("→ Using Homebrew")
		fmt.Println("Running: brew install --cask powershell")
		cmd := exec.Command("brew", "install", "--cask", "powershell")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via Homebrew")
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
	fmt.Println("1. Visit: https://aka.ms/powershell-release?tag=stable")
	fmt.Println("2. Download the PKG installer for macOS")
	fmt.Println("3. Run the installer package")
	fmt.Println()

	return fmt.Errorf("automatic installation failed - please install Homebrew or install manually")
}

func installPowerShellLinux() error {
	fmt.Println("Attempting Linux installation methods...")
	fmt.Println()

	// Detect distribution
	distro := detectLinuxDistro()
	fmt.Printf("Detected distribution: %s\n", distro)
	fmt.Println()

	var err error

	// Try distribution-specific package managers
	if _, err := exec.LookPath("apt"); err == nil {
		err = installPowerShellDebian()
	} else if _, err := exec.LookPath("dnf"); err == nil {
		err = installPowerShellFedora()
	} else if _, err := exec.LookPath("yum"); err == nil {
		err = installPowerShellRHEL()
	} else if _, err := exec.LookPath("zypper"); err == nil {
		err = installPowerShellOpenSUSE()
	} else if _, err := exec.LookPath("snap"); err == nil {
		// Try snap as fallback
		fmt.Println("→ Using snap")
		cmd := exec.Command("sudo", "snap", "install", "powershell", "--classic")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			fmt.Println("✓ Installation completed via snap")
			return nil
		}
		err = fmt.Errorf("snap installation failed")
	} else {
		return fmt.Errorf("no supported package manager found (apt, dnf, yum, zypper, snap)")
	}

	if err != nil {
		fmt.Println()
		fmt.Println("⚠ Package manager installation failed")
		fmt.Println()
		fmt.Println("Alternative methods:")
		fmt.Println("1. Using snap: sudo snap install powershell --classic")
		fmt.Println("2. Manual installation: https://aka.ms/powershell-release?tag=stable")
		fmt.Println("3. Using binary archives: https://docs.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-linux")
		fmt.Println()
		return err
	}

	return nil
}

func detectLinuxDistro() string {
	// Try to read /etc/os-release
	data, err := os.ReadFile("/etc/os-release")
	if err == nil {
		content := string(data)
		if strings.Contains(content, "Ubuntu") {
			return "Ubuntu"
		} else if strings.Contains(content, "Debian") {
			return "Debian"
		} else if strings.Contains(content, "Fedora") {
			return "Fedora"
		} else if strings.Contains(content, "CentOS") {
			return "CentOS"
		} else if strings.Contains(content, "Red Hat") {
			return "RHEL"
		} else if strings.Contains(content, "openSUSE") {
			return "openSUSE"
		}
	}
	return "Unknown"
}

func installPowerShellDebian() error {
	fmt.Println("→ Using apt (Debian/Ubuntu)")
	fmt.Println()

	// Update package list
	fmt.Println("Updating package lists...")
	updateCmd := exec.Command("sudo", "apt", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	_ = updateCmd.Run()

	// Install prerequisites
	fmt.Println()
	fmt.Println("Installing prerequisites...")
	prereqCmd := exec.Command("sudo", "apt", "install", "-y", "wget", "apt-transport-https", "software-properties-common")
	prereqCmd.Stdout = os.Stdout
	prereqCmd.Stderr = os.Stderr
	if err := prereqCmd.Run(); err != nil {
		return fmt.Errorf("failed to install prerequisites: %w", err)
	}

	// Download Microsoft repository GPG keys
	fmt.Println()
	fmt.Println("Adding Microsoft repository...")
	downloadCmd := exec.Command("wget", "-q", "https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb")
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		// Try generic approach
		fmt.Println("Using snap as alternative...")
		snapCmd := exec.Command("sudo", "snap", "install", "powershell", "--classic")
		snapCmd.Stdout = os.Stdout
		snapCmd.Stderr = os.Stderr
		return snapCmd.Run()
	}

	// Install the repository configuration
	installRepoCmd := exec.Command("sudo", "dpkg", "-i", "packages-microsoft-prod.deb")
	installRepoCmd.Stdout = os.Stdout
	installRepoCmd.Stderr = os.Stderr
	_ = installRepoCmd.Run()

	// Clean up
	_ = os.Remove("packages-microsoft-prod.deb")

	// Update package list again
	updateCmd2 := exec.Command("sudo", "apt", "update")
	updateCmd2.Stdout = os.Stdout
	updateCmd2.Stderr = os.Stderr
	_ = updateCmd2.Run()

	// Install PowerShell
	fmt.Println()
	fmt.Println("Installing PowerShell...")
	installCmd := exec.Command("sudo", "apt", "install", "-y", "powershell")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install PowerShell: %w", err)
	}

	fmt.Println("✓ Installation completed via apt")
	return nil
}

func installPowerShellFedora() error {
	fmt.Println("→ Using dnf (Fedora)")
	fmt.Println()

	// Register the Microsoft repository
	fmt.Println("Adding Microsoft repository...")
	repoCmd := exec.Command("sudo", "rpm", "--import", "https://packages.microsoft.com/keys/microsoft.asc")
	repoCmd.Stdout = os.Stdout
	repoCmd.Stderr = os.Stderr
	_ = repoCmd.Run()

	// Add repository
	curlCmd := exec.Command("bash", "-c", "curl https://packages.microsoft.com/config/rhel/8/prod.repo | sudo tee /etc/yum.repos.d/microsoft.repo")
	curlCmd.Stdout = os.Stdout
	curlCmd.Stderr = os.Stderr
	if err := curlCmd.Run(); err != nil {
		return fmt.Errorf("failed to add repository: %w", err)
	}

	// Install PowerShell
	fmt.Println()
	fmt.Println("Installing PowerShell...")
	installCmd := exec.Command("sudo", "dnf", "install", "-y", "powershell")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install PowerShell: %w", err)
	}

	fmt.Println("✓ Installation completed via dnf")
	return nil
}

func installPowerShellRHEL() error {
	fmt.Println("→ Using yum (CentOS/RHEL)")
	fmt.Println()

	// Register the Microsoft repository
	fmt.Println("Adding Microsoft repository...")
	repoCmd := exec.Command("sudo", "rpm", "--import", "https://packages.microsoft.com/keys/microsoft.asc")
	repoCmd.Stdout = os.Stdout
	repoCmd.Stderr = os.Stderr
	_ = repoCmd.Run()

	// Add repository
	curlCmd := exec.Command("bash", "-c", "curl https://packages.microsoft.com/config/rhel/8/prod.repo | sudo tee /etc/yum.repos.d/microsoft.repo")
	curlCmd.Stdout = os.Stdout
	curlCmd.Stderr = os.Stderr
	if err := curlCmd.Run(); err != nil {
		return fmt.Errorf("failed to add repository: %w", err)
	}

	// Install PowerShell
	fmt.Println()
	fmt.Println("Installing PowerShell...")
	installCmd := exec.Command("sudo", "yum", "install", "-y", "powershell")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install PowerShell: %w", err)
	}

	fmt.Println("✓ Installation completed via yum")
	return nil
}

func installPowerShellOpenSUSE() error {
	fmt.Println("→ Using zypper (openSUSE)")
	fmt.Println()

	// Register the Microsoft repository
	fmt.Println("Adding Microsoft repository...")
	repoCmd := exec.Command("sudo", "rpm", "--import", "https://packages.microsoft.com/keys/microsoft.asc")
	repoCmd.Stdout = os.Stdout
	repoCmd.Stderr = os.Stderr
	_ = repoCmd.Run()

	// Add repository
	curlCmd := exec.Command("bash", "-c", "curl https://packages.microsoft.com/config/rhel/8/prod.repo | sudo tee /etc/zypp/repos.d/microsoft.repo")
	curlCmd.Stdout = os.Stdout
	curlCmd.Stderr = os.Stderr
	if err := curlCmd.Run(); err != nil {
		return fmt.Errorf("failed to add repository: %w", err)
	}

	// Install PowerShell
	fmt.Println()
	fmt.Println("Installing PowerShell...")
	installCmd := exec.Command("sudo", "zypper", "install", "-y", "powershell")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install PowerShell: %w", err)
	}

	fmt.Println("✓ Installation completed via zypper")
	return nil
}
