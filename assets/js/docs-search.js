let pagefind;
export default () => ({
  query: "",
  results: [],
  matches: [],
  timer: null,
  loading: false,
  status: "",
  sequence: 0,
  queueSearch(query, composing = false) {
    this.query = query;
    ++this.sequence;
    clearTimeout(this.timer);
    this.loading = Boolean(query.trim());
    if (!query.trim()) {
      this.results = [];
      this.matches = [];
      this.status = "";
    }
    if (query.trim() && !composing)
      this.timer = setTimeout(() => this.search(), 400);
  },
  onEnter(event) {
    if (event.isComposing || event.keyCode === 229) return;
    event.preventDefault();
    event.stopPropagation();
    if (event.metaKey || event.ctrlKey) this.ask();
    else this.search();
  },
  destroy() {
    clearTimeout(this.timer);
    ++this.sequence;
  },
  async search() {
    clearTimeout(this.timer);
    const sequence = ++this.sequence;
    const query = this.query.trim();
    if (!query) {
      this.results = [];
      this.matches = [];
      this.status = "";
      this.loading = false;
      return;
    }
    this.loading = true;
    if (!this.results.length) this.status = "Searching…";
    try {
      pagefind ??= import("/pagefind/pagefind.js")
        .then(async (api) => {
          await api.options({
            ranking: {
              termFrequency: 0,
              termSimilarity: 2,
              pageLength: 0,
              termSaturation: 1,
              metaWeights: { title: 10, description: 4, keywords: 6 },
            },
          });
          return api;
        })
        .catch((error) => {
          pagefind = undefined;
          throw error;
        });
      const api = await pagefind;
      if (sequence !== this.sequence) return;
      const response = await api.search(query);
      if (sequence !== this.sequence) return;
      const results = await Promise.all(
        response.results.slice(0, 12).map((result) => result.data()),
      );
      if (sequence !== this.sequence) return;
      this.matches = response.results;
      this.results = results;
      this.status = this.matches.length
        ? `Showing ${this.results.length} of ${this.matches.length} results`
        : "No results. Try another search or ask Docker.";
    } catch {
      if (sequence === this.sequence)
        this.status = "Search is unavailable. You can still ask Docker.";
    } finally {
      if (sequence === this.sequence) this.loading = false;
    }
  },
  async loadMore(focus = false) {
    const sequence = this.sequence;
    const offset = this.results.length;
    this.loading = true;
    try {
      const results = await Promise.all(
        this.matches.slice(offset, offset + 12).map((result) => result.data()),
      );
      if (sequence !== this.sequence) return;
      this.results.push(...results);
      this.status = this.matches.length
        ? `Showing ${this.results.length} of ${this.matches.length} results`
        : "No results. Try another search or ask Docker.";
      if (focus)
        this.$nextTick(() => {
          if (sequence !== this.sequence) return;
          const links = this.$root.querySelectorAll(".search-results li a");
          links[offset]?.focus();
        });
    } catch {
      if (sequence === this.sequence)
        this.status = "Could not load results. Try again or ask Docker.";
    } finally {
      if (sequence === this.sequence) this.loading = false;
    }
  },
  ask() {
    const query = this.query.trim();
    if (!query) return;
    document.getElementById("docs-search")?.close();
    this.$store.gordon.open(query, true);
  },
});
