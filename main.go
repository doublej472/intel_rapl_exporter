package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"io/ioutil"
	"strings"
	"flag"
)

func main() {
	port := flag.Int("port", 9999, "Port to listen on")
	address := flag.String("address", "0.0.0.0", "Address to listen on")
	flag.Parse()

	listenAddr := fmt.Sprint(*address, ":", *port)

	http.HandleFunc("/metrics", metrics)
	fmt.Println("listening on", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic(err)
	}
}

func metrics(res http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir("/sys/class/powercap/")
    	if err != nil {
    		panic(err)
    	}
	io.WriteString(res, "# HELP power_intel_rapl_energy_uj Total microwatts used since boot\n# TYPE power_intel_rapl_energy_uj counter\n" )

	for _, file := range files {
		if (!strings.Contains(file.Name(), ":")) {
			continue
		}
		dat, err := os.ReadFile("/sys/class/powercap/" + file.Name() + "/energy_uj")
		powerstring := strings.TrimRight(string(dat), "\r\n")
    		if err != nil {
    			panic(err)
    		}

		namedat, err := os.ReadFile("/sys/class/powercap/" + file.Name() + "/name")
		powername := strings.TrimRight(string(namedat), "\r\n")
    		if err != nil {
    			panic(err)
    		}

		io.WriteString(res, "power_intel_rapl_energy_uj{name=\"" + powername +"\", path=\"" + file.Name() + "\"} " + powerstring + "\n" )
    	}
}
