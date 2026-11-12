// The source of game.json and sets.json. See vitest.config.ts for the command.
// Each scenario is a list of database rows. The Worker answers GET /leetcode/api/game for them, and the answer is the expected value of the Go test.
import { env, exports } from 'cloudflare:workers';
import { applyD1Migrations, reset } from 'cloudflare:test';
import { expect, it } from 'vitest';
import { addCalendarDays } from '../../../../../shared/dates';
import type { Acceptance, Difficulty, Role, TaskState, Understanding } from '../../../../../shared/leetcode';
import { sets, taskSlug, tasksOf } from '../../../../../src/leetcode/sets';
import { login } from '../../../../../test/client';

type Value = string | number | null;
interface Row { table: string; values: Record<string, Value> }
interface Scenario { name: string; role: Role; rows: Row[] }

const now = new Date(), today = now.toISOString().slice(0, 10);
const day = (offset: number) => addCalendarDays(today, offset);
const at = (offset: number, time = '12:00:00.000') => `${day(offset)}T${time}Z`;

function builder() {
	const rows: Row[] = [];
	let count = 0;
	const id = () => `00000000-0000-4000-8000-${String(++count).padStart(12, '0')}`;
	const add = (table: string, values: Record<string, Value>) => { rows.push({ table, values }); };
	const b = {
		rows,
		problem(difficulty: Difficulty = 'easy', options: { slug?: string; archived?: boolean } = {}): string {
			const problemId = id();
			add('problems', { id: problemId, slug: options.slug ?? `problem-${count}`, title: `Problem ${count}`, difficulty, topics: '[]', summary: '', search_text: '', archived_at: options.archived ? at(-1) : null, created_at: at(-40), updated_at: at(-40) });
			return problemId;
		},
		attempt(problemId: string, options: { state?: 'draft' | 'saved'; acceptance?: Acceptance; understanding?: Understanding; savedAt?: string } = {}): string {
			const attemptId = id(), state = options.state ?? 'saved';
			add('attempts', { id: attemptId, problem_id: problemId, state, document: '{}', search_text: '', acceptance: options.acceptance ?? 'accepted', understanding: options.understanding ?? 'with_help', created_at: at(-30), updated_at: at(-30), saved_at: state === 'saved' ? options.savedAt ?? at(-30) : null });
			return attemptId;
		},
		// A problem with a saved accepted attempt.
		accepted(difficulty: Difficulty = 'easy', options: { slug?: string; archived?: boolean; understanding?: Understanding; savedAt?: string } = {}): { problemId: string; attemptId: string } {
			const problemId = b.problem(difficulty, options);
			return { problemId, attemptId: b.attempt(problemId, { understanding: options.understanding, savedAt: options.savedAt }) };
		},
		review(problemId: string, result: Understanding, on: string): void {
			add('reviews', { id: id(), problem_id: problemId, result, note: '', reviewed_on: on, next_review_date: null, created_at: `${on}T10:00:00.000Z` });
		},
		homework(due: string | null): string {
			const homeworkId = id();
			add('homework', { id: homeworkId, title: `Homework ${count}`, instructions: '', due_date: due, created_at: at(-35), updated_at: at(-35) });
			return homeworkId;
		},
		task(homeworkId: string, problemId: string, state: TaskState, options: { createdAt?: string; completedAt?: string } = {}): string {
			const taskId = id();
			add('homework_tasks', { id: taskId, homework_id: homeworkId, problem_id: problemId, state, created_at: options.createdAt ?? at(-34), updated_at: at(-1), completed_at: options.completedAt ?? null });
			return taskId;
		},
		chapterTask(homeworkId: string, kind: 'chapter' | 'quiz', setSlug: string, chapter: number, state: TaskState, completedAt: string | null): void {
			add('homework_tasks', { id: id(), homework_id: homeworkId, problem_id: null, state, created_at: at(-34), updated_at: at(-1), kind, set_slug: setSlug, chapter, completed_at: completedAt });
		},
		submission(taskId: string, attemptId: string, createdAt: string): void {
			add('submissions', { id: id(), task_id: taskId, attempt_id: attemptId, task_version: count, created_at: createdAt });
		},
		card(correct: 0 | 1, answeredOn: string): void {
			add('card_answers', { id: id(), card_id: `card-${count}`, correct, interval_days: 1, answered_on: answeredOn, next_date: addCalendarDays(answeredOn, 1), created_at: `${answeredOn}T09:00:00.000Z` });
		},
		quiz(deck: string, correct: number, total: number, passed: 0 | 1, createdAt: string): void {
			add('quiz_attempts', { id: id(), deck, correct, total, passed, created_at: createdAt });
		},
		run(problemId: string, attemptId: string, passed: number, total: number, createdAt: string): void {
			add('test_runs', { id: id(), problem_id: problemId, attempt_id: attemptId, approach_id: 'approach', code_hash: 'hash', passed, total, run_on: createdAt.slice(0, 10), created_at: createdAt });
		},
		visit(value: string): void { add('visits', { day: value }); },
		badge(badgeId: string): void { add('earned_badges', { id: badgeId }); },
	};
	return b;
}

const scenarios: Scenario[] = [];
function scenario(name: string, role: Role, build: (b: ReturnType<typeof builder>) => void): void {
	const b = builder();
	build(b);
	scenarios.push({ name, role, rows: b.rows });
}
const repeat = (count: number, action: (index: number) => void) => { for (let index = 0; index < count; index++) action(index); };
const setSlugs = (slug: string) => tasksOf(sets.find(set => set.slug === slug)!).map(taskSlug);

scenario('an empty notebook read by the student', 'student', () => {});
scenario('an empty notebook read by the parent', 'parent', () => {});
scenario('accepted problems of each difficulty', 'student', b => {
	b.accepted('easy');
	b.accepted('medium');
	b.accepted('hard');
	b.accepted('hard', { archived: true });
	b.attempt(b.problem('hard'), { state: 'draft', acceptance: 'accepted' });
	b.attempt(b.problem('medium'), { acceptance: 'not_accepted' });
	b.attempt(b.problem('easy'), { acceptance: 'not_submitted' });
	b.problem('hard');
});
scenario('independent understanding counts one time for each accepted problem', 'student', b => {
	b.accepted('easy', { understanding: 'independent' });
	b.review(b.accepted('easy').problemId, 'independent', day(-3));
	const both = b.accepted('medium', { understanding: 'independent' }).problemId;
	b.review(both, 'independent', day(-4));
	b.review(both, 'independent', day(-2));
	b.review(b.accepted('easy').problemId, 'with_help', day(-3));
	const unaccepted = b.problem('easy');
	b.attempt(unaccepted, { acceptance: 'not_accepted', understanding: 'independent' });
	b.review(unaccepted, 'independent', day(-3));
});
scenario('completed homework tasks and recalls', 'student', b => {
	const homework = b.homework(day(5)), first = b.accepted('easy').problemId, second = b.accepted('medium').problemId;
	b.task(homework, first, 'completed', { completedAt: at(-2) });
	b.task(homework, second, 'completed', { completedAt: at(-1) });
	b.task(homework, b.problem('easy'), 'assigned');
	b.task(homework, b.problem('easy'), 'cancelled');
	b.chapterTask(homework, 'chapter', 'two-pointers', 1, 'completed', at(-2));
	b.chapterTask(homework, 'quiz', 'two-pointers', 1, 'assigned', null);
	b.review(first, 'with_help', day(-6));
	b.review(first, 'needs_practice', day(-6));
	b.review(first, 'independent', day(-3));
	b.review(second, 'with_help', day(-6));
});
scenario('cards of single answers and of the best passed attempt of each deck', 'student', b => {
	repeat(3, index => b.card(1, day(-index - 1)));
	repeat(2, index => b.card(0, day(-index - 1)));
	b.quiz('permutations-and-combinations-1', 4, 5, 1, at(-5));
	b.quiz('permutations-and-combinations-1', 5, 5, 1, at(-4));
	b.quiz('permutations-and-combinations-1', 3, 5, 0, at(-3));
	b.quiz('permutations-and-combinations-2', 2, 5, 0, at(-3));
	b.quiz('two-pointers-1', 8, 10, 1, at(-2));
});
scenario('test runs: passed problems, failed runs for each day, and the first try', 'student', b => {
	const first = b.accepted('easy'), second = b.accepted('easy'), third = b.accepted('easy'), fourth = b.problem('medium'), fourthAttempt = b.attempt(fourth, { acceptance: 'not_accepted' });
	repeat(5, index => b.run(first.problemId, first.attemptId, 2, 5, at(-4, `10:0${index}:00.000`)));
	repeat(2, index => b.run(first.problemId, first.attemptId, 4, 5, at(-3, `10:0${index}:00.000`)));
	b.run(first.problemId, first.attemptId, 5, 5, at(-3, '11:00:00.000'));
	b.run(first.problemId, first.attemptId, 5, 5, at(-2, '11:00:00.000'));
	b.run(second.problemId, second.attemptId, 0, 3, at(-2, '09:00:00.000'));
	b.run(third.problemId, third.attemptId, 0, 0, at(-2, '09:00:00.000'));
	b.run(fourth, fourthAttempt, 6, 6, at(-5, '08:00:00.000'));
	b.run(fourth, fourthAttempt, 1, 6, at(-4, '08:00:00.000'));
});
scenario('late homework tasks', 'student', b => {
	const past = b.homework(day(-1)), current = b.homework(day(0)), open = b.homework(null), solved = b.accepted('easy');
	b.task(past, b.problem(), 'assigned', { createdAt: at(-5) });
	b.task(past, b.problem(), 'completed', { createdAt: at(-5), completedAt: at(-1, '23:59:59.000') });
	b.task(past, b.problem(), 'completed', { createdAt: at(-5), completedAt: at(0, '00:00:01.000') });
	b.submission(b.task(past, solved.problemId, 'submitted', { createdAt: at(-5) }), solved.attemptId, at(-2));
	const lateWork = b.accepted('easy');
	b.submission(b.task(past, lateWork.problemId, 'submitted', { createdAt: at(-5) }), lateWork.attemptId, at(0, '01:00:00.000'));
	b.task(past, b.problem(), 'cancelled', { createdAt: at(-5) });
	b.task(past, b.problem(), 'assigned', { createdAt: at(0, '02:00:00.000') });
	b.task(past, b.problem(), 'in_progress', { createdAt: at(-1, '22:00:00.000') });
	b.task(current, b.problem(), 'assigned', { createdAt: at(-5) });
	b.task(open, b.problem(), 'assigned', { createdAt: at(-5) });
});
scenario('missed days and the longest run of visits', 'student', b => {
	for (const offset of [-12, -11, -10, -9, -8, -7, -6, -3, -2]) b.visit(day(offset));
});
scenario('visits read by the parent', 'parent', b => {
	for (const offset of [-4, -3, -1]) b.visit(day(offset));
});
scenario('a visit of today that is recorded before the read', 'student', b => {
	for (const offset of [-2, -1, 0]) b.visit(day(offset));
});
scenario('saved attempts on three days in a row', 'student', b => {
	for (const offset of [-5, -4, -3]) b.attempt(b.problem(), { acceptance: 'not_accepted', savedAt: at(offset) });
});
scenario('saved attempts on two days in a row and one later day', 'student', b => {
	for (const offset of [-5, -4, -2]) b.attempt(b.problem(), { acceptance: 'not_accepted', savedAt: at(offset) });
});
scenario('level one ends at 98 XP', 'parent', b => {
	repeat(9, () => b.accepted('easy'));
	repeat(4, index => b.card(1, day(-index - 1)));
});
scenario('level two starts at 100 XP', 'parent', b => { repeat(10, () => b.accepted('easy')); });
scenario('level five starts at 1000 XP', 'parent', b => { repeat(25, () => b.accepted('hard')); });
scenario('stored awards stay earned', 'parent', b => {
	for (const badge of ['first-hard', 'on-time-five', 'set-cleared', 'month-no-miss']) b.badge(badge);
});
scenario('all problems of one learning set are accepted', 'parent', b => {
	for (const slug of setSlugs('two-pointers')) b.accepted('medium', { slug });
});
scenario('all problems of two learning sets are accepted', 'parent', b => {
	for (const slug of [...setSlugs('two-pointers'), ...setSlugs('permutations-and-combinations')]) b.accepted('easy', { slug });
});
scenario('all problems of a learning set except one are accepted', 'parent', b => {
	for (const slug of setSlugs('two-pointers').slice(1)) b.accepted('medium', { slug });
	b.attempt(b.problem('medium', { slug: setSlugs('two-pointers')[0] }), { acceptance: 'not_accepted' });
});
scenario('five homework tasks done on or before the due date', 'parent', b => {
	const homework = b.homework(day(-2));
	repeat(3, () => b.task(homework, b.problem(), 'completed', { createdAt: at(-9), completedAt: at(-2, '20:00:00.000') }));
	repeat(2, () => { const work = b.accepted('easy'); b.submission(b.task(homework, work.problemId, 'submitted', { createdAt: at(-9) }), work.attemptId, at(-3)); });
	b.task(homework, b.problem(), 'completed', { createdAt: at(-9), completedAt: at(-1) });
	b.task(homework, b.problem(), 'cancelled', { createdAt: at(-9), completedAt: at(-3) });
});
scenario('four homework tasks done on or before the due date', 'parent', b => {
	const homework = b.homework(day(-2));
	repeat(4, () => b.task(homework, b.problem(), 'completed', { createdAt: at(-9), completedAt: at(-2, '20:00:00.000') }));
	b.task(homework, b.problem(), 'completed', { createdAt: at(-9), completedAt: at(-1) });
});
scenario('penalties larger than the gains give zero XP', 'student', b => {
	b.accepted('easy');
	const homework = b.homework(day(-3));
	repeat(3, () => b.task(homework, b.problem(), 'assigned', { createdAt: at(-9) }));
	b.visit(day(-10));
});
scenario('an accepted Codeforces problem', 'parent', b => { b.accepted('medium', { slug: 'cf-1-a' }); });
scenario('many achievements', 'parent', b => {
	const problems = Array.from({ length: 50 }, (_, index) => b.accepted('easy', { understanding: index < 5 ? 'independent' : 'with_help' }).problemId);
	repeat(10, index => b.review(problems[index], 'with_help', day(-index - 1)));
	repeat(50, index => b.card(1, day(-index % 20 - 1)));
});

async function load(rows: Row[]): Promise<void> {
	for (const { table, values } of rows) {
		const keys = Object.keys(values);
		await env.DB.prepare(`INSERT INTO ${table} (${keys.join(',')}) VALUES (${keys.map(() => '?').join(',')})`).bind(...keys.map(key => values[key])).run();
	}
}

it('writes the golden files', async () => {
	const results = [];
	for (const { name, role, rows } of scenarios) {
		await reset();
		await applyD1Migrations(env.DB, env.TEST_MIGRATIONS);
		await load(rows);
		const response = await exports.default.fetch('https://redlamp.me/leetcode/api/game', { headers: { Cookie: await login(role) } });
		expect(response.status, name).toBe(200);
		const visits = (await env.DB.prepare('SELECT day FROM visits ORDER BY day').all<{ day: string }>()).results.map(row => row.day);
		const badges = (await env.DB.prepare('SELECT id FROM earned_badges ORDER BY id').all<{ id: string }>()).results.map(row => row.id);
		results.push({ name, role, rows, game: await response.json(), visits, badges });
	}
	// The scenarios use dates relative to today, so the whole run must be on one UTC day.
	expect(new Date().toISOString().slice(0, 10)).toBe(today);
	console.log(`GOLDEN game ${JSON.stringify({ now: now.toISOString(), scenarios: results })}`);
	console.log(`GOLDEN sets ${JSON.stringify(sets.map(set => ({ slug: set.slug, title: set.title, chapters: set.chapters.map(chapter => ({ title: chapter.title, tasks: chapter.tasks.map(task => ({ url: task.url, number: task.number ?? null, title: task.title, difficulty: task.difficulty, topics: task.topics, note: task.note })) })) })))}`);
});
