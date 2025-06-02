# server.py
from flask import Flask, render_template
from flask_socketio import SocketIO, emit
import socket

app = Flask(__name__)
app.config['SECRET_KEY'] = 'secret!'
socketio = SocketIO(app)

# Получить локальный IP сервера
def get_local_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(('10.255.255.255', 1))
        ip = s.getsockname()[0]
    except Exception:
        ip = '127.0.0.1'
    finally:
        s.close()
    return ip

@app.route('/')
def index():
    return render_template('index.html')

@socketio.on('connect')
def handle_connect():
    print('Клиент подключён')

@socketio.on('disconnect')
def handle_disconnect():
    print('Клиент отключён')

@socketio.on('send_message')
def handle_message(data):
    print(f"Получено сообщение: {data}")
    emit('receive_message', data, broadcast=True)

if __name__ == '__main__':
    local_ip = get_local_ip()
    print(f"Сервер запущен по адресу: http://{local_ip}:5000")
    socketio.run(app, host='0.0.0.0', port=5000)