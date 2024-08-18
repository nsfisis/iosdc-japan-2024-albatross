import { PlayerInfo } from "../../models/PlayerInfo";
import PlayerProfile from "../PlayerProfile";

type Props = {
	gameDisplayName: string;
	playerInfo: Omit<PlayerInfo, "code">;
};

export default function GolfPlayAppWaiting({
	gameDisplayName,
	playerInfo,
}: Props) {
	return (
		<>
			<div className="min-h-screen bg-gray-100 flex flex-col font-bold text-center">
				<div className="text-white bg-iosdc-japan p-10">
					<div className="text-4xl">{gameDisplayName}</div>
				</div>
				<div className="grow grid mx-auto text-black">
					<PlayerProfile playerInfo={playerInfo} label="You" />
				</div>
			</div>
			<style>
				{`
        @keyframes changeHeight {
          0% { height: 20%; }
          50% { height: 100%; }
          100% { height: 20%; }
        }
      `}
			</style>
			<div
				style={{
					position: "fixed",
					bottom: 0,
					width: "100%",
					display: "flex",
					justifyContent: "center",
					alignItems: "flex-end",
					height: "100px",
					margin: "0 2px",
				}}
			>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "2.0s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.9s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.8s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.7s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.6s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.5s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.4s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.3s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.2s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.1s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.0s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.9s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.8s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.7s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.6s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.5s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.4s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.3s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.2s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.1s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.5s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.4s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.3s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.2s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.1s",
					}}
				></div>

				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.1s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.2s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.3s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.4s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.5s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.1s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.2s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.3s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.4s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.5s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.6s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.7s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.8s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "0.9s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.0s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.1s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.2s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.3s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.4s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.5s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.6s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.7s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.8s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "1.9s",
					}}
				></div>
				<div
					style={{
						width: "2%",
						margin: "0 2px",
						background:
							"linear-gradient(345deg, rgb(230, 36, 136) 0%, rgb(240, 184, 106) 100%)",
						display: "inline-block",
						animation: "changeHeight 1s infinite ease-in-out",
						animationDelay: "2.0s",
					}}
				></div>
			</div>
		</>
	);
}
