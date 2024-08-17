import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { ensureUserNotLoggedIn } from "../.server/auth";
import BorderedContainer from "../components/BorderedContainer";
import NavigateLink from "../components/NavigateLink";

export const meta: MetaFunction = () => [
	{ title: "iOSDC Japan 2024 Albatross.swift" },
];

export async function loader({ request }: LoaderFunctionArgs) {
	await ensureUserNotLoggedIn(request);
	return null;
}

export default function Index() {
	return (
		<div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center gap-y-6">
			<img
				src="/iosdc-japan/2024/code-battle/favicon.svg"
				alt="iOSDC Japan 2024"
				className="w-24 h-24"
			/>
			<div className="text-center">
				<div className="font-bold text-transparent bg-clip-text bg-iosdc-japan flex flex-col gap-y-2">
					<div className="text-3xl">iOSDC Japan 2024</div>
					<div className="text-6xl">Swift Code Battle</div>
				</div>
			</div>
			<div className="mx-2">
				<BorderedContainer>
					<p className="text-gray-900 max-w-prose">
						Swift コードバトルは指示された動作をする Swift
						コードをより短く書けた方が勝ち、という 1 対 1
						の対戦コンテンツです。8/22（木）day0 前夜祭では 8/12
						に実施された予選を勝ち抜いたプレイヤーによるトーナメント形式での
						Swift
						コードバトルを実施します。ここでは短いコードが正義です！可読性も保守性も放り投げた、イベントならではのコードをお楽しみください！
					</p>
				</BorderedContainer>
			</div>
			<div>
				<NavigateLink to="/login">ログイン</NavigateLink>
			</div>
		</div>
	);
}
