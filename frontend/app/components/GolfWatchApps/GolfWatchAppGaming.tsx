import type { PlayerInfo } from "../../types/PlayerInfo";
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
	gameResult: "winA" | "winB" | "draw" | null;
};

export default function GolfWatchAppGaming({
	gameDisplayName,
	gameDurationSeconds,
	leftTimeSeconds,
	playerInfoA,
	playerInfoB,
	problemTitle,
	problemDescription,
	gameResult,
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

	const topBg = gameResult
		? gameResult === "winA"
			? "bg-orange-400"
			: gameResult === "winB"
				? "bg-purple-400"
				: "bg-pink-500"
		: "bg-iosdc-japan";

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className={`text-white ${topBg} grid grid-cols-3 px-4 py-2`}>
				<div className="font-bold flex justify-between my-auto">
					<div className="flex gap-6">
						{playerInfoA.profile.iconPath && (
							<UserIcon
								iconPath={playerInfoA.profile.iconPath}
								displayName={playerInfoA.profile.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
						<div>
							<div className="text-gray-100">Player 1</div>
							<div className="text-2xl">{playerInfoA.profile.displayName}</div>
						</div>
					</div>
					<div className="text-6xl">{playerInfoA.state.score}</div>
				</div>
				<div className="font-bold text-center">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-3xl">
						{gameResult
							? gameResult === "winA"
								? `勝者 ${playerInfoA.profile.displayName}`
								: gameResult === "winB"
									? `勝者 ${playerInfoB.profile.displayName}`
									: "引き分け"
							: leftTime}
					</div>
				</div>
				<div className="font-bold flex justify-between my-auto">
					<div className="text-6xl">{playerInfoB.state.score}</div>
					<div className="flex gap-6 text-end">
						<div>
							<div className="text-gray-100">Player 2</div>
							<div className="text-2xl">{playerInfoB.profile.displayName}</div>
						</div>
						{playerInfoB.profile.iconPath && (
							<UserIcon
								iconPath={playerInfoB.profile.iconPath}
								displayName={playerInfoB.profile.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
					</div>
				</div>
			</div>
			<ScoreBar
				scoreA={playerInfoA.state.score}
				scoreB={playerInfoB.state.score}
				bgA="bg-orange-400"
				bgB="bg-purple-400"
			/>
			<div className="grow grid grid-cols-3 p-4 gap-4">
				<CodeBlock code={playerInfoA.state.code ?? ""} language="swift" />
				<div className="flex flex-col gap-4">
					<div className="grid grid-cols-2 gap-4">
						<SubmitResult result={playerInfoA.state.submitResult} />
						<SubmitResult result={playerInfoB.state.submitResult} />
					</div>
					<div>
						<div className="mb-2 text-center text-xl font-bold">
							{problemTitle}
						</div>
						<BorderedContainer>{problemDescription}</BorderedContainer>
					</div>
				</div>
				<CodeBlock code={playerInfoB.state.code ?? ""} language="swift" />
			</div>
		</div>
	);
}
