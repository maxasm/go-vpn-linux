package main

import (
	"os/exec"
	"bytes"
	"log"
	"os"
	"errors"
	"strings"
)

// Helper function to run a command
// return error, std_out, std_error
func run_cmd(cmd_name string, cmd_args []string) (error,string,string) {
	// run the command
	cmd := exec.Command(cmd_name, cmd_args...)

	// create the buffers where the output is stored
	std_out := bytes.Buffer{}
	std_err := bytes.Buffer{}

	cmd.Stdout = &std_out
	cmd.Stderr = &std_err

	// run the command
	err__run_cmd := cmd.Run()

	if err__run_cmd != nil {
		return err__run_cmd, std_out.String(),std_err.String()
	}

	return nil, std_out.String(), std_err.String()
}

func create_tun_device(dev_name string) (error) {
	// TODO: verify that dev_name is a valid device-name
	if len(strings.Trim(dev_name, " \t\r\n")) == 0 {
		return errors.New("device name can not be empty")
	}
	err__cmd_create_dev, _, cmd__create_dev_stderr := run_cmd("ip", []string{"tuntap", "add", "mode", "tun", "name", dev_name})
	if err__cmd_create_dev != nil {
		return errors.New(cmd__create_dev_stderr)
	}
	
	return nil
}

func main() {
	err__create_device := create_tun_device("test")
	if err__create_device != nil {
		log.Printf("[ERROR]:%s\n", err__create_device)
		os.Exit(1)
	}

	log.Printf("[SUCCESS]: created tun interface successfully.")
}
