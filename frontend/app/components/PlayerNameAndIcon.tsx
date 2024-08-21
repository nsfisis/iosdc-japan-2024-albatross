import { PlayerProfile } from "../types/PlayerProfile";
import UserIcon from "./UserIcon";

type Props = {
	label: string;
	profile: PlayerProfile;
};

export default function PlayerNameAndIcon({ label, profile }: Props) {
	return (
		<div className="flex flex-col gap-6 my-auto">
			<div className="flex flex-col gap-2">
				<div className="text-4xl">{label}</div>
				<div className="text-6xl">{profile.displayName}</div>
			</div>
			{profile.iconPath && (
				<UserIcon
					iconPath={profile.iconPath}
					displayName={profile.displayName}
					className="w-48 h-48"
				/>
			)}
		</div>
	);
}
