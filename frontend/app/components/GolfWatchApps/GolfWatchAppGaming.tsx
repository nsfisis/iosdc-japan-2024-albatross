type Props = {
	problem: string;
	playerInfoA: PlayerInfo;
	playerInfoB: PlayerInfo;
	leftTimeSeconds: number;
};

export type PlayerInfo = {
	displayName: string | null;
	iconPath: string | null;
	score: number | null;
	code: string | null;
	submissionResult?: SubmissionResult;
};

type SubmissionResult = {
	status: string;
	nextScore: number;
	executionResults: ExecutionResult[];
};

type ExecutionResult = {
	status: string;
	label: string;
	output: string;
};

export default function GolfWatchAppGaming({
	problem,
	playerInfoA,
	playerInfoB,
	leftTimeSeconds,
}: Props) {
	const leftTime = (() => {
		const m = Math.floor(leftTimeSeconds / 60);
		const s = leftTimeSeconds % 60;
		return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
	})();
	const scoreRatio = (() => {
		const scoreA = playerInfoA.score ?? 0;
		const scoreB = playerInfoB.score ?? 0;
		const totalScore = scoreA + scoreB;
		return totalScore === 0 ? 50 : (scoreA / totalScore) * 100;
	})();

	return (
		<div className="grid h-full w-full grid-rows-[auto_auto_1fr_auto]">
			<div className="grid grid-cols-[1fr_auto_1fr]">
				<div className="grid justify-start bg-red-500 p-2 text-white">
					{playerInfoA.displayName}
				</div>
				<div className="grid justify-center p-2">{leftTime}</div>
				<div className="grid justify-end bg-blue-500 p-2 text-white">
					{playerInfoB.displayName}
				</div>
			</div>
			<div className="grid grid-cols-[auto_1fr_auto]">
				<div className="grid justify-start bg-red-500 p-2 text-lg font-bold text-white">
					{playerInfoA.score ?? "-"}
				</div>
				<div className="w-full bg-blue-500">
					<div
						className="h-full bg-red-500"
						style={{ width: `${scoreRatio}%` }}
					></div>
				</div>
				<div className="grid justify-end bg-blue-500 p-2 text-lg font-bold text-white">
					{playerInfoB.score ?? "-"}
				</div>
			</div>
			<div className="grid grid-cols-[2fr_1fr_2fr_1fr] p-2">
				<div>
					<pre>
						<code>{playerInfoA.code}</code>
					</pre>
				</div>
				<div>
					<div>
						{playerInfoA.submissionResult?.status}(
						{playerInfoA.submissionResult?.nextScore})
					</div>
					<div>
						<ol>
							{playerInfoA.submissionResult?.executionResults.map(
								(result, idx) => (
									<li key={idx}>
										<div>
											<div>
												{result.status} {result.label}
											</div>
											<div>{result.output}</div>
										</div>
									</li>
								),
							)}
						</ol>
					</div>
				</div>
				<div>
					<pre>
						<code>{playerInfoB.code}</code>
					</pre>
				</div>
				<div>
					<div>
						{playerInfoB.submissionResult?.status}(
						{playerInfoB.submissionResult?.nextScore})
					</div>
					<div>
						<ol>
							{playerInfoB.submissionResult?.executionResults.map(
								(result, idx) => (
									<li key={idx}>
										<div>
											<div>
												{result.status} {result.label}
											</div>
											<div>{result.output}</div>
										</div>
									</li>
								),
							)}
						</ol>
					</div>
				</div>
			</div>
			<div className="grid justify-center p-2 bg-slate-300">{problem}</div>
		</div>
	);
}
