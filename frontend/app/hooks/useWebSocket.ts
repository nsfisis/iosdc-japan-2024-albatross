import useWebSocketOriginal, { ReadyState } from "react-use-websocket";

export { ReadyState };

// Typed version of useWebSocket() hook.
export default function useWebSocket<ReceiveMessage, SendMessage>(
	url: string,
): {
	sendJsonMessage: (message: SendMessage) => void;
	lastJsonMessage: ReceiveMessage;
	readyState: ReadyState;
} {
	return useWebSocketOriginal(url);
}
