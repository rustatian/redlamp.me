import { expect, it } from 'vitest';
import type { BankEntry, Difficulty } from '../shared/leetcode';
import { draw, freeDifficulties, type Picked } from './contestDraw';

// Math and Greedy have problems in all three difficulties, Geometry only at 1200. Difficulty and points follow the rating: 800 easy 3, 900-1100 medium 4, 1200 hard 6.
const bank: BankEntry[] = [
	{ slug: 'cf-1-a', rating: 800, topics: ['Math'], difficulty: 'easy', points: 3 },
	{ slug: 'cf-2-a', rating: 800, topics: ['Greedy'], difficulty: 'easy', points: 3 },
	{ slug: 'cf-3-b', rating: 900, topics: ['Math'], difficulty: 'medium', points: 4 },
	{ slug: 'cf-4-b', rating: 1100, topics: ['Math', 'Greedy'], difficulty: 'medium', points: 4 },
	{ slug: 'cf-5-c', rating: 1000, topics: ['Greedy'], difficulty: 'medium', points: 4 },
	{ slug: 'cf-6-d', rating: 1200, topics: ['Math'], difficulty: 'hard', points: 6 },
	{ slug: 'cf-7-d', rating: 1200, topics: ['Greedy', 'Geometry'], difficulty: 'hard', points: 6 },
];
// The picks of the bank problems with these slugs, each from its first topic.
const picks = (...slugs: string[]): Picked[] => slugs.map(slug => {
	const { topics, difficulty, rating, points } = bank.find(entry => entry.slug === slug)!;
	return { slug, topic: topics[0], difficulty, rating, points };
});

const frees: { name: string; picked: Picked[]; topic: string; free: Difficulty[] }[] = [
	{ name: 'hard place taken leaves easy and medium', picked: picks('cf-6-d'), topic: 'Math', free: ['easy', 'medium'] },
	{ name: 'two medium picks close medium', picked: picks('cf-3-b', 'cf-5-c'), topic: 'Math', free: ['easy', 'hard'] },
	{ name: 'topic with problems only in taken places', picked: picks('cf-6-d'), topic: 'Geometry', free: [] },
];
it.each(frees)('$name', row => { expect(freeDifficulties(bank, row.picked, row.topic)).toEqual(row.free); });

const draws: { name: string; picked: Picked[]; topic: string; random: number; drawn: Picked | null }[] = [
	{ name: 'topic with problems only in taken places draws nothing', picked: picks('cf-6-d'), topic: 'Geometry', random: 0, drawn: null },
	// Greedy still has cf-2-a, cf-4-b and cf-7-d, but the 4 places are full.
	{ name: 'four picks', picked: picks('cf-1-a', 'cf-3-b', 'cf-5-c', 'cf-6-d'), topic: 'Greedy', random: 0, drawn: null },
	// 0.4 selects medium of the 3 free difficulties, then the first of the medium Math problems: cf-3-b, if it were not picked.
	{ name: 'picked slug is not drawn again', picked: picks('cf-3-b'), topic: 'Math', random: 0.4, drawn: { slug: 'cf-4-b', topic: 'Math', difficulty: 'medium', rating: 1100, points: 4 } },
	{ name: 'random 0.99 picks the last free difficulty', picked: [], topic: 'Math', random: 0.99, drawn: { slug: 'cf-6-d', topic: 'Math', difficulty: 'hard', rating: 1200, points: 6 } },
];
it.each(draws)('$name', row => { expect(draw(bank, row.picked, row.topic, () => row.random)).toEqual(row.drawn); });
