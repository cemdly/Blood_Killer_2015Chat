// static/chat.js
const socket = io();

const messagesDiv = document.getElementById('messages');
const input = document.getElementById('messageInput');

socket.on('receive_message', function(data) {
    const msgElement = document.createElement('div');
    msgElement.className = 'message';
    msgElement.textContent = data;
    messagesDiv.appendChild(msgElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
});

function sendMessage() {
    const message = input.value;
    if (message.trim() !== '') {
        socket.emit('send_message', message);
        input.value = '';
    }
}