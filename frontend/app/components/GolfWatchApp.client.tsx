import { useEffect, useState } from "react";
import type { components } from "../.server/api/schema";
import useWebSocket, { ReadyState } from "../hooks/useWebSocket";
import GolfWatchAppConnecting from "./GolfWatchApps/GolfWatchAppConnecting";
import GolfWatchAppFinished from "./GolfWatchApps/GolfWatchAppFinished";
import GolfWatchAppGaming, {
	PlayerInfo,
} from "./GolfWatchApps/GolfWatchAppGaming";
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

	const playerA = game.players[0];
	const playerB = game.players[1];

	const [playerInfoA, setPlayerInfoA] = useState<PlayerInfo>({
		displayName: playerA?.display_name ?? null,
		iconPath: playerA?.icon_path ?? null,
		score: null,
		code: "",
		submissionResult: undefined,
	});
	const [playerInfoB, setPlayerInfoB] = useState<PlayerInfo>({
		displayName: playerB?.display_name ?? null,
		iconPath: playerB?.icon_path ?? null,
		score: null,
		code: "",
		submissionResult: undefined,
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
					const { player_id, preliminary_score } = lastJsonMessage.data;
					const setter =
						player_id === playerA?.user_id ? setPlayerInfoA : setPlayerInfoB;
					setter((prev) => ({
						...prev,
						submissionResult: {
							status: "running",
							preliminaryScore: preliminary_score,
							verificationResults: game.verification_steps.map((v) => ({
								testcase_id: v.testcase_id,
								status: "running",
								label: v.label,
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
						if (ret.submissionResult === undefined) {
							return ret;
						}
						ret.submissionResult = {
							...ret.submissionResult,
							verificationResults: ret.submissionResult.verificationResults.map(
								(v) =>
									v.testcase_id === testcase_id && v.status === "running"
										? {
												...v,
												status,
												stdout,
												stderr,
											}
										: v,
							),
						};
						return ret;
					});
				} else if (lastJsonMessage.type === "watcher:s2c:submitresult") {
					const { player_id, status } = lastJsonMessage.data;
					const setter =
						player_id === playerA?.user_id ? setPlayerInfoA : setPlayerInfoB;
					setter((prev) => {
						const ret = { ...prev };
						if (ret.submissionResult === undefined) {
							return ret;
						}
						ret.submissionResult = {
							...ret.submissionResult,
							status,
						};
						if (status === "success") {
							if (
								ret.score === null ||
								ret.submissionResult.preliminaryScore < ret.score
							) {
								ret.score = ret.submissionResult.preliminaryScore;
							}
						} else {
							ret.submissionResult.verificationResults =
								ret.submissionResult.verificationResults.map((v) =>
									v.status === "running" ? { ...v, status: "canceled" } : v,
								);
						}
						return ret;
					});
				}
			} else {
				setGameState("waiting");
			}
		}
	}, [
		game.verification_steps,
		lastJsonMessage,
		readyState,
		gameState,
		playerA?.user_id,
		playerB?.user_id,
	]);

	if (gameState === "connecting") {
		return <GolfWatchAppConnecting />;
	} else if (gameState === "waiting") {
		return <GolfWatchAppWaiting />;
	} else if (gameState === "starting") {
		return <GolfWatchAppStarting leftTimeSeconds={leftTimeSeconds!} />;
	} else if (gameState === "gaming") {
		return (
			<GolfWatchAppGaming
				problem={game.problem!.description}
				playerInfoA={playerInfoA}
				playerInfoB={playerInfoB}
				leftTimeSeconds={leftTimeSeconds!}
			/>
		);
	} else if (gameState === "finished") {
		return <GolfWatchAppFinished />;
	} else {
		return null;
	}
}
