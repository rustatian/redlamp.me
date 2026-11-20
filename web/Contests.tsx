import { useEffect, useRef, useState } from 'react';
import { Skull } from 'lucide-react';
import type { BankEntry, ContestDetail, ContestList, SessionUser } from '../shared/leetcode';
import { api, message } from './api';
import { draw, freeDifficulties, type Picked } from './contestDraw';
import { clock } from './contestTime';
import { go } from './navigate';
import { difficultyLabels } from './Problems';
import { plural } from './text';
import { useResource } from './useResource';

export const points = (count: number) => plural(count, ['очко', 'очка', 'очков']);

export function Contests({ user }: { user: SessionUser }) {
	const { data, error } = useResource<ContestList>('/contests');
	const [starting, setStarting] = useState(false);
	if (!data) return <div className="stack"><h1>Мясорубка</h1><p role="status">{error || 'Открываем контесты...'}</p></div>;
	const { best, active } = data, finished = data.items.filter(c => !c.active);
	const stats = [{ label: 'Рейтинг', value: data.rating }, { label: 'Лучший рейтинг', value: data.bestRating }, { label: best ? `Очки лучшего контеста, финиш ${clock(best.finishSeconds)}` : 'Очки лучшего контеста', value: best?.score ?? '-' }, { label: 'Контестов', value: finished.length }];
	return <div className="stack"><div className="heading"><div><h1>Мясорубка</h1><p>90 минут на задачи Codeforces. Время не остановить.</p></div>{active ? <a className="button primary" href={`/leetcode/contests/${active.id}`} data-tip="Открыть контест, который идёт сейчас.">Идёт контест</a> : user.role === 'student' && <button className="primary" onClick={() => setStarting(true)} data-tip="Выбрать задачи и начать контест на 90 минут.">В мясорубку</button>}</div>
		<div className="stats">{stats.map(item => <div className="stat" key={item.label}><strong>{item.value}</strong><span>{item.label}</span></div>)}</div>
		<section className="card"><header className="card-head"><h2>Законченные контесты</h2></header>{finished.length ? <div className="results"><table><thead><tr><th scope="col">Дата</th><th scope="col">Решено</th><th scope="col">Очки</th><th scope="col">Финиш</th><th scope="col">Изменение</th><th scope="col">Рейтинг</th></tr></thead>
			<tbody>{finished.map(c => <tr key={c.id}><td><a href={`/leetcode/contests/${c.id}`}>{new Date(c.startedAt).toLocaleString('ru-RU', { dateStyle: 'short', timeStyle: 'short' })}</a>{c.best && <span className="tag good">Лучший</span>}</td><td>{c.solved} из {c.total}</td><td>{c.score}</td><td>{clock(c.finishSeconds)}</td><td className={c.ratingChange! > 0 ? 'up' : c.ratingChange! < 0 ? 'down' : undefined}>{c.ratingChange!.toLocaleString('ru-RU', { signDisplay: 'exceptZero' })}</td><td>{c.ratingAfter}</td></tr>)}</tbody></table></div> : <p className="empty">Законченных контестов пока нет.</p>}</section>
		{starting && <Start close={() => setStarting(false)} />}
	</div>;
}

function Start({ close }: { close: () => void }) {
	const dialog = useRef<HTMLDialogElement>(null), [ready, setReady] = useState(false);
	useEffect(() => {
		const previous = document.activeElement, element = dialog.current;
		element?.showModal();
		return () => { element?.close(); if (previous instanceof HTMLElement) previous.focus(); };
	}, []);
	// The setup removes the focused button, so the focus goes to the dialog and the screen reader reads its new title.
	function agree() { setReady(true); dialog.current?.focus(); }
	return <dialog ref={dialog} className={ready ? undefined : 'contest-skull'} tabIndex={-1} onCancel={close} aria-labelledby="contest-start-title">{ready ? <Setup close={close} /> : <div className="stack"><Skull size={96} aria-hidden="true" /><h2 id="contest-start-title">Готов умереть?</h2><div className="row"><button className="primary" onClick={agree} data-tip="Перейти к выбору задач.">Да</button><button onClick={close} data-tip="Закрыть окно.">Нет</button></div></div>}</dialog>;
}

function Setup({ close }: { close: () => void }) {
	const index = useResource<{ topics: string[]; problems: BankEntry[] }>('/contests/bank');
	const [picked, setPicked] = useState<Picked[]>([]), [error, setError] = useState(''), [busy, setBusy] = useState(false);
	const bank = index.data?.problems ?? [];
	function add(topic: string) {
		const pick = draw(bank, picked, topic, Math.random);
		if (pick) setPicked([...picked, pick]);
	}
	async function start() {
		setBusy(true); setError('');
		try {
			const contest = await api<ContestDetail>('/contests', { method: 'POST', body: JSON.stringify({ problems: picked.map(({ slug, topic }) => ({ slug, topic })) }) });
			go(`/leetcode/contests/${contest.id}`);
		} catch (error) { setError(message(error)); setBusy(false); }
	}
	return <div className="stack"><h2 id="contest-start-title">Выбор задач</h2><p>Каждое нажатие на тему добавляет случайную задачу. В контесте до 4 задач: 1 лёгкая, 2 средние и 1 сложная.</p>
		{(error || index.error) && <p className="error" role="alert">{error || index.error}</p>}
		{index.data ? <div className="topics" role="group" aria-label="Темы">{index.data.topics.map(topic => { const count = picked.filter(p => p.topic === topic).length; return <button key={topic} disabled={!freeDifficulties(bank, picked, topic).length} onClick={() => add(topic)} data-tip="Добавить случайную задачу по этой теме.">{topic} <span className={count ? 'count picked' : 'count'}>x{count}</span></button>; })}</div> : !index.error && <p role="status">Открываем банк задач...</p>}
		<h3>Задачи: {picked.length} из 4</h3>
		{picked.length ? <ol className="card picks">{picked.map((p, i) => <li className="list-row" key={p.slug}><span className="grow">{p.topic}</span><span className={`tag ${p.difficulty}`}>{difficultyLabels[p.difficulty]}</span><span className="small">{points(p.points)}</span><button onClick={() => setPicked(picked.filter(item => item !== p))} aria-label={`Убрать задачу ${i + 1}: ${p.topic}`} data-tip="Убрать задачу и освободить её место.">Убрать</button></li>)}</ol> : <p className="empty">Нажмите на тему, чтобы добавить задачу.</p>}
		<div className="row"><button className="primary" disabled={!picked.length || busy} onClick={() => void start()} data-tip="Начать контест на 90 минут.">Начать</button><button onClick={close} data-tip="Закрыть окно. Выбор задач пропадёт.">Закрыть</button></div>
	</div>;
}
