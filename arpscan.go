package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mdlayher/arp"
)

type MacVendorResponse struct {
	Vendor string `json:"vendor"`
}

func getMacVendor(mac string) (string, error) {
	url := fmt.Sprintf("https://api.macvendors.com/%s", mac)
	resp, err := http.Get(url)
	if err != nil {
		return "Bilinmiyor", err
	}
	defer resp.Body.Close()

	var result MacVendorResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Bilinmiyor", err
	}

	return result.Vendor, nil
}

func getDefaultInterface() (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("uygun ağ arayüzü bulunamadi")
}

func parseIPRange(rangeStr string) ([]string, error) {
	parts := strings.Split(rangeStr, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("gecersiz IP araligi formati")
	}

	base := strings.Join(parts[:3], ".")
	var ips []string
	for i := 1; i <= 254; i++ {
		ips = append(ips, fmt.Sprintf("%s.%d", base, i))
	}

	return ips, nil
}

func main() {
	// Flag
	rangeFlag := flag.String("range", "192.168.1.0/24", "Taranacak IP araligi belirle (örn: 192.168.1.0/24)")
	showVendor := flag.Bool("vendor", false, "MAC üretici bilgisini göster")
	help := flag.Bool("h", false, "Kullanim bilgilerini göster")
	flag.Parse()

	// Yardım
	if *help {
		fmt.Println("Kullanim: ./arp_scan [seçenekler]")
		fmt.Println("Seçenekler:")
		fmt.Println("  -vendor  MAC üretici bilgisini göster")
		fmt.Println("  -range   Taranacak IP araligi belirle (örn: 192.168.1.0/24)")
		fmt.Println("  -h       Kullanim bilgilerini göster")
		os.Exit(0)
	}

	// Ağ arayüzü
	iface, err := getDefaultInterface()
	if err != nil {
		log.Fatalf("Ag arayuzu bulunamadi: %v", err)
	}
	fmt.Printf("Kullanilan arayuz: %s\n", iface.Name)

	// ARP istemcisi
	client, err := arp.New(iface)
	if err != nil {
		log.Fatalf("ARP istemcisi olusturulamadi: %v", err)
	}
	defer client.Close()

	// IP araligi parse
	ips, err := parseIPRange(*rangeFlag)
	if err != nil {
		log.Fatalf("Gecersiz IP araligi: %v", err)
	}

	// IP araligi
	for _, ip := range ips {
		parsedIP := net.ParseIP(ip)
		mac, err := client.Resolve(parsedIP)
		if err != nil {
			continue
		}

		//MAC adresi ve IP araligi
		if *showVendor {
			vendor, err := getMacVendor(mac.String())
			if err != nil {
				vendor = "Bilinmiyor"
			}
			fmt.Printf("IP: %s - MAC: %s - Üretici: %s\n", ip, mac, vendor)
		} else {
			fmt.Printf("IP: %s - MAC: %s\n", ip, mac)
		}

		time.Sleep(50 * time.Millisecond)
	}
}
