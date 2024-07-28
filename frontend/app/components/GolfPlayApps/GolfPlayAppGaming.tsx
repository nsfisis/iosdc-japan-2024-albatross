export default function GolfPlayAppGaming({
  problem,
  onCodeChange,
  currentScore,
}: {
  problem: string;
  onCodeChange: (code: string) => void;
  currentScore: number | null;
}) {
  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onCodeChange(e.target.value);
  };

  return (
    <div style={{ display: "flex" }}>
      <div style={{ flex: 1, padding: "10px", borderRight: "1px solid #ccc" }}>
        <div>{problem}</div>
        <div>
          {currentScore == null ? "Score: -" : `Score: ${currentScore}`}
        </div>
      </div>
      <div style={{ flex: 1, padding: "10px" }}>
        <textarea
          style={{ width: "100%", height: "100%" }}
          onChange={handleTextChange}
        />
      </div>
    </div>
  );
}
