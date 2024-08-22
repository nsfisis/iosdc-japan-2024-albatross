import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { ClientLoaderFunctionArgs, useLoaderData } from "@remix-run/react";
import { useHydrateAtoms } from "jotai/utils";
import { apiGetGame, apiGetToken } from "../.server/api/client";
import { ensureUserLoggedIn } from "../.server/auth";
import GolfWatchAppWithAudioPlayRequest from "../components/GolfWatchAppWithAudioPlayRequest.client";
import GolfWatchAppConnecting from "../components/GolfWatchApps/GolfWatchAppConnecting";
import {
	codeAAtom,
	codeBAtom,
	scoreAAtom,
	scoreBAtom,
	setCurrentTimestampAtom,
	setDurationSecondsAtom,
	submitResultAAtom,
	submitResultBAtom,
} from "../states/watch";
import { PlayerState } from "../types/PlayerState";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Watching ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
			: "Golf Watching | iOSDC Japan 2024 Albatross.swift",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token } = await ensureUserLoggedIn(request);

	const fetchGame = async () => {
		return (await apiGetGame(token, Number(params.gameId))).game;
	};
	const fetchSockToken = async () => {
		return (await apiGetToken(token)).token;
	};

	const [game, sockToken] = await Promise.all([fetchGame(), fetchSockToken()]);

	if (game.game_type !== "1v1") {
		throw new Response("Not Found", { status: 404 });
	}

	const playerStateA: PlayerState = {
		code: "",
		score: null,
		submitResult: {
			status: "waiting_submission",
			execResults: game.exec_steps.map((r) => ({
				testcase_id: r.testcase_id,
				status: "waiting_submission",
				label: r.label,
				stdout: "",
				stderr: "",
			})),
		},
	};
	const playerStateB = structuredClone(playerStateA);

	return {
		game,
		sockToken,
		playerStateA,
		playerStateB,
	};
}

export async function clientLoader({ serverLoader }: ClientLoaderFunctionArgs) {
	const data = await serverLoader<typeof loader>();

	const playerIdA = data.game.players[0]?.user_id;
	const playerIdB = data.game.players[1]?.user_id;

	if (playerIdA !== null) {
		const baseKeyA = `watcherState:${data.game.game_id}:${playerIdA}`;

		const localCodeA = (() => {
			const rawValue = window.localStorage.getItem(`${baseKeyA}:code`);

			if (rawValue === null) {
				return null;
			}
			return rawValue;
		})();

		const localScoreA = (() => {
			const rawValue = window.localStorage.getItem(`${baseKeyA}:score`);
			if (rawValue === null || rawValue === "") {
				return null;
			}
			return Number(rawValue);
		})();

		const localSubmissionResultA = (() => {
			const rawValue = window.localStorage.getItem(
				`${baseKeyA}:submissionResult`,
			);
			if (rawValue === null) {
				return null;
			}
			const parsed = JSON.parse(rawValue);
			if (typeof parsed !== "object") {
				return null;
			}
			return parsed;
		})();

		if (localCodeA !== null) {
			data.playerStateA.code = localCodeA;
		}
		if (localScoreA !== null) {
			data.playerStateA.score = localScoreA;
		}
		if (localSubmissionResultA !== null) {
			data.playerStateA.submitResult = localSubmissionResultA;
		}
	}

	if (playerIdB !== null) {
		const baseKeyB = `watcherState:${data.game.game_id}:${playerIdB}`;

		const localCodeB = (() => {
			const rawValue = window.localStorage.getItem(`${baseKeyB}:code`);
			if (rawValue === null) {
				return null;
			}
			return rawValue;
		})();

		const localScoreB = (() => {
			const rawValue = window.localStorage.getItem(`${baseKeyB}:score`);
			if (rawValue === null || rawValue === "") {
				return null;
			}
			return Number(rawValue);
		})();

		const localSubmissionResultB = (() => {
			const rawValue = window.localStorage.getItem(
				`${baseKeyB}:submissionResult`,
			);
			if (rawValue === null) {
				return null;
			}
			const parsed = JSON.parse(rawValue);
			if (typeof parsed !== "object") {
				return null;
			}
			return parsed;
		})();

		if (localCodeB !== null) {
			data.playerStateB.code = localCodeB;
		}
		if (localScoreB !== null) {
			data.playerStateB.score = localScoreB;
		}
		if (localSubmissionResultB !== null) {
			data.playerStateB.submitResult = localSubmissionResultB;
		}
	}

	return data;
}
clientLoader.hydrate = true;

export function HydrateFallback() {
	return <GolfWatchAppConnecting />;
}

export default function GolfWatch() {
	const { game, sockToken, playerStateA, playerStateB } =
		useLoaderData<typeof loader>();

	useHydrateAtoms([
		[setCurrentTimestampAtom, undefined],
		[setDurationSecondsAtom, game.duration_seconds],
		[codeAAtom, playerStateA.code],
		[codeBAtom, playerStateB.code],
		[scoreAAtom, playerStateA.score],
		[scoreBAtom, playerStateB.score],
		[submitResultAAtom, playerStateA.submitResult],
		[submitResultBAtom, playerStateB.submitResult],
	]);

	return <GolfWatchAppWithAudioPlayRequest game={game} sockToken={sockToken} />;
}
