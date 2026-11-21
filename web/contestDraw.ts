import type { BankEntry, Difficulty } from '../shared/leetcode';

export type Picked = { slug: string; topic: string; difficulty: Difficulty; rating: number; points: number };

// The places of a contest for each difficulty, as the server checks them.
const places: Record<Difficulty, number> = { easy: 1, medium: 2, hard: 1 };

const unpicked = (bank: BankEntry[], picked: Picked[], topic: string, difficulty: Difficulty) => bank.filter(entry => entry.difficulty === difficulty && entry.topics.includes(topic) && !picked.some(p => p.slug === entry.slug));

// The difficulties with a free place and an unpicked bank problem of the topic. 4 picks fill all places.
export function freeDifficulties(bank: BankEntry[], picked: Picked[], topic: string): Difficulty[] {
	return (Object.keys(places) as Difficulty[]).filter(difficulty => picked.filter(p => p.difficulty === difficulty).length < places[difficulty] && unpicked(bank, picked, topic, difficulty).length > 0);
}

// A random free difficulty with equal chance, then a random unpicked problem of the topic with that difficulty. Null when no difficulty is free.
export function draw(bank: BankEntry[], picked: Picked[], topic: string, random: () => number): Picked | null {
	const free = freeDifficulties(bank, picked, topic);
	if (!free.length) return null;
	const difficulty = free[Math.floor(random() * free.length)], candidates = unpicked(bank, picked, topic, difficulty);
	const { slug, rating, points } = candidates[Math.floor(random() * candidates.length)];
	return { slug, topic, difficulty, rating, points };
}
