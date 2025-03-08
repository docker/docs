# Errors


> API response structure

```json
{
    "error": "Descriptive error message"
}
```

The Docker Cloud API uses the following error codes:


Error Code | Meaning
---------- | -------
400 | Bad Request -- There's a problem in the content of your request. Retrying the same request will fail.
401 | Unauthorized -- Your API key is wrong or your account has been deactivated.
402 | Payment Required -- You need to provide billing information to perform this request.
403 | Forbidden -- Quota limit exceeded. Contact support to request a quota increase.
404 | Not Found -- The requested object cannot be found.
405 | Method Not Allowed -- The endpoint requested does not implement the method sent.
409 | Conflict -- The object cannot be created or updated because another object exists with the same unique fields
415 | Unsupported Media Type -- Make sure you are using `Accept` and `Content-Type` headers as `application/json` and that the data your are `POST`-ing or `PATCH`-ing is in valid JSON format.
429 | Too Many Requests -- You are being throttled because of too many requests in a short period of time.
500 | Internal Server Error -- There was a server error while processing your request. Try again later, or contact support.
503 | Service Unavailable -- We're temporarily offline for maintenance. Try again later.
504 | Gateway Timeout -- Our API servers are at full capacity. Try again later.