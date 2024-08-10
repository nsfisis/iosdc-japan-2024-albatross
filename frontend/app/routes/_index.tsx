import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";
import "@fortawesome/fontawesome-svg-core/styles.css";
import { ensureUserNotLoggedIn } from "../.server/auth";

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
				<div className="font-bold text-transparent bg-clip-text bg-gradient-to-r from-orange-400 via-pink-500 to-purple-400 flex flex-col gap-y-2">
					<div className="text-3xl">iOSDC Japan 2024</div>
					<div className="text-6xl">
						Swift <wbr />
						Code Battle
					</div>
				</div>
			</div>
			<p className="text-gray-900 max-w-prose bg-white p-4 rounded-xl border-2 border-pink-600 mx-2">
				Swift コードバトルは指示された動作をする Swift
				コードをより短く書けた方が勝ち、という 1 対 1
				の対戦コンテンツです。8/22（木）day0 前夜祭では 8/12
				に実施された予選を勝ち抜いたプレイヤーによるトーナメント形式での Swift
				コードバトルを実施します。ここでは短いコードが正義です！可読性も保守性も放り投げた、イベントならではのコードをお楽しみください！
			</p>
			<div>
				<Link
					to="/login"
					className="text-lg text-white bg-pink-600 px-4 py-2 rounded transition duration-300 hover:bg-pink-500"
				>
					ログイン
				</Link>
			</div>
		</div>
	);
}
