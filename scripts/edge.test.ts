import { expect, test } from 'bun:test';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

interface Result { uri?: string; statusCode?: number; headers?: { location?: { value: string } } }

// The code of the viewer request function is the FunctionCode block of the edge template.
const template = readFileSync(resolve(import.meta.dir, '../infra/edge.yaml'), 'utf8');
const code = /FunctionCode: \|\n((?: {8}.*\n|\n)+)/.exec(template)?.[1] ?? '';
const handler = new Function(`${code}\nreturn handler;`)() as (event: unknown) => Promise<Result>;

test.each([
	{ name: 'the main page', host: 'redlamp.me', uri: '/', query: undefined, file: '/index.html', location: undefined },
	{ name: 'a page without its extension', host: 'redlamp.me', uri: '/examples', query: undefined, file: '/examples.html', location: undefined },
	{ name: 'a page with its extension', host: 'redlamp.me', uri: '/examples.html', query: undefined, file: '/examples.html', location: undefined },
	{ name: 'a page in a directory', host: 'redlamp.me', uri: '/mathtools/primefactorization', query: undefined, file: '/mathtools/primefactorization.html', location: undefined },
	{ name: 'a file of the public site', host: 'redlamp.me', uri: '/main_page_style.css', query: undefined, file: '/main_page_style.css', location: undefined },
	{ name: 'the notebook', host: 'redlamp.me', uri: '/leetcode', query: undefined, file: '/leetcode', location: undefined },
	{ name: 'the notebook with the old extension', host: 'redlamp.me', uri: '/leetcode.html', query: undefined, file: '/leetcode.html', location: undefined },
	{ name: 'a notebook page', host: 'redlamp.me', uri: '/leetcode/problems', query: undefined, file: '/leetcode/problems', location: undefined },
	{ name: 'a notebook asset', host: 'redlamp.me', uri: '/leetcode/assets/main.js', query: undefined, file: '/leetcode/assets/main.js', location: undefined },
	{ name: 'www on the main page', host: 'www.redlamp.me', uri: '/', query: undefined, file: undefined, location: 'https://redlamp.me/' },
	{ name: 'www on a public page', host: 'www.redlamp.me', uri: '/examples', query: undefined, file: undefined, location: 'https://redlamp.me/examples' },
	{ name: 'www on a notebook page with a query', host: 'www.redlamp.me', uri: '/leetcode/login', query: 'return=%2Fleetcode%2Fhomework', file: undefined, location: 'https://redlamp.me/leetcode/login?return=%2Fleetcode%2Fhomework' },
])('$name', async ({ host, uri, query, file, location }) => {
	const result = await handler({ request: { uri, headers: { host: { value: host } }, rawQueryString: () => query } });
	expect(result.uri).toBe(file);
	expect(result.headers?.location?.value).toBe(location);
	expect(result.statusCode).toBe(location ? 302 : undefined);
});
