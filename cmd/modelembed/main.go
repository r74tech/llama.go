// Copyright (c) 2017-2025 The qitmeer developers

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/wrapper"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(20)

	var modelPath string
	var prompt string
	var interactive bool
	var contextSize int
	var nGpuLayers int
	var nPredict int
	var seed uint
	var useMemory bool
	var useMmap bool

	// Default model path
	defaultModel := filepath.Join("models", "tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf")

	flag.StringVar(&modelPath, "model", defaultModel, "Path to the GGUF model file")
	flag.StringVar(&prompt, "p", "Hello! How can I help you today?", "Initial prompt")
	flag.BoolVar(&interactive, "i", false, "Interactive mode")
	flag.IntVar(&contextSize, "ctx", 2048, "Context size")
	flag.IntVar(&nGpuLayers, "ngl", 0, "Number of GPU layers")
	flag.IntVar(&nPredict, "n", 512, "Number of tokens to predict")
	flag.UintVar(&seed, "seed", 42, "Random seed")
	flag.BoolVar(&useMemory, "memory", false, "Load model from memory buffer")
	flag.BoolVar(&useMmap, "mmap", true, "Use memory-mapped file loading (default)")
	flag.Parse()

	// If memory is explicitly requested, disable mmap
	if useMemory {
		useMmap = false
	}

	// Check if model exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Model file not found: %s\n", modelPath)
		fmt.Fprintf(os.Stderr, "Please ensure the model file exists or use -model flag to specify a different path.\n")
		os.Exit(1)
	}

	// Create configuration
	cfg := &config.Config{
		Model:       modelPath,
		Interactive: interactive,
		Prompt:      prompt,
		CtxSize:     contextSize,
		NGpuLayers:  nGpuLayers,
		NPredict:    nPredict,
		Seed:        seed,
	}

	// Print system information
	fmt.Printf("=== Memory Loading Example ===\n")
	fmt.Printf("System: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPU cores: %d\n", runtime.NumCPU())
	fmt.Printf("Go version: %s\n\n", runtime.Version())

	// Print configuration
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Model: %s\n", modelPath)
	fmt.Printf("  Context size: %d\n", contextSize)
	fmt.Printf("  GPU layers: %d\n", nGpuLayers)
	fmt.Printf("  Mode: ")
	if useMmap {
		fmt.Println("Memory-mapped file (efficient, no size limit)")
	} else if useMemory {
		fmt.Println("Memory buffer")
	} else {
		fmt.Println("Standard file")
	}
	fmt.Println()

	// Load the model based on the selected method
	var err error

	if useMmap {
		// Use memory-mapped file loading
		fmt.Println("Loading model using memory-mapped file...")
		addr, data, err := wrapper.MmapModel(modelPath)
		if err != nil {
			log.Fatalf("Failed to mmap model: %v", err)
		}
		defer wrapper.UnmapModel(data)

		fmt.Printf("✓ Model mapped at address: 0x%x\n", addr)
		fmt.Printf("✓ Model size: %.2f MB\n\n", float64(len(data))/(1024*1024))

		// Load from mmap
		err = wrapper.LoadFromMmap(addr, data, cfg)
		if err != nil {
			log.Fatalf("Failed to load model from mmap: %v", err)
		}
	} else if useMemory {
		// Load entire file into memory first
		fmt.Println("Loading model into memory buffer...")
		modelData, err := os.ReadFile(modelPath)
		if err != nil {
			log.Fatalf("Failed to read model file: %v", err)
		}

		fmt.Printf("✓ Model size: %.2f MB\n", float64(len(modelData))/(1024*1024))

		// Verify GGUF magic
		if len(modelData) < 4 || modelData[0] != 'G' || modelData[1] != 'G' ||
			modelData[2] != 'U' || modelData[3] != 'F' {
			log.Fatal("Invalid GGUF file format")
		}
		fmt.Println("✓ GGUF magic verified")

		err = wrapper.LoadFromMemory(modelData, cfg)
		if err != nil {
			log.Fatalf("Failed to load model from memory: %v", err)
		}
	} else {
		// Standard file-based loading
		fmt.Println("Loading model from file (standard mode)...")
		err = wrapper.LlamaStart(cfg)
		if err != nil {
			log.Fatalf("Failed to load model: %v", err)
		}
	}

	// Ensure cleanup on exit
	defer func() {
		fmt.Println("\nShutting down...")
		wrapper.LlamaStop()
	}()

	fmt.Println("Model loaded successfully!")
	fmt.Println("===============================")

	if interactive {
		// Interactive mode
		fmt.Println("\nInteractive mode. Type 'exit' to quit.")
		fmt.Println("-------------------------------")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("\n> ")
			if !scanner.Scan() {
				break
			}

			input := scanner.Text()
			if input == "exit" || input == "quit" || input == "/exit" {
				break
			}

			response, err := wrapper.LlamaGenerate(input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println(response)
		}
	} else {
		// Single generation mode
		fmt.Printf("\nPrompt: %s\n", prompt)
		fmt.Println("-------------------------------")

		response, err := wrapper.LlamaGenerate(prompt)
		if err != nil {
			log.Fatalf("Generation failed: %v", err)
		}

		fmt.Println(response)
	}

	// Print memory statistics
	printMemoryStats()
}

// printMemoryStats prints memory usage statistics
func printMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("\nMemory Statistics:\n")
	fmt.Printf("  Allocated: %.2f MB\n", float64(m.Alloc)/(1024*1024))
	fmt.Printf("  Total: %.2f MB\n", float64(m.TotalAlloc)/(1024*1024))
	fmt.Printf("  System: %.2f MB\n", float64(m.Sys)/(1024*1024))
	fmt.Printf("  GC cycles: %d\n", m.NumGC)
}
