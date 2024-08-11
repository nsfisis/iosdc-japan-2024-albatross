import type { SubmitResultStatus } from "../models/SubmitResult";

type Props = {
	status: SubmitResultStatus;
};

export default function SubmitStatusLabel({ status }: Props) {
	switch (status) {
		case "running":
			return <span>Running...</span>;
		case "success":
			return <span>Accepted</span>;
		case "wrong_answer":
			return <span>Wrong Answer</span>;
		case "timeout":
			return <span>Time Limit Exceeded</span>;
		case "compile_error":
			return <span>Compile Error</span>;
		case "runtime_error":
			return <span>Runtime Error</span>;
		case "internal_error":
			return <span>Internal Error</span>;
	}
}
