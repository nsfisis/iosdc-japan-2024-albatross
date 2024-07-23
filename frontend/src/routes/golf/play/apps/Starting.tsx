type Props = {
  gameId: number;
  playerId: number;
  timeLeft: number | null;
};

export default function Starting({ timeLeft }: Props) {
  return (
    <>
      <div>対戦相手が見つかりました。{timeLeft} 秒後にゲームを開始します。</div>
    </>
  );
}
