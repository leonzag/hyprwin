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
	Grouped        []string `json:"grouped"`
	Class          string   `json:"class"`
	Swallowing     string   `json:"swallowing"`
	InitialTitle   string   `json:"initialTitle"`
	InitialClass   string   `json:"inittialClass"`
	Address        string   `json:"address"`
	Title          string   `json:"title"`
	Workspace      wsObj    `json:"workspace"`
	Size           [2]int   `json:"size"`
	At             [2]int   `json:"at"`
	Monitor        int      `json:"monitor"`
	FocusHistoryID int      `json:"fucusHistoryID"`
	Pid            int      `json:"pid"`
	FullscreenMode int      `json:"fullscreenMode"`
	Pinned         bool     `json:"pinned"`
	Fullscreen     bool     `json:"fullscreen"`
	Xwayland       bool     `json:"xwayland"`
	FakeFullscreen bool     `json:"fakeFullscreen"`
	Hidden         bool     `json:"hidden"`
	Mapped         bool     `json:"mapped"`
	Floating       bool     `json:"floating"`
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
