const roomSelect = document.getElementById('roomSelect');
const joinButton = document.getElementById('joinButton');
const messageContainer = document.getElementById('messageContainer');
const messageInput = document.getElementById('messageInput');
const sendButton = document.getElementById('sendButton');
let socket = null;

joinButton.addEventListener('click', function() {
    // Fermer la connexion WebSocket existante avant d'en créer une nouvelle
    if (socket !== null) {
        socket.close();
    }

    // Créer une nouvelle connexion WebSocket avec la room sélectionnée
    socket = new WebSocket('ws://localhost:8080/ws?id=' + roomSelect.value);

    // Ajouter des écouteurs d'événements à la nouvelle connexion WebSocket
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
});

sendButton.addEventListener('click', function() {
    // Envoyer un message uniquement si la connexion WebSocket est ouverte
    if (socket !== null && socket.readyState === WebSocket.OPEN) {
        const message = {
            roomId: roomSelect.value,
            email: "test@example.com",
            username: "user1",
            message: messageInput.value
        };
        socket.send(JSON.stringify(message));
        messageInput.value = '';
    }
});
