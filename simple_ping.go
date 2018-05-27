package main

import (
	"fmt"
	"time"
	"strings"
	"os/exec"

)
func main(){
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			clock, _ := exec.Command("sh", "-c", "date +\"%Y-%m-%d %H:%M:%S %Z\"").Output()
			buffer := string(clock)
			buffer=strings.TrimSuffix(buffer,"\n")
			fmt.Printf("%s\n",buffer)

			out, _ := exec.Command("ping", "www.google.com", "-c 1", "-i 1", "-w 1").Output()//jika paketnya mau nambah tinggal ganti -c 1 menjadi  2 atau 3 atau nilai lainya
			if strings.Contains(string(out), "Destination Host Unreachable") {
				fmt.Println("TANGO DOWN")
			} else {
				fmt.Println("IT'S ALIVEEE")
				fmt.Printf("out:%s\n",out)
			}
		}
	}

}
