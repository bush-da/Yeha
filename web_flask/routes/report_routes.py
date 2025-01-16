import requests
from flask import Blueprint, request, jsonify, session, flash, redirect, url_for

report_bp = Blueprint('report', __name__, url_prefix='/report')

# Base URL for the Yeha API
API_BASE_URL = 'http://localhost:8080/api'


@report_bp.route('/post/<post_id>/report', methods=['POST'])
def report_post(post_id):
    """Report a post as inappropriate."""

    token = session.get('jwt_token')

    # Check if the post exists
    headers = {'Authorization': f'Bearer {token}'}

    # Validate user session
    user_id = session.get('user_id')
    if not user_id:
        return jsonify({'error': 'Unauthorized access'}), 401

    # Handle both JSON and form-encoded data
    if request.content_type == 'application/json':
        data = request.get_json()
    else:
        data = request.form  # Handle form submissions

    reason = data.get('reason')

    # Validate input
    valid_reasons = ['Spam', 'Harassment', 'Misinformation', 'Violence', 'Hate Speech', 'Other']
    if not reason or reason not in valid_reasons:
        return jsonify({'error': 'Invalid reason'}), 400

    # Check if the post exists using the API
    post_response = requests.get(f'{API_BASE_URL}/posts/posts/{post_id}')
    author_id = post_response.json().get('post').get('author').get('ID')
    if post_response.status_code == 404:
        return jsonify({'error': 'Post not found'}), 404

    # Check for duplicate report by the same user
    reports_response = requests.get(f'{API_BASE_URL}/report/reports/', headers=headers)
    print(reports_response.status_code)
    if reports_response.status_code == 200:
        existing_reports = reports_response.json()
        for report in existing_reports:
            if report['post_id'] == post_id:
                return jsonify({'success': 'Post already reported'})

    # Send a POST request to create the report
    report_payload = {
        'user_id': user_id,
        'author_id': author_id,
        'post_id': post_id,
        'reason': reason
    }
    create_report_response = requests.post(f'{API_BASE_URL}/report/reports/', json=report_payload, headers=headers)

    if create_report_response.status_code == 201:
        return jsonify({'success': 'Thank you for reporting the post. Our team will review it shortly'})
    else:
        return jsonify({'error': 'Failed to create report'}), create_report_response.status_code


@report_bp.route('/reports', methods=['GET'])
def get_all_reports():
    """Fetch all reports."""
    response = requests.get(f'{API_BASE_URL}/report/reports/')
    if response.status_code == 200:
        return jsonify(response.json())
    return jsonify({'error': 'Failed to fetch reports'}), response.status_code


@report_bp.route('/reports/<report_id>/action', methods=['POST'])
def take_action_on_report(report_id):
    """Take action on a specific report (admin only)."""
    user_id = session.get('user_id')
    if not user_id or not session.get('is_admin', False):
        return jsonify({'error': 'Unauthorized access'}), 401

    data = request.get_json() or {}
    action = data.get('action')  # Action can be "delete" or "ignore"
    if action not in ['delete', 'ignore']:
        return jsonify({'error': 'Invalid action'}), 400

    # Send the action to the API
    response = requests.post(f'{API_BASE_URL}/report/reports/{report_id}/action', json={'action': action})
    if response.status_code == 200:
        return jsonify({'message': 'Action performed successfully'})
    return jsonify({'error': 'Failed to perform action'}), response.status_code


@report_bp.route('/reports/user/<user_id>', methods=['DELETE'])
def delete_user(user_id):
    """Delete a user (admin only)."""
    admin_id = session.get('user_id')
    if not admin_id or not session.get('is_admin', False):
        return jsonify({'error': 'Unauthorized access'}), 401

    # Send DELETE request to the API
    response = requests.delete(f'{API_BASE_URL}/report/reports/user/{user_id}')
    if response.status_code == 200:
        return jsonify({'message': 'User deleted successfully'})
    return jsonify({'error': 'Failed to delete user'}), response.status_code
