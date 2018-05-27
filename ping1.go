package main

import (
	"fmt"
	"os/exec"
	"github.com/sparrc/go-ping"
)

func main() {
	permission, _ := exec.Command("sh", "-c", "sudo sysctl -w net.ipv4.ping_group_range="+"0  "+"2147483647").Output()
	fmt.Printf("%s",permission)	
	
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
			panic(err)
	}

	pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
					pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
					stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
					stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()//block until finish jadi kalo misalkan banyak function dan routine yang kamu gunakan sebaiknya pake count pinger.Count = 3 seperti example githubnya
}
