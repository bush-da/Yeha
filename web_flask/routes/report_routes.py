from flask import Blueprint, request, jsonify, session, flash, redirect, url_for
from models import storage
from models.report import Report
from models.post import Post

report_bp = Blueprint('report', __name__, url_prefix='/report')


@report_bp.route('/post/<post_id>/report', methods=['POST'])
def report_post(post_id):
    """Report a post as inappropriate."""
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

    # Check if the post exists
    post = storage.all(Post).get(f"Post.{post_id}")
    if not post:
        return jsonify({'error': 'Post not found'}), 404

    # Check for duplicate report by the same user
    existing_reports = storage.all(Report).values()
    for report in existing_reports:
        if report.user_id == user_id and report.post_id == post_id:
            flash('post already reported.', 'info')
            return redirect(url_for('home.index'))

    # Create and save the report
    new_report = Report(user_id=user_id, post_id=post_id, reason=reason)
    storage.new(new_report)
    storage.save()
    flash('Thank you for reporting the post. Our team will review it shortly.', 'success')
    return redirect(url_for('home.index'))
