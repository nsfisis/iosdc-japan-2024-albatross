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
	status:
		| "running"
		| "success"
		| "wrong_answer"
		| "timeout"
		| "compile_error"
		| "runtime_error"
		| "internal_error";
	preliminaryScore: number;
	verificationResults: VerificationResult[];
};

type VerificationResult = {
	testcase_id: number | null;
	status:
		| "running"
		| "success"
		| "wrong_answer"
		| "timeout"
		| "compile_error"
		| "runtime_error"
		| "internal_error"
		| "canceled";
	label: string;
	stdout: string;
	stderr: string;
};

function submissionResultStatusToLabel(
	status: SubmissionResult["status"] | null,
) {
	switch (status) {
		case null:
			return "-";
		case "running":
			return "Running...";
		case "success":
			return "Accepted";
		case "wrong_answer":
			return "Wrong Answer";
		case "timeout":
			return "Time Limit Exceeded";
		case "compile_error":
			return "Compile Error";
		case "runtime_error":
			return "Runtime Error";
		case "internal_error":
			return "Internal Error";
	}
}

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
		return totalScore === 0 ? 50 : (scoreB / totalScore) * 100;
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
			<div className="grid grid-cols-[3fr_2fr_3fr_2fr] p-2">
				<div>
					<pre>
						<code>{playerInfoA.code}</code>
					</pre>
				</div>
				<div>
					<div>
						{submissionResultStatusToLabel(
							playerInfoA.submissionResult?.status ?? null,
						)}{" "}
						({playerInfoA.submissionResult?.preliminaryScore})
					</div>
					<div>
						<ol>
							{playerInfoA.submissionResult?.verificationResults.map(
								(result, idx) => (
									<li key={idx}>
										<div>
											<div>
												{result.status} {result.label}
											</div>
											<div>
												{result.stdout}
												{result.stderr}
											</div>
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
						{submissionResultStatusToLabel(
							playerInfoB.submissionResult?.status ?? null,
						)}{" "}
						({playerInfoB.submissionResult?.preliminaryScore ?? "-"})
					</div>
					<div>
						<ol>
							{playerInfoB.submissionResult?.verificationResults.map(
								(result, idx) => (
									<li key={idx}>
										<div>
											<div>
												{result.status} {result.label}
											</div>
											<div>
												{result.stdout}
												{result.stderr}
											</div>
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
