// The source of sets.json: the whole content of the learning sets of the Worker. See vitest.config.ts for the command.
// The Go tests use it as the content of the server and compare the Go content with it.
import { it } from 'vitest';
import { sets } from '../../../../../src/leetcode/sets';

it('writes the learning sets', () => {
	console.log(`GOLDEN sets ${JSON.stringify(sets)}`);
});
