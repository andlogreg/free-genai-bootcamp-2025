# Description

Simple ollama and webui setup.

# Quick Start

## Ollama

Ollama API will be available at localhost:11434.

Examples to test the correct installation:


**Pull model:**

```bash
curl -X POST http://localhost:11434/api/pull -d '{"name": "llama3.2:1b"}'
```

**Text generation:**

```bash
curl -X POST http://localhost:11434/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Write a short poem about machine learning",
  "stream": false
}'
```


## WebUI

WebUI will be available at localhost:3000.

To access it, you need to open the browser and go to http://localhost:3000.

From there, you can use the WebUI to interact Ollama to pull models and generate text.
