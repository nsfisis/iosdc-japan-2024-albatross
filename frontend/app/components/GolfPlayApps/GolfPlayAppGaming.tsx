import React, { useRef } from "react";

type Props = {
	problem: string;
	onCodeChange: (code: string) => void;
	onCodeSubmit: (code: string) => void;
	currentScore: number | null;
	lastExecStatus: string | null;
};

export default function GolfPlayAppGaming({
	problem,
	onCodeChange,
	onCodeSubmit,
	currentScore,
	lastExecStatus,
}: Props) {
	const textareaRef = useRef<HTMLTextAreaElement>(null);

	const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
		onCodeChange(e.target.value);
	};

	const handleSubmitButtonClick = () => {
		if (textareaRef.current) {
			onCodeSubmit(textareaRef.current.value);
		}
	};

	return (
		<div className="min-h-screen flex">
			<div className="mx-auto flex min-h-full flex-grow">
				<div className="flex w-1/2 flex-col justify-between p-4">
					<div>
						<div className="mb-2 text-xl font-bold">TODO</div>
						<div className="text-gray-700">{problem}</div>
					</div>
					<div className="mb-4 mt-auto">
						<div className="mb-2">
							<div className="font-semibold text-green-500">
								Score: {currentScore ?? "-"} ({lastExecStatus ?? "-"})
							</div>
						</div>
						<button
							onClick={handleSubmitButtonClick}
							className="focus:shadow-outline rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
						>
							Submit
						</button>
					</div>
				</div>
				<div className="w-1/2 p-4 flex">
					<div className="flex-grow">
						<textarea
							ref={textareaRef}
							onChange={handleTextChange}
							className="h-full w-full rounded-lg border border-gray-300 p-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
						></textarea>
					</div>
				</div>
			</div>
		</div>
	);
}
