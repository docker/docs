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
      label: "macOS · Installer",
      url: "/get-docker/desktop/",
      downloads: [
        {
          label: "Apple silicon",
          url: "https://desktop.docker.com/mac/main/arm64/Docker.dmg",
        },
        {
          label: "Intel",
          url: "https://desktop.docker.com/mac/main/amd64/Docker.dmg",
        },
      ],
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
      label: "Windows · Installer",
      url: "/get-docker/desktop/",
      downloads: [
        {
          label: "x86_64",
          url: "https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe",
        },
        {
          label: "Arm (Early Access)",
          url: "https://desktop.docker.com/win/main/arm64/Docker%20Desktop%20Installer.exe",
        },
      ],
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
    if (install.downloads) {
      const downloads = row.querySelector(".home-install-downloads");
      for (const download of install.downloads) {
        const item = document.createElement("li");
        const link = document.createElement("a");
        link.className = "button";
        link.href = download.url;
        link.textContent = `Download for ${download.label}`;
        link.setAttribute(
          "aria-label",
          `Download ${install.product} for ${download.label}`,
        );
        link.setAttribute(
          "marlin-label",
          platform === "macos" ? "download_dmg_click" : "download_exe_click",
        );
        item.append(link);
        downloads.append(item);
      }
      downloads.hidden = false;
      return;
    }
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
