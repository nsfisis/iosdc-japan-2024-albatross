import type { SubmitResultStatus } from "../models/SubmitResult";

type Props = {
	status: SubmitResultStatus;
};

export default function SubmitStatusLabel({ status }: Props) {
	switch (status) {
		case "waiting_submission":
			return "提出待ち";
		case "running":
			return "実行中...";
		case "success":
			return "成功";
		case "wrong_answer":
			return "テスト失敗";
		case "timeout":
			return "時間切れ";
		case "compile_error":
			return "コンパイルエラー";
		case "runtime_error":
			return "実行時エラー";
		case "internal_error":
			return "！内部エラー！";
	}
}
