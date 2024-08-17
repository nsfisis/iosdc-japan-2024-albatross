type Props = {
	scoreA: number | null;
	scoreB: number | null;
	bgA: string;
	bgB: string;
};

export default function ScoreBar({ scoreA, scoreB, bgA, bgB }: Props) {
	let scoreRatio;
	if (scoreA === null && scoreB === null) {
		scoreRatio = 50;
	} else if (scoreA === null) {
		scoreRatio = 0;
	} else if (scoreB === null) {
		scoreRatio = 100;
	} else {
		scoreRatio = (scoreB / (scoreA + scoreB)) * 100;
	}

	return (
		<div className={`w-full ${bgB}`}>
			<div className={`h-6 ${bgA}`} style={{ width: `${scoreRatio}%` }}></div>
		</div>
	);
}
