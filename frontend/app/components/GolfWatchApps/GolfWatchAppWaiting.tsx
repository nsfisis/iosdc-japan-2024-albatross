import type { PlayerProfile } from "../../types/PlayerProfile";
import PlayerNameAndIcon from "../PlayerNameAndIcon";

type Props = {
	gameDisplayName: string;
	playerProfileA: PlayerProfile;
	playerProfileB: PlayerProfile;
};

export default function GolfWatchAppWaiting({
	gameDisplayName,
	playerProfileA,
	playerProfileB,
}: Props) {
	return (
		<div className="min-h-screen bg-gray-100 flex flex-col font-bold text-center">
			<div className="text-white bg-iosdc-japan p-10">
				<div className="text-4xl">{gameDisplayName}</div>
			</div>
			<div className="grow grid grid-cols-3 gap-10 mx-auto text-black">
				<PlayerNameAndIcon label="Player 1" profile={playerProfileA} />
				<div className="text-8xl my-auto">vs.</div>
				<PlayerNameAndIcon label="Player 2" profile={playerProfileB} />
			</div>
		</div>
	);
}
