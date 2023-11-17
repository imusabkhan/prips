package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [CIDR] [output_path_optional]\n\n", os.Args[0])
		// flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n%s 192.168.0.0/24 /tmp/output.txt\n", os.Args[0])
	}
	flag.Parse()

	var cidr string
	switch flag.NArg() {
	case 0:
		// Read from standard input
		_, err := fmt.Scanln(&cidr)
		if err != nil {
			fmt.Println("Error reading CIDR input:", err)
			return
		}
	case 1:
		// Read CIDR from command-line argument
		cidr = flag.Arg(0)
	case 2:
		// Read CIDR and output file from command-line arguments
		cidr = flag.Arg(0)
		outputPath := flag.Arg(1)
		ipList, err := expandCIDR(cidr)
		if err != nil {
			fmt.Println("Error expanding CIDR:", err)
			return
		}
		err = printAndSaveToFile(outputPath, ipList)
		if err != nil {
			fmt.Println("Error printing and saving IPs:", err)
			return
		}
		return
	default:
		flag.Usage()
		return
	}

	ipList, err := expandCIDR(cidr)
	if err != nil {
		fmt.Println("Error expanding CIDR:", err)
		return
	}

	// Print IPs to standard output
	for _, ip := range ipList {
		fmt.Println(ip)
	}
}

func expandCIDR(cidr string) ([]string, error) {
	ipList := []string{}
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ipList = append(ipList, ip.String())
	}

	// Remove network address and broadcast address
	ipList = ipList[1 : len(ipList)-1]

	return ipList, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func printAndSaveToFile(filePath string, ips []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, ip := range ips {
		fmt.Println(ip)
		_, err := file.WriteString(ip + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
