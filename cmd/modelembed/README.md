# Model Memory-Mapped Loading Example

This example demonstrates efficient model loading using **memory-mapped files (mmap)**, which is the recommended approach for large language models.

## Why Memory-Mapped Files?

Memory-mapped files (mmap) provide the best solution for loading large models:
- **No size limitations** - Works with models of any size (unlike go:embed's 2GB limit)
- **Efficient memory usage** - OS manages paging, only loads needed parts into RAM
- **Fast loading** - Near-instant model loading
- **Shared between processes** - Multiple processes can share the same mapped model
- **Zero-copy** - Direct access to file data without copying

## Prerequisites

Place your GGUF model file in the `models/` directory:
```bash
cp /path/to/your/model.gguf models/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf
```

Or specify a different model path with the `-model` flag.

## Building

```bash
# From llama.go root directory
./scripts/build.sh
```

This will build the `modelembed` binary in `build/bin/`.

## Usage

### Memory-Mapped Loading (Default, Recommended)
```bash
# Mmap is the default mode
./build/bin/modelembed -p "What is the capital of France?"

# Or explicitly specify mmap
./build/bin/modelembed -mmap -p "What is the capital of France?"
```

### Memory Buffer Loading (Alternative)
For comparison, you can also use traditional memory buffer loading:
```bash
./build/bin/modelembed -memory -mmap=false -p "What is the capital of France?"
```

### Standard File Loading
```bash
./build/bin/modelembed -mmap=false -p "What is the capital of France?"
```

### Interactive Mode
Start an interactive chat session:
```bash
# Using mmap (default)
./build/bin/modelembed -i

# Using memory buffer
./build/bin/modelembed -memory -mmap=false -i
```

## Command Line Options

- `-model <path>` - Path to GGUF model file (default: models/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf)
- `-p <prompt>` - Initial prompt (default: "Hello! How can I help you today?")
- `-i` - Interactive mode
- `-mmap` - Use memory-mapped file loading (default: true)
- `-memory` - Load model into memory buffer (default: false)
- `-ctx <size>` - Context size (default: 2048)
- `-ngl <layers>` - Number of GPU layers (default: 0)
- `-n <tokens>` - Number of tokens to predict (default: 512)
- `-seed <seed>` - Random seed (default: 42)

## Performance Comparison

The example shows memory statistics after each run to compare different loading methods:

### Memory-Mapped (mmap) ✨
- **Allocated**: Minimal, only active pages
- **System**: Efficient, managed by OS
- **Load time**: Near-instant
- **Suitable for**: All model sizes, especially large models

### Memory Buffer
- **Allocated**: Full model size
- **System**: Higher usage
- **Load time**: Slower (needs to read entire file)
- **Suitable for**: Small models with plenty of RAM

### Standard File
- **Allocated**: Variable
- **System**: Standard
- **Load time**: Standard
- **Suitable for**: Traditional usage

## Implementation Details

The mmap implementation (`wrapper/mmap_*.go`) provides:
- Cross-platform support (Darwin/Linux/Windows)
- Page-aligned memory mapping
- GGUF format validation
- Efficient memory management

## Example Output

```
=== Memory Loading Example ===
System: darwin/arm64
CPU cores: 8
Go version: go1.21.0

Configuration:
  Model: models/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf
  Context size: 2048
  GPU layers: 0
  Mode: Memory-mapped file (efficient, no size limit)

Loading model using memory-mapped file...
✓ Model mapped at address: 0x12a000000
✓ Model size: 637.54 MB

Model loaded successfully!
===============================

Prompt: What is the capital of France?
-------------------------------
Response: The capital of France is Paris.

Memory Statistics:
  Allocated: 12.45 MB
  Total: 15.23 MB
  System: 48.67 MB
  GC cycles: 3
```

## Advantages for Large Models

When working with large language models (7B, 13B, 70B+ parameters):

1. **Memory Efficiency**: Only active parts of the model are loaded into RAM
2. **Fast Startup**: No need to load the entire model before starting
3. **Multi-Process**: Multiple instances can share the same mapped model
4. **OS Optimization**: Operating system handles caching and paging efficiently
5. **No Size Limits**: Works with models much larger than available RAM

## Troubleshooting

### "Model file not found"
- Ensure the model exists at the specified path
- Check file permissions

### "mmap failed"
- Check available system resources
- Ensure the file is not corrupted
- Try with a smaller model first

### High memory usage with mmap
- This is normal initially as the OS maps the file
- Actual RAM usage will be much lower than the file size
- The OS will page out unused portions automatically