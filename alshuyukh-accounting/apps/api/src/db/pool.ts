import pg from 'pg';

/** Creates a pool whose sessions always use UTC. */
export function createPool(connectionString: string): pg.Pool {
  return new pg.Pool({
    connectionString,
    max: 20,
    options: '-c timezone=UTC',
  });
}

// Return NUMERIC as string so money is never converted to a JS float.
// Business logic converts with a decimal library (introduced in Phase 2).
pg.types.setTypeParser(pg.types.builtins.NUMERIC, (v) => v);
// Return bigint (e.g. COUNT(*)) as string; callers convert explicitly.
pg.types.setTypeParser(pg.types.builtins.INT8, (v) => v);
// Return DATE as 'YYYY-MM-DD' text. The default converts it to a JS Date in
// the server's local time zone, which can shift a calendar date by a day.
pg.types.setTypeParser(pg.types.builtins.DATE, (v) => v);
