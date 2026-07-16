#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <base-url> [sample-url ...]" >&2
  exit 1
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required" >&2
  exit 1
fi

if ! command -v rg >/dev/null 2>&1; then
  echo "rg is required" >&2
  exit 1
fi

BASE_URL="${1%/}"
shift || true
SAMPLE_SIZE="${SAMPLE_SIZE:-12}"
LLMS_SAMPLE_SIZE="${LLMS_SAMPLE_SIZE:-2}"
CHECK_TOOL_MANIFESTS="${CHECK_TOOL_MANIFESTS:-1}"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

count_matches() {
  local pattern="$1"
  local file="$2"
  rg -o "$pattern" "$file" 2>/dev/null | wc -l | tr -d ' ' || true
}

header_value() {
  local header_file="$1"
  local name="$2"
  awk -F': ' -v target="$name" '
    tolower($1) == tolower(target) { value = $2 }
    END {
      gsub(/\r/, "", value)
      print value
    }
  ' "$header_file"
}

normalize_text() {
  printf '%s' "$1" \
    | tr '[:upper:]' '[:lower:]' \
    | sed -E 's/[[:space:]]+/ /g; s/^[[:space:]]+//; s/[[:space:]]+$//; s/ \| docker docs$//'
}

code_fence_stats() {
  local file="$1"
  awk '
    BEGIN { in_block = 0; total = 0; tagged = 0 }
    /^```/ {
      line = $0
      sub(/^```[[:space:]]*/, "", line)
      if (!in_block) {
        total++
        if (line != "") {
          tagged++
        }
        in_block = 1
      } else {
        in_block = 0
      }
    }
    END {
      printf "%d\t%d\n", total, tagged
    }
  ' "$file"
}

resource_probe() {
  local url="$1"
  local label="$2"
  local body="$TMPDIR/resource-body"
  local headers="$TMPDIR/resource-headers"
  local status
  local content_type
  local bytes

  status="$(curl -sS -L -o "$body" -D "$headers" -w '%{http_code}' "$url" || true)"
  content_type="$(header_value "$headers" "content-type")"
  bytes="$(wc -c < "$body" | tr -d ' ')"

  printf '%s\t%s\t%s\t%s\t%s\n' "$label" "$url" "$status" "$content_type" "$bytes"
}

page_probe() {
  local url="$1"
  local html="$TMPDIR/page-html"
  local html_headers="$TMPDIR/page-html-headers"
  local md="$TMPDIR/page-md"
  local md_headers="$TMPDIR/page-md-headers"
  local direct_md="$TMPDIR/page-direct-md"
  local direct_md_headers="$TMPDIR/page-direct-md-headers"
  local alt_md="$TMPDIR/page-alt-md"
  local alt_md_headers="$TMPDIR/page-alt-md-headers"
  local status
  local content_type
  local final_url
  local h1_count
  local main_count
  local article_count
  local canonical_count
  local jsonld_count
  local md_alt
  local md_alt_url
  local direct_md_url
  local html_title
  local html_h1
  local md_h1
  local md_status
  local md_content_type
  local md_bytes
  local direct_md_status
  local direct_md_content_type
  local md_alt_status="na"
  local md_alt_content_type="na"
  local title_md_h1_match="no"
  local html_h1_md_h1_match="no"
  local code_blocks_total
  local code_blocks_tagged

  status="$(
    curl -sS -L -o "$html" -D "$html_headers" \
      -w '%{http_code}\t%{url_effective}' "$url" || true
  )"
  content_type="$(header_value "$html_headers" "content-type")"
  final_url="${status#*$'\t'}"
  status="${status%%$'\t'*}"

  h1_count="$(count_matches '<h1[ >]' "$html")"
  main_count="$(count_matches '<main[ >]' "$html")"
  article_count="$(count_matches '<article[ >]' "$html")"
  canonical_count="$(count_matches 'rel=canonical' "$html")"
  jsonld_count="$(count_matches 'application/ld\+json' "$html")"
  md_alt="$(
    rg -o 'type=text/markdown href=[^ >]+|href=[^ >]+[^>]*type=text/markdown' \
      "$html" -m 1 2>/dev/null | sed -E 's/.*href=([^ >]+).*/\1/' || true
  )"

  md_status="$(curl -sS -L -H 'Accept: text/markdown' -o "$md" -D "$md_headers" -w '%{http_code}' "$url" || true)"
  md_content_type="$(header_value "$md_headers" "content-type")"
  md_bytes="$(wc -c < "$md" | tr -d ' ')"
  direct_md_url="$(printf '%s' "$final_url" | sed 's#/$##').md"
  direct_md_status="$(curl -sS -L -o "$direct_md" -D "$direct_md_headers" -w '%{http_code}' "$direct_md_url" || true)"
  direct_md_content_type="$(header_value "$direct_md_headers" "content-type")"
  html_title="$(rg -o '<title>[^<]+' "$html" -m 1 2>/dev/null | sed 's/<title>//' || true)"
  html_h1="$(rg -o '<h1[^>]*>[^<]+' "$html" -m 1 2>/dev/null | sed -E 's/<h1[^>]*>//' || true)"
  md_h1="$(awk '/^# / { sub(/^# /, ""); print; exit }' "$md" || true)"

  if [[ -n "$html_title" && -n "$md_h1" ]]; then
    if [[ "$(normalize_text "$html_title")" == "$(normalize_text "$md_h1")" ]]; then
      title_md_h1_match="yes"
    fi
  fi

  if [[ -n "$html_h1" && -n "$md_h1" ]]; then
    if [[ "$(normalize_text "$html_h1")" == "$(normalize_text "$md_h1")" ]]; then
      html_h1_md_h1_match="yes"
    fi
  fi

  IFS=$'\t' read -r code_blocks_total code_blocks_tagged < <(code_fence_stats "$md")

  if [[ -n "$md_alt" ]]; then
    if [[ "$md_alt" =~ ^https?:// ]]; then
      md_alt_url="$md_alt"
    elif [[ "$md_alt" == /* ]]; then
      md_alt_url="${BASE_URL}${md_alt}"
    else
      md_alt_url="${BASE_URL}/${md_alt}"
    fi
    md_alt_status="$(curl -sS -L -o "$alt_md" -D "$alt_md_headers" -w '%{http_code}' "$md_alt_url" || true)"
    md_alt_content_type="$(header_value "$alt_md_headers" "content-type")"
  else
    md_alt_url="na"
  fi

  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
    "$url" \
    "$status" \
    "$content_type" \
    "$final_url" \
    "$h1_count" \
    "$main_count" \
    "$article_count" \
    "$canonical_count" \
    "$jsonld_count" \
    "$md_status" \
    "$md_content_type" \
    "$md_bytes" \
    "$direct_md_url" \
    "$direct_md_status" \
    "$direct_md_content_type" \
    "$md_alt_url" \
    "$md_alt_status" \
    "$md_alt_content_type" \
    "$title_md_h1_match" \
    "$html_h1_md_h1_match" \
    "$code_blocks_total" \
    "$code_blocks_tagged"
}

llms_urls() {
  local llms="$TMPDIR/llms-sample.txt"
  local llms_status

  llms_status="$(curl -sS -L -o "$llms" -w '%{http_code}' "$BASE_URL/llms.txt" || true)"
  if [[ "$llms_status" == "200" ]]; then
    rg -o '\(https?://[^)]+\)' "$llms" 2>/dev/null \
      | tr -d '()' \
      | rg "^${BASE_URL//./\\.}" \
      | rg -v '/404\.html$|/search/?$|\.xml$|\.txt$' \
      | awk -v limit="$LLMS_SAMPLE_SIZE" '!seen[$0]++ && NR <= limit { print }'
  fi
}

sitemap_urls() {
  local sitemap="$TMPDIR/sitemap.xml"
  local sitemap_status

  sitemap_status="$(curl -sS -L -o "$sitemap" -w '%{http_code}' "$BASE_URL/sitemap.xml" || true)"
  if [[ "$sitemap_status" == "200" ]]; then
    rg -o '<loc>[^<]+' "$sitemap" \
      | sed 's/<loc>//' \
      | rg "^${BASE_URL//./\\.}" \
      | rg -v '/404\.html$|/search/?$|\.xml$|\.txt$' \
      | awk '!seen[$0]++ { print }'
  fi
}

sample_urls() {
  if [[ $# -gt 0 ]]; then
    printf '%s\n' "$@"
    return
  fi

  local sample_file="$TMPDIR/sampled-urls.txt"

  {
    llms_urls
    sitemap_urls
  } | awk -v limit="$SAMPLE_SIZE" '
    !seen[$0]++ {
      print
      count++
      if (count >= limit) {
        exit
      }
    }
  ' > "$sample_file"

  if [[ ! -s "$sample_file" ]]; then
    printf '%s/\n' "$BASE_URL"
  else
    cat "$sample_file"
  fi
}

printf 'META\tbase-url\t%s\n' "$BASE_URL"
if [[ $# -gt 0 ]]; then
  printf 'META\tsample-source\texplicit\n'
else
  printf 'META\tsample-source\tllms-and-sitemap-or-homepage\n'
fi
printf '\nSITEWIDE\n'
printf 'label\turl\tstatus\tcontent-type\tbytes\n'
resource_probe "$BASE_URL/llms.txt" "llms.txt"
resource_probe "$BASE_URL/llms-full.txt" "llms-full.txt"
resource_probe "$BASE_URL/robots.txt" "robots.txt"
resource_probe "$BASE_URL/sitemap.xml" "sitemap.xml"
if [[ "$CHECK_TOOL_MANIFESTS" == "1" ]]; then
  resource_probe "$BASE_URL/.well-known/ai-plugin.json" "ai-plugin.json"
  resource_probe "$BASE_URL/.well-known/agent.json" "agent.json"
  resource_probe "$BASE_URL/.well-known/agents.json" "agents.json"
fi

printf '\nPAGES\n'
printf 'url\tstatus\tcontent-type\tfinal-url\th1\tmain\tarticle\tcanonical\tjsonld\tmd-negotiate-status\tmd-negotiate-content-type\tmd-bytes\tmd-direct-url\tmd-direct-status\tmd-direct-content-type\tmd-alt-url\tmd-alt-status\tmd-alt-content-type\ttitle-md-h1-match\th1-md-h1-match\tcode-blocks-total\tcode-blocks-tagged\n'
while IFS= read -r page_url; do
  [[ -z "$page_url" ]] && continue
  page_probe "$page_url"
done < <(sample_urls "$@")
