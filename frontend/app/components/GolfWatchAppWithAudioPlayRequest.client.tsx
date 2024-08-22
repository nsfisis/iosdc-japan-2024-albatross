import { useAtom } from "jotai";
import { AudioController } from "../.client/audio/AudioController";
import { audioControllerAtom } from "../states/watch";
import GolfWatchApp, { type Props } from "./GolfWatchApp.client";
import SubmitButton from "./SubmitButton";

export default function GolfWatchAppWithAudioPlayRequest({
	game,
	sockToken,
}: Omit<Props, "audioController">) {
	const [audioController, setAudioController] = useAtom(audioControllerAtom);
	const audioPlayPermitted = audioController !== null;

	if (audioPlayPermitted) {
		return <GolfWatchApp game={game} sockToken={sockToken} />;
	} else {
		return (
			<div className="min-h-screen bg-gray-100 flex items-center justify-center">
				<div className="text-center">
					<SubmitButton
						onClick={async () => {
							const audioController = new AudioController();
							await audioController.loadAll();
							await audioController.playDummySoundEffect();
							setAudioController(audioController);
						}}
					>
						開始
					</SubmitButton>
				</div>
			</div>
		);
	}
}
