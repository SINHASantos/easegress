package main

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"syscall"

	"github.com/megaease/easegateway/pkg/api"
	"github.com/megaease/easegateway/pkg/cluster"
	"github.com/megaease/easegateway/pkg/logger"
	"github.com/megaease/easegateway/pkg/model"
	"github.com/megaease/easegateway/pkg/option"
	"github.com/megaease/easegateway/pkg/stat"
	"github.com/megaease/easegateway/pkg/store"
	"github.com/megaease/easegateway/pkg/version"
)

func main() {
	defer logger.Close()

	logger.Infof("%s", version.Long)

	cpuProfileDone := setupCPUProfile()
	memProfileDone := setupMemoryoryProfile()

	cluster, err := cluster.New()
	if err != nil {
		logger.Errorf("new cluster failed: %v", err)
		os.Exit(1)
	}
	store, err := store.New(cluster)
	if err != nil {
		logger.Errorf("new store failed: %v", err)
		os.Exit(1)
	}
	model, err := model.NewModel(store)
	if err != nil {
		logger.Errorf("new model failed: %v", err)
		os.Exit(1)
	}

	stat := stat.NewStat(cluster, model)

	api := api.MustNewAPIServer(cluster)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	go func() {
		sig := <-sigChan
		logger.Infof("%s signal received, closing easegateway immediately", sig)
		os.Exit(255)
	}()

	logger.Infof("%s signal received, closing easegateway", sig)

	api.Close()
	stat.Close()
	model.Close()
	store.Close()
	cluster.Close()

	if cpuProfileDone != nil {
		cpuProfileDone <- struct{}{}
		<-cpuProfileDone
	}
	if memProfileDone != nil {
		memProfileDone <- struct{}{}
		<-memProfileDone
	}
}

func setupCPUProfile() chan struct{} {
	if option.Global.CPUProfileFile == "" {
		return nil
	}

	done := make(chan struct{}, 1)

	f, err := os.Create(option.Global.CPUProfileFile)
	if err != nil {
		logger.Errorf("create cpu profile failed: %v", err)
		os.Exit(1)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		logger.Errorf("start cpu profile failed: %v", err)
		os.Exit(1)
	}

	logger.Infof("cpu profile: %s", option.Global.CPUProfileFile)
	go func() {
		<-done
		pprof.StopCPUProfile()
		err := f.Close()
		if err != nil {
			logger.Errorf("close %s failed: %v", option.Global.CPUProfileFile, err)
		}
		close(done)
	}()

	return done
}

func setupMemoryoryProfile() chan struct{} {
	if option.Global.MemoryProfileFile == "" {
		return nil
	}

	done := make(chan struct{}, 1)

	// to include every allocated block in the profile
	runtime.MemProfileRate = 1

	go func() {
		<-done
		logger.Infof("memory profile: %s", option.Global.MemoryProfileFile)
		f, err := os.Create(option.Global.MemoryProfileFile)
		if err != nil {
			logger.Errorf("create memory profile failed: %v", err)
			return
		}

		runtime.GC()         // get up-to-date statistics
		debug.FreeOSMemory() // help developer when using outside monitor tool

		if err := pprof.WriteHeapProfile(f); err != nil {
			logger.Errorf("write memory file failed: %v", err)
			return
		}
		if err := f.Close(); err != nil {
			logger.Errorf("close memory file failed: %v", err)
			return
		}
		close(done)
	}()

	return done
}