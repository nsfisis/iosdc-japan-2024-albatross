type Props = {
	iconPath: string;
	displayName: string;
	className: string;
};

export default function UserIcon({ iconPath, displayName, className }: Props) {
	return (
		<img
			src={
				process.env.NODE_ENV === "development"
					? `http://localhost:8002/iosdc-japan/2024/code-battle${iconPath}`
					: `/iosdc-japan/2024/code-battle${iconPath}`
			}
			alt={`${displayName} のアイコン`}
			className={`rounded-full border-4 border-white ${className}`}
		/>
	);
}
