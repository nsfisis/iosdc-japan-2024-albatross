import { useState, useEffect } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useDebouncedCallback } from "use-debounce";
import Connecting from "./apps/Connecting.tsx";
import Waiting from "./apps/Waiting.tsx";
import Starting from "./apps/Starting.tsx";
import Gaming from "./apps/Gaming.tsx";
import Finished from "./apps/Finished.tsx";
import Failed from "./apps/Failed.tsx";
import type { GameState } from "./GameState.ts";

type Props = {
  gameId: number;
  playerId: number;
};

type WebSocketMessage =
  | { type: "connect" }
  | { type: "prepare"; data: { problem: string } }
  | { type: "ready" }
  | { type: "start"; data: { startTime: string } }
  | {
      type: "finish";
      data: { yourScore: number | null; opponentScore: number | null };
    }
  | { type: "score"; data: { score: number } }
  | { type: "code"; data: { code: string } }
  | {
      type: "watch";
      data: {
        problem: string;
        scoreA: number | null;
        codeA: string;
        scoreB: number | null;
        codeB: string;
      };
    };

export default ({ gameId, playerId }: Props) => {
  // const socketUrl = `wss://t.nil.ninja/iosdc/2024/sock/golf/${gameId}/${playerId}/`;
  const socketUrl = `ws://localhost:8002/sock/golf/${gameId}/${playerId}/`;

  const { sendJsonMessage, lastJsonMessage, readyState } =
    useWebSocket<WebSocketMessage>(socketUrl);

  const [gameState, setGameState] = useState<GameState>("connecting");

  const [problem, setProblem] = useState<string | null>(null);

  // in seconds
  const [timeLeft, setTimeLeft] = useState<number | null>(null);
  useEffect(() => {
    if (gameState === "starting" && timeLeft !== null) {
      const timer = setInterval(() => {
        setTimeLeft((prevTime) => {
          // `prevTime` is not null because `timeLeft` is not null.
          prevTime = prevTime!;
          if (prevTime <= 1) {
            clearInterval(timer);
            setGameState("gaming");
            return 0;
          }
          return prevTime - 1;
        });
      }, 1000);

      return () => clearInterval(timer);
    }
  }, [gameState]);

  const [score, setScore] = useState<number | null>(null);

  const [result, setResult] = useState<{
    yourScore: number | null;
    opponentScore: number | null;
  } | null>(null);

  const onCodeChange = useDebouncedCallback((data) => {
    sendJsonMessage({ type: "code", data });
  }, 1000);

  useEffect(() => {
    if (readyState === ReadyState.UNINSTANTIATED) {
      setGameState("failed");
    } else if (
      readyState === ReadyState.CLOSING ||
      readyState === ReadyState.CLOSED
    ) {
      if (gameState !== "finished") {
        setGameState("failed");
      }
    } else if (readyState === ReadyState.CONNECTING) {
      setGameState("connecting");
    } else if (readyState === ReadyState.OPEN) {
      if (lastJsonMessage !== null) {
        if (lastJsonMessage.type === "prepare") {
          const { problem } = lastJsonMessage.data;
          setProblem(problem);
          sendJsonMessage({ type: "ready", data: {} });
        } else if (lastJsonMessage.type === "start") {
          const { startTime } = lastJsonMessage.data;
          const startTimeMs = Date.parse(startTime);
          setTimeLeft(
            Math.max(0, Math.floor((startTimeMs - Date.now()) / 1000)),
          );
          setGameState("starting");
        } else if (lastJsonMessage.type === "finish") {
          const result = lastJsonMessage.data;
          setResult(result);
          setGameState("finished");
        } else if (lastJsonMessage.type === "score") {
          const { score } = lastJsonMessage.data;
          setScore(score);
        } else {
          setGameState("failed");
        }
      } else {
        setGameState("waiting");
        sendJsonMessage({ type: "connect", data: {} });
      }
    }
  }, [readyState, lastJsonMessage]);

  return (
    <div>
      <h1>
        Game #{gameId} (playerId #{playerId})
      </h1>
      <div>
        {gameState === "connecting" ? (
          <Connecting gameId={gameId} playerId={playerId} />
        ) : gameState === "waiting" ? (
          <Waiting gameId={gameId} playerId={playerId} />
        ) : gameState === "starting" ? (
          <Starting gameId={gameId} playerId={playerId} timeLeft={timeLeft} />
        ) : gameState === "gaming" ? (
          <Gaming
            gameId={gameId}
            playerId={playerId}
            problem={problem}
            score={score}
            onCodeChange={onCodeChange}
          />
        ) : gameState === "finished" ? (
          <Finished gameId={gameId} playerId={playerId} result={result!} />
        ) : (
          <Failed gameId={gameId} playerId={playerId} />
        )}
      </div>
    </div>
  );
};
