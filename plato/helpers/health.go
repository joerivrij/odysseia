package helpers

import (
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"runtime"
	"time"
)

func GetHealthOfApp(elasticClient elastic.Client) models.Health {
	currentTime := time.Now()
	memUsage := GetMemoryUsage()
	elasticHealth := elasticClient.Health().Info()
	overallHealth := false

	if elasticHealth.Healthy {
		overallHealth = true
	}

	health := models.Health{
		Healthy:  overallHealth,
		Time:     currentTime.String(),
		Database: elasticHealth,
		Memory:   memUsage,
	}

	return health
}

func GetHealthWithVault(vaultHealth bool) models.Health {
	currentTime := time.Now()
	memUsage := GetMemoryUsage()

	health := models.Health{
		Healthy: vaultHealth,
		Time:    currentTime.String(),
		Memory:  memUsage,
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
