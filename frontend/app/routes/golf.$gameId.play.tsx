import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { ClientLoaderFunctionArgs, useLoaderData } from "@remix-run/react";
import { useHydrateAtoms } from "jotai/utils";
import { apiGetGame, apiGetToken } from "../.server/api/client";
import { ensureUserLoggedIn } from "../.server/auth";
import GolfPlayApp from "../components/GolfPlayApp.client";
import GolfPlayAppConnecting from "../components/GolfPlayApps/GolfPlayAppConnecting";
import {
	scoreAtom,
	setCurrentTimestampAtom,
	setDurationSecondsAtom,
	submitResultAtom,
} from "../states/play";
import { PlayerState } from "../types/PlayerState";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
	{
		title: data
			? `Golf Playing ${data.game.display_name} | iOSDC Japan 2024 Albatross.swift`
			: "Golf Playing | iOSDC Japan 2024 Albatross.swift",
	},
];

export async function loader({ params, request }: LoaderFunctionArgs) {
	const { token, user } = await ensureUserLoggedIn(request);

	const fetchGame = async () => {
		return (await apiGetGame(token, Number(params.gameId))).game;
	};
	const fetchSockToken = async () => {
		return (await apiGetToken(token)).token;
	};

	const [game, sockToken] = await Promise.all([fetchGame(), fetchSockToken()]);

	const playerState: PlayerState = {
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

	return {
		game,
		player: user,
		sockToken,
		playerState,
	};
}

export async function clientLoader({ serverLoader }: ClientLoaderFunctionArgs) {
	const data = await serverLoader<typeof loader>();
	const baseKey = `playerState:${data.game.game_id}:${data.player.user_id}`;

	const localCode = (() => {
		const rawValue = window.localStorage.getItem(`${baseKey}:code`);
		if (rawValue === null) {
			return null;
		}
		return rawValue;
	})();

	const localScore = (() => {
		const rawValue = window.localStorage.getItem(`${baseKey}:score`);
		if (rawValue === null || rawValue === "") {
			return null;
		}
		return Number(rawValue);
	})();

	const localSubmissionResult = (() => {
		const rawValue = window.localStorage.getItem(`${baseKey}:submissionResult`);
		if (rawValue === null) {
			return null;
		}
		const parsed = JSON.parse(rawValue);
		if (typeof parsed !== "object") {
			return null;
		}
		return parsed;
	})();

	if (localCode !== null) {
		data.playerState.code = localCode;
	}
	if (localScore !== null) {
		data.playerState.score = localScore;
	}
	if (localSubmissionResult !== null) {
		data.playerState.submitResult = localSubmissionResult;
	}

	return data;
}
clientLoader.hydrate = true;

export function HydrateFallback() {
	return <GolfPlayAppConnecting />;
}

export default function GolfPlay() {
	const { game, player, sockToken, playerState } =
		useLoaderData<typeof loader>();

	useHydrateAtoms([
		[setCurrentTimestampAtom, undefined],
		[setDurationSecondsAtom, game.duration_seconds],
		[scoreAtom, playerState.score],
		[submitResultAtom, playerState.submitResult],
	]);

	return (
		<GolfPlayApp
			game={game}
			player={player}
			initialCode={playerState.code}
			sockToken={sockToken}
		/>
	);
}
