package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
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
	return nil, fmt.Errorf("uygun ag arayuzu bulunamadi")
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
	showVendor := flag.Bool("vendor", false, "MAC uretici bilgisini goster")
	rangeFlag := flag.String("range", "192.168.1.0/24", "Taranacak IP araligi")
	help := flag.Bool("h", false, "Kullanim bilgilerini goster")
	flag.Parse()

	if *help {
		fmt.Println("Kullanim: ./arp_scan [secenekler]")
		fmt.Println("Secenekler:")
		fmt.Println("  -vendor  MAC uretici bilgisini goster")
		fmt.Println("  -range   Taranacak IP araligini belirle (orn: 192.168.1.0/24)")
		fmt.Println("  -h       Kullanim bilgilerini goster")
		os.Exit(0)
	}

	//ag arayuzu bul
	iface, err := getDefaultInterface()
	if err != nil {
		log.Fatalf("Ag arayuzu bulunamadi: %v", err)
	}

	fmt.Printf("Kullanilan arayuz: %s\n", iface.Name)

	// Packet Listener olustur
	packetConn, err := net.ListenPacket("ethernet", iface.Name)
	if err != nil {
		log.Fatalf("Packet listener oluşturulamadı: %v", err)
	}
	defer packetConn.Close()

	// ARP istemcisi olustur
	client, err := arp.New(iface, packetConn)
	if err != nil {
		log.Fatalf("ARP istemcisi oluşturulamadı: %v", err)
	}
	defer client.Close()

	// Kullanici gelen IP araligi
	ips, err := parseIPRange(*rangeFlag)
	if err != nil {
		log.Fatalf("Geçersiz IP aralığı: %v", err)
	}

	// IP araligi
	for _, ip := range ips {
		// IP'yi `netip.Addr` formatina cevir
		parsedAddr, err := netip.ParseAddr(ip)
		if err != nil {
			log.Printf("Geçersiz IP adresi: %s\n", ip)
			continue
		}

		// ARP istegi
		mac, err := client.Resolve(parsedAddr)
		if err != nil {
			continue
		}

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
