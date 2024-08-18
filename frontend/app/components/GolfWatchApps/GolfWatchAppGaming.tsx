import { PlayerInfo } from "../../models/PlayerInfo";
import BorderedContainer from "../BorderedContainer";
import CodeBlock from "../Gaming/CodeBlock";
import ScoreBar from "../Gaming/ScoreBar";
import SubmitResult from "../Gaming/SubmitResult";
import UserIcon from "../UserIcon";

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
	problemTitle,
	problemDescription,
}: Props) {
	const leftTime = (() => {
		const k = gameDurationSeconds + leftTimeSeconds;
		if (k <= 0) {
			return "00:00";
		}
		const m = Math.floor(k / 60);
		const s = k % 60;
		return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
	})();

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className="text-white bg-iosdc-japan grid grid-cols-3 px-4 py-2">
				<div className="font-bold flex justify-between my-auto">
					<div className="flex gap-4">
						{playerInfoA.iconPath && (
							<UserIcon
								iconPath={playerInfoA.iconPath}
								displayName={playerInfoA.displayName!}
								className="w-12 h-12 my-auto"
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
							<UserIcon
								iconPath={playerInfoB.iconPath}
								displayName={playerInfoB.displayName!}
								className="w-12 h-12 my-auto"
							/>
						)}
					</div>
				</div>
			</div>
			<ScoreBar
				scoreA={playerInfoA.score}
				scoreB={playerInfoB.score}
				bgA="bg-orange-400"
				bgB="bg-purple-400"
			/>
			<div className="grow grid grid-cols-3 p-4 gap-4">
				<CodeBlock code={playerInfoA.code ?? ""} language="swift" />
				<div className="flex flex-col gap-4">
					<div className="grid grid-cols-2 gap-4">
						<SubmitResult result={playerInfoA.submitResult} />
						<SubmitResult result={playerInfoB.submitResult} />
					</div>
					<div>
						<div className="mb-2 text-center text-xl font-bold">
							{problemTitle}
						</div>
						<BorderedContainer>{problemDescription}</BorderedContainer>
					</div>
				</div>
				<CodeBlock code={playerInfoB.code ?? ""} language="swift" />
			</div>
		</div>
	);
}
