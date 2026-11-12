export class ApiError extends Error {
	constructor(public status: number, public code: string, message: string, public fields?: Record<string, string>) { super(message); }
}

type Init = Omit<RequestInit, 'body'> & { body?: string | Uint8Array<ArrayBuffer> };

// Sends a request to the notebook API. A text body is JSON.
// CloudFront signs each request to the function URL, and Lambda takes no unsigned body, so a request with a body method has the SHA-256 of its body:
// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/private-content-restricting-access-to-lambda.html
export async function send(path: string, init: Init = {}): Promise<Response> {
	const headers = new Headers(init.headers), method = init.method ?? 'GET';
	if (typeof init.body === 'string') headers.set('Content-Type', 'application/json');
	if (method !== 'GET' && method !== 'HEAD') {
		const body = typeof init.body === 'string' ? new TextEncoder().encode(init.body) : init.body ?? new Uint8Array();
		headers.set('x-amz-content-sha256', Array.from(new Uint8Array(await crypto.subtle.digest('SHA-256', body)), byte => byte.toString(16).padStart(2, '0')).join(''));
	}
	return fetch(`/leetcode/api${path}`, { ...init, headers, cache: 'no-store' });
}

export async function api<T>(path: string, init: Init = {}): Promise<T> {
	const response = await send(path, init);
	if (response.status === 204) return undefined as T;
	if (!response.ok) throw await apiError(response);
	return response.json().catch(() => null);
}

export async function apiError(response: Response): Promise<ApiError> {
	const value = await response.json().catch(() => null);
	return new ApiError(response.status, value?.error?.code ?? 'request_failed', value?.error?.message ?? 'Запрос не выполнен. Попробуйте снова.', value?.error?.fields);
}

export function message(error: unknown): string { return error instanceof Error ? error.message : 'Запрос не выполнен. Попробуйте снова.'; }
