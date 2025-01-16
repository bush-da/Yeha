from flask import Blueprint, request, render_template, redirect, url_for, flash
from flask import session
from werkzeug.security import generate_password_hash, check_password_hash
import requests

auth_bp = Blueprint('auth', __name__)

import requests

@auth_bp.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        username = request.form.get('username')
        email = request.form.get('email')
        password = request.form.get('password')
        conf_password = request.form.get('confirm_password')
        gender = request.form.get('gender')

        if not username or not email or not password:
            flash('All fields are required!')
            return redirect(url_for('auth.register'))
        if conf_password != password:
            flash('Passwords do not match!')
            return redirect(url_for('auth.register'))

        # Check if the user already exists by calling the Go API
        existing_user = check_if_user_exists(email, username)
        if existing_user:
            flash('User already exists!')
            return redirect(url_for('auth.register'))

        # Prepare data for registration
        registration_data = {
            'username': username,
            'email': email,
            'password': password,
            'gender': gender
        }

        # Send POST request to the Go API to create the user
        response = requests.post('http://localhost:8080/api/users/register', json=registration_data)

        if response.status_code == 201:
            flash('Registration successful!')
            return redirect(url_for('auth.login'))
        else:
            # If registration fails on the Go API
            error_message = response.json().get('error', 'Registration failed.')
            flash(error_message)
            return redirect(url_for('auth.register'))

    return render_template('register.html')

def check_if_user_exists(email, username):
    """
    Helper function to check if the user already exists by checking the Go API.
    """
    response = requests.get(f'http://localhost:8080/api/users/check', params={'email': email, 'username': username})
    if response.status_code == 200:
        data = response.json()
        return data.get('exists', False)
    return False

@auth_bp.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        email = request.form.get('email')
        password = request.form.get('password')

        # Make request to Go API for login
        login_data = {
            'email': email,
            'password': password
        }
        response = requests.post('http://localhost:8080/api/users/login', json=login_data)

        if response.status_code != 200:
            flash('Invalid credentials!')
            return redirect(url_for('auth.login'))

        # If login successful, get token and user info
        data = response.json()
        token = data.get('token')
        user_info = data.get('user')

        # Save token and user info in session
        session['user_id'] = user_info['id']
        session['is_admin'] = user_info['isAdmin']
        session['jwt_token'] = token  # Save the JWT token

        flash('Login successful!')
        if user_info['isAdmin']:
            return redirect(url_for('admin.admin_dashboard'))
        return redirect(url_for('home.index'))

    return render_template('login.html')

@auth_bp.route('/logout')
def logout():
    session.pop('user_id', None)
    session.pop('is_admin', None)
    session.pop('jwt_token', None)
    flash('You have been logged out.')
    return redirect(url_for('home.index'))
