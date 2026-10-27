---
title: "Calculator Tool"
description: "Calculate arithmetic expressions with exact results instead of relying on the model's mental arithmetic."
keywords: docker agent, ai agents, tools, toolsets, calculator, arithmetic
linkTitle: "Calculator"
weight: 146
canonical: https://docs.docker.com/ai/docker-agent/tools/calculator/
---

The calculator toolset gives agents one tool, `calculate`, for reliable arithmetic.
It uses exact rational numbers rather than floating-point arithmetic: `0.1 + 0.2`
returns exactly `0.3`, and large integers do not silently overflow.

Calculations run locally, with no code execution, filesystem access, network
access, or side effects. The tool is read-only and available in native and browser
(WASM) builds.

## Configuration

```yaml
agents:
  root:
    model: openai/gpt-5-mini
    instruction: Use calculate for all numerical calculations.
    toolsets:
      - type: calculator
```

No toolset-specific configuration is required. See
[`examples/calculator.yaml`](https://github.com/docker/docker-agent/blob/main/examples/calculator.yaml)
for a complete agent.

## Arguments

`calculate` takes one required argument:

| Argument | Type | Description |
| --- | --- | --- |
| `expression` | string | An arithmetic expression, such as `(19.99 * 3) * (1 + 8.25 / 100)`. |

Supported syntax:

- Decimal integers and fractions, including `.5`, `1.`, and scientific notation
  (`1.5e3`). Leading zeros are decimal, not octal.
- Addition (`+`), subtraction (`-`), multiplication (`*`), and division (`/`).
- Integer remainder (`%`), with the sign of the dividend: `-17 % 5` is `-2`.
  **This is not a percentage operator**; write `15 / 100` for 15%.
- Powers (`^`) with integer exponents, including negative exponents (`2^-3`).
- Parentheses, whitespace, and unary `+` and `-`.

Powers bind before unary signs, then multiplication/division/remainder, then
addition/subtraction. Powers associate right-to-left: `2^3^2` is `512`;
`-2^2` is `-4`, while `(-2)^2` is `4`.

Variables, functions, units, commas, implicit multiplication, and programming
language syntax are not supported. Convert quantities to compatible units before
calculating; the calculator cannot verify that the expression models the user's
problem correctly.

## Results and precision

For terminating decimals, both `result` and `decimal` are exact:

```json
{"result":"0.3","decimal":"0.3","approximate":false}
```

For repeating decimals, `result` is an exact reduced fraction and `decimal` is
rounded to 20 decimal places:

```json
{"result":"1/3","decimal":"0.33333333333333333333","approximate":true}
```

All numbers are returned as **strings**, preserving precision across JSON clients.
Use `result`, not a rounded `decimal`, in subsequent expressions. Intermediate
calculations remain exact: `1 / 3 * 3` returns `1`.

## Errors and limits

Malformed expressions, division or remainder by zero, non-integer remainder
operands or power exponents, and `0^0` return tool errors rather than guessed
answers. To bound resource use, expressions are limited to 8192 bytes, nesting to
64 levels, scientific and power exponents to -1024 through 1024, and intermediate
numerators and denominators to 4096 bits. Oversized powers can be rejected before
calculation using a conservative size estimate.
