import { PlayerInfo } from "../../models/PlayerInfo";
import PlayerProfile from "../PlayerProfile";

type Props = {
	gameDisplayName: string;
	playerInfo: Omit<PlayerInfo, "code">;
};

export default function GolfPlayAppWaiting({
	gameDisplayName,
	playerInfo,
}: Props) {
	return (
		<div className="min-h-screen bg-gray-100 flex flex-col font-bold text-center">
			<div className="text-white bg-iosdc-japan p-10">
				<div className="text-4xl">{gameDisplayName}</div>
			</div>
			<div className="grow grid mx-auto text-black">
				<PlayerProfile playerInfo={playerInfo} label="You" />
			</div>
		</div>
	);
}
