import App from "./play/App.tsx";

export default function GolfPlay() {
  return (
    <div>
      <h1>Golf Play</h1>
      <App gameId={0} playerId={0} />
    </div>
  );
}
