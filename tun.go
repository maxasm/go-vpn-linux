package main

// the name of the `TUN` device being made
const TUN_DEV_NAME = "tunmax"

// This function creates a new `TUN` device with the name TUN_DEV_NAME
// and assignes it the IP 172.8.8.8/8
func init_tun() {
	// 1. Check if the `TUN` module is installed
	err__run_lsmod, lsmod := run_cmd(nil, "lsmod", nil)
	if err__run_lsmod != nil {
		dl.Fatalf("error running lsmod. %s\n", err__run_lsmod)
	}

	// pipe in the output of lsmod to `grep tun`
	err__run_grep, grep_tun := run_cmd(lsmod, "grep", []string{"tun"}) 
	if err__run_grep != nil {
		dl.Fatalf("error runnig grep command.%s\n", err__run_grep)
	}

	// check if there was any output `lsmod | grep tun`
	if len(string(grep_tun)) == 0 {
		dl.Fatalf("the tun module is not loaded.\n")
	}
	dl.Printf("tun module is loaded.\n")
	
	// 2. Create the TUN device
	err__create_device, _ := run_cmd(nil, "ip", []string{"tuntap", "add", "mode", "tun", "name", TUN_DEV_NAME})
	if err__create_device != nil {
		dl.Fatalf("%s\n", err__create_device)
	}

	dl.Printf("created tun interface [%s] successfully.\n", TUN_DEV_NAME)
	// 3. Bring the device 'UP'
	err__up_tun_device, _ := run_cmd(nil, "ip", []string{"link", "set", TUN_DEV_NAME, "up"})
	if err__up_tun_device != nil {
		dl.Fatalf("%s\n", err__up_tun_device)
	}

	dl.Printf("tun device %s is UP.\n", TUN_DEV_NAME)

	// 4. Assign an IP address and a Subnet Mask to the interface
	err__assign_ip,_ := run_cmd(nil, "ip", []string{"address", "add", "171.8.8.8/8", "dev", TUN_DEV_NAME})
	if err__assign_ip != nil {
		dl.Fatalf("%s\n", err__assign_ip)
	}

	dl.Printf("assigned ip & subnet 172.8.8.8/8 to tun device [%s]\n", TUN_DEV_NAME)
}
