import { Link } from "@remix-run/react";
import React, { useRef } from "react";
import SubmitButton from "../../components/SubmitButton";
import type { PlayerInfo } from "../../models/PlayerInfo";
import BorderedContainer from "../BorderedContainer";
import SubmitResult from "../Gaming/SubmitResult";

type Props = {
	gameDisplayName: string;
	gameDurationSeconds: number;
	leftTimeSeconds: number;
	playerInfo: Omit<PlayerInfo, "code">;
	problemTitle: string;
	problemDescription: string;
	onCodeChange: (code: string) => void;
	onCodeSubmit: (code: string) => void;
};

export default function GolfPlayAppGaming({
	gameDisplayName,
	gameDurationSeconds,
	leftTimeSeconds,
	playerInfo,
	problemTitle,
	problemDescription,
	onCodeChange,
	onCodeSubmit,
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

	const leftTime = (() => {
		const k = gameDurationSeconds + leftTimeSeconds;
		const m = Math.floor(k / 60);
		const s = k % 60;
		return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
	})();

	return (
		<div className="min-h-screen bg-gray-100 flex flex-col">
			<div className="text-white bg-iosdc-japan flex flex-row justify-between px-4 py-2">
				<div className="font-bold">
					<div className="text-gray-100">{gameDisplayName}</div>
					<div className="text-2xl">{leftTime}</div>
				</div>
				<Link to={"/dashboard"}>
					<div className="flex gap-4 my-auto font-bold">
						<div className="text-6xl">{playerInfo.score}</div>
						<div className="text-end">
							<div className="text-gray-100">Player 1</div>
							<div className="text-2xl">{playerInfo.displayName}</div>
						</div>
						{playerInfo.iconPath && (
							<img
								src={
									process.env.NODE_ENV === "development"
										? `http://localhost:8002/iosdc-japan/2024/code-battle${playerInfo.iconPath}`
										: `/iosdc-japan/2024/code-battle${playerInfo.iconPath}`
								}
								alt={`${playerInfo.displayName} のアイコン`}
								className="w-12 h-12 rounded-full my-auto border-4 border-white"
							/>
						)}
					</div>
				</Link>
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
					<SubmitResult
						result={playerInfo.submitResult}
						submitButton={
							<SubmitButton onClick={handleSubmitButtonClick}>
								提出
							</SubmitButton>
						}
					/>
				</div>
			</div>
		</div>
	);
}
