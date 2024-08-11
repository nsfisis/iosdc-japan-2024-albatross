import { useEffect, useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import type { components } from "../.server/api/schema";
import useWebSocket, { ReadyState } from "../hooks/useWebSocket";
import GolfPlayAppConnecting from "./GolfPlayApps/GolfPlayAppConnecting";
import GolfPlayAppFinished from "./GolfPlayApps/GolfPlayAppFinished";
import GolfPlayAppGaming from "./GolfPlayApps/GolfPlayAppGaming";
import GolfPlayAppStarting from "./GolfPlayApps/GolfPlayAppStarting";
import GolfPlayAppWaiting from "./GolfPlayApps/GolfPlayAppWaiting";

type GamePlayerMessageS2C = components["schemas"]["GamePlayerMessageS2C"];
type GamePlayerMessageC2S = components["schemas"]["GamePlayerMessageC2S"];

type Game = components["schemas"]["Game"];

type GameState = "connecting" | "waiting" | "starting" | "gaming" | "finished";

export default function GolfPlayApp({
	game,
	sockToken,
}: {
	game: Game;
	sockToken: string;
}) {
	const socketUrl =
		process.env.NODE_ENV === "development"
			? `ws://localhost:8002/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/play?token=${sockToken}`
			: `wss://t.nil.ninja/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/play?token=${sockToken}`;

	const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket<
		GamePlayerMessageS2C,
		GamePlayerMessageC2S
	>(socketUrl);

	const [gameState, setGameState] = useState<GameState>("connecting");

	const [startedAt, setStartedAt] = useState<number | null>(null);

	const [leftTimeSeconds, setLeftTimeSeconds] = useState<number | null>(null);

	useEffect(() => {
		if (gameState === "starting" && startedAt !== null) {
			const timer1 = setInterval(() => {
				setLeftTimeSeconds((prev) => {
					if (prev === null) {
						return null;
					}
					if (prev <= 1) {
						clearInterval(timer1);
						setGameState("gaming");
						return 0;
					}
					return prev - 1;
				});
			}, 1000);

			const timer2 = setInterval(() => {
				const nowSec = Math.floor(Date.now() / 1000);
				const finishedAt = startedAt + game.duration_seconds;
				if (nowSec >= finishedAt) {
					clearInterval(timer2);
					setGameState("finished");
				}
			}, 1000);

			return () => {
				clearInterval(timer1);
				clearInterval(timer2);
			};
		}
	}, [gameState, startedAt, game.duration_seconds]);

	const [currentScore, setCurrentScore] = useState<number | null>(null);

	const [lastExecStatus, setLastExecStatus] = useState<string | null>(null);

	const onCodeChange = useDebouncedCallback((code: string) => {
		console.log("player:c2s:code");
		sendJsonMessage({
			type: "player:c2s:code",
			data: { code },
		});
	}, 1000);

	const onCodeSubmit = useDebouncedCallback((code: string) => {
		console.log("player:c2s:submit");
		sendJsonMessage({
			type: "player:c2s:submit",
			data: { code },
		});
	}, 1000);

	if (readyState === ReadyState.UNINSTANTIATED) {
		throw new Error("WebSocket is not connected");
	}

	useEffect(() => {
		if (readyState === ReadyState.CLOSING || readyState === ReadyState.CLOSED) {
			if (gameState !== "finished") {
				setGameState("connecting");
			}
		} else if (readyState === ReadyState.CONNECTING) {
			setGameState("connecting");
		} else if (readyState === ReadyState.OPEN) {
			if (lastJsonMessage !== null) {
				console.log(lastJsonMessage.type);
				if (lastJsonMessage.type === "player:s2c:start") {
					if (
						gameState !== "starting" &&
						gameState !== "gaming" &&
						gameState !== "finished"
					) {
						const { start_at } = lastJsonMessage.data;
						setStartedAt(start_at);
						const nowSec = Math.floor(Date.now() / 1000);
						setLeftTimeSeconds(start_at - nowSec);
						setGameState("starting");
					}
				} else if (lastJsonMessage.type === "player:s2c:execresult") {
					const { status, score } = lastJsonMessage.data;
					if (
						score !== null &&
						(currentScore === null || score < currentScore)
					) {
						setCurrentScore(score);
					}
					setLastExecStatus(status);
				}
			} else {
				if (game.started_at) {
					const nowSec = Math.floor(Date.now() / 1000);
					if (game.started_at <= nowSec) {
						// The game has already started.
						if (gameState !== "gaming" && gameState !== "finished") {
							setStartedAt(game.started_at);
							setLeftTimeSeconds(0);
							setGameState("gaming");
						}
					} else {
						// The game is starting.
						if (
							gameState !== "starting" &&
							gameState !== "gaming" &&
							gameState !== "finished"
						) {
							setStartedAt(game.started_at);
							setLeftTimeSeconds(game.started_at - nowSec);
							setGameState("starting");
						}
					}
				} else {
					setGameState("waiting");
				}
			}
		}
	}, [
		game.started_at,
		sendJsonMessage,
		lastJsonMessage,
		readyState,
		gameState,
		currentScore,
	]);

	if (gameState === "connecting") {
		return <GolfPlayAppConnecting />;
	} else if (gameState === "waiting") {
		return <GolfPlayAppWaiting />;
	} else if (gameState === "starting") {
		return <GolfPlayAppStarting leftTimeSeconds={leftTimeSeconds!} />;
	} else if (gameState === "gaming") {
		return (
			<GolfPlayAppGaming
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				onCodeChange={onCodeChange}
				onCodeSubmit={onCodeSubmit}
				currentScore={currentScore}
				lastExecStatus={lastExecStatus}
			/>
		);
	} else if (gameState === "finished") {
		return <GolfPlayAppFinished />;
	} else {
		return null;
	}
}
