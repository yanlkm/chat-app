const socket = new WebSocket('ws://localhost:8080/ws');

const messageContainer = document.getElementById('messageContainer');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');

sendButton.addEventListener('click', function() {
    const message = {
        email: "test@example.com",
        username: "user1",
        message: messageInput.value
    };

    socket.send(JSON.stringify(message));
    messageInput.value = "";
});

socket.addEventListener('open', function(event) {
    console.log('Connected to WebSocket server');
});

socket.addEventListener('message', function(event) {
    const message = JSON.parse(event.data);
    const messageElement = document.createElement('div');
    messageElement.className = 'message';
    messageElement.textContent = `${message.username}: ${message.message}`;
    messageContainer.appendChild(messageElement);
});

socket.addEventListener('close', function(event) {
    console.log('Connection closed');
});
