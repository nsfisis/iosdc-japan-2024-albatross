import { useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import type { components } from "../.server/api/schema";
import GolfWatchAppConnecting from "./GolfWatchApps/GolfWatchAppConnecting";
import GolfWatchAppFinished from "./GolfWatchApps/GolfWatchAppFinished";
import GolfWatchAppGaming from "./GolfWatchApps/GolfWatchAppGaming";
import GolfWatchAppStarting from "./GolfWatchApps/GolfWatchAppStarting";
import GolfWatchAppWaiting from "./GolfWatchApps/GolfWatchAppWaiting";

type WebSocketMessage = components["schemas"]["GameWatcherMessageS2C"];

type Game = components["schemas"]["Game"];

type GameState = "connecting" | "waiting" | "starting" | "gaming" | "finished";

export default function GolfWatchApp({
  game,
  sockToken,
}: {
  game: Game;
  sockToken: string;
}) {
  // const socketUrl = `wss://t.nil.ninja/iosdc-japan/2024/sock/golf/${game.game_id}/watch?token=${sockToken}`;
  const socketUrl =
    process.env.NODE_ENV === "development"
      ? `ws://localhost:8002/sock/golf/${game.game_id}/watch?token=${sockToken}`
      : `ws://api-server/sock/golf/${game.game_id}/watch?token=${sockToken}`;

  const { lastJsonMessage, readyState } = useWebSocket<WebSocketMessage>(
    socketUrl,
    {},
  );

  const [gameState, setGameState] = useState<GameState>("connecting");

  const [startedAt, setStartedAt] = useState<number | null>(null);

  const [timeLeftSeconds, setTimeLeftSeconds] = useState<number | null>(null);

  useEffect(() => {
    if (gameState === "starting" && startedAt !== null) {
      const timer1 = setInterval(() => {
        setTimeLeftSeconds((prev) => {
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

  const [scoreA, setScoreA] = useState<number | null>(null);
  const [scoreB, setScoreB] = useState<number | null>(null);
  const [codeA, setCodeA] = useState<string>("");
  const [codeB, setCodeB] = useState<string>("");

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
            setTimeLeftSeconds(start_at - nowSec);
            setGameState("starting");
          }
        } else if (lastJsonMessage.type === "watcher:s2c:code") {
          const { player_id, code } = lastJsonMessage.data;
          setCodeA(code);
        } else if (lastJsonMessage.type === "watcher:s2c:execresult") {
          const { score } = lastJsonMessage.data;
          if (score !== null && (scoreA === null || score < scoreA)) {
            setScoreA(score);
          }
        }
      } else {
        setGameState("waiting");
      }
    }
  }, [lastJsonMessage, readyState, gameState, scoreA]);

  if (gameState === "connecting") {
    return <GolfWatchAppConnecting />;
  } else if (gameState === "waiting") {
    return <GolfWatchAppWaiting />;
  } else if (gameState === "starting") {
    return <GolfWatchAppStarting timeLeft={timeLeftSeconds!} />;
  } else if (gameState === "gaming") {
    return (
      <GolfWatchAppGaming
        problem={game.problem!.description}
        codeA={codeA}
        scoreA={scoreA}
        codeB={codeB}
        scoreB={scoreB}
      />
    );
  } else if (gameState === "finished") {
    return <GolfWatchAppFinished />;
  } else {
    return null;
  }
}
