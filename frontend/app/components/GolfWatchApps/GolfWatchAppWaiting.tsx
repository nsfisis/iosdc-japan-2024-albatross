import { PlayerInfo as FullPlayerInfo } from "../../models/PlayerInfo";
import PlayerProfile from "../PlayerProfile";

type PlayerInfo = Pick<FullPlayerInfo, "displayName" | "iconPath">;

type Props = {
	gameDisplayName: string;
	playerInfoA: PlayerInfo;
	playerInfoB: PlayerInfo;
};

export default function GolfWatchAppWaiting({
	gameDisplayName,
	playerInfoA,
	playerInfoB,
}: Props) {
	return (
		<div className="min-h-screen bg-gray-100 flex flex-col font-bold text-center">
			<div className="text-white bg-iosdc-japan p-10">
				<div className="text-4xl">{gameDisplayName}</div>
			</div>
			<div className="grow grid grid-cols-3 gap-10 mx-auto text-black">
				<PlayerProfile playerInfo={playerInfoA} label="Player 1" />
				<div className="text-8xl my-auto">vs.</div>
				<PlayerProfile playerInfo={playerInfoB} label="Player 2" />
			</div>
		</div>
	);
}
