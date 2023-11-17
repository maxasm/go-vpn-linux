package main

import (
	"os/exec"
	"log"
	"errors"
)

// the name of the `TUN` device being made
const TUN_DEV_NAME = "tunmax"

// TODO: create an helper function to run the commands
func run_cmd(pipein []byte,command string, args []string) (error,[]byte) {
	cmd := exec.Command(command, args...)
	cmd_stdin_pipe, err__get_cmd_stdin_pipe := cmd.StdinPipe()
	if err__get_cmd_stdin_pipe != nil {
		return err__get_cmd_stdin_pipe, nil
	}

	cmd_stdin_pipe.Write(pipein)
	cmd_stdin_pipe.Close()

	co, err__run_cmd := cmd.CombinedOutput()
	if err__run_cmd != nil {
		return errors.New(string(co)), co
	}

	return nil, co	
}

func init_tun() {
	// 1. Check if the `TUN` module is installed
	err__run_lsmod, lsmod := run_cmd(nil, "lsmod", nil)
	if err__run_lsmod != nil {
		log.Printf("[ERROR]: %s\n", err__run_lsmod)
	}

	// pipe in the output of lsmod to `grep tun`
	err__run_grep, grep_tun := run_cmd(lsmod, "grep", []string{"tun"}) 
	if err__run_grep != nil {
		log.Fatalf("[ERROR]: %s\n", err__run_grep)
	}

	// check if there was any output `lsmod | grep tun`
	if len(string(grep_tun)) == 0 {
		log.Fatalf("[ERROR]: the tun module is not loaded.\n")
	}
	log.Printf("[DEBUG]: tun module is loaded.\n")
	
	// 2. Create the TUN device
	err__create_device, _ := run_cmd(nil, "ip", []string{"tuntap", "add", "mode", "tun", "name", TUN_DEV_NAME})
	if err__create_device != nil {
		log.Fatalf("[ERROR]: %s\n", err__create_device)
	}

	log.Printf("[DEBUG]: created tun interface [%s] successfully.\n", TUN_DEV_NAME)

	// 3. Bring the device 'UP'
	err__up_tun_device, _ := run_cmd(nil, "ip", []string{"link", "set", TUN_DEV_NAME, "up"})
	if err__up_tun_device != nil {
		log.Fatalf("[ERROR]: %s\n", err__up_tun_device)
	}

	log.Printf("[DEBUG]: tun device %s is UP.\n", TUN_DEV_NAME)

	// 4. Assign an IP address and a Subnet Mask to the interface
	err__assign_ip,_ := run_cmd(nil, "ip", []string{"address", "add", "171.8.8.8/8", "dev", TUN_DEV_NAME})
	if err__assign_ip != nil {
		log.Fatalf("[ERROR]: %s\n", err__assign_ip)
	}

	log.Printf("[DEBUG]: assigned ip & subnet 172.8.8.8/8 to tun device [%s]\n", TUN_DEV_NAME)
}

func main() {
	err__check_tun, _ := run_cmd(nil, "ip", []string{"link", "show", "dev", TUN_DEV_NAME})
	if err__check_tun != nil {
		log.Printf("[DEBUG]: %s\n", err__check_tun)
		log.Printf("[DEBUG]: creating TUN device [%s]\n", TUN_DEV_NAME)
		init_tun()
	}

	log.Printf("[DEBUG]: tun device [%s] is created.\n", TUN_DEV_NAME)
}
