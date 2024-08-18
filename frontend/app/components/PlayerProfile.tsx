import { PlayerInfo as FullPlayerInfo } from "../models/PlayerInfo";
import UserIcon from "./UserIcon";

type PlayerInfo = Pick<FullPlayerInfo, "displayName" | "iconPath">;

type Props = {
	playerInfo: PlayerInfo;
	label: string;
};

export default function PlayerProfile({ playerInfo, label }: Props) {
	return (
		<div className="flex flex-col gap-6 my-auto">
			<div className="flex flex-col gap-2">
				<div className="text-4xl">{label}</div>
				<div className="text-6xl">{playerInfo.displayName}</div>
			</div>
			{playerInfo.iconPath && (
				<UserIcon
					iconPath={playerInfo.iconPath}
					displayName={playerInfo.displayName!}
					className="w-48 h-48"
				/>
			)}
		</div>
	);
}
