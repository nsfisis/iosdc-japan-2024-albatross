import { faArrowDown } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { PlayerInfo } from "../../models/PlayerInfo";
import BorderedContainer from "../BorderedContainer";
import ExecStatusIndicatorIcon from "../ExecStatusIndicatorIcon";
import SubmitStatusLabel from "../SubmitStatusLabel";

type Props = {
	gameDisplayName: string;
	gameDurationSeconds: number;
	leftTimeSeconds: number;
	playerInfoA: PlayerInfo;
	playerInfoB: PlayerInfo;
	problemTitle: string;
	problemDescription: string;
};

export default function GolfWatchAppGaming({
	gameDisplayName,
	gameDurationSeconds,
	leftTimeSeconds,
	playerInfoA,
	playerInfoB,
}: Props) {
	const leftTime = (() => {
		const k = gameDurationSeconds + leftTimeSeconds;
		const m = Math.floor(k / 60);
		const s = k % 60;
		return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
	})();

	const scoreRatio = (() => {
		const scoreA = playerInfoA.score;
		const scoreB = playerInfoB.score;
		if (scoreA === null && scoreB === null) {
			return 50;
		} else if (scoreA === null) {
			return 0;
		} else if (scoreB === null) {
			return 100;
		} else {
			return (scoreB / (scoreA + scoreB)) * 100;
		}
	})();

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className="text-white bg-iosdc-japan grid grid-cols-3 px-4 py-2">
				<div className="font-bold flex justify-between my-auto">
					<div className="flex gap-4">
						{playerInfoA.iconPath && (
							<img
								src={
									process.env.NODE_ENV === "development"
										? `http://localhost:8002/iosdc-japan/2024/code-battle${playerInfoA.iconPath}`
										: `/iosdc-japan/2024/code-battle${playerInfoA.iconPath}`
								}
								alt={`${playerInfoA.displayName} のアイコン`}
								className="w-12 h-12 rounded-full my-auto border-4 border-white"
							/>
						)}
						<div>
							<div className="text-gray-100">Player 1</div>
							<div className="text-2xl">{playerInfoA.displayName}</div>
						</div>
					</div>
					<div className="text-6xl">{playerInfoA.score}</div>
				</div>
				<div className="font-bold text-center">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-3xl">{leftTime}</div>
				</div>
				<div className="font-bold flex justify-between my-auto">
					<div className="text-6xl">{playerInfoB.score}</div>
					<div className="flex gap-4 text-end">
						<div>
							<div className="text-gray-100">Player 2</div>
							<div className="text-2xl">{playerInfoB.displayName}</div>
						</div>
						{playerInfoB.iconPath && (
							<img
								src={
									process.env.NODE_ENV === "development"
										? `http://localhost:8002/iosdc-japan/2024/code-battle${playerInfoB.iconPath}`
										: `/iosdc-japan/2024/code-battle${playerInfoB.iconPath}`
								}
								alt={`${playerInfoB.displayName} のアイコン`}
								className="w-12 h-12 rounded-full my-auto border-4 border-white"
							/>
						)}
					</div>
				</div>
			</div>
			<div className="w-full bg-purple-400">
				<div
					className="h-6 bg-orange-400"
					style={{ width: `${scoreRatio}%` }}
				></div>
			</div>
			<div className="grow grid grid-cols-10 p-4 gap-4">
				<div className="col-span-3">
					<pre className="bg-white resize-none h-full w-full rounded-lg border border-gray-300 p-2">
						<code>{playerInfoA.code}</code>
					</pre>
				</div>
				<div className="col-span-2 flex flex-col gap-4">
					<div className="flex">
						<div className="grow font-bold text-xl text-center">
							<SubmitStatusLabel status={playerInfoA.submitResult.status} />
						</div>
					</div>
					<ul className="flex flex-col gap-2">
						{playerInfoA.submitResult.execResults.map((r, idx) => (
							<li key={r.testcase_id ?? -1} className="flex gap-2">
								<div className="flex flex-col gap-2 p-2">
									<div className="w-6">
										<ExecStatusIndicatorIcon status={r.status} />
									</div>
									{idx !== playerInfoA.submitResult.execResults.length - 1 && (
										<div>
											<FontAwesomeIcon
												icon={faArrowDown}
												fixedWidth
												className="text-gray-500"
											/>
										</div>
									)}
								</div>
								<div className="grow p-2 overflow-x-scroll">
									<BorderedContainer>
										<div className="font-semibold">{r.label}</div>
										<div>
											<code>
												{r.stdout}
												{r.stderr}
											</code>
										</div>
									</BorderedContainer>
								</div>
							</li>
						))}
					</ul>
				</div>
				<div className="col-span-2 flex flex-col gap-4">
					<div className="flex">
						<div className="grow font-bold text-xl text-center">
							<SubmitStatusLabel status={playerInfoB.submitResult.status} />
						</div>
					</div>
					<ul className="flex flex-col gap-2">
						{playerInfoB.submitResult.execResults.map((r, idx) => (
							<li key={r.testcase_id ?? -1} className="flex gap-2">
								<div className="flex flex-col gap-2 p-2">
									<div className="w-6">
										<ExecStatusIndicatorIcon status={r.status} />
									</div>
									{idx !== playerInfoB.submitResult.execResults.length - 1 && (
										<div>
											<FontAwesomeIcon
												icon={faArrowDown}
												fixedWidth
												className="text-gray-500"
											/>
										</div>
									)}
								</div>
								<div className="grow p-2 overflow-x-scroll">
									<BorderedContainer>
										<div className="font-semibold">{r.label}</div>
										<div>
											<code>
												{r.stdout}
												{r.stderr}
											</code>
										</div>
									</BorderedContainer>
								</div>
							</li>
						))}
					</ul>
				</div>
				<div className="col-span-3">
					<pre className="bg-white resize-none h-full w-full rounded-lg border border-gray-300 p-2">
						<code>{playerInfoB.code}</code>
					</pre>
				</div>
			</div>
		</div>
	);
}
