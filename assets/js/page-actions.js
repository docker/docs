export default () => ({
  status: "",
  async copy(url) {
    this.status = "Copying…";
    try {
      const response = await fetch(url);
      if (!response.ok) throw new Error("Markdown unavailable");
      await navigator.clipboard.writeText(await response.text());
      this.status = "Markdown copied";
    } catch {
      this.status = "Could not copy. Use View Markdown to open the page.";
    }
  },
});
