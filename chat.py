from flask import Flask, render_template, request, redirect, session, url_for
from flask_sqlalchemy import SQLAlchemy
from flask_socketio import SocketIO, emit
import os

app = Flask(__name__)
app.config['SECRET_KEY'] = 'supersecretkey'
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///chat.db'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

db = SQLAlchemy(app)
socketio = SocketIO(app)

# --- Модели ---
class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password = db.Column(db.String(120), nullable=False)

class Message(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    content = db.Column(db.Text, nullable=False)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'))
    user = db.relationship('User', backref='messages')

class Sticker(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    src = db.Column(db.Text, nullable=False)

@app.route('/login', methods=['GET', 'POST'])
def login():
    error_username = None
    error_password = None

    if request.method == 'POST':
        username = request.form['username']
        password = request.form['password']

        user = User.query.filter_by(username=username).first()

        if not user or user.password != password:
            error_username = "Неверное имя пользователя или пароль"
            error_password = "Неверное имя пользователя или пароль"

        if user and user.password == password:
            session['user_id'] = user.id
            return redirect(url_for('index'))

    return render_template('login.html', register=False, invalid_username=error_username, invalid_password=error_password)


@app.route('/register', methods=['GET', 'POST'])
def register():
    error_username = None
    error_password = None

    if request.method == 'POST':
        username = request.form['username']
        password = request.form['password']

        if User.query.filter_by(username=username).first():
            error_username = "Пользователь уже существует"

        if len(password) < 6:
            error_password = "Пароль должен быть не короче 6 символов"

        if not error_username and not error_password:
            new_user = User(username=username, password=password)
            db.session.add(new_user)
            db.session.commit()
            return redirect(url_for('login'))

    return render_template('login.html', register=True, invalid_username=error_username, invalid_password=error_password)

# --- Socket.IO ---
@socketio.on('send_message')
def handle_send_message(data):
    user_id = session.get('user_id')
    if not user_id:
        return

    message = Message(content=data['message'], user_id=user_id)
    db.session.add(message)
    db.session.commit()

    username = User.query.get(user_id).username
    emit('receive_message', {'message': data['message'], 'username': username}, broadcast=True)

# --- Запуск ---
if __name__ == '__main__':
    with app.app_context():
        db.create_all()
    socketio.run(app, host='0.0.0.0', port=5000, debug=True)