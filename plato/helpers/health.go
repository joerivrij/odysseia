package helpers

import (
	"fmt"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetHealthOfApp(elasticClient elastic.Client) models.Health {
	currentTime := time.Now()
	memUsage := GetMemoryUsage()
	cpuUsage := GetCPUSample()
	elasticHealth := elasticClient.Health().Info()
	overallHealth := false

	cpuPercentage := strconv.FormatUint(cpuUsage, 10)

	if elasticHealth.Healthy {
		overallHealth = true
	}

	health := models.Health{
		Healthy:       overallHealth,
		Time:          currentTime.String(),
		Database:      elasticHealth,
		Memory:        memUsage,
		CPUPercentage: fmt.Sprintf("%s%%", cpuPercentage),
	}

	return health
}

func GetHealthWithVault(vaultHealth bool) models.Health {
	currentTime := time.Now()
	memUsage := GetMemoryUsage()
	cpuUsage := GetCPUSample()

	cpuPercentage := strconv.FormatUint(cpuUsage, 10)

	health := models.Health{
		Healthy:       vaultHealth,
		Time:          currentTime.String(),
		Memory:        memUsage,
		CPUPercentage: fmt.Sprintf("%s%%", cpuPercentage),
	}

	return health
}

func GetMemoryUsage() models.Memory {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats

	return models.Memory{
		Free:       bToMb(m.Frees),      //Frees is the cumulative count of heap objects freed.
		Alloc:      bToMb(m.Alloc),      //Alloc is bytes of allocated heap objects.
		TotalAlloc: bToMb(m.TotalAlloc), //TotalAlloc is cumulative bytes allocated for heap objects.
		Sys:        bToMb(m.Sys),        //Sys is the total bytes of memory obtained from the OS.
		Unit:       "mb",
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func GetCPUSample() uint64 {
	var idle uint64
	var total uint64

	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return 100 * (total - idle) / total
		}
	}
	return 100 * (total - idle) / total
}
