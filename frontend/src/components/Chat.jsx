import { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Allapi from '../common';
import io from 'socket.io-client';

const Chat = () => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const { user, logout } = useAuth();
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();
  const socketRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    if (!user) {
      navigate('/');
      return;
    }

    // Initialize socket connection
    socketRef.current = io("http://localhost:4000", {
      withCredentials: true,
      transports: ['websocket'],
      reconnectionAttempts: 5,
      forceNew: true,
    });

    socketRef.current.on("connect", () => {
      console.log("Connected to WebSocket server!");
    });

    // Listen for incoming messages
    socketRef.current.on("receiveMessage", (msg) => {
      console.log("Received message:", msg);
      setMessages(prevMessages => [...prevMessages, msg]);
    });

    // Fetch initial messages
    fetchMessages();

    return () => {
      if (socketRef.current) {
        socketRef.current.off("receiveMessage");
        socketRef.current.disconnect();
      }
    };
  }, [user, navigate]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const fetchMessages = async () => {
    try {
      const response = await axios.get(Allapi.messages.url, { withCredentials: true });
      setMessages(response.data);
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to fetch messages');
    }
  };

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !socketRef.current) return;

    const messageData = {
      content: newMessage,
      username: user.username,
      user_id: user.userid,
      timestamp: new Date().toISOString()
    };

    try {
      // Emit the socket event
      socketRef.current.emit("sendMessage", messageData);
      
      // Make the HTTP request
      await axios({
        method: Allapi.sendMessage.method,
        url: Allapi.sendMessage.url,
        data: { content: newMessage, username: user.username },
        withCredentials: true
      });

      setNewMessage('');
    } catch (error) {
      toast.error(error.response?.data?.error || 'Failed to send message');
    }
  };

  const handleLogout = async () => {
    const success = await logout();
    if (success) {
      if (socketRef.current) {
        socketRef.current.disconnect();
      }
      toast.success('Logged out successfully');
      navigate('/');
    } else {
      toast.error('Failed to logout');
    }
  };

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      <div className="bg-white shadow-md p-4 flex justify-between items-center">
        <h1 className="text-xl font-bold">Group Chat</h1>
        <div className="flex items-center gap-4">
          <span className="text-gray-600">Welcome, {user?.username}</span>
          <button
            onClick={handleLogout}
            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
          >
            Logout
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((message, index) => (
          <div
            key={`${message.ID || message.timestamp}-${index}`}
            className={`flex ${
              message.user_id === user?.userid ? 'justify-end' : 'justify-start'
            }`}
          >
            <div
              className={`max-w-[70%] rounded-lg p-3 ${
                message.user_id === user?.userid
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200'
              }`}
            >
              <p className="text-sm font-semibold mb-1">
                {message.user_id === user?.userid ? 'You' : message.username}
              </p>
              <p>{message.content}</p>
              <p className="text-xs mt-1 opacity-75">
                {new Date(message.timestamp).toLocaleTimeString()}
              </p>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form onSubmit={handleSendMessage} className="p-4 bg-white shadow-md">
        <div className="flex gap-2">
          <input
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type your message..."
            className="flex-1 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-200"
          />
          <button
            type="submit"
            className="bg-blue-500 text-white px-6 py-2 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          >
            Send
          </button>
        </div>
      </form>
    </div>
  );
};

export default Chat;