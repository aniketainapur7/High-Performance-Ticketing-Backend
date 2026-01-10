import pkg from "pg"
import { env } from "./env";

const { Pool } = pkg;

export const pgPool = new Pool({
  host: env.postgres.host,
  user: env.postgres.user,
  password: env.postgres.password,
  database: env.postgres.database,
});