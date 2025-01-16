from flask import Blueprint, render_template, redirect, url_for, request, flash, jsonify, session
from functools import wraps
import requests

admin_bp = Blueprint('admin', __name__, url_prefix='/admin')

API_BASE_URL = 'http://localhost:8080/api'

def admin_required(f):
    """Decorator to enforce admin-only access."""
    @wraps(f)
    def decorated_function(*args, **kwargs):
        token = session.get('jwt_token')
        headers = {'Authorization': f'Bearer {token}'}

        user_id = session.get('user_id')  # Retrieve logged-in user ID from session
        if not user_id:  # No user logged in
            flash('Please log in to access this page.', 'warning')
            return redirect(url_for('auth.login', next=request.url))

        user = requests.get(f"{API_BASE_URL}/users/{user_id}", headers=headers).json().get('user')
        if not user or not user.get('IsAdmin'):  # User doesn't exist or isn't an admin
            flash('You are not authorized to access this page.', 'danger')
            return redirect(url_for('home.index'))  # Redirect to home or login page

        return f(*args, **kwargs)  # Allow access if user is admin

    return decorated_function


@admin_bp.route('/', methods=['GET'])
@admin_required
def admin_dashboard():
    """Admin Dashboard Overview Page."""
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    user_response = requests.get(f'{API_BASE_URL}/users', headers=headers)
    users = user_response.json()
    total_users = len(users)

    post_response = requests.get(f'{API_BASE_URL}/posts', headers=headers)
    posts = post_response.json()
    total_posts = len(posts)

    report_response = requests.get(f'{API_BASE_URL}/report/reports', headers=headers)
    pending_reports = len([report for report in report_response.json() if not report.get('reviewed')])
    # Calculate tag popularity
    tags_response = requests.get(f'{API_BASE_URL}/posts/tags', headers=headers)

    tags = tags_response.json()

    tag_counts = [(tag.get('name'), tag.get('post_count')) for tag in tags.get('tags') if tag.get('post_count') > 0]  # Count posts per tag
    popular_tags = sorted(tag_counts, key=lambda x: x[1], reverse=True)[:5]  # Sort by count (descending)

    return render_template('admin/overview.html',
                           total_users=total_users,
                           total_posts=total_posts,
                           pending_reports=pending_reports,
                           popular_tags=[tag for tag in popular_tags])


@admin_bp.route('/users', methods=['GET'])
@admin_required
def admin_users():
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    user_response = requests.get(f'{API_BASE_URL}/users', headers=headers)
    users = user_response.json()

    return render_template('admin/users.html', users=users)



@admin_bp.route('/user/<user_id>', methods=['GET'])
@admin_required
def user_details(user_id):
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    user_response = requests.get(f'{API_BASE_URL}/users/{user_id}', headers=headers)
    user = user_response.json().get('user')

    if not user:
        flash('User not found', 'error')
        return redirect(url_for('admin.admin_users'))
    return render_template('admin/user_details.html', user=user)



@admin_bp.route('/user/<user_id>/edit', methods=['GET', 'POST'])
@admin_required
def edit_user(user_id):
    """Edit user details, including role."""
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    user_response = requests.get(f'{API_BASE_URL}/users/{user_id}', headers=headers)
    user = user_response.json().get('user')

    if not user:
        flash('User not found', 'error')
        return redirect(url_for('admin.admin_users'))

    if request.method == 'POST':
        # Update basic fields
        role = request.form.get('role')
        is_admin = True if role == 'admin' else False  # Update based on role selection
        url = f"{API_BASE_URL}/users/{user_id}"


        payload = {
        "username": request.form.get('name'),
        "email": request.form.get('email'),
        "is_admin": is_admin,
        }

        response = requests.put(url, json=payload, headers=headers)

        flash('User updated successfully', 'success')
        return redirect(url_for('admin.admin_users'))

    return render_template('admin/edit_user.html', user=user)


@admin_bp.route('/user/<user_id>/delete', methods=['POST'])
@admin_required
def delete_user(user_id):
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    user_response = requests.get(f'{API_BASE_URL}/users/{user_id}', headers=headers)
    user = user_response.json().get('user')
    if not user:
        flash('User not found', 'error')
    else:
        user_id = user.get('ID')
        response = requests.delete(f"{API_BASE_URL}/users/{user_id}", headers=headers)
        if response.status_code == 200:
            return jsonify({'success': 'user deleted successfully'})
        else:
            return jsonify({'error': 'Failed to delete user'}), response.status_code
    return redirect(url_for('admin.admin_users'))


@admin_bp.route('/reports', methods=['GET'])
@admin_required
def admin_reports():
    """Fetch all reported posts and comments for review."""
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    # Make an API call to get all reports
    response = requests.get(f'{API_BASE_URL}/report/reports')
    if response.status_code != 200:
        return render_template('admin/reports.html', error="Failed to fetch reports.")

    reports_data = response.json()
    results = []

    for report in reports_data:

        # Extract target data from the report
        target_type = 'post' if report['post_id'] else 'comment'
        target = report['post'] if target_type == 'post' else None  # Using post data directly from the report


        # Fetch the reporter data
        user_id = report['user_id']
        user_response = requests.get(f'{API_BASE_URL}/users/{user_id}', headers=headers)

        if user_response.status_code != 200:
            continue  # Skip this report if the user cannot be fetched
        reporter = user_response.json().get('user')

        # Prepare the report data for the template
        results.append({
            'post_id': report['post_id'],
            'id': report['id'],
            'title': target['title'] if target else 'N/A',  # For posts
            'user': {
                'name': reporter['Username']
            },
            'post': {
                'id': target['id'] if target else None
            },
            'reason': report['reason'],
            'reviewed': report['reviewed']
        })

    return render_template('admin/reports.html', reports=results)

@admin_bp.route('/report/<report_id>/action', methods=['POST'])
@admin_required
def review_flag(report_id):
    """Review flagged content."""
    token = session.get('jwt_token')
    headers = {'Authorization': f'Bearer {token}'}

    action = request.form.get('action') or request.form.get('_method')

    # Validate action
    if action not in ['delete', 'ignore']:
        return jsonify({'error': 'Invalid action'}), 400

    # Check if the report exists via the API
    report_response = requests.get(f'{API_BASE_URL}/report/reports/', headers=headers)
    if report_response.status_code == 200:
        reports = [report for report in report_response.json()]
        report = next((r for r in reports if r['id'] == report_id), None)
        if not report:
            return jsonify({'error': 'Reported content not found'}), 404
    else:
        return jsonify({'error': 'Failed to fetch reports'}), report_response.status_code

    # Send action to the API
    action_response = requests.post(
        f'{API_BASE_URL}/report/reports/{report_id}/action',
        json={'action': action},
        headers=headers
    )

    if action_response.status_code == 200:
        return jsonify({'success': 'Action updated successfully'})
    else:
        return jsonify({'error': 'Failed to update action'}), action_response.status_code
