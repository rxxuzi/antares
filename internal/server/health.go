package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func healthHandler(w http.ResponseWriter, r *http.Request, staticFS fs.FS) {
	status := getHealthStatus()

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
		return
	}

	tmplContent, err := fs.ReadFile(staticFS, "health.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read template: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("health").Funcs(template.FuncMap{
		"toLower":     strings.ToLower,
		"formatBytes": formatBytes,
	}).Parse(string(tmplContent))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, status)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		status := getHealthStatus()
		data := map[string]interface{}{
			"CPUUse":    status.CPUUse,
			"MemoryUse": status.MemoryUse,
			"Uptime":    status.Uptime,
		}
		conn.WriteJSON(data)
		time.Sleep(time.Second)
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

type HealthStatus struct {
	Status        string    `json:"status"`
	Uptime        string    `json:"uptime"`
	MemoryUse     float64   `json:"memory_use"`
	CPUUse        float64   `json:"cpu_use"`
	DiskUse       float64   `json:"disk_use"`
	NetworkIn     uint64    `json:"network_in"`
	NetworkOut    uint64    `json:"network_out"`
	DiskStatus    string    `json:"disk_status"`
	NetworkStatus string    `json:"network_status"`
	StartTime     time.Time `json:"start_time"`
}

func getHealthStatus() HealthStatus {
	uptime := time.Since(startTime).Round(time.Second)

	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := net.IOCounters(false)

	var netIn, netOut uint64
	if len(n) > 0 {
		netIn, netOut = n[0].BytesRecv, n[0].BytesSent
	}

	diskStatus := "Healthy"
	if d.UsedPercent > 90 {
		diskStatus = "Warning"
	}

	networkStatus := "Operational"
	if netIn == 0 && netOut == 0 {
		networkStatus = "No Traffic"
	}

	return HealthStatus{
		Status:        "healthy",
		Uptime:        uptime.String(),
		MemoryUse:     v.UsedPercent,
		CPUUse:        c[0],
		DiskUse:       d.UsedPercent,
		NetworkIn:     netIn,
		NetworkOut:    netOut,
		DiskStatus:    diskStatus,
		NetworkStatus: networkStatus,
		StartTime:     startTime,
	}
}
