package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type bymin struct {
	Requests int64 `json:"requests"`
	total    int64 `json:"total"`
	Avg      int64 `json:"avg"`
	Longest  int64 `json:"longest"`
	Shortest int64 `json:"shortest"`
}

type stat struct {
	sync.RWMutex
	data map[string]bymin
}

var stats stat

func main() {
	stats = stat{
		data: make(map[string]bymin),
	}

	http.HandleFunc("/stat", func(rw http.ResponseWriter, req *http.Request) {
		log.Println("received request from ", req.RemoteAddr, "at ", req.RequestURI)
		stats.RLock()
		defer stats.RUnlock()
		bt, err := json.Marshal(stats.data)
		if err != nil {
			log.Println(err)
		}
		rw.Write(bt)
		return
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		s := time.Now()
		log.Println("received request from ", req.RemoteAddr, "at ", req.RequestURI)
		rw.Write([]byte("received request from " + req.RemoteAddr + "\n"))
		log.Println("acquiring lock")
		stats.Lock()
		defer stats.Unlock()
		v, k := stats.data[s.Format("02-01-2006 15:04")]
		lat := time.Now().Sub(s).Nanoseconds()
		if !k {
			stats.data[s.Format("02-01-2006 15:04")] = bymin{
				Requests: 1,
				total:    lat,
				Avg:      lat,
				Longest:  lat,
				Shortest: lat,
			}
			return
		}
		v.Requests++
		v.total += lat
		v.Avg = (v.total / v.Requests)
		if lat > v.Longest {
			v.Longest = lat
		}
		if lat < v.Shortest {
			v.Shortest = lat
		}
		stats.data[s.Format("02-01-2006 15:04")] = v
		return
	})

	
	http.HandleFunc("/createfile", func(rw http.ResponseWriter, req *http.Request) {
	    if req.Method == http.MethodPost {
	      log.Println("received request from ", req.RemoteAddr, "at ", req.RequestURI)
	      hostname, err := os.Hostname()
	      if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	      }

	      file, err := os.OpenFile("/data/"+hostname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	      if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	      }

	      defer file.Close()

	      _, err = file.WriteString("data written from " + hostname + " at " + time.Now().String() + "\n")
	      if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Internal error:"+err.Error(), http.StatusInternalServerError)
		return
	      }
	      return

	    }
	    http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	    return
	  })
	
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
