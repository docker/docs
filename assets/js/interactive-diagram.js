const svgNamespace = "http://www.w3.org/2000/svg";

document.querySelectorAll("[data-interactive-diagram]").forEach((root) => {
  const configElement = root.querySelector("[data-interactive-diagram-config]");
  if (!configElement) return;

  let config;
  try {
    config = JSON.parse(configElement.textContent);
  } catch {
    return;
  }

  const stage = root.querySelector("[data-diagram-stage]");
  if (!stage || !config.canvas || !config.steps?.length) return;

  const state = buildDiagram(stage, config);
  const previousButton = root.querySelector("[data-step-previous]");
  const nextButton = root.querySelector("[data-step-next]");
  const playButton = root.querySelector("[data-step-play]");
  const progress = root.querySelector("[data-step-progress]");
  let currentStep = 0;
  let playTimer;

  const stepButtons = config.steps.map((step, index) => {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "interactive-diagram__progress-step";
    button.setAttribute("aria-label", `Show step ${index + 1}: ${step.label}`);
    button.addEventListener("click", () => {
      stopPlaying();
      showStep(index);
    });
    progress?.append(button);
    return button;
  });

  function showStep(index) {
    currentStep = Math.max(0, Math.min(index, config.steps.length - 1));
    const step = config.steps[currentStep];
    const activeNodes = new Set(step.activeNodes ?? []);
    const activeEdges = new Set(step.activeEdges ?? []);

    state.nodes.forEach((element, id) => {
      element.classList.toggle("is-active", activeNodes.has(id));
    });
    state.edges.forEach((element, id) => {
      element.classList.toggle("is-active", activeEdges.has(id));
    });
    positionToken(state, config, step.token);

    setText(
      root,
      "[data-step-number]",
      `${String(currentStep + 1).padStart(2, "0")} / ${String(config.steps.length).padStart(2, "0")}`,
    );
    setText(root, "[data-step-title]", step.label);
    setText(root, "[data-step-body]", step.body);
    setText(root, "[data-state-label]", step.stateLabel);
    setText(root, "[data-state-value]", step.stateValue);

    if (previousButton) previousButton.disabled = currentStep === 0;
    if (nextButton)
      nextButton.disabled = currentStep === config.steps.length - 1;
    stepButtons.forEach((button, buttonIndex) => {
      const selected = buttonIndex === currentStep;
      button.classList.toggle("is-active", selected);
      button.setAttribute("aria-current", selected ? "step" : "false");
    });
  }

  function stopPlaying() {
    window.clearInterval(playTimer);
    playTimer = undefined;
    playButton?.setAttribute("aria-pressed", "false");
    if (playButton) playButton.textContent = "Play";
  }

  previousButton?.addEventListener("click", () => {
    stopPlaying();
    showStep(currentStep - 1);
  });
  nextButton?.addEventListener("click", () => {
    stopPlaying();
    showStep(currentStep + 1);
  });
  playButton?.addEventListener("click", () => {
    if (playTimer) {
      stopPlaying();
      return;
    }
    if (currentStep === config.steps.length - 1) showStep(0);
    playButton.setAttribute("aria-pressed", "true");
    playButton.textContent = "Pause";
    playTimer = window.setInterval(() => {
      if (currentStep === config.steps.length - 1) {
        stopPlaying();
      } else {
        showStep(currentStep + 1);
      }
    }, 2600);
  });

  root.classList.add("is-enhanced");
  showStep(0);
});

function buildDiagram(stage, config) {
  const svg = createSvgElement("svg", {
    viewBox: `0 0 ${config.canvas.width} ${config.canvas.height}`,
    role: "img",
    "aria-label": config.description,
  });
  svg.classList.add("interactive-diagram__svg");

  const definitions = createSvgElement("defs");
  const marker = createSvgElement("marker", {
    id: `arrow-${Math.random().toString(36).slice(2)}`,
    viewBox: "0 0 10 10",
    refX: "8",
    refY: "5",
    markerWidth: "7",
    markerHeight: "7",
    orient: "auto-start-reverse",
  });
  marker.append(createSvgElement("path", { d: "M 0 0 L 10 5 L 0 10 z" }));
  definitions.append(marker);
  svg.append(definitions);

  config.boundaries?.forEach((boundary) => {
    const group = createSvgElement("g");
    group.classList.add(
      "interactive-diagram__boundary",
      `interactive-diagram__boundary--${boundary.kind}`,
    );
    group.append(
      createSvgElement("rect", {
        x: boundary.x,
        y: boundary.y,
        width: boundary.width,
        height: boundary.height,
        rx: "8",
      }),
      createSvgText(
        boundary.label,
        boundary.x + 16,
        boundary.y + 25,
        "interactive-diagram__boundary-label",
      ),
    );
    svg.append(group);
  });

  const nodesById = new Map(config.nodes.map((node) => [node.id, node]));
  const edges = new Map();
  config.edges.forEach((edge) => {
    const from = nodesById.get(edge.from);
    const to = nodesById.get(edge.to);
    if (!from || !to) return;
    const points = edgePoints(from, to, edge.offset);
    const line = createSvgElement("line", {
      x1: points.x1,
      y1: points.y1,
      x2: points.x2,
      y2: points.y2,
      "marker-end": `url(#${marker.id})`,
    });
    line.classList.add("interactive-diagram__edge");
    svg.append(line);
    edges.set(edge.id, line);
  });

  const nodes = new Map();
  config.nodes.forEach((node) => {
    const group = createSvgElement("g");
    group.classList.add(
      "interactive-diagram__node",
      `interactive-diagram__node--${node.kind}`,
    );
    group.append(
      createSvgElement("rect", {
        x: node.x,
        y: node.y,
        width: node.width,
        height: node.height,
        rx: "7",
      }),
      createSvgText(
        node.label,
        node.x + 14,
        node.y + 29,
        "interactive-diagram__node-label",
      ),
      createSvgText(
        node.description,
        node.x + 14,
        node.y + 51,
        "interactive-diagram__node-description",
      ),
    );
    svg.append(group);
    nodes.set(node.id, group);
  });

  const token = createSvgElement("g");
  token.classList.add("interactive-diagram__token");
  const tokenRect = createSvgElement("rect", { height: "28", rx: "14" });
  const tokenText = createSvgText("", 0, 0, "interactive-diagram__token-label");
  token.append(tokenRect, tokenText);
  svg.append(token);
  stage.append(svg);

  return {
    nodes,
    edges,
    nodesById,
    token,
    tokenRect,
    tokenText,
    tokenAnimation: undefined,
  };
}

function positionToken(state, config, tokenConfig) {
  state.tokenAnimation?.remove();
  state.tokenAnimation = undefined;

  if (!tokenConfig) {
    state.token.hidden = true;
    return;
  }

  let x;
  let y;
  let travelStart;
  if (tokenConfig.node) {
    const node = state.nodesById.get(tokenConfig.node);
    if (node) {
      x = node.x + node.width / 2;
      y = node.y + node.height / 2;
    }
  } else if (tokenConfig.edge) {
    const edge = config.edges.find(
      (candidate) => candidate.id === tokenConfig.edge,
    );
    const from = state.nodesById.get(edge?.from);
    const to = state.nodesById.get(edge?.to);
    if (from && to) {
      const points = edgePoints(from, to, edge.offset);
      const progress = 0.62;
      x = points.x1 + (points.x2 - points.x1) * progress;
      y = points.y1 + (points.y2 - points.y1) * progress;
      travelStart = { x: points.x1, y: points.y1 };
    }
  }
  if (x === undefined || y === undefined) return;

  const width = Math.max(70, tokenConfig.label.length * 7.2 + 24);
  state.token.hidden = false;
  state.token.setAttribute(
    "transform",
    `translate(${x - width / 2} ${y - 14})`,
  );
  state.tokenRect.setAttribute("width", width);
  state.tokenText.setAttribute("x", width / 2);
  state.tokenText.setAttribute("y", "18");
  state.tokenText.textContent = tokenConfig.label;
  state.token.classList.remove("is-entering");
  if (travelStart && !prefersReducedMotion()) {
    animateTokenTravel(state, travelStart, { x, y }, width);
  } else {
    window.requestAnimationFrame(() =>
      state.token.classList.add("is-entering"),
    );
  }
}

function animateTokenTravel(state, from, to, width) {
  const animation = createSvgElement("animateTransform", {
    attributeName: "transform",
    type: "translate",
    from: `${from.x - width / 2} ${from.y - 14}`,
    to: `${to.x - width / 2} ${to.y - 14}`,
    dur: "650ms",
    calcMode: "spline",
    keyTimes: "0;1",
    keySplines: "0.2 0 0 1",
  });
  state.token.prepend(animation);
  state.tokenAnimation = animation;
  if (typeof animation.beginElement === "function") {
    animation.beginElement();
    window.setTimeout(() => {
      animation.remove();
      if (state.tokenAnimation === animation) state.tokenAnimation = undefined;
    }, 700);
  } else {
    animation.remove();
    state.tokenAnimation = undefined;
  }
}

function prefersReducedMotion() {
  return (
    window.matchMedia?.("(prefers-reduced-motion: reduce)").matches ?? false
  );
}

function edgePoints(from, to, offset = 0) {
  const fromCenter = {
    x: from.x + from.width / 2,
    y: from.y + from.height / 2,
  };
  const toCenter = { x: to.x + to.width / 2, y: to.y + to.height / 2 };
  const dx = toCenter.x - fromCenter.x;
  const dy = toCenter.y - fromCenter.y;
  if (Math.abs(dx) >= Math.abs(dy)) {
    return {
      x1: dx >= 0 ? from.x + from.width : from.x,
      y1: fromCenter.y + offset,
      x2: dx >= 0 ? to.x : to.x + to.width,
      y2: toCenter.y + offset,
    };
  }
  return {
    x1: fromCenter.x + offset,
    y1: dy >= 0 ? from.y + from.height : from.y,
    x2: toCenter.x + offset,
    y2: dy >= 0 ? to.y : to.y + to.height,
  };
}

function createSvgElement(name, attributes = {}) {
  const element = document.createElementNS(svgNamespace, name);
  Object.entries(attributes).forEach(([key, value]) =>
    element.setAttribute(key, value),
  );
  return element;
}

function createSvgText(value, x, y, className) {
  const text = createSvgElement("text", { x, y });
  text.classList.add(className);
  text.textContent = value;
  return text;
}

function setText(root, selector, value) {
  const element = root.querySelector(selector);
  if (element) element.textContent = value;
}
