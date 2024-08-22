import { useAtomValue } from "jotai";
import {
	codeAAtom,
	codeBAtom,
	gamingLeftTimeSecondsAtom,
	scoreAAtom,
	scoreBAtom,
	submitResultAAtom,
	submitResultBAtom,
} from "../../states/watch";
import type { PlayerProfile } from "../../types/PlayerProfile";
import BorderedContainer from "../BorderedContainer";
import CodeBlock from "../Gaming/CodeBlock";
import ScoreBar from "../Gaming/ScoreBar";
import SubmitResult from "../Gaming/SubmitResult";
import UserIcon from "../UserIcon";

type Props = {
	gameDisplayName: string;
	playerProfileA: PlayerProfile;
	playerProfileB: PlayerProfile;
	problemTitle: string;
	problemDescription: string;
	gameResult: "winA" | "winB" | "draw" | null;
};

export default function GolfWatchAppGaming({
	gameDisplayName,
	playerProfileA,
	playerProfileB,
	problemTitle,
	problemDescription,
	gameResult,
}: Props) {
	const leftTimeSeconds = useAtomValue(gamingLeftTimeSecondsAtom)!;
	const codeA = useAtomValue(codeAAtom);
	const codeB = useAtomValue(codeBAtom);
	const scoreA = useAtomValue(scoreAAtom);
	const scoreB = useAtomValue(scoreBAtom);
	const submitResultA = useAtomValue(submitResultAAtom);
	const submitResultB = useAtomValue(submitResultBAtom);

	const leftTime = (() => {
		const m = Math.floor(leftTimeSeconds / 60);
		const s = leftTimeSeconds % 60;
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
						{playerProfileA.iconPath && (
							<UserIcon
								iconPath={playerProfileA.iconPath}
								displayName={playerProfileA.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
						<div>
							<div className="text-gray-100">Player 1</div>
							<div className="text-2xl">{playerProfileA.displayName}</div>
						</div>
					</div>
					<div className="text-6xl">{scoreA}</div>
				</div>
				<div className="font-bold text-center">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-3xl">
						{gameResult
							? gameResult === "winA"
								? `勝者 ${playerProfileA.displayName}`
								: gameResult === "winB"
									? `勝者 ${playerProfileB.displayName}`
									: "引き分け"
							: leftTime}
					</div>
				</div>
				<div className="font-bold flex justify-between my-auto">
					<div className="text-6xl">{scoreB}</div>
					<div className="flex gap-6 text-end">
						<div>
							<div className="text-gray-100">Player 2</div>
							<div className="text-2xl">{playerProfileB.displayName}</div>
						</div>
						{playerProfileB.iconPath && (
							<UserIcon
								iconPath={playerProfileB.iconPath}
								displayName={playerProfileB.displayName}
								className="w-12 h-12 my-auto"
							/>
						)}
					</div>
				</div>
			</div>
			<ScoreBar
				scoreA={scoreA}
				scoreB={scoreB}
				bgA="bg-orange-400"
				bgB="bg-purple-400"
			/>
			<div className="grow grid grid-cols-3 p-4 gap-4">
				<CodeBlock code={codeA} language="swift" />
				<div className="flex flex-col gap-4">
					<div className="grid grid-cols-2 gap-4">
						<SubmitResult result={submitResultA} />
						<SubmitResult result={submitResultB} />
					</div>
					<div>
						<div className="mb-2 text-center text-xl font-bold">
							{problemTitle}
						</div>
						<BorderedContainer>{problemDescription}</BorderedContainer>
					</div>
				</div>
				<CodeBlock code={codeB} language="swift" />
			</div>
		</div>
	);
}
