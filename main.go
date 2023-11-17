package main

import (
	"os/exec"
	"log"
	"errors"
	"strings"
)

// the name of the `TUN` device being made
const TUN_DEV_NAME = "tunmax"

func create_tun_device() (error) {
	create_tun_cmd := exec.Command("ip", "tuntap", "add", "mode", "tun", "name", TUN_DEV_NAME)
	co, err__run_create_tun_cmd := create_tun_cmd.CombinedOutput()
	if err__run_create_tun_cmd != nil {
		return errors.New(string(co)) 
	}
		
	return nil
}

// checks if the `TUN` module is loaded. 
func check_tun_module() (error) {
	lsmod_cmd := exec.Command("lsmod")
	lsmod_cmd_stdout, err__lsmod_cmd := lsmod_cmd.Output()
	if err__lsmod_cmd != nil {
		return err__lsmod_cmd
	}

	// create the `grep` command
	grep_cmd := exec.Command("grep", "tun")

	grep_stdin_pipe, err__grep_stdin_pipe := grep_cmd.StdinPipe()
	if err__grep_stdin_pipe != nil {
		return err__grep_stdin_pipe
	}

	// write the `lsmod` output to it
	grep_stdin_pipe.Write([]byte(lsmod_cmd_stdout))
	grep_stdin_pipe.Close()

	// get the output of the `grep` command
	grep_cmd_stdout, err__run_grep_cmd := grep_cmd.Output()
	if err__run_grep_cmd != nil {
		return err__run_grep_cmd
	}

	// check if the outout of the grep command is empty
	if len(strings.Trim(string(grep_cmd_stdout), " ")) == 0 {
		return errors.New("tun module is not loaded.")
	} 

	return nil
}

// bring `UP` the created TUN device
func up_tun_device() error {
	up_cmd := exec.Command("ip", "link", "set", TUN_DEV_NAME, "up")
	co, err__run_up_cmd := up_cmd.CombinedOutput()
	if err__run_up_cmd != nil {
		return errors.New(string(co))
	}
	return nil
}

func init_tun() {
	// 1. Check if the `TUN` module is installed
	err__check_tun_module := check_tun_module()
	if err__check_tun_module != nil {
		log.Fatalf("[ERROR]: %s\n", err__check_tun_module)
	}
	log.Printf("[CHECK]: tun module is loaded.\n")

	// 2. Create the TUN device
	err__create_device := create_tun_device()
	if err__create_device != nil {
		log.Fatalf("[ERROR]: %s\n", err__create_device)
	}

	log.Printf("[SUCCESS]: created tun interface [%s] successfully.\n", TUN_DEV_NAME)

	// 3. Bring the device 'UP'
	err__up_tun_device := up_tun_device()
	if err__up_tun_device != nil {
		log.Fatalf("[ERROR]: %s\n", err__up_tun_device)
	}

	log.Printf("[SUCCESS]: tun device %s is UP.\n", TUN_DEV_NAME)
	// Assign an IP address and a Subnet Mask to the interface
	// TODO: research on private IPs and how I can check if a certain range is available to use.
}

func main() {
	init_tun()
}
