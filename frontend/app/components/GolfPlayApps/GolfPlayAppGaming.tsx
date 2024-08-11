import { Link } from "@remix-run/react";
import React, { useRef } from "react";
import SubmitButton from "../../components/SubmitButton";
import BorderedContainer from "../BorderedContainer";

type Props = {
	gameDisplayName: string;
	playerDisplayName: string;
	problemTitle: string;
	problemDescription: string;
	onCodeChange: (code: string) => void;
	onCodeSubmit: (code: string) => void;
	currentScore: number | null;
	lastExecStatus: string | null;
};

export default function GolfPlayAppGaming({
	gameDisplayName,
	playerDisplayName,
	problemTitle,
	problemDescription,
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
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className="text-white bg-iosdc-japan flex flex-row justify-between px-4 py-2">
				<div className="font-bold">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-2xl">03:21</div>
				</div>
				<div>
					<Link to={"/dashboard"} className="font-bold text-xl">
						{playerDisplayName}
					</Link>
				</div>
			</div>
			<div className="grow grid grid-cols-3 divide-x divide-gray-300">
				<div className="p-4">
					<div className="mb-2 text-xl font-bold">{problemTitle}</div>
					<div className="p-2">
						<BorderedContainer>
							<div className="text-gray-700">{problemDescription}</div>
						</BorderedContainer>
					</div>
				</div>
				<div className="p-4">
					<textarea
						ref={textareaRef}
						onChange={handleTextChange}
						className="resize-none h-full w-full rounded-lg border border-gray-300 p-2 focus:outline-none focus:ring-2 focus:ring-gray-400 transition duration-300"
					></textarea>
				</div>
				<div className="p-4">
					<SubmitButton onClick={handleSubmitButtonClick}>提出</SubmitButton>
					<div className="mb-2 mt-auto">
						<div className="font-semibold text-green-500">
							Score: {currentScore ?? "-"} ({lastExecStatus ?? "-"})
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
