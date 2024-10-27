"use client";

import { redirect } from "next/navigation";
import React, { useEffect, useRef, useState } from "react";

const topics = ["anime", "books", "games", "movies", "music"];

const Page = ({ params }: { params: Promise<{ topic: string }> }) => {
  const { topic } = React.use(params);
  if (!topics.includes(topic)) {
    redirect("/");
  }

  const [message, setMessage] = useState<string>("");
  const [messageHistory, setMessageHistory] = useState<string[]>([]);
  const socket = React.useRef<WebSocket | null>(null);
  const formRef = useRef<HTMLFormElement>(null);
  const chatRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    socket.current = new WebSocket("ws://localhost:3100/ws/" + topic);
    socket.current.onmessage = (e) => {
      setMessageHistory((prevMessages) => [...prevMessages, e.data]);
    };

    return () => {
      socket.current?.close();
    };
  }, []);

  useEffect(() => {
    if (!chatRef.current) {
      return;
    }

    chatRef.current.scrollIntoView({ behavior: "smooth" });
  }, [messageHistory]);

  const sendMessage = () => {
    if (!socket.current || !message) {
      return;
    }

    socket.current.send(message);
  };

  const handleMessage = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
  };

  const handleEnter = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key != "Enter") {
      return;
    }

    e.preventDefault();
    formRef.current?.requestSubmit();
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    sendMessage();
    setMessage("");
  };

  return (
    <div className="flex h-screen justify-center p-8">
      <div className="flex flex-col justify-end h-full min-w-[400px] w-[600px] rounded-2xl overflow-hidden bg-gray-900">
        <div className="flex justify-center items-center p-2 mb-auto text-2xl bg-gray-800">
          {topic.charAt(0).toUpperCase() + topic.slice(1)} Chatroom
        </div>
        <div className="overflow-y-auto">
          <ul>
            {messageHistory.map((m, i) => (
              <li key={i}>{m}</li>
            ))}
          </ul>
          <div ref={chatRef} />
        </div>
        <form
          ref={formRef}
          onSubmit={handleSubmit}
          className="h-[10%] flex justify-center items-center p-8 bg-gray-800"
        >
          <textarea
            value={message}
            onChange={handleMessage}
            onKeyDown={handleEnter}
            wrap="soft"
            className="w-full rounded-3xl px-2 overflow-hidden resize-none bg-gray-600"
          />
        </form>
      </div>
    </div>
  );
};

export default Page;
