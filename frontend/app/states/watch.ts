import { atom } from "jotai";
import { AudioController } from "../.client/audio/AudioController";
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

export const codeAAtom = atom("");
export const codeBAtom = atom("");
export const scoreAAtom = atom<number | null>(null);
export const scoreBAtom = atom<number | null>(null);
export const submitResultAAtom = atom<SubmitResult>({
	status: "waiting_submission",
	execResults: [],
});
export const submitResultBAtom = atom<SubmitResult>({
	status: "waiting_submission",
	execResults: [],
});

type GameWatcherMessageS2CSubmitPayload =
	components["schemas"]["GameWatcherMessageS2CSubmitPayload"];
type GameWatcherMessageS2CCodePayload =
	components["schemas"]["GameWatcherMessageS2CCodePayload"];
type GameWatcherMessageS2CExecResultPayload =
	components["schemas"]["GameWatcherMessageS2CExecResultPayload"];
type GameWatcherMessageS2CSubmitResultPayload =
	components["schemas"]["GameWatcherMessageS2CSubmitResultPayload"];

export const handleWsCodeMessageAtom = atom(
	null,
	(
		_,
		set,
		data: GameWatcherMessageS2CCodePayload,
		getTarget: <T>(player_id: number, atomA: T, atomB: T) => T,
		callback: (player_id: number, code: string) => void,
	) => {
		const { player_id, code } = data;
		const codeAtom = getTarget(player_id, codeAAtom, codeBAtom);
		set(codeAtom, code);
		callback(player_id, code);
	},
);

export const handleWsSubmitMessageAtom = atom(
	null,
	(
		get,
		set,
		data: GameWatcherMessageS2CSubmitPayload,
		getTarget: <T>(player_id: number, atomA: T, atomB: T) => T,
		callback: (player_id: number, submissionResult: SubmitResult) => void,
	) => {
		const { player_id } = data;
		const submitResultAtom = getTarget(
			player_id,
			submitResultAAtom,
			submitResultBAtom,
		);
		const prev = get(submitResultAtom);
		const newResult = {
			status: "running" as const,
			execResults: prev.execResults.map((r) => ({
				...r,
				status: "running" as const,
				stdout: "",
				stderr: "",
			})),
		};
		set(submitResultAtom, newResult);
		callback(player_id, newResult);
	},
);

export const handleWsExecResultMessageAtom = atom(
	null,
	(
		get,
		set,
		data: GameWatcherMessageS2CExecResultPayload,
		getTarget: <T>(player_id: number, atomA: T, atomB: T) => T,
		callback: (player_id: number, submissionResult: SubmitResult) => void,
	) => {
		const { player_id, testcase_id, status, stdout, stderr } = data;
		const submitResultAtom = getTarget(
			player_id,
			submitResultAAtom,
			submitResultBAtom,
		);
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
		callback(player_id, newResult);
	},
);

export const handleWsSubmitResultMessageAtom = atom(
	null,
	(
		get,
		set,
		data: GameWatcherMessageS2CSubmitResultPayload,
		getTarget: <T>(player_id: number, atomA: T, atomB: T) => T,
		callback: (
			player_id: number,
			submissionResult: SubmitResult,
			score: number | null,
		) => void,
	) => {
		const { player_id, status, score } = data;
		const submitResultAtom = getTarget(
			player_id,
			submitResultAAtom,
			submitResultBAtom,
		);
		const scoreAtom = getTarget(player_id, scoreAAtom, scoreBAtom);
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
		callback(player_id, newResult, score);
	},
);

export const audioControllerAtom = atom<AudioController | null>(null);
