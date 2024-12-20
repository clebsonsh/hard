package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	hardPath string
	wwwPath  string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hard",
	Short: "Laravel Hard CLI",
	Long:  `Laravel Hard CLI is a tool for managing Laravel projects with Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayHelp()
	},
}

func init() {
	// add flags
	upCmd.Flags().BoolP("detach", "d", false, "Run in detached mode")

	// Add commands
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(bashCmd)
	rootCmd.AddCommand(laravelCmd)
	rootCmd.AddCommand(phpCmd)
	rootCmd.AddCommand(composerCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(npmCmd)
	rootCmd.AddCommand(yarnCmd)
	rootCmd.AddCommand(updateCmd)
}

func init() {
	homeDir, _ := os.UserHomeDir()
	hardPath = filepath.Join(homeDir, ".hard")

	// Load .env file
	if err := loadEnv(filepath.Join(hardPath, ".env")); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	wwwPath = os.Getenv("WWW_PATH")
	if wwwPath == "" {
		fmt.Println("ERROR: WWW_PATH not set in .env file")
		os.Exit(1)
	}
}

func loadEnv(filename string) error {
	// For simplicity, we'll just check if the file exists and load it if it does
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("ERROR: Couldn't find '.env' in the hard directory")
	}
	err := godotenv.Load(filename)
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v", err)
	}

	return nil
}

func dockerComposeCmd() *exec.Cmd {
	dockerComposePath := filepath.Join(hardPath, "docker-compose.yml")
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("ERROR: Docker is not installed or not in PATH")
		os.Exit(1)
	}

	if _, err := exec.LookPath("docker-compose"); err == nil {
		return exec.Command("docker-compose", "-f", dockerComposePath)
	}

	return exec.Command("docker", "compose", "-f", dockerComposePath)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Long:  "Update the Laravel Hard CLI",
	Short: "Update the Laravel Hard CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating Laravel Hard CLI...")

		exec.Command("git", "-C", hardPath, "pull").Run()

		// Drop container images
		fmt.Println("Drop container images...")
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "down", "--rmi", "all", "--volumes")
		dcCmd.Run()

		// Rebuild container images
		fmt.Println("Rebuild container images...")
		dcCmd.Args = append(dcCmd.Args, "build", "--no-cache")
		dcCmd.Run()

		// Start the environment
		fmt.Println("Starting the environment...")
		dcCmd.Args = append(dcCmd.Args, "up", "-d")
		dcCmd.Run()

		fmt.Println("Laravel Hard CLI updated successfully!")
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the environment",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "up")
		if detach, _ := cmd.Flags().GetBool("detach"); detach {
			dcCmd.Args = append(dcCmd.Args, "-d")
		}
		dcCmd.Args = append(dcCmd.Args, args...)
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Run()
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the environment",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "down")
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Run()
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the environment",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "restart")
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Run()
	},
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Display the status of all containers",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "ps")
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Run()
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the environment",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "build")
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Run()
	},
}

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Start a shell session within the app container",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "exec", "app", "bash")
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Stdin = os.Stdin
		dcCmd.Run()
	},
}

var laravelCmd = &cobra.Command{
	Use:   "laravel",
	Short: "Run Laravel commands",
	Run: func(cmd *cobra.Command, args []string) {
		dcCmd := dockerComposeCmd()
		dcCmd.Args = append(dcCmd.Args, "exec", "app", "laravel")
		dcCmd.Args = append(dcCmd.Args, args...)
		dcCmd.Stdout = os.Stdout
		dcCmd.Stderr = os.Stderr
		dcCmd.Stdin = os.Stdin
		dcCmd.Run()
	},
}

func runInProject(projectCmd string, args []string) {
	project := ""
	if len(args) > 0 {
		if _, err := os.Stat(filepath.Join(wwwPath, args[0])); err == nil {
			project = args[0]
			args = args[1:]
		}
	}

	if project == "" {
		dir, _ := os.Getwd()
		project = filepath.Base(dir)
	}

	dcCmd := dockerComposeCmd()
	dcCmd.Args = append(dcCmd.Args, "exec", "-it", "--workdir", "/var/www/"+project, "app", projectCmd)
	dcCmd.Args = append(dcCmd.Args, args...)
	dcCmd.Stdout = os.Stdout
	dcCmd.Stderr = os.Stderr
	dcCmd.Stdin = os.Stdin
	dcCmd.Run()
}

var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "Run PHP commands",
	Run: func(cmd *cobra.Command, args []string) {
		runInProject("php", args)
	},
}

var composerCmd = &cobra.Command{
	Use:   "composer",
	Short: "Run Composer commands",
	Run: func(cmd *cobra.Command, args []string) {
		runInProject("composer", args)
	},
}

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Run Node commands",
	Run: func(cmd *cobra.Command, args []string) {
		runInProject("node", args)
	},
}

var npmCmd = &cobra.Command{
	Use:   "npm",
	Short: "Run npm commands",
	Run: func(cmd *cobra.Command, args []string) {
		runInProject("npm", args)
	},
}

var yarnCmd = &cobra.Command{
	Use:   "yarn",
	Short: "Run Yarn commands",
	Run: func(cmd *cobra.Command, args []string) {
		runInProject("yarn", args)
	},
}

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func displayHelp() {
	fmt.Println("Laravel Hard")
	fmt.Println()
	fmt.Printf("%s\n", yellow("docker-compose Commands:"))
	fmt.Printf("  %s       Build the environment\n", green("hard build"))
	fmt.Printf("  %s        Start the environment\n", green("hard up"))
	fmt.Printf("  %s     Start the environment in the background\n", green("hard up -d"))
	fmt.Printf("  %s      Stop the environment\n", green("hard down"))
	fmt.Printf("  %s   Restart the environment\n", green("hard restart"))
	fmt.Printf("  %s        Display the status of all containers\n", green("hard ps"))
	fmt.Printf("  %s      Start a shell session within the app container\n", green("hard bash"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("Laravel Commands:"))
	fmt.Printf("  %s Run the Laravel command\n", green("hard laravel ..."))
	fmt.Printf("  %s\n", green("hard laravel new awesome-project"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("PHP Commands:"))
	fmt.Printf("  %s Run a snippet of PHP code\n", green("hard php ..."))
	fmt.Printf("  %s\n", green("hard awesome-project php artisan migrate"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("Composer Commands:"))
	fmt.Printf("  %s Run a Composer command\n", green("hard composer ..."))
	fmt.Printf("  %s\n", green("hard awesome-project composer require laravel/sanctum"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("Node Commands:"))
	fmt.Printf("  %s Run a Node command\n", green("hard node ..."))
	fmt.Printf("  %s\n", green("hard awesome-project node --version"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("NPM Commands:"))
	fmt.Printf("  %s Run a npm command\n", green("hard npm ..."))
	fmt.Printf("  %s\n", green("hard awesome-project npm run prod"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("Yarn Commands:"))
	fmt.Printf("  %s Run a Yarn command\n", green("hard yarn ..."))
	fmt.Printf("  %s\n", green("hard awesome-project yarn prod"))
	fmt.Println()
	fmt.Printf("%s\n", yellow("Update Commands:"))
	fmt.Printf("  %s Update the Laravel Hard CLI\n", green("hard update"))
	os.Exit(0)
}
