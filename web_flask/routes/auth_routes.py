from flask import Blueprint, request, render_template, redirect, url_for, flash
from flask import session
from werkzeug.security import generate_password_hash, check_password_hash
from models import storage
from models.user import User

auth_bp = Blueprint('auth', __name__)

@auth_bp.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        username = request.form.get('username')
        email = request.form.get('email')
        password = request.form.get('password')
        gender = request.form.get('gender')

        if not username or not email or not password:
            flash('All fields are required!')
            return redirect(url_for('auth.register'))

        existing_user = storage.all(User).get(f"User.{email}")
        if existing_user:
            flash('User already exists!')
            return redirect(url_for('auth.register'))

        hashed_password = generate_password_hash(password)
        new_user = User(username=username, email=email, password=hashed_password, bio='', )
        storage.new(new_user)
        storage.save()

        flash('Registration successful!')
        return redirect(url_for('auth.login'))

    return render_template('register.html')

@auth_bp.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        email = request.form.get('email')
        password = request.form.get('password')

        users = storage.all(User)
        user = None
        for u in users.values():
            if u.email == email:
                user = u
                break

        #user = storage.all(User).get(f"User.{email}")
        if not user or not check_password_hash(user.password, password):
            flash('Invalid credentials!')
            return redirect(url_for('auth.login'))

        # Assuming you use session management to track logged-in users
        session['user_id'] = user.id
        flash('Login successful!')
        return redirect(url_for('home.index'))

    return render_template('login.html')

@auth_bp.route('/logout')
def logout():
    session.pop('user_id', None)
    flash('You have been logged out.')
    return redirect(url_for('home.index'))
