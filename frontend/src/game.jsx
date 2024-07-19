import { createRoot } from 'react-dom/client';
import App from './game/App.jsx';

const url = new URL(window.location.href);
const path = url.pathname;
const match = path.match(/\/golf\/(\d+)\/(a|b)\/$/);
if (match) {
  const gameId = match[1];
  const team = match[2];

  createRoot(document.getElementById('app')).render(<App gameId={gameId} team={team} />);
}
