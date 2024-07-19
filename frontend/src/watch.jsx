import { createRoot } from 'react-dom/client';
import App from './watch/App.jsx';

const url = new URL(window.location.href);
const path = url.pathname;
const match = path.match(/\/golf\/(\d+)\/watch\/$/);
if (match) {
  const gameId = match[1];

  createRoot(document.getElementById('app')).render(<App gameId={gameId} />);
}
