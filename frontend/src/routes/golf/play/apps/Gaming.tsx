import React from "react";

type Props = {
  gameId: number;
  playerId: number;
  problem: string | null;
  onCodeChange: (data: { code: string }) => void;
  score: number | null;
};

export default function Gaming({ problem, onCodeChange, score }: Props) {
  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onCodeChange({ code: e.target.value });
  };

  return (
    <div style={{ display: "flex" }}>
      <div style={{ flex: 1, padding: "10px", borderRight: "1px solid #ccc" }}>
        <div>{problem}</div>
        <div>{score == null ? "Score: -" : `Score: ${score} byte`}</div>
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
