package hyprwin

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var (
	BufSize              = 8192
	ErrIpcHyprSignNotSet = errors.New("failed connect to ipc: $HYPRLAND_INSTANCE_SIGNATURE env var not set")
	ErrIpcSocketNotFound = errors.New("failed connect to ipc: socket not found")
)

type WinObj struct {
	Address          string   `json:"address"`
	Mapped           bool     `json:"mapped"`
	Hidden           bool     `json:"hidden"`
	At               [2]int   `json:"at"`
	Size             [2]int   `json:"size"`
	Workspace        wsObj    `json:"workspace"`
	Floating         bool     `json:"floating"`
	Pseudo           bool     `json:"pseudo"`
	Monitor          int      `json:"monitor"`
	Class            string   `json:"class"`
	Title            string   `json:"title"`
	InitialClass     string   `json:"inittialClass"`
	InitialTitle     string   `json:"initialTitle"`
	Pid              int      `json:"pid"`
	Xwayland         bool     `json:"xwayland"`
	Pinned           bool     `json:"pinned"`
	Fullscreen       int      `json:"fullscreen"`
	FullscreenClient int      `json:"fullscreenClient"`
	Grouped          []string `json:"grouped"`
	Tags             []string `json:"tags"`
	Swallowing       string   `json:"swallowing"`
	FocusHistoryID   int      `json:"fucusHistoryID"`
}

type wsObj struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type IPC interface {
	Hyprctl(commands ...string) ([]byte, error)
	ActiveWindow() (*WinObj, error)
}

type ipc struct {
	addr *net.UnixAddr
}

func InitIPC() (IPC, error) {
	sign := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if sign == "" {
		return nil, ErrIpcHyprSignNotSet
	}

	socketPath := ""

	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	tmpDir := os.TempDir()
	dirs := []string{tmpDir, runtimeDir}

	for _, socketHome := range dirs {
		fpath := filepath.Join(socketHome, "hypr", sign, ".socket.sock")
		if finfo, err := os.Stat(fpath); err == nil && !finfo.IsDir() {
			socketPath = fpath
		}
	}
	if socketPath == "" {
		return nil, ErrIpcSocketNotFound
	}

	return &ipc{
		&net.UnixAddr{
			Name: socketPath,
			Net:  "unix",
		},
	}, nil
}

func (c *ipc) ActiveWindow() (*WinObj, error) {
	jsonStr, err := c.Hyprctl("activewindow")
	if err != nil {
		return nil, err
	}
	win := &WinObj{}
	err = json.Unmarshal([]byte(jsonStr), win)
	return win, err
}

// Hyprctl executes commands and returns response. If only one is passed, then sets json flag (-j)
func (c *ipc) Hyprctl(commands ...string) (resp []byte, err error) {
	if len(commands) == 0 {
		return nil, errors.New("attempt to write to a socket with empty string")
	}
	cmd := "j/" + commands[0]
	if len(commands) > 1 {
		cmd = "[[BATCH]] " + strings.Join(commands, "; ")
	}

	conn, err := net.DialUnix("unix", nil, c.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(cmd)); err != nil {
		return nil, err
	}

	var response []byte
	buf := make([]byte, BufSize)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		response = append(response, buf[:n]...)

		if n < BufSize {
			break
		}
	}

	return response, nil
}
