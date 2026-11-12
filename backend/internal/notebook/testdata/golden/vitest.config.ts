// One-time generator of the golden files for the Go tests. It runs the TypeScript code of the Worker and writes the results.
// Run from the repository root:
//   bun run build
//   bunx vitest run --config backend/internal/notebook/testdata/golden/vitest.config.ts
// It writes ../game.json and ../sets.json. The dates of game.json are relative to the day of the run, and the file records that moment.
import { writeFileSync } from 'node:fs';
import path from 'node:path';
import { cloudflareTest, readD1Migrations } from '@cloudflare/vitest-plugin';
import { defineConfig } from 'vitest/config';
import { testAccountSecret } from '../../../../../test/fixtures.ts';

const root = path.resolve(import.meta.dirname, '../../../../..');
export default defineConfig({
	root,
	// The generator runs in the runtime of the Worker tests, with the D1 schema and the test accounts.
	plugins: [cloudflareTest({
		wrangler: { configPath: path.join(root, 'wrangler.toml') },
		remoteBindings: false,
		miniflare: { bindings: { LEETCODE_ACCOUNTS: testAccountSecret, TUTOR_TOKEN_KEY: btoa('t'.repeat(32)), TEST_MIGRATIONS: await readD1Migrations(path.join(root, 'migrations')) } },
	})],
	test: {
		include: ['backend/internal/notebook/testdata/golden/*.golden.ts'],
		silent: false,
		// The Worker runtime has no file system, so the test prints each file as one line and this process writes it.
		onConsoleLog(log) {
			const files = [...log.matchAll(/^GOLDEN (\w+) (.+)$/gm)];
			for (const [, name, json] of files) writeFileSync(path.join(import.meta.dirname, '..', `${name}.json`), `${JSON.stringify(JSON.parse(json), null, '\t')}\n`);
			if (files.length) return false;
		},
	},
});
