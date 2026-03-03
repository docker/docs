# docker model run

<!---MARKER_GEN_START-->
Run a model and interact with it using a submitted prompt or chat mode

### Options

| Name             | Type     | Default | Description                                          |
|:-----------------|:---------|:--------|:-----------------------------------------------------|
| `--color`        | `string` | `no`    | Use colored output (auto\|yes\|no)                   |
| `--debug`        | `bool`   |         | Enable debug logging                                 |
| `-d`, `--detach` | `bool`   |         | Load the model in the background without interaction |
| `--openaiurl`    | `string` |         | OpenAI-compatible API endpoint URL to chat with      |


<!---MARKER_GEN_END-->

## Description

When you run a model, Docker calls an inference server API endpoint hosted by the Model Runner through Docker Desktop. The model stays in memory until another model is requested, or until a pre-defined inactivity timeout is reached (currently 5 minutes).

You do not have to use Docker model run before interacting with a specific model from a host process or from within a container. Model Runner transparently loads the requested model on-demand, assuming it has been pulled and is locally available.

You can also use chat mode in the Docker Desktop Dashboard when you select the model in the **Models** tab.

## Examples

### One-time prompt

```console
docker model run ai/smollm2 "Hi"
```

Output:

```console
Hello! How can I assist you today?
```

### Interactive chat

```console
docker model run ai/smollm2
```

Output:

```console
> Hi
Hi there! It's SmolLM, AI assistant. How can I help you today?
> /bye
```

### Pre-load a model

```console
docker model run --detach ai/smollm2
```

This loads the model into memory without interaction, ensuring maximum performance for subsequent requests.
