import {
	faBan,
	faCircleCheck,
	faCircleExclamation,
	faRotate,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { VerificationResultStatus } from "../models/VerificationResult";

type Props = {
	status: VerificationResultStatus;
};

export default function ExecStatusIndicatorIcon({ status }: Props) {
	switch (status) {
		case "running":
			return (
				<FontAwesomeIcon
					icon={faRotate}
					spin
					fixedWidth
					className="text-gray-700"
				/>
			);
		case "success":
			return (
				<FontAwesomeIcon
					icon={faCircleCheck}
					fixedWidth
					className="text-green-500"
				/>
			);
		case "canceled":
			return (
				<FontAwesomeIcon icon={faBan} fixedWidth className="text-gray-400" />
			);
		default:
			return (
				<FontAwesomeIcon
					icon={faCircleExclamation}
					fixedWidth
					className="text-red-500"
				/>
			);
	}
}
