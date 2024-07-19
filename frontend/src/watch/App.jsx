import { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import Connecting from './apps/Connecting.jsx';
import Waiting from './apps/Waiting.jsx';
import Gaming from './apps/Gaming.jsx';
import Failed from './apps/Failed.jsx';
import { WATCH_STATE_CONNECTING, WATCH_STATE_WAITING, WATCH_STATE_GAMING, WATCH_STATE_FINISHED, WATCH_STATE_FAILED } from './WatchState.js';

export default ({ gameId }) => {
  // const socketUrl = `wss://t.nil.ninja/iosdc/2024/sock/golf/${gameId}/watch/`;
  const socketUrl = `ws://localhost:8002/sock/golf/${gameId}/watch/`;

  const { lastJsonMessage, readyState } = useWebSocket(socketUrl);

  const [watchState, setWatchState] = useState(WATCH_STATE_CONNECTING);

  const [problem, setProblem] = useState(null);

  const [scoreA, setScoreA] = useState(null);
  const [codeA, setCodeA] = useState(null);
  const [scoreB, setScoreB] = useState(null);
  const [codeB, setCodeB] = useState(null);

  useEffect(() => {
    if (readyState === ReadyState.UNINSTANTIATED) {
      setWatchState(WATCH_STATE_FAILED);
    } else if (readyState === ReadyState.CLOSING || readyState === ReadyState.CLOSED) {
      if (watchState !== WATCH_STATE_FINISHED) {
        setWatchState(WATCH_STATE_FAILED);
      }
    } else if (readyState === ReadyState.CONNECTING) {
      setWatchState(WATCH_STATE_CONNECTING);
    } else if (readyState === ReadyState.OPEN) {
      if (lastJsonMessage !== null) {
        if (lastJsonMessage.type === 'watch') {
          const { problem, scoreA: scoreA_, codeA: codeA_, scoreB: scoreB_, codeB: codeB_ } = lastJsonMessage.data;
          setProblem(problem);
          setScoreA(scoreA_);
          setCodeA(codeA_);
          setScoreB(scoreB_);
          setCodeB(codeB_);
          setWatchState(WATCH_STATE_GAMING);
        } else {
          setWatchState(WATCH_STATE_FAILED);
        }
      } else {
        setWatchState(WATCH_STATE_WAITING);
      }
    }
  }, [readyState, lastJsonMessage]);

  return (
    <div>
      <h1>Game #{gameId} watching</h1>
      <div>
        { watchState === WATCH_STATE_CONNECTING ? (<Connecting gameId={gameId} />)
          : watchState === WATCH_STATE_WAITING ? (<Waiting gameId={gameId} />)
          : watchState === WATCH_STATE_GAMING || watchState === WATCH_STATE_FINISHED ? (<Gaming gameId={gameId} problem={problem} scoreA={scoreA} codeA={codeA} scoreB={scoreB} codeB={codeB} />)
          : (<Failed />) }
      </div>
    </div>
  );
};
