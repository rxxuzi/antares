package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rxxuzi/antares/internal/global"
	"github.com/rxxuzi/antares/internal/server"
)

func main() {
	var (
		port int
		root string
	)

	flag.IntVar(&port, "port", 0, "Port to run the server on (overrides config file)")
	flag.StringVar(&root, "root", "", "Root directory to serve (overrides config file)")

	flag.Usage = func() {
		fmt.Println("Usage of antares:")
		fmt.Println("  -port <int>")
		fmt.Println("        Port to run the server on (overrides config file)")
		fmt.Println("  -root <string>")
		fmt.Println("        Root directory to serve (overrides config file)")
		fmt.Println("\nExample:")
		fmt.Println("  antares -port 8080 -root public/")
	}

	// 引数の解析
	flag.Parse()

	// 設定ファイルの読み込み
	config, err := global.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// コマンドライン引数で上書き
	if port != 0 {
		config.Port = port
	}
	if root != "" {
		config.Root = root
	}

	// rootが相対パスの場合、絶対パスに変換
	if !filepath.IsAbs(config.Root) {
		absRoot, err := filepath.Abs(config.Root)
		if err != nil {
			log.Fatalf("Failed to get absolute path for root directory: %v", err)
		}
		config.Root = absRoot
	}

	_, err = os.Stat(config.Root)
	if os.IsNotExist(err) {
		fmt.Println(config.Root + " does not exist.")
		return
	}

	serverConfig := &server.Config{
		Port:    config.Port,
		RootDir: config.Root,
		LogFlag: false,
	}

	// サーバーの作成
	svr, err := server.CreateServer(serverConfig)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// アクセス情報の表示
	server.PrintAccessInfo(serverConfig)

	// サーバーの起動
	go func() {
		log.Printf("Starting server on port %d, serving directory: %s\n", config.Port, config.Root)
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// シグナル処理の設定
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// シグナル待機
	<-sigChan

	// シャットダウン処理
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}
