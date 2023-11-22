package main

import(
	"log"
	"os"
)

var dl *log.Logger = log.New(os.Stdout, "[DEBUG]: ", log.Lshortfile)

func main() {
	// check if the TUN device named TUN_DEV_NAME exists
	err__check_tun, _ := run_cmd(nil, "ip", []string{"link", "show", "dev", TUN_DEV_NAME})
	if err__check_tun != nil {
		dl.Printf("TUN device [%s] does not exist.\n", TUN_DEV_NAME)
		dl.Printf("creating TUN device [%s]\n", TUN_DEV_NAME)
		init_tun()
	}

	dl.Printf("tun device [%s] is created.\n", TUN_DEV_NAME)
	dl.Printf("Opening TUN device @ /dev/net/tun\n")

	open_tun_device()
}
