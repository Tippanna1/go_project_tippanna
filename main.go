package main

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// instances
	r := gin.Default()

	// New endpoint for sequential sorting
	r.POST("/process-single", func(c *gin.Context) {
		var inputData struct {
			ToSort [][]int `json:"to_sort"`
		}

		if err := c.BindJSON(&inputData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		startTime := time.Now()

		sortedArrays := make([][]int, len(inputData.ToSort))
		for i, arr := range inputData.ToSort {
			sortedArr := make([]int, len(arr))
			copy(sortedArr, arr)
			sort.Ints(sortedArr)
			sortedArrays[i] = sortedArr
		}

		elapsedTime := time.Since(startTime)

		c.JSON(http.StatusOK, gin.H{
			"sorted_arrays": sortedArrays,
			"time_ns":       elapsedTime.Nanoseconds(),
		})
	})

	// New endpoint for concurrent sorting
	r.POST("/process-concurrent", func(c *gin.Context) {
		var inputData struct {
			ToSort [][]int `json:"to_sort"`
		}

		if err := c.BindJSON(&inputData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		startTime := time.Now()

		sortedArrays := make([][]int, len(inputData.ToSort))
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i, arr := range inputData.ToSort {
			wg.Add(1)
			go func(index int, array []int) {
				defer wg.Done()

				sortedArr := make([]int, len(array))
				copy(sortedArr, array)
				sort.Ints(sortedArr)

				mu.Lock()
				sortedArrays[index] = sortedArr
				mu.Unlock()
			}(i, arr)
		}

		wg.Wait()

		elapsedTime := time.Since(startTime)

		c.JSON(http.StatusOK, gin.H{
			"sorted_arrays": sortedArrays,
			"time_ns":       elapsedTime.Nanoseconds(),
		})
	})

	r.Run()
}
