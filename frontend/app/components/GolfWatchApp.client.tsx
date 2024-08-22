import { useAtomValue, useSetAtom } from "jotai";
import { useCallback, useEffect } from "react";
import { useTimer } from "react-use-precision-timer";
import { AudioController } from "../.client/audio/AudioController";
import type { components } from "../.server/api/schema";
import useWebSocket, { ReadyState } from "../hooks/useWebSocket";
import {
	gameStartAtom,
	gameStateKindAtom,
	handleWsCodeMessageAtom,
	handleWsConnectionClosedAtom,
	handleWsExecResultMessageAtom,
	handleWsSubmitMessageAtom,
	handleWsSubmitResultMessageAtom,
	setCurrentTimestampAtom,
	setGameStateConnectingAtom,
	setGameStateWaitingAtom,
} from "../states/watch";
import GolfWatchAppConnecting from "./GolfWatchApps/GolfWatchAppConnecting";
import GolfWatchAppGaming from "./GolfWatchApps/GolfWatchAppGaming";
import GolfWatchAppStarting from "./GolfWatchApps/GolfWatchAppStarting";
import GolfWatchAppWaiting from "./GolfWatchApps/GolfWatchAppWaiting";

type GameWatcherMessageS2C = components["schemas"]["GameWatcherMessageS2C"];
type GameWatcherMessageC2S = never;

type Game = components["schemas"]["Game"];

export type Props = {
	game: Game;
	sockToken: string;
	audioController: AudioController;
};

export default function GolfWatchApp({
	game,
	sockToken,
	audioController,
}: Props) {
	const socketUrl =
		process.env.NODE_ENV === "development"
			? `ws://localhost:8002/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/watch?token=${sockToken}`
			: `wss://t.nil.ninja/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/watch?token=${sockToken}`;

	const gameStateKind = useAtomValue(gameStateKindAtom);
	const setCurrentTimestamp = useSetAtom(setCurrentTimestampAtom);
	const gameStart = useSetAtom(gameStartAtom);
	const setGameStateConnecting = useSetAtom(setGameStateConnectingAtom);
	const setGameStateWaiting = useSetAtom(setGameStateWaitingAtom);
	const handleWsConnectionClosed = useSetAtom(handleWsConnectionClosedAtom);
	const handleWsCodeMessage = useSetAtom(handleWsCodeMessageAtom);
	const handleWsSubmitMessage = useSetAtom(handleWsSubmitMessageAtom);
	const handleWsExecResultMessage = useSetAtom(handleWsExecResultMessageAtom);
	const handleWsSubmitResultMessage = useSetAtom(
		handleWsSubmitResultMessageAtom,
	);

	useTimer({ delay: 1000, startImmediately: true }, setCurrentTimestamp);

	const { lastJsonMessage, readyState } = useWebSocket<
		GameWatcherMessageS2C,
		GameWatcherMessageC2S
	>(socketUrl);

	const playerA = game.players[0]!;
	const playerB = game.players[1]!;

	const getTargetAtomByPlayerId: <T>(
		player_id: number,
		atomA: T,
		atomB: T,
	) => T = useCallback(
		(player_id, atomA, atomB) =>
			player_id === playerA.user_id ? atomA : atomB,
		[playerA.user_id],
	);

	const playerProfileA = {
		displayName: playerA.display_name,
		iconPath: playerA.icon_path ?? null,
	};
	const playerProfileB = {
		displayName: playerB.display_name,
		iconPath: playerB.icon_path ?? null,
	};

	if (readyState === ReadyState.UNINSTANTIATED) {
		throw new Error("WebSocket is not connected");
	}

	useEffect(() => {
		if (readyState === ReadyState.CLOSING || readyState === ReadyState.CLOSED) {
			handleWsConnectionClosed();
		} else if (readyState === ReadyState.CONNECTING) {
			setGameStateConnecting();
		} else if (readyState === ReadyState.OPEN) {
			if (lastJsonMessage !== null) {
				console.log(lastJsonMessage.type);
				if (lastJsonMessage.type === "watcher:s2c:start") {
					const { start_at } = lastJsonMessage.data;
					gameStart(start_at);
				} else if (lastJsonMessage.type === "watcher:s2c:code") {
					handleWsCodeMessage(
						lastJsonMessage.data,
						getTargetAtomByPlayerId,
						(player_id, code) => {
							const baseKey = `watcherState:${game.game_id}:${player_id}`;
							window.localStorage.setItem(`${baseKey}:code`, code);
						},
					);
				} else if (lastJsonMessage.type === "watcher:s2c:submit") {
					handleWsSubmitMessage(
						lastJsonMessage.data,
						getTargetAtomByPlayerId,
						(player_id, submissionResult) => {
							const baseKey = `watcherState:${game.game_id}:${player_id}`;
							window.localStorage.setItem(
								`${baseKey}:submissionResult`,
								JSON.stringify(submissionResult),
							);
						},
					);
				} else if (lastJsonMessage.type === "watcher:s2c:execresult") {
					handleWsExecResultMessage(
						lastJsonMessage.data,
						getTargetAtomByPlayerId,
						(player_id, submissionResult) => {
							const baseKey = `watcherState:${game.game_id}:${player_id}`;
							window.localStorage.setItem(
								`${baseKey}:submissionResult`,
								JSON.stringify(submissionResult),
							);
						},
					);
				} else if (lastJsonMessage.type === "watcher:s2c:submitresult") {
					handleWsSubmitResultMessage(
						lastJsonMessage.data,
						getTargetAtomByPlayerId,
						(player_id, submissionResult, score) => {
							const baseKey = `watcherState:${game.game_id}:${player_id}`;
							window.localStorage.setItem(
								`${baseKey}:submissionResult`,
								JSON.stringify(submissionResult),
							);
							window.localStorage.setItem(
								`${baseKey}:score`,
								score === null ? "" : score.toString(),
							);
						},
					);
				}
			} else {
				if (game.started_at) {
					gameStart(game.started_at);
				} else {
					setGameStateWaiting();
				}
			}
		}
	}, [
		game.started_at,
		game.game_id,
		lastJsonMessage,
		readyState,
		gameStart,
		getTargetAtomByPlayerId,
		handleWsCodeMessage,
		handleWsConnectionClosed,
		handleWsExecResultMessage,
		handleWsSubmitMessage,
		handleWsSubmitResultMessage,
		setGameStateConnecting,
		setGameStateWaiting,
	]);

	if (gameStateKind === "connecting") {
		return <GolfWatchAppConnecting />;
	} else if (gameStateKind === "waiting") {
		return (
			<GolfWatchAppWaiting
				gameDisplayName={game.display_name}
				playerProfileA={playerProfileA}
				playerProfileB={playerProfileB}
			/>
		);
	} else if (gameStateKind === "starting") {
		return <GolfWatchAppStarting gameDisplayName={game.display_name} />;
	} else if (gameStateKind === "gaming" || gameStateKind === "finished") {
		return (
			<GolfWatchAppGaming
				gameDisplayName={game.display_name}
				playerProfileA={playerProfileA}
				playerProfileB={playerProfileB}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				gameResult={null /* TODO */}
			/>
		);
	} else {
		return null;
	}
}
