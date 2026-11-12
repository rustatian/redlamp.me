import { test as base, expect, type Page } from '@playwright/test';
import { execFileSync, spawn } from 'node:child_process';
import { createHash, randomBytes } from 'node:crypto';
import { once } from 'node:events';
import { createServer, type AddressInfo } from 'node:net';
import { Client } from 'pg';
import type { Attempt, Problem, Role, TutorKind } from '../shared/leetcode';

// The Go commands that e2e/global-setup.ts builds.
const bin = '.playwright/bin';
const passwords: Record<Role, string> = { parent: 'parent-test-password-001', student: 'student-test-password-01' };
// The LEETCODE_ACCOUNTS secret of the test accounts.
const accounts = JSON.stringify(Object.fromEntries(Object.entries(passwords).map(([role, password]) => [role, { username: `${role}-test`, passwordHash: createHash('sha256').update(password).digest('hex') }])));

export async function signIn(page: Page, role: Role = 'student') {
	await page.goto('/leetcode/login');
	await page.getByLabel('Имя пользователя').fill(`${role}-test`);
	await page.getByLabel('Пароль').fill(passwords[role]);
	await page.getByRole('button', { name: 'Войти', exact: true }).click();
	await expect(page.getByRole('button', { name: 'Выйти' })).toBeVisible();
}

export async function addTask(db: Client, problemId: string, state = 'assigned', dueDate: string | null = null) {
	const homeworkId = crypto.randomUUID(), taskId = crypto.randomUUID();
	await db.query("INSERT INTO homework (id,title,instructions,due_date,created_at,updated_at) VALUES ($1,'Week 41','Explain',$2,'now','now')", [homeworkId, dueDate]);
	await db.query("INSERT INTO homework_tasks (id,homework_id,problem_id,state,created_at,updated_at) VALUES ($1,$2,$3,$4,'now','now')", [taskId, homeworkId, problemId, state]);
	return { homeworkId, taskId };
}

// One minute between messages, so the order is fixed.
export async function addMessages(db: Client, kind: TutorKind, ref: string, messages: { author: 'student' | 'assistant'; body: string; quote?: string }[], start = '2026-10-09T10:00:00.000Z') {
	for (const [i, m] of messages.entries()) await db.query('INSERT INTO tutor_messages (id,context_kind,context_ref,author,body,quote,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)', [crypto.randomUUID(), kind, ref, m.author, m.body, m.quote ?? null, new Date(Date.parse(start) + i * 60_000).toISOString()]);
}

async function admin(sql: string) {
	const client = new Client({ connectionString: process.env.TEST_DATABASE_URL });
	await client.connect();
	try { await client.query(sql); } finally { await client.end(); }
}

async function freePort() {
	const server = createServer().listen(0, '127.0.0.1');
	await once(server, 'listening');
	const { port } = server.address() as AddressInfo;
	server.close();
	await once(server, 'close');
	return port;
}

// Each test gets a new schema and its own server, so the sign-in limit of the server (10 attempts in 60 seconds for each address) applies to one test only.
export const test = base.extend<{ server: { origin: string; schema: string }; db: Client; notebook: { problem: Problem; attempt: Attempt } }>({
	server: async ({ browserName }, use, testInfo) => {
		if (browserName !== 'chromium') throw new Error('Use Chromium for these checks.');
		const schema = `e2e_${randomBytes(8).toString('hex')}`, port = await freePort(), origin = `http://localhost:${port}`;
		// pgx sends the unknown URL parameters as runtime parameters. Aurora DSQL has only REPEATABLE READ.
		const url = `${process.env.TEST_DATABASE_URL}${process.env.TEST_DATABASE_URL!.includes('?') ? '&' : '?'}search_path=${schema}&default_transaction_isolation=repeatable%20read`;
		await admin(`CREATE SCHEMA ${schema}`);
		let log = '';
		try {
			execFileSync(`${bin}/migrate`, ['-database-url', url], { stdio: 'pipe' });
			const env: Record<string, string | undefined> = { ...process.env, PUBLIC_ORIGIN: origin, LEETCODE_ACCOUNTS: accounts, HTTPS_PROXY: 'http://127.0.0.1:9' };
			// Without a key the server has no tutor client. The proxy on the discard port refuses each request to leetcode.com.
			for (const name of ['ANTHROPIC_API_KEY', 'ANTHROPIC_API_KEY_PARAMETER', 'ANTHROPIC_AUTH_TOKEN', 'ANTHROPIC_PROFILE', 'ANTHROPIC_BASE_URL']) delete env[name];
			const api = spawn(`${bin}/api`, ['-listen', `127.0.0.1:${port}`, '-database-url', url, '-assets', 'dist'], { env, stdio: ['ignore', 'ignore', 'pipe'] }), exited = once(api, 'exit');
			api.stderr.setEncoding('utf8').on('data', (chunk: string) => { log += chunk; });
			try {
				for (let i = 0; ; i++) {
					if ((api.exitCode ?? api.signalCode) !== null) throw new Error(`The server stopped:\n${log}`);
					if (await fetch(`http://127.0.0.1:${port}/leetcode/login`).then(() => true, () => false)) break;
					if (i === 100) throw new Error(`The server does not answer:\n${log}`);
					await new Promise(resolve => setTimeout(resolve, 50));
				}
				await use({ origin, schema });
			} finally {
				api.kill();
				await exited;
			}
		} finally {
			if (testInfo.status !== testInfo.expectedStatus) await testInfo.attach('server log', { body: log, contentType: 'text/plain' });
			await admin(`DROP SCHEMA ${schema} CASCADE`);
		}
	},
	baseURL: async ({ server }, use) => { await use(server.origin); },
	db: async ({ server }, use) => {
		const db = new Client({ connectionString: process.env.TEST_DATABASE_URL, options: `-c search_path=${server.schema}` });
		await db.connect();
		try { await use(db); } finally { await db.end(); }
	},
	notebook: async ({ page, baseURL }, use) => {
		await signIn(page);
		const headers = { Origin: new URL(baseURL!).origin };
		const created = await page.request.post('/leetcode/api/problems', { headers, data: { url: 'https://leetcode.com/problems/two-sum/', number: 1, title: 'Two Sum', difficulty: 'easy', topics: ['Arrays'], summary: 'Find two values with the target sum.' } });
		expect(created.status()).toBe(201);
		const problem = await created.json() as Problem;
		const response = await page.request.post(`/leetcode/api/problems/${problem.id}/attempts`, { headers, data: {} });
		expect(response.status()).toBe(201);
		const attempt = await response.json() as Attempt;
		// The server gets the two-sum examples and statement from leetcode.com. The routes cover every page of the context. A test that needs them adds its own page route, which comes first.
		await page.context().route('**/leetcode/api/problems/*/tests', route => route.fulfill({ json: { spec: null, source: 'none', standard: [], own: [] } }));
		await page.context().route('**/leetcode/api/problems/*/statement', route => route.fulfill({ json: { html: null, text: null } }));
		await page.goto(`/leetcode/problems/${problem.id}`);
		await expect(page.getByLabel('Основная идея', { exact: true })).toBeVisible();
		await use({ problem, attempt });
	},
});
export { expect };
