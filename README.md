# llama.go
Go bindings to llama.cpp

### Installation
***make sure you have `git golang cmake gcc make` installed on the system to build.***
* Build from source
```bash
~ git clone https://github.com/Qitmeer/llama.go.git
~ cd llama.go
~ make
```

### Get model
* Manually download the model:[Hugging Face Qwen3-8B-GGUF](https://huggingface.co/ggml-org/Qwen3-8B-GGUF/tree/main)

### Local startup

```bash
~ ./llama --model=./qwen2.5-0.5b-q8_0.gguf --prompt=天空为什么是蓝的
```
Or enable interactive mode to run:
```bash
~ ./llama --model=./qwen2.5-0.5b-q8_0.gguf -i
```


### As the startup of the server

```bash
~ ./llama --model=./qwen2.5-0.5b-q8_0.gguf
```

* Supports two access modes: grpc and regular REST API:
```bash
~ curl -s -k -X POST -H 'Content-Type: application/json' --data '{"prompt":"天空为什么是蓝的"}' http://127.0.0.1:8081/v1/generate
```

### Embedding

```bash
~ ./llama --model=./qwen2.5-0.5b-q8_0.gguf --prompt=天空为什么是蓝的 --output-file=./embs.json embedding
```