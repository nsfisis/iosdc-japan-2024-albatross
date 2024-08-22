import { SoundEffect, getFileUrl } from "./SoundEffect";

export class AudioController {
	audioElements: Record<SoundEffect, HTMLAudioElement | null>;

	constructor() {
		this.audioElements = {
			finish: null,
			winner_1: null,
			winner_2: null,
			good_1: null,
			good_2: null,
			good_3: null,
			good_4: null,
			new_score_1: null,
			new_score_2: null,
			new_score_3: null,
			compile_error_1: null,
			compile_error_2: null,
		};
	}

	loadAll(): Promise<void> {
		return new Promise((resolve) => {
			const files = Object.keys(this.audioElements).map(
				(se) => [se as SoundEffect, getFileUrl(se as SoundEffect)] as const,
			);
			const totalCount = files.length;
			let loadedCount = 0;

			files.forEach(([se, fileUrl]) => {
				const audio = new Audio(fileUrl);

				audio.addEventListener(
					"canplaythrough",
					() => {
						loadedCount++;
						this.audioElements[se] = audio;
						if (loadedCount === totalCount) {
							resolve();
						}
					},
					{ once: true },
				);

				audio.addEventListener("error", () => {
					console.log(`Failed to load audio file: ${fileUrl}`);
					// Ignore the error and continue loading other files.
				});
			});
		});
	}

	async playDummySoundEffect(): Promise<void> {
		const audio = this.audioElements["good_1"];
		if (!audio) {
			return;
		}
		audio.muted = true;
		audio.currentTime = 0;
		await audio.play();
		audio.muted = false;
	}

	async playSoundEffect(soundEffect: SoundEffect): Promise<void> {
		const audio = this.audioElements[soundEffect];
		if (!audio) {
			return;
		}
		audio.currentTime = 0;
		await audio.play();
	}

	async playSoundEffectFinish(): Promise<void> {
		await this.playSoundEffect("finish");
	}

	async playSoundEffectWinner(winner: 1 | 2): Promise<void> {
		await this.playSoundEffect(`winner_${winner}`);
	}

	async playSoundEffectGood(): Promise<void> {
		const variant = Math.floor(Math.random() * 4) + 1;
		if (variant !== 1 && variant !== 2 && variant !== 3 && variant !== 4) {
			return; // unreachable
		}
		return await this.playSoundEffect(`good_${variant}`);
	}

	async playSoundEffectNewScore(): Promise<void> {
		const variant = Math.floor(Math.random() * 3) + 1;
		if (variant !== 1 && variant !== 2 && variant !== 3) {
			return; // unreachable
		}
		return await this.playSoundEffect(`new_score_${variant}`);
	}

	async playSoundEffectCompileError(): Promise<void> {
		const variant = Math.floor(Math.random() * 2) + 1;
		if (variant !== 1 && variant !== 2) {
			return; // unreachable
		}
		return await this.playSoundEffect(`compile_error_${variant}`);
	}
}
