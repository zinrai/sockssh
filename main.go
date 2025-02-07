package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Defaults struct {
		SocksServer string `yaml:"socks_server"`
		Port        int    `yaml:"port"`
	} `yaml:"defaults"`
	Environments map[string]struct {
		SocksServer string `yaml:"socks_server"`
		Port        int    `yaml:"port"`
	} `yaml:"environments"`
}

type SocksConfig struct {
	server string
	port   int
}

func main() {
	var (
		envFlag     string
		verboseFlag bool
	)

	flag.StringVar(&envFlag, "env", "", "Environment to use")
	flag.BoolVar(&verboseFlag, "verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [-env <environment>] [-verbose] [--] <target> [ssh-options...]

Options:
`, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  %s user@example.com
  %s -env prod -- user@example.com -i ~/.ssh/id_rsa
  %s user@example.com -v
`, os.Args[0], os.Args[0], os.Args[0])
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Target host not specified")
		flag.Usage()
		os.Exit(1)
	}

	if args[0] == "--" {
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No command specified")
		os.Exit(1)
	}

	config, err := loadConfig(envFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	target := args[0]

	var sshOptions []string
	if len(args) > 1 {
		sshOptions = args[1:]
	}

	// Build base SSH arguments with required options
	sshArgs := []string{
		"-o", fmt.Sprintf("ProxyCommand=nc -x %s:%d %%h %%p", config.server, config.port),
		"-o", "ForwardAgent=yes",
	}

	// Add any additional SSH options
	if len(sshOptions) > 0 {
		sshArgs = append(sshArgs, sshOptions...)
	}

	// Add the target host last
	sshArgs = append(sshArgs, target)

	if verboseFlag {
		fmt.Fprintf(os.Stderr, "sockssh: Using SOCKS proxy %s:%d\n", config.server, config.port)
		fmt.Fprintf(os.Stderr, "Command: ssh %s\n", strings.Join(sshArgs, " "))
	}

	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if verboseFlag {
			fmt.Fprintf(os.Stderr, "sockssh: Error executing SSH command: %v\n", err)
		}
		os.Exit(1)
	}
}

func loadConfig(env string) (*SocksConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "sockssh.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	result := &SocksConfig{
		server: config.Defaults.SocksServer,
		port:   config.Defaults.Port,
	}

	if env != "" {
		if envConfig, ok := config.Environments[env]; ok {
			if envConfig.SocksServer != "" {
				result.server = envConfig.SocksServer
			}
			if envConfig.Port != 0 {
				result.port = envConfig.Port
			}
		} else {
			return nil, fmt.Errorf("environment '%s' not found in config", env)
		}
	}

	if result.server == "" {
		return nil, fmt.Errorf("socks_server not configured")
	}

	return result, nil
}
