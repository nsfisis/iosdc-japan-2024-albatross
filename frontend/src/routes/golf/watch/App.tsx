import { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import Connecting from './apps/Connecting.jsx';
import Waiting from './apps/Waiting.jsx';
import Gaming from './apps/Gaming.jsx';
import Failed from './apps/Failed.jsx';
import type { WatchState } from './WatchState.js';

type Props = {
  gameId: number;
};

type WebSocketMessage =
  { type: 'connect' }
  | { type: 'prepare'; data: { problem: string } }
  | { type: 'ready' }
  | { type: 'start'; data: { startTime: string } }
  | { type: 'finish'; data: { yourScore: number | null; opponentScore: number | null } }
  | { type: 'score'; data: { score: number } }
  | { type: 'code'; data: { code: string } }
  | { type: 'watch'; data: { problem: string; scoreA: number | null; codeA: string; scoreB: number | null; codeB: string } }

export default ({ gameId }: Props) => {
  // const socketUrl = `wss://t.nil.ninja/iosdc/2024/sock/golf/${gameId}/watch/`;
  const socketUrl = `ws://localhost:8002/sock/golf/${gameId}/watch/`;

  const { lastJsonMessage, readyState } = useWebSocket<WebSocketMessage>(socketUrl);

  const [watchState, setWatchState] = useState<WatchState>('connecting');

  const [problem, setProblem] = useState<string | null>(null);

  const [scoreA, setScoreA] = useState<number | null>(null);
  const [codeA, setCodeA] = useState<string | null>(null);
  const [scoreB, setScoreB] = useState<number | null>(null);
  const [codeB, setCodeB] = useState<string | null>(null);

  useEffect(() => {
    if (readyState === ReadyState.UNINSTANTIATED) {
      setWatchState('failed');
    } else if (readyState === ReadyState.CLOSING || readyState === ReadyState.CLOSED) {
      if (watchState !== 'finished') {
        setWatchState('failed');
      }
    } else if (readyState === ReadyState.CONNECTING) {
      setWatchState('connecting');
    } else if (readyState === ReadyState.OPEN) {
      if (lastJsonMessage !== null) {
        if (lastJsonMessage.type === 'watch') {
          const { problem, scoreA: scoreA_, codeA: codeA_, scoreB: scoreB_, codeB: codeB_ } = lastJsonMessage.data;
          setProblem(problem);
          setScoreA(scoreA_);
          setCodeA(codeA_);
          setScoreB(scoreB_);
          setCodeB(codeB_);
          setWatchState('gaming');
        } else {
          setWatchState('failed');
        }
      } else {
        setWatchState('waiting');
      }
    }
  }, [readyState, lastJsonMessage]);

  return (
    <div>
      <h1>Game #{gameId} watching</h1>
      <div>
        {watchState === 'connecting' ? (<Connecting gameId={gameId} />)
          : watchState === 'waiting' ? (<Waiting gameId={gameId} />)
            : watchState === 'gaming' || watchState === 'finished' ? (<Gaming gameId={gameId} problem={problem} scoreA={scoreA} codeA={codeA} scoreB={scoreB} codeB={codeB} />)
              : (<Failed gameId={gameId} />)}
      </div>
    </div>
  );
};
