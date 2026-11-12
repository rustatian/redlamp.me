import { execFileSync } from 'node:child_process';

// Builds the Go server one time for a run. Each test starts it with its own schema in e2e/fixtures.ts.
export default function globalSetup() {
	if (!process.env.TEST_DATABASE_URL) throw new Error('Set TEST_DATABASE_URL to a PostgreSQL connection string, for example postgres://postgres:test@127.0.0.1:54329/postgres?sslmode=disable.');
	execFileSync('go', ['build', '-buildvcs=false', '-o', '../.playwright/bin/', './cmd/api', './cmd/migrate'], { cwd: 'backend', stdio: 'inherit' });
}
