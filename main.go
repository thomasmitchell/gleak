package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//It's not thread safe but... who cares?
var blackhole []uint8

// HelloHandler says hello and gives you the loop that the program is on
type HelloHandler struct {
	//Basically just here as a vessel for ListenAndServe
}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := uint8(random.Uint64())
	w.Write([]byte(fmt.Sprintf("Hello world: index %d is %d", i, blackhole[i])))
}

func useResources(bytes uint64) {
	blackhole = make([]uint8, bytes)

	for i := uint64(0); i < bytes; i++ {
		blackhole[i] = uint8(i)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//Shuffle it to use CPU
	for {
		for curIndex := uint64(0); curIndex < bytes; curIndex++ {
			swapWith := r.Uint64() % (bytes - curIndex)
			blackhole[curIndex], blackhole[swapWith] = blackhole[swapWith], blackhole[curIndex]
		}
	}
}

func main() {
	fmt.Printf("PORT: %s\n", os.Getenv("PORT"))
	var numMegabytes uint64 = 1024 //Default
	var err error
	if os.Getenv("MEGABYTES") != "" {
		numMegabytes, err = strconv.ParseUint(os.Getenv("MEGABYTES"), 10, 64)
		if err != nil {
			fmt.Printf("MEGABYTES envvar could not be parsed as uint: %s\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("MEGABYTES: %d\n", numMegabytes)

	numBytes := 1024 * 1024 * numMegabytes
	go useResources(numBytes)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), HelloHandler{})
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
