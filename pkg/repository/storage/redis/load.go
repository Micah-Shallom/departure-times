// package redis

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/Micah-Shallom/departure-times/external/external_models"
// 	"github.com/Micah-Shallom/departure-times/external/requests"
// 	agencyService "github.com/Micah-Shallom/departure-times/services/agency"
// 	stopService "github.com/Micah-Shallom/departure-times/services/stops"
// 	"github.com/Micah-Shallom/departure-times/utility"
// )

// var (
// 	cacheMutex = &sync.Mutex{}
// )

// var cache = external_models.Cache{
// 	RouteListCache:   make(map[string]external_models.GetRoutesResponse),
// 	RouteConfigCache: make(map[string]external_models.GetRoutesConfigResponse),
// 	StopListCache:    make(map[string]map[string][]external_models.Stop),
// }

// func LoadCache(logger *utility.Logger) {
// 	var (
// 		extReq = requests.ExternalRequest{Logger: logger}
// 	)

// 	cacheMutex.Lock()
// 	defer cacheMutex.Unlock()

// 	agencyList, err := agencyService.GetAgencyList(logger, extReq)
// 	if err != nil {
// 		logger.Error("unable to fetch agency list", err)
// 		log.Printf("unable to fetch agency list", err)
// 		return
// 	}
// 	logger.Info("agency list fetched successfully")

// 	for _, agency := range agencyList.Agencies {
// 		routeList, err := agencyService.GetRouteList(logger, extReq, agency.Tag)
// 		if err != nil {
// 			logger.Error("error fetching route list for agency %s", agency.Tag, err)
// 			log.Printf("error fetching route list for agency %s", agency.Tag, err)
// 			continue
// 		}
// 		logger.Info("route list fetched successfully")
// 		cache.RouteListCache[agency.Tag] = routeList
// 		logger.Info("route list cached successfully")

// 		for _, route := range routeList.Routes {
// 			routeConfigs, err := stopService.GetRouteConfigurations(logger, extReq, agency.Tag, route.Tag)
// 			if err != nil {
// 				logger.Error("unable to fetch route configurations for agency %s, route %s", agency.Tag, route.Tag, err)
// 				log.Printf("unable to fetch route configurations for agency %s, route %s", agency.Tag, route.Tag, err)
// 				continue
// 			}

// 			logger.Info("route configurations fetched successfully")
// 			cache.RouteConfigCache[route.Tag] = routeConfigs
// 			log.Printf("routeConfigs.Config.Stops: %+v\n", routeConfigs.Config.Stops)

// 			//lets initialize the inner map if it doesn't exist
// 			if _, ok := cache.StopListCache[agency.Tag]; !ok {
// 				cache.StopListCache[agency.Tag] = make(map[string][]external_models.Stop)
// 			}

// 			cache.StopListCache[agency.Tag][route.Tag] = routeConfigs.Config.Stops
// 			logger.Info("route configurations cached successfully")
// 		}
// 	}

// }

package redis

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Micah-Shallom/departure-times/external/external_models"
	"github.com/Micah-Shallom/departure-times/external/requests"
	agencyService "github.com/Micah-Shallom/departure-times/services/agency"
	stopService "github.com/Micah-Shallom/departure-times/services/stops"
	"github.com/Micah-Shallom/departure-times/utility"
)

const (
	requestDelay = 100 * time.Millisecond
)

var (
	numConcurrentWorkers = int(math.Min(float64(runtime.NumCPU()), 10))
	cacheMutex           = &sync.Mutex{}
	cache                = external_models.Cache{
		RouteListCache:   make(map[string]external_models.GetRoutesResponse),
		RouteConfigCache: make(map[string]external_models.GetRoutesConfigResponse),
		StopListCache:    make(map[string]map[string][]external_models.Stop),
	}

	// statistics counters
	totalTasks     uint64
	completedTasks uint64
	errorCount     uint64
)

type routeTask struct {
	agencyTag string
	routeTag  string
	taskID    int
}

func LoadCache(logger *utility.Logger) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	startTime := time.Now()
	log.Printf("\n🚀 Starting cache loading process at %v\n", startTime.Format("15:04:05"))

	var (
		extReq = requests.ExternalRequest{Logger: logger}
		wg     sync.WaitGroup
	)

	// Get agency list
	log.Printf("\n📋 Fetching agency list...")
	agencyList, err := agencyService.GetAgencyList(logger, extReq)
	if err != nil {
		log.Printf("❌ Failed to fetch agency list: %v\n", err)
		return
	}
	log.Printf("✅ Successfully fetched %d agencies\n", len(agencyList.Agencies))

	rateLimiter := make(chan struct{}, numConcurrentWorkers)
	routeTasks := make(chan routeTask)
	errorChan := make(chan error, 100)

	// Start worker pool
	log.Printf("\n👷 Starting %d workers...\n", numConcurrentWorkers)
	for i := 1; i <= numConcurrentWorkers; i++ {
		go worker(i, logger, extReq, routeTasks, rateLimiter, &wg, errorChan)
	}

	taskID := 0
	// Process each agency
	for _, agency := range agencyList.Agencies {
		log.Printf("\n🏢 Processing agency: %s\n", agency.Tag)

		routeList, err := agencyService.GetRouteList(logger, extReq, agency.Tag)
		if err != nil {
			log.Printf("❌ Error fetching routes for agency %s: %v\n", agency.Tag, err)
			continue
		}

		log.Printf("📍 Found %d routes for agency %s\n", len(routeList.Routes), agency.Tag)
		atomic.AddUint64(&totalTasks, uint64(len(routeList.Routes)))

		cacheMutex.Lock()
		cache.RouteListCache[agency.Tag] = routeList
		cache.StopListCache[agency.Tag] = make(map[string][]external_models.Stop)
		cacheMutex.Unlock()

		for _, route := range routeList.Routes {
			taskID++
			wg.Add(1)
			routeTasks <- routeTask{
				agencyTag: agency.Tag,
				routeTag:  route.Tag,
				taskID:    taskID,
			}
		}
	}

	log.Printf("\n📊 Total tasks queued: %d\n", atomic.LoadUint64(&totalTasks))
	close(routeTasks)

	// Start progress monitoring goroutine
	go func() {
		for {
			time.Sleep(5 * time.Second)
			completed := atomic.LoadUint64(&completedTasks)
			total := atomic.LoadUint64(&totalTasks)
			errors := atomic.LoadUint64(&errorCount)

			if completed >= total {
				return
			}

			log.Printf("\n📈 Progress: %d/%d (%.1f%%) completed, %d errors\n",
				completed, total, float64(completed)/float64(total)*100, errors)
		}
	}()

	wg.Wait()
	close(errorChan)

	// Final statistics
	duration := time.Since(startTime)
	log.Printf("\n🏁 Cache loading completed in %v\n", duration)
	log.Printf("📊 Final Statistics:\n")
	log.Printf("   - Total tasks: %d\n", atomic.LoadUint64(&totalTasks))
	log.Printf("   - Completed: %d\n", atomic.LoadUint64(&completedTasks))
	log.Printf("   - Errors: %d\n", atomic.LoadUint64(&errorCount))
	log.Printf("   - Average time per task: %v\n", duration/time.Duration(atomic.LoadUint64(&completedTasks)))
}

func worker(
	workerID int,
	logger *utility.Logger,
	extReq requests.ExternalRequest,
	tasks <-chan routeTask,
	rateLimiter chan struct{},
	wg *sync.WaitGroup,
	errorChan chan<- error,
) {
	for task := range tasks {
		// startTime := time.Now()
		// log.Printf("👨‍💼 Worker %d starting task %d (%s-%s)\n",
		//     workerID, task.taskID, task.agencyTag, task.routeTag)

		rateLimiter <- struct{}{}
		time.Sleep(requestDelay)

		routeConfigs, err := stopService.GetRouteConfigurations(
			logger,
			extReq,
			task.agencyTag,
			task.routeTag,
		)

		<-rateLimiter

		if err != nil {
			atomic.AddUint64(&errorCount, 1)
			errorChan <- fmt.Errorf(
				"error fetching route config for agency %s, route %s: %v",
				task.agencyTag,
				task.routeTag,
				err,
			)
			log.Printf("❌ Worker %d failed task %d (%s-%s): %v\n",
			    workerID, task.taskID, task.agencyTag, task.routeTag, err)
			wg.Done()
			continue
		}

		cacheMutex.Lock()
		cache.RouteConfigCache[task.routeTag] = routeConfigs
		cache.StopListCache[task.agencyTag][task.routeTag] = routeConfigs.Config.Stops
		cacheMutex.Unlock()

		atomic.AddUint64(&completedTasks, 1)
		// log.Printf("✅ Worker %d completed task %d (%s-%s) in %v\n",
		//     workerID, task.taskID, task.agencyTag, task.routeTag, time.Since(startTime))

		wg.Done()
	}
	log.Printf("👋 Worker %d shutting down\n", workerID)
}
