import { PlayerInfo } from "../../models/PlayerInfo";
import { submissionResultStatusToLabel } from "../../models/SubmissionResult";
import ExecStatusIndicatorIcon from "../ExecStatusIndicatorIcon";

type Props = {
	problem: string;
	playerInfoA: PlayerInfo;
	playerInfoB: PlayerInfo;
	leftTimeSeconds: number;
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
					{playerInfoA.score}
				</div>
				<div className="w-full bg-blue-500">
					<div
						className="h-full bg-red-500"
						style={{ width: `${scoreRatio}%` }}
					></div>
				</div>
				<div className="grid justify-end bg-blue-500 p-2 text-lg font-bold text-white">
					{playerInfoB.score}
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
						)}
					</div>
					<div>
						<ol>
							{playerInfoA.submissionResult?.verificationResults.map(
								(result) => (
									<li key={result.testcase_id ?? -1}>
										<div>
											<div>
												<ExecStatusIndicatorIcon status={result.status} />{" "}
												{result.label}
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
						)}
					</div>
					<div>
						<ol>
							{playerInfoB.submissionResult?.verificationResults.map(
								(result, idx) => (
									<li key={idx}>
										<div>
											<div>
												<ExecStatusIndicatorIcon status={result.status} />{" "}
												{result.label}
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
