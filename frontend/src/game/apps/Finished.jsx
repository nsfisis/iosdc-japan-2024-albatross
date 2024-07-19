export default ({ result }) => {
  const { yourScore, opponentScore } = result;
  const yourScoreToCompare = yourScore ?? Infinity;
  const opponentScoreToCompare = opponentScore ?? Infinity;
  const resultText = yourScoreToCompare === opponentScoreToCompare ? '引き分け' : (yourScoreToCompare < opponentScoreToCompare ? 'あなたの勝ち' : 'あなたの負け');
  return (
    <>
      <div>
        対戦終了
      </div>
      <div>
        <div>
          {resultText}
        </div>
        <div>
          あなたのスコア: {yourScore ?? 'なし'}
        </div>
        <div>
          相手のスコア: {opponentScore ?? 'なし'}
        </div>
      </div>
    </>
  );
}
