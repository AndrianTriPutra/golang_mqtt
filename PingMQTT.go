/*
 *golang
 * ping google and send data to server with mqtt
 * created ATP
*/
package main

import (
	"fmt"
	"time"
	"bytes"
	"strconv"	
	"strings"
	"os/exec"
	"github.com/sparrc/go-ping"	
	"github.com/eclipse/paho.mqtt.golang"
)

var paket string
var Clock string
var bufferpaket string
var DataRTS string

var opts *mqtt.ClientOptions
var buffer bytes.Buffer

func main() {		
	go Ping()
	go clock()
	MQTT()
	fmt.Scanf("%f\n")
}

func Ping(){
	var state_network bool = true
	
	permission, _ := exec.Command("sh", "-c", "sudo sysctl -w net.ipv4.ping_group_range="+"0  "+"2147483647").Output()
	fmt.Printf("%s",permission)	

	for{//for
		pinger, err := ping.NewPinger("www.google.com")
		if err != nil {
			if state_network{			
				paket = "DOWN"
				fmt.Printf("paket:%s",paket)
				fmt.Println()
				state_network=false
			}			
		}else{
			pinger.OnRecv = func(pkt *ping.Packet) {
			buffer.Reset()				
			buffer.WriteString("UP")		
			buffer.WriteString(",")
					
			//paket
			s := strconv.Itoa(pkt.Nbytes)
			buffer.WriteString(s)
			buffer.WriteString(",")
								
			Rtt := int64(pkt.Rtt/time.Millisecond) 
			t := strconv.FormatInt(Rtt,10)
			buffer.WriteString(t)		
					
			//statistic
			stats   	 := pinger.Statistics()
					
			PacketsSent  := stats.PacketsSent
			u := strconv.Itoa(PacketsSent)
			buffer.WriteString(",")
			buffer.WriteString(u)
					
			PacketsRecv  := stats.PacketsRecv
			v := strconv.Itoa(PacketsRecv)
			buffer.WriteString(",")
			buffer.WriteString(v)
					
			PacketLoss  := stats.PacketLoss
			w := strconv.FormatFloat(PacketLoss, 'f',2,64)
			buffer.WriteString(",")
			buffer.WriteString(w)					
			paket = buffer.String()
			}

			fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
			pinger.Run()			
		}		

	}//for
}

func MQTT(){
	ticker := time.NewTicker(1 * time.Second)
    quit := make(chan struct{})
    
    go func() {
		opts = mqtt.NewClientOptions().AddBroker("tcp://test.mosquitto.org:1883").SetClientID("ATP").SetMaxReconnectInterval(1 * time.Minute).SetAutoReconnect(true)	
		c := mqtt.NewClient(opts)
		token := c.Connect()

		if token.Wait() && token.Error() != nil {
			fmt.Println("error")
		} else {
			fmt.Println("Connected")
		}
        for {
            select {
            case <-ticker.C:
				buffer.Reset()				
				buffer.WriteString(Clock)
				buffer.WriteString("\t")    
				if bufferpaket != paket && paket!="DOWN" {		
					bufferpaket = paket
					buffer.WriteString(paket)   
				}else{
					paket="DOWN"
					buffer.WriteString(paket) 
				}   			 
				DataRTS=buffer.String()   
				token2 := c.Publish("BIR/ATP",0,false,DataRTS)
				if token2.Wait() && token2.Error() != nil {
					fmt.Println("error")
				} else {
					fmt.Printf("Send paket:%s\n",DataRTS)
				} 
            case <-quit:
                ticker.Stop()
            }
        }
    }()	
}

func clock() {
	 for{
		out, _ := exec.Command("sh", "-c", "date +\"%Y-%m-%d %H:%M:%S %Z\"").Output()
		buffer := string(out)
		buffer=strings.TrimSuffix(buffer,"\n")
		Clock=buffer	 
	 }
}


