type Props = {
  gameId: number;
  problem: string | null;
  scoreA: number | null;
  codeA: string | null;
  scoreB: number | null;
  codeB: string | null;
};

export default function Gaming({
  problem,
  scoreA,
  codeA,
  scoreB,
  codeB,
}: Props) {
  return (
    <>
      <div style={{ display: "flex", flexDirection: "column" }}>
        <div style={{ display: "flex", flex: 1, justifyContent: "center" }}>
          {problem}
        </div>
        <div style={{ display: "flex", flex: 3 }}>
          <div style={{ display: "flex", flex: 3, flexDirection: "column" }}>
            <div style={{ flex: 1, justifyContent: "center" }}>{scoreA}</div>
            <div style={{ flex: 3 }}>
              <pre>
                <code>{codeA}</code>
              </pre>
            </div>
          </div>
          <div style={{ display: "flex", flex: 3, flexDirection: "column" }}>
            <div style={{ flex: 1, justifyContent: "center" }}>{scoreB}</div>
            <div style={{ flex: 3 }}>
              <pre>
                <code>{codeB}</code>
              </pre>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
