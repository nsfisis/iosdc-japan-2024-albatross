type Props = {
	timeLeft: number;
};

export default function GolfPlayAppStarting({ timeLeft }: Props) {
	return (
		<div className="min-h-screen bg-gray-100 flex items-center justify-center">
			<div className="text-center">
				<h1 className="text-4xl font-bold text-black-600 mb-4">
					Starting... ({timeLeft} s)
				</h1>
			</div>
		</div>
	);
}
