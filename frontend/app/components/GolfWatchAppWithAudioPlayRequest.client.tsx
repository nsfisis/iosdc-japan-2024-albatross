import { useState } from "react";
import { AudioController } from "../.client/audio/AudioController";
import GolfWatchApp, { type Props } from "./GolfWatchApp.client";

export default function GolfWatchAppWithAudioPlayRequest({
	game,
	sockToken,
}: Omit<Props, "audioController">) {
	const [audioController, setAudioController] =
		useState<AudioController | null>(null);
	const audioPlayPermitted = audioController !== null;

	if (audioPlayPermitted) {
		return (
			<GolfWatchApp
				game={game}
				sockToken={sockToken}
				audioController={audioController}
			/>
		);
	} else {
		return (
			<div>
				<button
					onClick={async () => {
						const audioController = new AudioController();
						await audioController.loadAll();
						setAudioController(audioController);
					}}
				>
					Enable Audio Play
				</button>
			</div>
		);
	}
}
