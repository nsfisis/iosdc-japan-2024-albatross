type Props = {
  problem: string;
  onCodeChange: (code: string) => void;
  currentScore: number | null;
};

export default function GolfPlayAppGaming({
  problem,
  onCodeChange,
  currentScore,
}: Props) {
  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    onCodeChange(e.target.value);
  };

  return (
    <div className="min-h-screen flex">
      <div className="mx-auto flex min-h-full flex-grow">
        <div className="flex w-1/2 flex-col justify-between p-4">
          <div>
            <div className="mb-2 text-xl font-bold">TODO</div>
            <div className="text-gray-700">{problem}</div>
          </div>
          <div className="mb-4 mt-auto">
            <div className="font-semibold text-green-500">
              Score: {currentScore == null ? "-" : `${currentScore}`}
            </div>
          </div>
        </div>
        <div className="w-1/2 p-4 flex">
          <div className="flex-grow">
            <textarea
              className="h-full w-full rounded-lg border border-gray-300 p-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              onChange={handleTextChange}
            ></textarea>
          </div>
        </div>
      </div>
    </div>
  );
}
