# ARP Scan - MAC Adresi Taraması
Bu Go programı, belirli bir IP aralığında ARP taraması yaparak, her cihazın MAC adresini tespit eder. Ayrıca, isteğe bağlı olarak, her MAC adresinin üretici bilgisini (vendor) de sorgular ve ekrana yazdırır. Program, ağınızdaki aktif cihazları kolayca tespit etmenize olanak tanır.

### Özellikler
ARP Taraması: Verilen IP aralığındaki cihazların MAC adreslerini çözümler.
MAC Vendor Bilgisi: -vendor bayrağı ile MAC adresinin üretici bilgisini sorgular.
Esnek IP Aralığı: IP aralığını komut satırından belirtme imkânı sunar.
Yardım Bayrağı: Kullanıcıya komut satırında kullanım talimatlarını gösterir.
Kullanım
Program, terminal üzerinden çalıştırılabilir ve çeşitli bayraklarla kontrol edilebilir.

### Komut Yapısı:

./arp_scan -h 
Seçenekler:
1. -range <IP aralığı>: Taranacak IP aralığını belirtir. Varsayılan olarak 192.168.1.0/24 kullanılır.
2. -vendor: MAC adresinin üretici bilgilerini gösterir. Bu bayrak belirtilmezse, sadece IP ve MAC adresi yazdırılır.
3. -h: Kullanım bilgilerini gösterir.
Örnek Kullanımlar:
Sadece IP ve MAC Adreslerini Gösterme: Belirtilen IP aralığındaki cihazların IP ve MAC adreslerini görmek için:


```
./arp_scan -range 192.168.1.0/24
Bu komut, IP aralığındaki her cihazın IP ve MAC adreslerini listeler.
```
MAC Vendor Bilgisiyle Gösterme: Eğer her cihazın MAC adresinin üretici bilgisini de görmek isterseniz:

```

./arp_scan -range 192.168.1.0/24 -vendor
Bu komut, her cihazın IP ve MAC adresinin yanı sıra, MAC adresinin üreticisini (vendor) de listeler.
```


./arp_scan -h
Bu komut, programın nasıl kullanılacağı hakkında bilgi verir.
```
### Örnek Çıktılar:
MAC Vendor Bilgisi ile Çıktı:

```

IP: 192.168.1.1 - MAC: 00:11:22:33:44:55 - Üretici: Cisco Systems
IP: 192.168.1.2 - MAC: 66:77:88:99:AA:BB - Üretici: TP-Link Technologies

```
```
IP: 192.168.1.1 - MAC: 00:11:22:33:44:55
IP: 192.168.1.2 - MAC: 66:77:88:99:AA:BB
```
### Kurulum

```
go get github.com/mdlayher/arp

go build arp_scan.go

./arp_scan -range 192.168.1.0/24 -vendor
```
