package utils

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

// 執行DOCKER 內命令
func RunCmd(params []interface{}, cmd string, statusChan chan string) {
	cmdStr := fmt.Sprintf(cmd, params...)
	fmt.Println("-- run --\n", cmdStr)

	command := exec.Command("/bin/sh", "-c", cmdStr)

	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatal("--", err)
	}

	if err := command.Start(); err != nil {
		log.Fatal("--", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		statusChan <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("--", err)
	}
	if err := command.Wait(); err != nil {
		log.Fatal("--", err)
	}
}

// 執行命令
func ExecCmd(cmd *exec.Cmd) string {
	// 获取命令的输出
	fmt.Println(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed to run: %s\n", err)
		return ""
	}
	// 将命令的输出内容存储在变量中
	outputStr := string(output)
	fmt.Println(outputStr)
	return outputStr
}
