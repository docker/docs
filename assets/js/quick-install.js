// Browser platform detection selects examples, not installation eligibility.
export function detectPlatform(platform = "", userAgent = "", touchPoints = 0) {
  if (/Android|iPhone|iPad|iPod/i.test(userAgent)) return null;
  if (/Mac/i.test(platform) && touchPoints > 1) return null;
  if (/Win/i.test(platform)) return "windows";
  if (/Mac/i.test(platform)) return "macos";
  if (/Linux/i.test(platform)) return "linux";
  return null;
}

export const quickInstalls = {
  macos: [
    {
      product: "Docker Sandboxes",
      label: "macOS · Apple silicon · Homebrew",
      url: "/get-docker/sandboxes/",
      command: "brew trust docker/tap && brew install docker/tap/sbx",
    },
    {
      product: "Docker Desktop",
      label: "macOS · Homebrew",
      url: "/get-docker/desktop/",
      command: "brew install --cask docker-desktop",
    },
  ],
  windows: [
    {
      product: "Docker Sandboxes",
      label: "Windows 11 · x64 · WinGet",
      url: "/get-docker/sandboxes/",
      command: "winget install -h Docker.sbx",
    },
    {
      product: "Docker Desktop",
      label: "Windows · WinGet",
      url: "/get-docker/desktop/",
      command: "winget install -e --id Docker.DockerDesktop",
    },
  ],
  linux: [
    {
      product: "Docker Sandboxes",
      label: "Ubuntu 24.04+",
      url: "/get-docker/sandboxes/",
      command:
        "curl -fsSL https://get.docker.com | sudo REPO_ONLY=1 sh && sudo apt install docker-sbx",
    },
    {
      product: "Docker Engine",
      label: "Linux · dev/test",
      url: "/get-docker/engine/",
      command: "curl -fsSL https://get.docker.com | sudo sh",
    },
  ],
};

if (typeof document !== "undefined") {
  const platform = detectPlatform(
    navigator.userAgentData?.platform || navigator.platform,
    navigator.userAgent,
    navigator.maxTouchPoints,
  );
  const rows = document.querySelectorAll(".home-install-row");
  quickInstalls[platform]?.forEach((install, index) => {
    const row = rows[index];
    if (!row) return;
    row.querySelector("a").href = install.url;
    row.querySelector(".home-install-product").textContent = install.product;
    row.querySelector(".home-install-platform").textContent = install.label;
    row.querySelector("code").textContent = install.command;
    row.querySelector(".home-install-command").hidden = false;
    const button = row.querySelector("[data-copy-text]");
    button.dataset.copyText = install.command;
    button.setAttribute(
      "aria-label",
      `Copy ${install.product} install command`,
    );
  });
}
