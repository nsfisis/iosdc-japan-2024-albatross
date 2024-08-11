import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement>;

export default function SubmitButton(props: ButtonProps) {
	return (
		<button
			{...props}
			className="text-lg text-white bg-pink-600 px-4 py-2 rounded transition duration-300 hover:bg-pink-500 focus:ring focus:ring-pink-400 focus:outline-none"
		/>
	);
}
