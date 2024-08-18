type Props = {
	leftTimeSeconds: number;
};

export default function GolfWatchAppStarting({
	leftTimeSeconds: timeLeft,
}: Props) {
	return (
		<div className="min-h-screen bg-gray-100 flex items-center justify-center">
			<div className="text-center text-black font-black text-10xl animate-ping">
				{timeLeft}
			</div>
		</div>
	);
}
