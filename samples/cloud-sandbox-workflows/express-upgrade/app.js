import express from "express";

const app = express();

app.get("/api/status", (req, res) => {
  res.json({ status: "ok" });
});

app.get("*", (req, res) => {
  res
    .type("html")
    .send("<!doctype html><title>Client app</title><h1>Client app</h1>");
});

export default app;
