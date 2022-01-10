package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func AttackByRate(attacker vegeta.Attacker, targeter vegeta.Targeter, attackRate int, duration time.Duration) vegeta.Metrics {
	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: attackRate, Per: time.Second}
	for res := range attacker.Attack(targeter, rate, duration, "Big Ban") {
		metrics.Add(res)
	}
	metrics.Close()

	return metrics
}

func binarrySearch(attacker vegeta.Attacker, targeter vegeta.Targeter, maxRate int, duration time.Duration, delay time.Duration) vegeta.Metrics {

	var metrics vegeta.Metrics
	minRate := 0
	midRate := (maxRate + minRate) / 2
	for {
		if (maxRate - minRate) < 10 {
			break
		}

		metrics = AttackByRate(attacker, targeter, midRate, duration)
		if *debug {
			fmt.Printf("Max: %d, Min: %d, Mid: %d, success: %f\n", maxRate, minRate, midRate, metrics.Success)
		}

		if metrics.Success == 1.0 {
			minRate = midRate
		} else {
			maxRate = midRate
		}
		time.Sleep(delay)
		midRate = (maxRate + minRate) / 2
	}

	return metrics
}

var (
	duration = flag.Duration("duration", 10*time.Second, "Benchmark duration")
	delay    = flag.Duration("delay", 10*time.Second, "Benchmark interval")
	maxRate  = flag.Int("rate", 10000, "Max Benchmark Rate(rps)")
	insecure = flag.Bool("insecure", false, "Ignore TLS errors")
	debug    = flag.Bool("debug", false, "Output Debug Logs")
	output   = flag.String("output", "text", "Output format [text, json]")
)

func formatText(metrics vegeta.Metrics) {
	fmt.Printf("Requests [total, rate, throughput]: %d, %f, %f\n", metrics.Requests, metrics.Rate, metrics.Throughput)
	fmt.Printf("Durationt [total, attack, wait]: %.2fs, %.2fs, %.2fs\n", (metrics.Duration + metrics.Wait).Seconds(), metrics.Duration.Seconds(), metrics.Wait.Seconds())
	fmt.Printf("Success [ratio] %.2f%%\n", metrics.Success*100)
}

func formatJson(metrics vegeta.Metrics) {
	bytes, _ := json.Marshal(&metrics)
	fmt.Println(string(bytes))
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Need host for target")
		return
	}
	host := flag.Args()[0]
	_, err := url.ParseRequestURI(host)
	if err != nil {
		panic(err)
	}

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    host,
	})

	tlsc := tls.Config{InsecureSkipVerify: *insecure}
	attacker := vegeta.NewAttacker(vegeta.TLSConfig(&tlsc))

	metrics := binarrySearch(*attacker, targeter, *maxRate, *duration, *delay)
	if *output == "json" {
		formatJson(metrics)
	} else {
		formatText(metrics)
	}
}
