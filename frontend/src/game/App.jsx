import { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { useDebouncedCallback } from 'use-debounce';
import Connecting from './apps/Connecting.jsx';
import Waiting from './apps/Waiting.jsx';
import Starting from './apps/Starting.jsx';
import Gaming from './apps/Gaming.jsx';
import Finished from './apps/Finished.jsx';
import Failed from './apps/Failed.jsx';
import { GAME_STATE_CONNECTING, GAME_STATE_WAITING, GAME_STATE_STARTING, GAME_STATE_GAMING, GAME_STATE_FINISHED, GAME_STATE_FAILED } from './GameState.js';

export default ({ gameId, team }) => {
  // const socketUrl = `wss://t.nil.ninja/iosdc/2024/sock/golf/${gameId}/${team}/`;
  const socketUrl = `ws://localhost:8002/sock/golf/${gameId}/${team}/`;

  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(socketUrl);

  const [gameState, setGameState] = useState(GAME_STATE_CONNECTING);

  const [problem, setProblem] = useState(null);

  // in seconds
  const [timeLeft, setTimeLeft] = useState(null);
  useEffect(() => {
    if (gameState === GAME_STATE_STARTING && timeLeft !== null) {
      const timer = setInterval(() => {
        setTimeLeft(prevTime => {
          if (prevTime <= 1) {
            clearInterval(timer);
            setGameState(GAME_STATE_GAMING);
            return 0;
          }
          return prevTime - 1;
        });
      }, 1000);

      return () => clearInterval(timer);
    }
  }, [gameState]);

  const [score, setScore] = useState(null);

  const [result, setResult] = useState(null);

  const onCodeChange = useDebouncedCallback((data) => {
    sendJsonMessage({ type: 'code', data });
  }, 1000);

  useEffect(() => {
    if (readyState === ReadyState.UNINSTANTIATED) {
      setGameState(GAME_STATE_FAILED);
    } else if (readyState === ReadyState.CLOSING || readyState === ReadyState.CLOSED) {
      if (gameState !== GAME_STATE_FINISHED) {
        setGameState(GAME_STATE_FAILED);
      }
    } else if (readyState === ReadyState.CONNECTING) {
      setGameState(GAME_STATE_CONNECTING);
    } else if (readyState === ReadyState.OPEN) {
      if (lastJsonMessage !== null) {
        if (lastJsonMessage.type === 'prepare') {
          const { problem } = lastJsonMessage.data;
          setProblem(problem);
          sendJsonMessage({ type: 'ready', data: {} });
        } else if (lastJsonMessage.type === 'start') {
          const { startTime } = lastJsonMessage.data;
          const startTimeMs = Date.parse(startTime);
          setTimeLeft(Math.max(0, Math.floor((startTimeMs - Date.now()) / 1000)));
          setGameState(GAME_STATE_STARTING);
        } else if (lastJsonMessage.type === 'finish') {
          const result = lastJsonMessage.data;
          setResult(result);
          setGameState(GAME_STATE_FINISHED);
        } else if (lastJsonMessage.type === 'score') {
          const { score } = lastJsonMessage.data;
          setScore(score);
        } else {
          setGameState(GAME_STATE_FAILED);
        }
      } else {
        setGameState(GAME_STATE_WAITING);
        sendJsonMessage({ type: 'connect', data: {} });
      }
    }
  }, [readyState, lastJsonMessage]);

  return (
    <div>
      <h1>Game #{gameId} (team #{team})</h1>
      <div>
        { gameState === GAME_STATE_CONNECTING ? (<Connecting gameId={gameId} team={team} />)
          : gameState === GAME_STATE_WAITING ? (<Waiting gameId={gameId} team={team} />)
          : gameState === GAME_STATE_STARTING ? (<Starting gameId={gameId} team={team} timeLeft={timeLeft} />)
          : gameState === GAME_STATE_GAMING ? (<Gaming gameId={gameId} team={team} problem={problem} score={score} onCodeChange={onCodeChange} />)
          : gameState === GAME_STATE_FINISHED ? (<Finished gameId={gameId} team={team} result={result} />)
          : (<Failed />) }
      </div>
    </div>
  );
};
