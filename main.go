package main

import (
	"bufio"
	"fmt"
	"game-lag-watcher/traceroute"
	"github.com/sparrc/go-ping"
	"log"
	"os"
	"time"
)

func main() {
	hosts := make([]string, 0)
	file, err := os.Open("hosts.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text()
		if len(host) < 3 {
			continue
		}
		hosts = append(hosts, host)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Hosts defined: %s", hosts)

	// print the time and traceroute command output
	logFile, err := os.OpenFile("traceroute.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()

	// put parsed output into csv
	csvFile, err := os.OpenFile("traceroute.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	// put parsed output into csv
	pingCsvFile, err := os.OpenFile("pings.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer pingCsvFile.Close()

	for {
		for _, host := range hosts {
			log.Printf("Running ping for host %s", host)
			pinger, err := ping.NewPinger(host)
			if err != nil {
				log.Println(err)
				continue
			}
			pinger.SetPrivileged(true)
			pinger.Count = 3
			pinger.Run() // blocks until finished
			stats := pinger.Statistics() // get send/receive/rtt stats
			log.Printf("ping %+v", stats)

			pingRow := fmt.Sprintf("%s, %s, %s, %s, %s, %s\n", time.Now().Format(time.RFC1123), host, stats.MinRtt, stats.MaxRtt, stats.AvgRtt, stats.StdDevRtt)
			if _, err := pingCsvFile.WriteString(pingRow); err != nil {
				log.Println(err)
			}
			log.Printf("Finished ping for host %s", host)

			log.Printf("Running tracert for host %s", host)
			res, err := traceroute.Run(host)
			if err != nil {
				log.Println(err)
				continue
			}

			if _, err := logFile.WriteString(fmt.Sprintf("\n%s \n %s\n", time.Now(), res.Output)); err != nil {
				log.Println(err)
			}

			for _, v := range res.Hops {
				row := fmt.Sprintf("%s, %s, %d, %s, %s, %s, %s\n", time.Now().Format(time.RFC1123), host, v.Id, v.T1, v.T2, v.T3, v.Dest)
				if _, err := csvFile.WriteString(row); err != nil {
					log.Println(err)
				}
			}
			log.Printf("Finished tracert for host %s", host)
		}

		<- time.After(5 * time.Minute)
	}
}
