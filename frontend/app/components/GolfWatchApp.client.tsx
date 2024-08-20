import { useEffect, useState } from "react";
import type { components } from "../.server/api/schema";
import useWebSocket, { ReadyState } from "../hooks/useWebSocket";
import type { PlayerInfo } from "../models/PlayerInfo";
import GolfWatchAppConnecting from "./GolfWatchApps/GolfWatchAppConnecting";
import GolfWatchAppGaming from "./GolfWatchApps/GolfWatchAppGaming";
import GolfWatchAppStarting from "./GolfWatchApps/GolfWatchAppStarting";
import GolfWatchAppWaiting from "./GolfWatchApps/GolfWatchAppWaiting";

type GameWatcherMessageS2C = components["schemas"]["GameWatcherMessageS2C"];
type GameWatcherMessageC2S = never;

type Game = components["schemas"]["Game"];

type GameState = "connecting" | "waiting" | "starting" | "gaming" | "finished";

export default function GolfWatchApp({
	game,
	sockToken,
}: {
	game: Game;
	sockToken: string;
}) {
	const socketUrl =
		process.env.NODE_ENV === "development"
			? `ws://localhost:8002/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/watch?token=${sockToken}`
			: `wss://t.nil.ninja/iosdc-japan/2024/code-battle/sock/golf/${game.game_id}/watch?token=${sockToken}`;

	const { lastJsonMessage, readyState } = useWebSocket<
		GameWatcherMessageS2C,
		GameWatcherMessageC2S
	>(socketUrl);

	const [gameState, setGameState] = useState<GameState>("connecting");

	const [startedAt, setStartedAt] = useState<number | null>(null);

	const [leftTimeSeconds, setLeftTimeSeconds] = useState<number | null>(null);

	useEffect(() => {
		if (
			(gameState === "starting" || gameState === "gaming") &&
			startedAt !== null
		) {
			const timer = setInterval(() => {
				setLeftTimeSeconds((prev) => {
					if (prev === null) {
						return null;
					}
					if (prev <= 1) {
						const nowSec = Math.floor(Date.now() / 1000);
						const finishedAt = startedAt + game.duration_seconds;
						if (nowSec >= finishedAt) {
							clearInterval(timer);
							setGameState("finished");
						} else {
							setGameState("gaming");
						}
					}
					return prev - 1;
				});
			}, 1000);

			return () => {
				clearInterval(timer);
			};
		}
	}, [gameState, startedAt, game.duration_seconds]);

	const playerA = game.players[0];
	const playerB = game.players[1];

	const [playerInfoA, setPlayerInfoA] = useState<PlayerInfo>({
		displayName: playerA?.display_name ?? null,
		iconPath: playerA?.icon_path ?? null,
		score: null,
		code: "",
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
	});
	const [playerInfoB, setPlayerInfoB] = useState<PlayerInfo>({
		displayName: playerB?.display_name ?? null,
		iconPath: playerB?.icon_path ?? null,
		score: null,
		code: "",
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
	});

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
				if (lastJsonMessage.type === "watcher:s2c:start") {
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
				} else if (lastJsonMessage.type === "watcher:s2c:code") {
					const { player_id, code } = lastJsonMessage.data;
					const setter =
						player_id === playerA?.user_id ? setPlayerInfoA : setPlayerInfoB;
					setter((prev) => ({ ...prev, code }));
				} else if (lastJsonMessage.type === "watcher:s2c:submit") {
					const { player_id } = lastJsonMessage.data;
					const setter =
						player_id === playerA?.user_id ? setPlayerInfoA : setPlayerInfoB;
					setter((prev) => ({
						...prev,
						submitResult: {
							status: "running",
							execResults: prev.submitResult.execResults.map((r) => ({
								...r,
								status: "running",
								stdout: "",
								stderr: "",
							})),
						},
					}));
				} else if (lastJsonMessage.type === "watcher:s2c:execresult") {
					const { player_id, testcase_id, status, stdout, stderr } =
						lastJsonMessage.data;
					const setter =
						player_id === playerA?.user_id ? setPlayerInfoA : setPlayerInfoB;
					setter((prev) => {
						const ret = { ...prev };
						ret.submitResult = {
							...prev.submitResult,
							execResults: prev.submitResult.execResults.map((r) =>
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
						return ret;
					});
				} else if (lastJsonMessage.type === "watcher:s2c:submitresult") {
					const { player_id, status, score } = lastJsonMessage.data;
					const setter =
						player_id === playerA?.user_id ? setPlayerInfoA : setPlayerInfoB;
					setter((prev) => {
						const ret = { ...prev };
						ret.submitResult = {
							...prev.submitResult,
							status,
						};
						if (status === "success") {
							if (score) {
								if (ret.score === null || score < ret.score) {
									ret.score = score;
								}
							}
						} else {
							ret.submitResult.execResults = prev.submitResult.execResults.map(
								(r) =>
									r.status === "running" ? { ...r, status: "canceled" } : r,
							);
						}
						return ret;
					});
				}
			} else {
				if (game.started_at) {
					const nowSec = Math.floor(Date.now() / 1000);
					if (game.started_at <= nowSec) {
						// The game has already started.
						if (gameState !== "gaming" && gameState !== "finished") {
							setStartedAt(game.started_at);
							setLeftTimeSeconds(game.started_at - nowSec);
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
		lastJsonMessage,
		readyState,
		gameState,
		playerA?.user_id,
		playerB?.user_id,
	]);

	if (gameState === "connecting") {
		return <GolfWatchAppConnecting />;
	} else if (gameState === "waiting") {
		return (
			<GolfWatchAppWaiting
				gameDisplayName={game.display_name}
				playerInfoA={playerInfoA}
				playerInfoB={playerInfoB}
			/>
		);
	} else if (gameState === "starting") {
		return (
			<GolfWatchAppStarting
				gameDisplayName={game.display_name}
				leftTimeSeconds={leftTimeSeconds!}
			/>
		);
	} else if (gameState === "gaming" || gameState === "finished") {
		return (
			<GolfWatchAppGaming
				gameDisplayName={game.display_name}
				gameDurationSeconds={game.duration_seconds}
				leftTimeSeconds={leftTimeSeconds!}
				playerInfoA={playerInfoA}
				playerInfoB={playerInfoB}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				gameResult={null /* TODO */}
			/>
		);
	} else {
		return null;
	}
}
