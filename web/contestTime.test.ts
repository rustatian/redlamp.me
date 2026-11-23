import { expect, it } from 'vitest';
import { clock, remaining } from './contestTime';

const endsAt = '2026-11-20T10:30:00.000Z', serverNow = '2026-11-20T10:00:00.000Z';
// The client clock is 10 minutes fast: it reads 10:10:00 when the server reads 10:00:00.
const fetchedAt = Date.parse('2026-11-20T10:10:00.000Z');

const remainders: { name: string; clientNow: number; seconds: number }[] = [
	{ name: 'client clock 10 minutes fast', clientNow: fetchedAt, seconds: 1800 },
	{ name: 'one minute after the fetch', clientNow: fetchedAt + 60_000, seconds: 1740 },
	{ name: 'half a second before the end rounds up', clientNow: fetchedAt + 1_799_500, seconds: 1 },
	{ name: 'after the end', clientNow: fetchedAt + 1_860_000, seconds: 0 },
];
it.each(remainders)('$name', row => { expect(remaining(endsAt, serverNow, fetchedAt, row.clientNow)).toBe(row.seconds); });

const clocks: { name: string; seconds: number; text: string }[] = [
	{ name: 'whole contest', seconds: 5400, text: '90:00' },
	{ name: 'one minute and one second', seconds: 61, text: '01:01' },
	{ name: 'end of the contest', seconds: 0, text: '00:00' },
];
it.each(clocks)('$name', row => { expect(clock(row.seconds)).toBe(row.text); });
