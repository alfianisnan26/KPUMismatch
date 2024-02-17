package model

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
	Active      uint64  `json:"active"`
	Inactive    uint64  `json:"inactive"`
}

func NewFromVM(vm *mem.VirtualMemoryStat) *MemoryInfo {
	if vm == nil {
		return nil
	}

	return &MemoryInfo{
		Total:       vm.Total,
		Available:   vm.Available,
		Used:        vm.Used,
		UsedPercent: vm.UsedPercent,
		Free:        vm.Free,
		Active:      vm.Active,
		Inactive:    vm.Inactive,
	}
}

type WebStats struct {
	start        time.Time
	targetLength int

	UploadID       uint64
	MemoryInfo     *MemoryInfo
	RPS            float32
	Percentage     float32
	DataCount      int
	Timestamp      time.Time
	Estimation     time.Duration
	ProcessingTime time.Duration
}

func NewWebStats(targetLength int) *WebStats {
	return &WebStats{
		start:        time.Now(),
		targetLength: targetLength,
	}
}

func (ws *WebStats) Update(newDataLength int) error {
	ws.RPS = float32(newDataLength) / float32(time.Now().Sub(ws.Timestamp).Seconds())
	if ws.RPS < 0 || ws.RPS > 100000 {
		ws.RPS = 0
	}

	ws.Percentage = float32(ws.DataCount) / float32(ws.targetLength) * 100
	ws.Estimation = estimateTime(ws.start, ws.Percentage, 100)
	ws.Timestamp = time.Now()
	ws.ProcessingTime = ws.Timestamp.Sub(ws.start)

	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	ws.MemoryInfo = NewFromVM(vm)
	return nil
}

func (ws *WebStats) String() string {
	return fmt.Sprintf("[%2.2f%%] %s Pt: (min) %2.2f | Est: (min) %2.2f | RPS: %2.2f | %d/%d", ws.Percentage, ws.Timestamp.Format(time.TimeOnly), ws.ProcessingTime.Minutes(), ws.Estimation.Minutes(), ws.RPS, ws.DataCount, ws.targetLength)
}

// estimateTime calculates the estimated time until reaching the target percentage.
func estimateTime(startTime time.Time, currentPercentage, targetPercentage float32) time.Duration {

	elapsed := time.Since(startTime)
	remainingPercentage := targetPercentage - currentPercentage
	ct := time.Duration(currentPercentage * 1000)
	rt := time.Duration(remainingPercentage * 1000)
	if ct == 0 || rt == 0 {
		return 0
	}

	remainingTime := elapsed / ct * rt
	return remainingTime
}
