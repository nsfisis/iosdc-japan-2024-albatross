import { useEffect, useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import type { components } from "../.server/api/schema";
import useWebSocket, { ReadyState } from "../hooks/useWebSocket";
import type { PlayerState } from "../types/PlayerState";
import GolfPlayAppConnecting from "./GolfPlayApps/GolfPlayAppConnecting";
import GolfPlayAppFinished from "./GolfPlayApps/GolfPlayAppFinished";
import GolfPlayAppGaming from "./GolfPlayApps/GolfPlayAppGaming";
import GolfPlayAppStarting from "./GolfPlayApps/GolfPlayAppStarting";
import GolfPlayAppWaiting from "./GolfPlayApps/GolfPlayAppWaiting";

type GamePlayerMessageS2C = components["schemas"]["GamePlayerMessageS2C"];
type GamePlayerMessageC2S = components["schemas"]["GamePlayerMessageC2S"];

type Game = components["schemas"]["Game"];
type User = components["schemas"]["User"];

type GameState = "connecting" | "waiting" | "starting" | "gaming" | "finished";

export default function GolfPlayApp({
	game,
	player,
	sockToken,
}: {
	game: Game;
	player: User;
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

	const playerProfile = {
		displayName: player.display_name,
		iconPath: player.icon_path ?? null,
	};
	const [playerState, setPlayerState] = useState<PlayerState>({
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
	});

	const onCodeChange = useDebouncedCallback((code: string) => {
		console.log("player:c2s:code");
		sendJsonMessage({
			type: "player:c2s:code",
			data: { code },
		});
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
		setPlayerState((prev) => ({
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
					const { testcase_id, status, stdout, stderr } = lastJsonMessage.data;
					setPlayerState((prev) => {
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
				} else if (lastJsonMessage.type === "player:s2c:submitresult") {
					const { status, score } = lastJsonMessage.data;
					setPlayerState((prev) => {
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
		sendJsonMessage,
		lastJsonMessage,
		readyState,
		gameState,
	]);

	if (gameState === "connecting") {
		return <GolfPlayAppConnecting />;
	} else if (gameState === "waiting") {
		return (
			<GolfPlayAppWaiting
				gameDisplayName={game.display_name}
				playerProfile={playerProfile}
			/>
		);
	} else if (gameState === "starting") {
		return (
			<GolfPlayAppStarting
				gameDisplayName={game.display_name}
				leftTimeSeconds={leftTimeSeconds!}
			/>
		);
	} else if (gameState === "gaming") {
		return (
			<GolfPlayAppGaming
				gameDisplayName={game.display_name}
				gameDurationSeconds={game.duration_seconds}
				leftTimeSeconds={leftTimeSeconds!}
				playerInfo={{
					profile: playerProfile,
					state: playerState,
				}}
				problemTitle={game.problem.title}
				problemDescription={game.problem.description}
				onCodeChange={onCodeChange}
				onCodeSubmit={onCodeSubmit}
			/>
		);
	} else if (gameState === "finished") {
		return <GolfPlayAppFinished />;
	} else {
		return null;
	}
}
