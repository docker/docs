# docker model bench

<!---MARKER_GEN_START-->
Benchmark a model's performance showing tokens per second at different concurrency levels.

This command runs a series of benchmarks with 1, 2, 4, and 8 concurrent requests by default,
measuring the tokens per second (TPS) that the model can generate.

### Options

| Name            | Type       | Default                                                                         | Description                           |
|:----------------|:-----------|:--------------------------------------------------------------------------------|:--------------------------------------|
| `--concurrency` | `intSlice` | `[1,2,4,8]`                                                                     | Concurrency levels to test            |
| `--duration`    | `duration` | `30s`                                                                           | Duration to run each concurrency test |
| `--json`        | `bool`     |                                                                                 | Output results in JSON format         |
| `--prompt`      | `string`   | `Write a comprehensive 100 word summary on whales and their impact on society.` | Prompt to use for benchmarking        |
| `--timeout`     | `duration` | `5m0s`                                                                          | Timeout for each individual request   |


<!---MARKER_GEN_END-->

