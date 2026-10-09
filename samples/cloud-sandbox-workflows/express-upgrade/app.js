import express from "express";

const app = express();
const tickets = [
  { id: 1, title: "Reset password", status: "open", assignee: "alex" },
  {
    id: 2,
    title: "Update billing address",
    status: "closed",
    assignee: "alex",
  },
  { id: 3, title: "Restore project access", status: "open", assignee: "sam" },
];

app.use(express.json());

function findTickets(filters = {}) {
  return tickets.filter(
    (ticket) =>
      (!filters.status || ticket.status === filters.status) &&
      (!filters.assignee || ticket.assignee === filters.assignee),
  );
}

app.get("/api/status", (req, res) => {
  res.json({ status: "ok" });
});

app.get("/api/tickets", (req, res) => {
  res.json({ tickets: findTickets(req.query.filters) });
});

app.post("/api/tickets/search", (req, res) => {
  const { filters } = req.body;
  res.json({ tickets: findTickets(filters) });
});

app.get("*", (req, res) => {
  res
    .type("html")
    .send("<!doctype html><title>Client app</title><h1>Client app</h1>");
});

export default app;
