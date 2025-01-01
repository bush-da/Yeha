from flask import Blueprint, render_template, redirect, url_for, request, flash, jsonify, session
from models import storage
from models.user import User
from models.post import Post
from models.comment import Comment
from models.tag import Tag
from models.report import Report
from functools import wraps

admin_bp = Blueprint('admin', __name__, url_prefix='/admin')

def admin_required(f):
    """Decorator to enforce admin-only access."""
    @wraps(f)
    def decorated_function(*args, **kwargs):
        user_id = session.get('user_id')  # Retrieve logged-in user ID from session
        if not user_id:  # No user logged in
            flash('Please log in to access this page.', 'warning')
            return redirect(url_for('auth.login', next=request.url))

        user = storage.all(User).get(f"User.{user_id}")
        if not user or not user.is_admin:  # User doesn't exist or isn't an admin
            flash('You are not authorized to access this page.', 'danger')
            return redirect(url_for('home.index'))  # Redirect to home or login page

        return f(*args, **kwargs)  # Allow access if user is admin

    return decorated_function


@admin_bp.route('/', methods=['GET'])
@admin_required
def admin_dashboard():
    """Admin Dashboard Overview Page."""
    total_users = len(storage.all(User))
    total_posts = len(storage.all(Post))
    pending_reports = len([report for report in storage.all(Report).values() if not report.reviewed])
    # Calculate tag popularity
    tags = storage.all(Tag).values()
    tag_counts = [(tag.name, len(tag.posts)) for tag in tags if len(tag.posts) > 0]  # Count posts per tag
    popular_tags = sorted(tag_counts, key=lambda x: x[1], reverse=True)[:5]  # Sort by count (descending)

    return render_template('admin/overview.html',
                           total_users=total_users,
                           total_posts=total_posts,
                           pending_reports=pending_reports,
                           popular_tags=[tag for tag in popular_tags])


@admin_bp.route('/users', methods=['GET'])
@admin_required
def admin_users():
    users = storage.all(User).values()
    return render_template('admin/users.html', users=users)




@admin_bp.route('/user/<user_id>', methods=['GET'])
@admin_required
def user_details(user_id):
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found', 'error')
        return redirect(url_for('admin.admin_users'))
    return render_template('admin/user_details.html', user=user)



@admin_bp.route('/user/<user_id>/edit', methods=['GET', 'POST'])
@admin_required
def edit_user(user_id):
    """Edit user details, including role."""
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found', 'error')
        return redirect(url_for('admin.admin_users'))

    if request.method == 'POST':
        # Update basic fields
        user.username = request.form['name']
        user.email = request.form['email']

        # Update role
        role = request.form['role']
        user.is_admin = True if role == 'admin' else False  # Update based on role selection

        storage.save()
        flash('User updated successfully', 'success')
        return redirect(url_for('admin.admin_users'))

    return render_template('admin/edit_user.html', user=user)


@admin_bp.route('/user/<user_id>/delete', methods=['POST'])
@admin_required
def delete_user(user_id):
    user = storage.all(User).get(f"User.{user_id}")
    if not user:
        flash('User not found', 'error')
    else:
        storage.delete(user)
        storage.save()
        flash('User deleted successfully', 'success')
    return redirect(url_for('admin.admin_users'))

@admin_bp.route('/post/<post_id>/edit', methods=['GET', 'POST'])
@admin_required
def edit_post(post_id):
    post = storage.all(Post).get(f"Post.{post_id}")
    if not post:
        flash('Post not found', 'error')
        return redirect(url_for('admin.admin_posts'))

    if request.method == 'POST':
        post.title = request.form['title']
        post.content = request.form['content']
        post.is_published = bool(int(request.form['is_published']))
        storage.save()
        flash('Post updated successfully', 'success')
        return redirect(url_for('admin.admin_posts'))

    return render_template('admin/edit_post.html', post=post)



@admin_bp.route('/post/<post_id>/delete', methods=['POST'])
@admin_required
def delete_post(post_id):
    post = storage.all(Post).get(f"Post.{post_id}")
    if not post:
        flash('Post not found', 'error')
    else:
        storage.delete(post)
        storage.save()
        flash('Post deleted successfully', 'success')
    return redirect(url_for('admin.admin_posts'))






@admin_bp.route('/reports', methods=['GET'])
@admin_required
def admin_reports():
    """Fetch all reported posts and comments for review."""
    reports = storage.all(Report)
    results = []

    for report in reports.values():
        # Determine if the report targets a post or a comment
        target_type = 'post' if report.post_id else 'comment'
        target_id = report.post_id or report.comment_id

        # Fetch the target and reporter data
        if target_type == 'post':
            target = storage.all(Post).get(f"Post.{target_id}")
        else:
            target = storage.all(Comment).get(f"Comment.{target_id}")

        reporter = storage.all(User).get(f"User.{report.user_id}")



        results.append({
            'post_id': target_id,
            'id': report.id,
            'title': target.title if target_type == 'post' else target.content[:30] + '...',
            'user': {
                'name': reporter.username
            },
            'post': {
                'id': target.id if target_type == 'post' else None
            },
            'reason': report.reason,
            'reviewed': report.reviewed
        })

    return render_template('admin/reports.html', reports=results)



@admin_bp.route('/report/<report_id>/action', methods=['POST'])
@admin_required
def review_flag(report_id):
    """Review flagged content."""
    action = request.form.get('action') or request.form.get('_method')

    # Validate action
    if action not in ['delete', 'ignore']:
        return jsonify({'error': 'Invalid action'}), 400

    # Retrieve flagged content
    report = storage.all(Report).get(f"Report.{report_id}")
    if not report:
        return jsonify({'error': 'Reported content not found'}), 404

    # Perform actions
    if action == 'delete' and report.post_id:
        post = storage.all(Post).get(f"Post.{report.post_id}")
        if post:
            storage.delete(post)
            storage.save()
            report.action_taken = 'deleted'
    else:
        report.action_taken = 'ignored'

    report.reviewed = True
    storage.save()

    return jsonify({'success': 'Action updated successfully'})
