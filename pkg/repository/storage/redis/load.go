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
	"github.com/redis/go-redis/v9"
)

const (
	requestDelay = 100 * time.Millisecond
)

var (
	numConcurrentWorkers = int(math.Max(float64(runtime.NumCPU()*10), 10))
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

func LoadCache(logger *utility.Logger, redisClient *redis.Client) {
	var (
		extReq = requests.ExternalRequest{Logger: logger}
		wg     sync.WaitGroup
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	startTime := time.Now()
	log.Printf("\n🚀 Starting cache loading process at %v\n", startTime.Format("15:04:05"))

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
		go worker(i, logger, redisClient, extReq, routeTasks, rateLimiter, &wg, errorChan)
	}

	taskID := 0
	for _, agency := range agencyList.Agencies {
		log.Printf("\n🏢 Processing agency: %s\n", agency.Tag)

		routeList, err := agencyService.GetRouteList(logger, extReq, agency.Tag)
		if err != nil {
			log.Printf("❌ Error fetching routes for agency %s: %v\n", agency.Tag, err)
			continue
		}

		log.Printf("📍 Found %d routes for agency %s\n", len(routeList.Routes), agency.Tag)
		atomic.AddUint64(&totalTasks, uint64(len(routeList.Routes)))

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

	// monitoring goroutine
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
	redisClient *redis.Client,
	extReq requests.ExternalRequest,
	tasks <-chan routeTask,
	rateLimiter chan struct{},
	wg *sync.WaitGroup,
	errorChan chan<- error,
) {
	for task := range tasks {
		startTime := time.Now()

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
		storeStopsInRedis(redisClient, task.agencyTag, task.routeTag, routeConfigs.Route.Stop)
		cacheMutex.Unlock()

		atomic.AddUint64(&completedTasks, 1)
		log.Printf("✅ Worker %d completed task %d (%s-%s) in %v\n",
			workerID, task.taskID, task.agencyTag, task.routeTag, time.Since(startTime))

		wg.Done()
	}
	log.Printf("👋 Worker %d shutting down\n", workerID)
}
