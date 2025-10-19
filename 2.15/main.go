package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var currentCmd *exec.Cmd //for ctrl+c

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for {
			<-sigChan
			if currentCmd != nil {
				currentCmd.Process.Signal(syscall.SIGINT)
			} else {
				fmt.Println("\nMinishell: ")
			}
		}
	}()

	fmt.Println("Welcome to Minishell: ")
	reader := bufio.NewReader(os.Stdin) //with NewScanner cant process Ctrl+d
	for {
		fmt.Print("Minishell: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nexit")
				return
			}
			fmt.Fprintln(os.Stderr, "Error:", err)
			return
		}
		// Убрать \n в конце
		command := input[:len(input)-1]
		commands := strings.Split(command, "|")
		if len(commands) == 1 { //1 command
			arg := strings.Fields(commands[0])
			if len(arg) == 0 {
				fmt.Println("Not found argument")
				continue
			}
			if arg[0] == "exit" {
				fmt.Println("exit")
				return
			}
			switch arg[0] {
			case "pwd":
				dir, err := os.Getwd()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				fmt.Println(dir)
			case "echo":
				fmt.Println(strings.Join(arg[1:], " "))
			case "cd":
				if len(arg) < 2 || arg[1] == "~" {
					home := os.Getenv("HOME")
					err := os.Chdir(home)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				} else {
					err := os.Chdir(arg[1])
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error:", err)
					}
				}
			case "ps":
				cmd := exec.Command("ps", "aux")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
			case "kill":
				if len(arg) < 2 {
					fmt.Fprintln(os.Stderr, "Not enough argument")
					continue
				}
				pid, err := strconv.Atoi(arg[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					continue
				}
				err = syscall.Kill(pid, syscall.SIGTERM)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
			default:
				currentCmd = exec.Command(arg[0], arg[1:]...)
				currentCmd.Stdout = os.Stdout
				currentCmd.Stderr = os.Stderr
				err := currentCmd.Run()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				currentCmd = nil //clean after finish
			}
		} else { //many command
			var cmdList []*exec.Cmd //commands

			for _, cmd := range commands {
				arg := strings.Fields(strings.TrimSpace(cmd))
				cmd := exec.Command(arg[0], arg[1:]...)
				cmdList = append(cmdList, cmd)
			}
			currentCmd = cmdList[0]
			//create pipes
			for i := 0; i < len(cmdList)-1; i++ {
				stdout, _ := cmdList[i].StdoutPipe()
				cmdList[i+1].Stdin = stdout
			}

			cmdList[0].Stdin = os.Stdin
			cmdList[len(cmdList)-1].Stdout = os.Stdout

			for _, cmd := range cmdList {
				cmd.Start()
			}

			for _, cmd := range cmdList {
				cmd.Wait()
			}
			currentCmd = nil
		}

	}
}
