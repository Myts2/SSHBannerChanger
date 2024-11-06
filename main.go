package main

import (
	"bytes"
	"fmt"
	"github.com/cakturk/go-netstat/netstat"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Get the local open SSH port, PID, and the corresponding sshd process file location
func getSSHPortAndProcess() (string, string, string, error) {
	fn := func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	}
	tabs, err := netstat.TCPSocks(fn)
	if err != nil {
		panic(err)
	}
	sshdPort := 0
	sshdPath := ""
	sshdPid := 0
	for _, tab := range tabs {
		if tab.Process != nil &&
			tab.Process.Name == "sshd" &&
			strings.Compare(tab.LocalAddr.IP.String(), "0.0.0.0") == 0 {
			sshdPort = int(tab.LocalAddr.Port)
			sshdPid = tab.Process.Pid
			sshdPath, err = os.Readlink(fmt.Sprintf("/proc/%d/exe", tab.Process.Pid))
		}
	}
	sshdPath = strings.Split(sshdPath, " ")[0]
	return strconv.Itoa(sshdPort), strconv.Itoa(sshdPid), sshdPath, nil
}

// ReplaceAndBackupBytesInFile replaces bytes in a file and backs up the original file
func ReplaceAndBackupBytesInFile(filePath, backupPath string, oldBytes, newBytes []byte) error {
	sourceFileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(backupPath, input, sourceFileInfo.Mode())
	if err != nil {
		return err
	}
	fmt.Printf("Successfully backup sshd executable to: %s\n", backupPath)

	// Use bytes.ReplaceAll to replace
	output := bytes.ReplaceAll(input, oldBytes, newBytes)

	tmpName := fmt.Sprintf("%s.%d.tmp", filePath, time.Now().Unix())
	// Write the replaced bytes back to the file
	err = ioutil.WriteFile(tmpName, output, sourceFileInfo.Mode())
	if err != nil {
		return err
	}

	return os.Rename(tmpName, filePath)
}

// Modify the banner in the sshd file
func modifySSHDBanner(banner, sshdPath string) error {
	// Backup the sshd executable
	backupPath := sshdPath + fmt.Sprintf(".%d", time.Now().Unix())
	replacement := "OpenSSH_fix"
	if len(replacement) < len(banner) {
		replacement = fmt.Sprintf("%s%s", replacement, strings.Repeat(" ", len(banner)-len(replacement)))
	}
	return ReplaceAndBackupBytesInFile(sshdPath, backupPath, []byte(banner), []byte(replacement))
}

func getSSHBanner(port string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Read the banner
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	// Filter to keep only visible characters
	var visibleChars []rune
	for _, r := range string(buffer[:n]) {
		if unicode.IsPrint(r) {
			visibleChars = append(visibleChars, r)
		}
	}
	input := string(visibleChars)
	start := strings.Index(input, "OpenSSH")
	result := input[start:]

	return result, nil
}

func main() {
	// Get SSH port, PID, and process file path
	sshPort, sshPID, sshdPath, err := getSSHPortAndProcess()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("SSH port: %s, PID: %s, SSHD Locate: %s\n", sshPort, sshPID, sshdPath)
	banner, err := getSSHBanner(sshPort)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Current Banner:", banner)
	// Modify the banner in the sshd file
	err = modifySSHDBanner(banner, sshdPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Successfully modified the SSHD banner")
}
