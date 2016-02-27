package main

import (
	"bytes"
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/config"
	"github.com/L7-MCPE/lav7/format"
	"github.com/L7-MCPE/lav7/gen"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/util"
)

func main() {

	cfg, err := os.Open("lav7.properties")
	if os.IsNotExist(err) {
		err := ioutil.WriteFile("lav7.properties", []byte(config.DefaultConfig), 0644)
		if err != nil {
			panic(err)
		}
		main()
	} else if err != nil {
		log.Fatalln("Error while opening properties file:", err)
	}

	ppr := flag.Bool("pprof", false, "starts pprof debug server on :8080")
	mutex := flag.Bool("mutex", false, "trace mutexes for debugging")
	flag.Parse()

	if *ppr {
		go http.ListenAndServe(":8080", nil)
		log.Println("debug: pprof server is running on :8080")
	}

	if *mutex {
		util.MutexDebug = true
		log.Println("debug: Mutex debugger is ON. Server may be slower.")
	} else {
		util.MutexDebug = false // I know a zero value for bool is false, but setting anyway for certainty
	}

	config.Parse(cfg)

	log.Printf("Starting lav7 version %s(git commit %s)", lav7.Version, lav7.GitCommit)
	log.Println("lav7 is licensed under the GPLv3 License; see http://rawgit.com/L7-MCPE/lav7/master/LICENSE.txt")
	log.Printf("Build timestamp: %s", lav7.BuildTime)

	start := time.Now()
	if runtime.NumCPU() < 2 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	initLevel(config.Generator, config.GeneratorArgs, config.Format)
	initRaknet()
	startLevel()
	startRouter(config.Port)

	log.Println("All done! Elapsed time:", time.Since(start).Seconds(), "seconds")
	log.Println("Server is ready. Type 'stop' to stop server.")
	lav7.HandleCommand()
}

func initLevel(genname string, arg string, lvformat string) {
	var g gen.Generator
	if lvformat == "img" {
		log.Print("* Using EXPERIMENTAL image chunk generator")
		file, err := os.Open(arg)
		if err != nil {
			log.Fatalln("Error while opening image:", err)
		}
		buf := new(bytes.Buffer)
		io.Copy(buf, file)
		bs := buf.Bytes()
		cfg, format, err := image.DecodeConfig(buf)
		if err != nil {
			log.Fatalln("Error while decoding image:", err)
		}
		log.Printf("* Image size: %d * %d, format detected: %s", cfg.Width, cfg.Height, format)
		if cfg.Width < 16 || cfg.Height < 16 {
			log.Fatalln("Fatal: Image size should be bigger than 16*16.")
		}
		img, _, err := image.Decode(bytes.NewBuffer(bs))
		if err != nil {
			log.Fatalln("Error while loading image:", err)
		}
		g = &gen.ImageGenerator{
			Image:  img,
			Width:  int32(cfg.Width),
			Height: int32(cfg.Height),
		}
	} else {
		g = gen.GetGenerator(strings.ToLower(genname))
		if g == nil {
			log.Fatalln("Fatal: cannot find given generator name:", genname)
		}
	}
	log.Println("Generator type:", reflect.TypeOf(g))
	g.Init()
	log.Println("Generator init done. Initializing level...")
	log.Println("Level format type:", lvformat)
	p := format.GetProvider(lvformat)
	if p == nil {
		log.Fatalln("Error: cannot find the format provider from server.")
	}
	p.Init(lav7.GetDefaultLevel().Name)
	lav7.GetDefaultLevel().Init(p)
	lav7.GetDefaultLevel().Gen = g.Gen
	log.Printf("Level init done.")
}

func initRaknet() {
	raknet.ServerName = config.ServerName
	atomic.StoreInt32(&raknet.MaxPlayers, config.MaxPlayers)
}

func startLevel() {
	go lav7.GetDefaultLevel().Process()
}

func startRouter(port uint16) {
	log.Println("Starting raknet router, version", raknet.Version)
	var r *raknet.Router
	var err error
	if r, err = raknet.CreateRouter(lav7.RegisterPlayer, lav7.UnregisterPlayer, port); err != nil {
		log.Fatalln("Error while creating router:", err)
	}
	r.Start()
}
