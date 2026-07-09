package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/ParsaImi/imiterm/internal/model"
)

// Exec replaces the current process with an SSH session to the given host.
// On success, this function never returns — the process becomes ssh.
func Exec(h model.Host) error {
	argv, err := buildArgs(h)
	if err != nil {
		return err
	}

	binary, err := exec.LookPath(argv[0])
	if err != nil {
		return fmt.Errorf("cannot find %s in PATH: %w", argv[0], err)
	}

	if os.Getenv("TMUX") != "" {
		exec.Command("tmux", "rename-window", h.Name).Run()
	}

	// Copy password to clipboard so sudo prompts can be answered with paste.
	// Silently skip if clipboard is unavailable (headless / no display).
	if h.Password != "" {
		if err := clipboard.WriteAll(h.Password); err == nil {
			fmt.Fprintf(os.Stderr, "\r\033[2m[sudo password copied to clipboard]\033[0m\n")
		}
	}

	return syscall.Exec(binary, argv, os.Environ())
}

func buildArgs(h model.Host) ([]string, error) {
	// accept-new: auto-accept unknown host keys, still reject changed ones.
	// Prevents the yes/no prompt that breaks sshpass-based connections.
	strictCheck := "-o StrictHostKeyChecking=accept-new"

	switch h.AuthMethod() {
	case model.AuthKey:
		expanded := expandTilde(h.KeyPath)
		return []string{
			"ssh",
			"-p", fmt.Sprintf("%d", h.Port),
			"-i", expanded,
			strictCheck,
			fmt.Sprintf("%s@%s", h.User, h.Hostname),
		}, nil

	case model.AuthPass:
		if _, err := exec.LookPath("sshpass"); err != nil {
			return nil, fmt.Errorf("password auth requires sshpass (apt install sshpass)")
		}
		return []string{
			"sshpass", "-p", h.Password,
			"ssh",
			"-p", fmt.Sprintf("%d", h.Port),
			strictCheck,
			fmt.Sprintf("%s@%s", h.User, h.Hostname),
		}, nil

	default: // AuthAgent
		return []string{
			"ssh",
			"-p", fmt.Sprintf("%d", h.Port),
			strictCheck,
			fmt.Sprintf("%s@%s", h.User, h.Hostname),
		}, nil
	}
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}
