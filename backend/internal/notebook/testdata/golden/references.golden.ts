// The source of references.json: the hidden test inputs and the reference answers of the learning set tasks. See vitest.config.ts for the command.
// The Go tests run the Go generators and the Go references and compare the results with this file.
import { it } from 'vitest';
import { parsePy, reprPy } from '../../../../../shared/answers';
import { mulberry32 } from '../../../../../src/leetcode/judge';
import { contentOf, sets, taskSlug, tasksOf } from '../../../../../src/leetcode/sets';

// The seeds of the hidden tests: the smallest, the largest and three other values.
const SEEDS = [0, 1, 20261112, 2147483648, 4294967295];

// A text of more than 512 characters is in the file as its SHA-256 and its length, so the file stays small.
async function compact(value: string): Promise<unknown> {
	if (value.length <= 512) return value;
	const digest = new Uint8Array(await crypto.subtle.digest('SHA-256', new TextEncoder().encode(value)));
	return { sha256: Array.from(digest, byte => byte.toString(16).padStart(2, '0')).join(''), length: value.length };
}

it('writes the hidden tests and the reference answers', async () => {
	const tasks: Record<string, unknown> = {};
	for (const task of sets.flatMap(tasksOf)) {
		const slug = taskSlug(task), content = contentOf(slug)!, { reference, generate } = content;
		if (!reference) continue;
		// The answer text of the acceptance check: stdout, or the Python literal of the return value.
		const answered = async (inputs: string[]) => Promise.all(inputs.map(async input => {
			const answer = reference.mode === 'stdio' ? reference.solve(input) : reprPy(reference.solve(...parsePy(input) as unknown[]));
			return { input: await compact(input), answer: await compact(answer) };
		}));
		tasks[slug] = {
			mode: reference.mode,
			arity: reference.mode === 'call' ? reference.solve.length : 0,
			cases: await answered(task.tests?.cases.map(item => item.input) ?? []),
			seeds: generate ? await Promise.all(SEEDS.map(async seed => ({ seed, tests: await answered(generate(mulberry32(seed))) }))) : [],
		};
	}
	console.log(`GOLDEN references ${JSON.stringify(tasks)}`);
	// The references run about 4,000 times, which takes longer than the default time limit of a test on a busy machine.
}, 120_000);
