interface ChatMessageInterface {
  message: string;
  end?: boolean;
}

const ChatMessage = ({ message, end = false }: ChatMessageInterface) => {
  const justify = end ? "justify-end" : "justify-start";

  return (
    <div className={`flex py-2 ${justify}`}>
      <div className="text-wrap break-words bg-gray-800 px-2 py-1 rounded-lg overflow-hidden">
        {message}
      </div>
    </div>
  );
};

export default ChatMessage;
