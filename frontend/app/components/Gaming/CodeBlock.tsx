import Prism, { highlight, languages } from "prismjs";
import "prismjs/components/prism-swift";
import "prismjs/themes/prism.min.css";

Prism.manual = true;

type Props = {
	code: string;
	language: string;
};

export default function CodeBlock({ code, language }: Props) {
	const highlighted = highlight(code, languages[language]!, language);

	return (
		<pre className="bg-white resize-none h-full w-full rounded-lg border border-gray-300 p-2">
			<code dangerouslySetInnerHTML={{ __html: highlighted }} />
		</pre>
	);
}
