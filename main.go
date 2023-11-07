package main

import (
	"fmt"
	"os"
	"os/exec"
	"udate/control"

	"github.com/spf13/cobra"
)

func main() {
	control.CHttp()
}

func ssn() {

	var rootCmd = &cobra.Command{Use: "sshls"}
	var host string
	var ldate string

	// 添加名为 "ls" 的命令
	var lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "SSH to a host and execute 'ls' command",
		Run: func(cmd *cobra.Command, args []string) {
			if host == "" {
				fmt.Println("Please provide both host and user")
				return
			}
			// 使用 SSH 命令执行 'ls' 在远程主机上
			datecmd := fmt.Sprintf("/root/get_date.sh")
			sshCmd := exec.Command("ssh", "root@"+host, datecmd)
			sshCmd.Stdout = os.Stdout
			sshCmd.Stderr = os.Stderr
			err := sshCmd.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		},
	}
	// 添加参数
	lsCmd.Flags().StringVarP(&host, "host", "H", "", "Remote host address (e.g., ubuntu@1.2.33.4)")
	lsCmd.Flags().StringVarP(&ldate, "ldate", "C", "", "linux command line")
	// 将 "ls" 命令添加到根命令
	rootCmd.AddCommand(lsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
