import type { EventKind } from './types';

// Fixed time zone so server render and browser hydration produce identical text.
const TZ = 'Asia/Singapore';

const money = new Intl.NumberFormat('en-SG', { style: 'currency', currency: 'SGD' });

export const formatPrice = (cents: number) => (cents === 0 ? 'Free' : money.format(cents / 100));

export const formatDay = (iso: string) =>
	new Intl.DateTimeFormat('en-SG', { weekday: 'short', day: 'numeric', month: 'short', timeZone: TZ }).format(
		new Date(iso)
	);

export const formatTime = (iso: string) =>
	new Intl.DateTimeFormat('en-SG', { hour: 'numeric', minute: '2-digit', timeZone: TZ }).format(new Date(iso));

export const formatRange = (start: string, end: string) => `${formatTime(start)} – ${formatTime(end)}`;

export const dateParts = (iso: string) => {
	const d = new Date(iso);
	const part = (o: Intl.DateTimeFormatOptions) => new Intl.DateTimeFormat('en-SG', { ...o, timeZone: TZ }).format(d);
	return { day: part({ day: 'numeric' }), month: part({ month: 'short' }), weekday: part({ weekday: 'short' }) };
};

export const durationMins = (start: string, end: string) =>
	Math.round((new Date(end).getTime() - new Date(start).getTime()) / 60000);

export const kindLabel: Record<EventKind, string> = {
	training: 'Training',
	coaching: 'Coaching',
	event: 'Event'
};

const sportGlyphs: Record<string, string> = {
	Running: '🏃',
	Tennis: '🎾',
	Badminton: '🏸',
	Football: '⚽',
	Basketball: '🏀',
	Fitness: '🏋️',
	Swimming: '🏊'
};
export const sportGlyph = (sport: string) => sportGlyphs[sport] ?? '🏅';
