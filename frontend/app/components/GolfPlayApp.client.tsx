import { useAtomValue, useSetAtom } from "jotai";
import { useEffect } from "react";
import { useTimer } from "react-use-precision-timer";
import { useDebouncedCallback } from "use-debounce";
import type { components } from "../.server/api/schema";
import useWebSocket, { ReadyState } from "../hooks/useWebSocket";
import {
	gameStartAtom,
	gameStateKindAtom,
	handleSubmitCodeAtom,
	handleWsConnectionClosedAtom,
	handleWsExecResultMessageAtom,
	handleWsSubmitResultMessageAtom,
	setCurrentTimestampAtom,
	setGameStateConnectingAtom,
	setGameStateWaitingAtom,
} from "../states/play";
import GolfPlayAppConnecting from "./GolfPlayApps/GolfPlayAppConnecting";
import GolfPlayAppFinished from "./GolfPlayApps/GolfPlayAppFinished";
import GolfPlayAppGaming from "./GolfPlayApps/GolfPlayAppGaming";
import GolfPlayAppStarting from "./GolfPlayApps/GolfPlayAppStarting";
import GolfPlayAppWaiting from "./GolfPlayApps/GolfPlayAppWaiting";

type GamePlayerMessageS2C = components["schemas"]["GamePlayerMessageS2C"];
type GamePlayerMessageC2S = components["schemas"]["GamePlayerMessageC2S"];

type Game = components["schemas"]["Game"];
type User = components["schemas"]["User"];

type Props = {
	game: Game;
	player: User;
	initialCode: string;
	sockToken: string;
};

export default function GolfPlayApp({
	game,
	player,
	initialCode,
	sockToken,
}: Props) {
	const socketUrl =
		process.env.NODE_ENV === "development"
			? `ws://localhost:8002/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/play?token=${sockToken}`
			: `wss://t.nil.ninja/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/play?token=${sockToken}`;

	const gameStateKind = useAtomValue(gameStateKindAtom);
	const setCurrentTimestamp = useSetAtom(setCurrentTimestampAtom);
	const gameStart = useSetAtom(gameStartAtom);
	const setGameStateConnecting = useSetAtom(setGameStateConnectingAtom);
	const setGameStateWaiting = useSetAtom(setGameStateWaitingAtom);
	const handleWsConnectionClosed = useSetAtom(handleWsConnectionClosedAtom);
	const handleWsExecResultMessage = useSetAtom(handleWsExecResultMessageAtom);
	const handleWsSubmitResultMessage = useSetAtom(
		handleWsSubmitResultMessageAtom,
	);
	const handleSubmitCode = useSetAtom(handleSubmitCodeAtom);

	useTimer({ delay: 1000, startImmediately: true }, setCurrentTimestamp);

	const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket<
		GamePlayerMessageS2C,
		GamePlayerMessageC2S
	>(socketUrl);

	const playerProfile = {
		displayName: player.display_name,
		iconPath: player.icon_path ?? null,
	};

	const onCodeChange = useDebouncedCallback((code: string) => {
		console.log("player:c2s:code");
		sendJsonMessage({
			type: "player:c2s:code",
			data: { code },
		});
		const baseKey = `playerState:${game.game_id}:${player.user_id}`;
		window.localStorage.setItem(`${baseKey}:code`, code);
	}, 1000);

	const onCodeSubmit = useDebouncedCallback((code: string) => {
		if (code === "") {
			return;
		}
		console.log("player:c2s:submit");
		sendJsonMessage({
			type: "player:c2s:submit",
			data: { code },
		});
		handleSubmitCode();
	}, 1000);

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
				console.log(lastJsonMessage.data);
				if (lastJsonMessage.type === "player:s2c:start") {
					const { start_at } = lastJsonMessage.data;
					gameStart(start_at);
				} else if (lastJsonMessage.type === "player:s2c:execresult") {
					handleWsExecResultMessage(
						lastJsonMessage.data,
						(submissionResult) => {
							const baseKey = `playerState:${game.game_id}:${player.user_id}`;
							window.localStorage.setItem(
								`${baseKey}:submissionResult`,
								JSON.stringify(submissionResult),
							);
						},
					);
				} else if (lastJsonMessage.type === "player:s2c:submitresult") {
					handleWsSubmitResultMessage(
						lastJsonMessage.data,
						(submissionResult, score) => {
							const baseKey = `playerState:${game.game_id}:${player.user_id}`;
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
		game.game_id,
		game.started_at,
		player.user_id,
		sendJsonMessage,
		lastJsonMessage,
		readyState,
		gameStart,
		handleWsConnectionClosed,
		handleWsExecResultMessage,
		handleWsSubmitResultMessage,
		setGameStateConnecting,
		setGameStateWaiting,
	]);

	if (gameStateKind === "connecting") {
		return <GolfPlayAppConnecting />;
	} else if (gameStateKind === "waiting") {
		return (
			<GolfPlayAppWaiting
				gameDisplayName={game.display_name}
				playerProfile={playerProfile}
			/>
		);
	} else if (gameStateKind === "starting") {
		return <GolfPlayAppStarting gameDisplayName={game.display_name} />;
	} else if (gameStateKind === "gaming") {
		return (
			<GolfPlayAppGaming
				gameDisplayName={game.display_name}
				playerProfile={playerProfile}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				initialCode={initialCode}
				onCodeChange={onCodeChange}
				onCodeSubmit={onCodeSubmit}
			/>
		);
	} else if (gameStateKind === "finished") {
		return <GolfPlayAppFinished />;
	} else {
		return null;
	}
}
