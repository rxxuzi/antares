// antares project

package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// Config  holds the configuration for the server
type Config struct {
	Port    int
	RootDir string
	LogFlag bool
}

// GetLocalIP returns the non-loopback local IPv4 address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unable to get local IP"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && !ipnet.IP.IsLinkLocalUnicast() {
				return ipnet.IP.String()
			}
		}
	}
	return "No suitable IP found"
}

// findAvailablePort tries to find an available port starting from the given port
func findAvailablePort(startPort int) (int, error) {
	for port := startPort; port < startPort+100; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		listener.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no available ports found in range %d-%d", startPort, startPort+99)
}

// PrintAccessInfo prints server access information
func PrintAccessInfo(config *Config) {
	localIP := GetLocalIP()
	fmt.Printf("Server Configuration:\n")
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Root Directory: %s\n", config.RootDir)
	fmt.Printf("  Logging: %v\n", config.LogFlag)
	fmt.Println()

	fmt.Println("Access URLs:")
	fmt.Printf("  http://%s:%d\n", localIP, config.Port)
	fmt.Printf("  http://localhost:%d\n", config.Port)

	fmt.Println("\nNote: Use Ctrl+C to stop the server")
}

// logRequest is a middleware that logs each request
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s - [%s] \"%s %s %s\" %s",
			r.RemoteAddr,
			time.Now().Format(time.RFC1123),
			r.Method,
			r.URL.Path,
			r.Proto,
			time.Since(start),
		)
	})
}
