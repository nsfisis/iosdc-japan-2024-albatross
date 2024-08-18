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
		<pre className="h-full w-full p-2 bg-gray-50 rounded-lg border border-gray-300 whitespace-pre-wrap break-words">
			<code dangerouslySetInnerHTML={{ __html: highlighted }} />
		</pre>
	);
}
