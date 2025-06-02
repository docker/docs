# docker model pull

<!---MARKER_GEN_START-->
Pull a model from Docker Hub or HuggingFace to your local environment


<!---MARKER_GEN_END-->

## Description

Pull a model to your local environment. Downloaded models also appear in the Docker Desktop Dashboard.

## Examples

### Pulling a model from Docker Hub

```console
docker model pull ai/smollm2
```

### Pulling from HuggingFace

You can pull GGUF models directly from [Hugging Face](https://huggingface.co/models?library=gguf).

```console
docker model pull hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF
```
