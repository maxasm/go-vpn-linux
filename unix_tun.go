package main

import (
	"golang.org/x/sys/unix"
)

const (
	// points to the TUN device location
	TUN_DEVICE = "/dev/net/tun"
	// This tells the Kernel that we are working with a TUN device and we dont need any additional packet information.
	TUN_TYPE = unix.IFF_TUN | unix.IFF_NO_PI
	// IOCTL access code. Basically referes to TUN/TAP devices
	IOCTL_ACCESS_CODE = unix.TUNSETIFF
)

// use `epoll` to read from the file descriptor
func epoll_fd(fd int) {
	// create an new `epoll` instance given by the kernel
	ep_fd, err__get_epoll := unix.EpollCreate(1)
	if err__get_epoll != nil {
		dl.Fatalf("failed to get a new epoll instance.\n")
	}

	dl.Printf("created a new epoll instance.\n")
	
	// create an epoll event instance used to add the fds you want to read from
	// and the events you want to `listen` to
	// in this case, I am adding the fd for the tun interface
	tun_fd_epoll_event := unix.EpollEvent{
		Events: unix.EPOLLIN, // when there is data to read from the fd
		Fd: int32(fd), // the fd that we are interested in reading from
	}

	// use EpollCtl to add the event to the instance.
	// This  system  call  is  used to add, modify, or remove entries in the interest list of the epoll(7) instance referred to by the file descriptor epfd.
	// It requests that the operation op be performed for the target file descriptor fd.
	err__add_fd := unix.EpollCtl(ep_fd, unix.EPOLL_CTL_ADD, fd, &tun_fd_epoll_event)
	if err__add_fd != nil {
		dl.Fatalf("failed to add a new fd to the Epoll interest list using epoll_ctl.\n")
	}

	dl.Printf("successfully added tun_fd to the interest list of the Epoll fd.\n")

	// create the `Event Store` where all triggered events are stored.
	// since we have only one fd that we are watching, MAX_EVENTS = 1
	MAX_EVENTS := 1
	event_store := make([]unix.EpollEvent, MAX_EVENTS)

	dl.Printf("waiting for data from tun device ...\n")

	for {
		n_events, err__wait := unix.EpollWait(ep_fd, event_store, -1)
		if err__wait != nil {
			dl.Fatalf("failed to wait for events on tun device\n")
		}

		dl.Printf("got %d events.\n", n_events)

		rfd := int(event_store[0].Fd)

		buffer := make([]byte, 1024)
		n_bytes, err__read := unix.Read(rfd, buffer)
		if err__read != nil {
			dl.Fatalf("failed to read data from the tun_device\n")
		}

		dl.Printf("EPOLL: read %d bytes from /dev/net/tun/\n", n_bytes)
	}
}

func open_tun_device() {
	// open the TUN device in nonblocking mode
	tun_fd, err__open_tun := unix.Open(TUN_DEVICE, unix.O_CREAT|unix.O_RDWR|unix.O_CLOEXEC|unix.O_NONBLOCK, 0644)
	if err__open_tun != nil {
		dl.Fatalf("failed to open TUN device @ %s.\n", TUN_DEVICE)
	}

	// NewIfreq creates an Ifreq with the input network interface name after validating the name does not exceed IFNAMSIZ-1 (trailing NULL required) bytes.
	ifreq, err__get_ifreq := unix.NewIfreq(TUN_DEV_NAME)
	if err__get_ifreq != nil {
		dl.Fatalf("failed to create a new IFREQ\n")
	}

	// set the flags
	// SetUint16 sets a C short/Go uint16 value as the Ifreq's union data
	ifreq.SetUint16(uint16(TUN_TYPE))
	
	// make the IOCTL unix to set configure the TUN device
	// IoctlIfreq performs an ioctl using an Ifreq structure for input and/or output. See the netdevice(7) man page for details.
	err__run_ioctl := unix.IoctlIfreq(tun_fd, IOCTL_ACCESS_CODE, ifreq)
	if err__run_ioctl != nil {
		dl.Fatalf("failed to run IOCTL#IFREQ\n")
	}

	dl.Printf("successfully run IOCTL#IFREQ\n")
	epoll_fd(tun_fd)
}
