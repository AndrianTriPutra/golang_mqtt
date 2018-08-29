//referensi https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
package main

import (
	"bytes"
	"log"
	"fmt"
	"os/exec"
	"time"
)

var data string
var backup bytes.Buffer
var status string

func main() {
	ticker := time.NewTicker(10 * time.Second)

	for{
		select {
			case <-ticker.C:
				cmd := exec.Command("ping", "-c 2", "-i 1", "8.8.8.8")
				var buf bytes.Buffer
				cmd.Stdout = &buf				
				cmd.Start()			
				done := make(chan error)
				go func() { done <- cmd.Wait() }()
				timeout := time.After(2 * time.Second)
		
				select {
				case <-timeout:
					cmd.Process.Kill()
					status = "Time Out"
				case err := <-done:					
					if err == nil {
						status = "UP"
					}else{
						status = "DOWN"
					}
		
				}

				TimeUTC := time.Now().UTC()
				hhmmss := fmt.Sprintf("%02d:%02d:%02d", TimeUTC.Hour(), TimeUTC.Minute(), TimeUTC.Second())
				yymmdd := fmt.Sprintf("%02d-%02d-%02d", TimeUTC.Year()%100, TimeUTC.Month(), TimeUTC.Day())
				
				year = yymmdd[0:2]
				month = yymmdd[3:5]
				date = yymmdd[6:8]
			
				hour = hhmmss[0:2]
				minute = hhmmss[3:5]
				second = hhmmss[6:8]

				backup.Reset()
				backup.WriteString(yymmdd)
				backup.WriteString(",")
				backup.WriteString(hhmmss)
				backup.WriteString(",")
				backup.WriteString(status)
				data = backup.String()

				log.Printf("data:%s",data)
		}		
	}

}
