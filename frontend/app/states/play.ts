import { atom } from "jotai";
import type { components } from "../.server/api/schema";
import type { SubmitResult } from "../types/SubmitResult";

type RawGameState =
	| {
			kind: "connecting";
			startedAtTimestamp: null;
	  }
	| {
			kind: "waiting";
			startedAtTimestamp: null;
	  }
	| {
			kind: "starting";
			startedAtTimestamp: number;
	  };

const rawGameStateAtom = atom<RawGameState>({
	kind: "connecting",
	startedAtTimestamp: null,
});

export type GameStateKind =
	| "connecting"
	| "waiting"
	| "starting"
	| "gaming"
	| "finished";

export const gameStateKindAtom = atom<GameStateKind>((get) => {
	const { kind: rawKind, startedAtTimestamp } = get(rawGameStateAtom);
	if (rawKind === "connecting" || rawKind === "waiting") {
		return rawKind;
	} else {
		const durationSeconds = get(rawDurationSecondsAtom);
		const finishedAtTimestamp = startedAtTimestamp + durationSeconds;
		const currentTimestamp = get(rawCurrentTimestampAtom);
		if (currentTimestamp < startedAtTimestamp) {
			return "starting";
		} else if (currentTimestamp < finishedAtTimestamp) {
			return "gaming";
		} else {
			return "finished";
		}
	}
});

export const gameStartAtom = atom(null, (get, set, value: number) => {
	const { kind } = get(rawGameStateAtom);
	if (kind === "starting") {
		return;
	}
	set(rawGameStateAtom, {
		kind: "starting",
		startedAtTimestamp: value,
	});
});
export const setGameStateConnectingAtom = atom(null, (_, set) =>
	set(rawGameStateAtom, { kind: "connecting", startedAtTimestamp: null }),
);
export const setGameStateWaitingAtom = atom(null, (_, set) =>
	set(rawGameStateAtom, { kind: "waiting", startedAtTimestamp: null }),
);

const rawCurrentTimestampAtom = atom(0);
export const setCurrentTimestampAtom = atom(null, (_, set) =>
	set(rawCurrentTimestampAtom, Math.floor(Date.now() / 1000)),
);

const rawDurationSecondsAtom = atom<number>(0);
export const setDurationSecondsAtom = atom(null, (_, set, value: number) =>
	set(rawDurationSecondsAtom, value),
);

export const startingLeftTimeSecondsAtom = atom<number | null>((get) => {
	const { startedAtTimestamp } = get(rawGameStateAtom);
	if (startedAtTimestamp === null) {
		return null;
	}
	const currentTimestamp = get(rawCurrentTimestampAtom);
	return Math.max(0, startedAtTimestamp - currentTimestamp);
});

export const gamingLeftTimeSecondsAtom = atom<number | null>((get) => {
	const { startedAtTimestamp } = get(rawGameStateAtom);
	if (startedAtTimestamp === null) {
		return null;
	}
	const durationSeconds = get(rawDurationSecondsAtom);
	const finishedAtTimestamp = startedAtTimestamp + durationSeconds;
	const currentTimestamp = get(rawCurrentTimestampAtom);
	return Math.min(
		durationSeconds,
		Math.max(0, finishedAtTimestamp - currentTimestamp),
	);
});

export const handleWsConnectionClosedAtom = atom(null, (get, set) => {
	const kind = get(gameStateKindAtom);
	if (kind !== "finished") {
		set(setGameStateConnectingAtom);
	}
});

export const scoreAtom = atom<number | null>(null);
export const submitResultAtom = atom<SubmitResult>({
	status: "waiting_submission",
	execResults: [],
});

export const handleSubmitCodeAtom = atom(null, (_, set) => {
	set(submitResultAtom, (prev) => ({
		status: "running",
		execResults: prev.execResults.map((r) => ({
			...r,
			status: "running",
			stdout: "",
			stderr: "",
		})),
	}));
});

type GamePlayerMessageS2CExecResultPayload =
	components["schemas"]["GamePlayerMessageS2CExecResultPayload"];
type GamePlayerMessageS2CSubmitResultPayload =
	components["schemas"]["GamePlayerMessageS2CSubmitResultPayload"];

export const handleWsExecResultMessageAtom = atom(
	null,
	(
		get,
		set,
		data: GamePlayerMessageS2CExecResultPayload,
		callback: (submissionResult: SubmitResult) => void,
	) => {
		const { testcase_id, status, stdout, stderr } = data;
		const prev = get(submitResultAtom);
		const newResult = {
			...prev,
			execResults: prev.execResults.map((r) =>
				r.testcase_id === testcase_id && r.status === "running"
					? {
							...r,
							status,
							stdout,
							stderr,
						}
					: r,
			),
		};
		set(submitResultAtom, newResult);
		callback(newResult);
	},
);

export const handleWsSubmitResultMessageAtom = atom(
	null,
	(
		get,
		set,
		data: GamePlayerMessageS2CSubmitResultPayload,
		callback: (submissionResult: SubmitResult, score: number | null) => void,
	) => {
		const { status, score } = data;
		const prev = get(submitResultAtom);
		const newResult = {
			...prev,
			status,
		};
		if (status !== "success") {
			newResult.execResults = prev.execResults.map((r) =>
				r.status === "running" ? { ...r, status: "canceled" } : r,
			);
		} else {
			newResult.execResults = prev.execResults.map((r) => ({
				...r,
				status: "success",
			}));
		}
		set(submitResultAtom, newResult);
		if (status === "success" && score !== null) {
			const currentScore = get(scoreAtom);
			if (currentScore === null || score < currentScore) {
				set(scoreAtom, score);
			}
		}
		callback(newResult, score);
	},
);
