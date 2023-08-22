import docsearch from '@docsearch/js';
import * as params from '@params'; // hugo dict

const { appid, apikey, indexname } = params

docsearch({
  container: '#docsearch',
  appId: appid,
  apiKey: apikey,
  indexName: indexname,
});
