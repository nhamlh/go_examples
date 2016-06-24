/*
Purpose: simple ipcalc cli command written in golang
Author : NhamLH <lehoainham@gmail.com>
Date   : Jun 24 2016
*/

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"math"
)

type InvalidFormat struct {
}

func (i *InvalidFormat) Error() string {
	return "Invalid IP format."
}

type IP struct {
	Addr string
	Net  int
}

func (ip *IP) _netmaskToIntArray() [32]int {
	var a [32]int
	for i := 0; i < ip.Net; i++ {
		a[i] = 1
	}
	return a
}

func (ip *IP) _netaddrToIntArray() [32]int {
	net := ip._netmaskToIntArray()

	netaddr := ip._addrToArray()
	// compare ip and netmask to get network address
	for i := 31; i >= 0; i-- {
		if net[i] == 0 {
			netaddr[i] = 0
		}
	}

	return netaddr
}

func (ip *IP) _addrToArray() [32]int {
	var octet [4]int

	// devide ip into octets of int
	re := regexp.MustCompile("([0-9]+)\\.([0-9]+)\\.([0-9]+)\\.([0-9]+)")
	match := re.FindStringSubmatch(ip.Addr)

	octet[0], _ = strconv.Atoi(match[1])
	octet[1], _ = strconv.Atoi(match[2])
	octet[2], _ = strconv.Atoi(match[3])
	octet[3], _ = strconv.Atoi(match[4])

	// convert octet to string
	var bin_octet string
	for i := 0; i < len(octet); i++ {
		bin_octet += IntToBin(octet[i])
	}

	// convert octet string to int array
	var addr [32]int
	for i := 0; i < len(bin_octet); i++ {
		if bin_octet[i] == 49 {
			addr[i] = 1
		} else {
			addr[i] = 0
		}
	}

	return addr
}

// return a dot-decimal string of netmask
func (ip *IP) NetmaskToDot() string {
	return IntArrayToDot(ip._netmaskToIntArray())
}

func (ip *IP) NetaddrToDot() string {
	return IntArrayToDot(ip._netaddrToIntArray())
}

func (ip *IP) AddrToDot() string {
	return IntArrayToDot(ip._addrToArray())
}

func (ip *IP) NetmaskToCIDR() string {
	return BinToCIDR(ip.NetmaskToDot())
}

func (ip *IP) NetaddrToCIDR() string {
	return BinToCIDR(ip.NetaddrToDot())
}

func NewIP(s string) (IP, error) {
	var ip IP
	re := regexp.MustCompile("([0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+)/([0-9]+)")
	match := re.FindStringSubmatch(s)

	if len(match) == 0 {
		return ip, &InvalidFormat{}
	}

	addr := match[1]
	net, _ := strconv.Atoi(match[2])

	if net < 0 || net > 32 {
		return ip, &InvalidFormat{}
	}

	ip = IP{addr, net}

	return ip, nil
}

// Helper func to return a string represent binary array of an integer
func IntToBin(i int) string {
	s := strconv.FormatInt(int64(i), 2)

	for len(s) < 8 {
		s = "0" + s
	}
	return s
}

// Helper func to return an integer from a string
// ex: 10000000 -> 128, 11000001 -> 193
func BinToInt(s string) int {
	var n int
	for i:=0; i<len(s); i++ {
		if s[i] == 49 {
			n += int(math.Pow(2, float64(7-i)))
		}
	}
	return n
}

// Helper func to return an cidr from a string
// ex: 11000000.10101000.00000001.00000000 => 192.168.1.0
func BinToCIDR(s string) string {

	re := regexp.MustCompile("([0-1]+)\\.([0-1]+)\\.([0-1]+)\\.([0-1]+)")
	match := re.FindStringSubmatch(s)

	var result string
	for i:=1; i<len(match); i++ {
		result += strconv.Itoa(BinToInt(match[i]))

		if i< len(match)-1 {
			result += "."
		}
	}

	return result
}

// Helper func to return a dot-decimal representation of a int array
// This will work when you have converted netmask or net address to int array
func IntArrayToDot(a [32]int) string {
	var s string

	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			bit := strconv.Itoa(a[j+(i)*8])
			s += bit
		}
		if i != 3 {
			s += "."
		}
	}
	return s
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ipcalc <IP>/<NETMASK>")
	}
	raw_ip := os.Args[1]

	ip, error := NewIP(raw_ip)
	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("IP Address      ", ip.AddrToDot())
	fmt.Println("Netmask         ", ip.NetmaskToDot())
	fmt.Println("Net Address     ", ip.NetaddrToDot())
	fmt.Println("Netmask         ", ip.NetmaskToCIDR())
	fmt.Println("Net Address     ", ip.NetaddrToCIDR())
}
