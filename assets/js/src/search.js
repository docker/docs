import docsearch from "@docsearch/js";
import * as params from "@params"; // hugo dict

const { appid, apikey, indexname } = params;

docsearch({
  container: "#docsearch",
  appId: appid,
  apiKey: apikey,
  indexName: indexname,
  transformItems(items) {
    return items.map((item) => ({
      ...item,
      url: item.url.replace("https://docs.docker.com", ""),
    }));
  },
});
