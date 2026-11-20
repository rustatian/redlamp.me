// The time left in whole seconds, from the server time of the response and the client time since the fetch, so a wrong client clock has no effect.
// It rounds up: 00:00 shows only at the end.
export function remaining(endsAt: string, serverNow: string, fetchedAt: number, clientNow: number): number {
	return Math.max(0, Math.ceil((Date.parse(endsAt) - Date.parse(serverNow) - (clientNow - fetchedAt)) / 1000));
}

export function clock(seconds: number): string {
	return `${String(Math.floor(seconds / 60)).padStart(2, '0')}:${String(seconds % 60).padStart(2, '0')}`;
}
