from flask import Blueprint, render_template, request, redirect, url_for, session
from models.post import Post
from models.tag import Tag
from models.follower import Follower
from models import storage

home_bp = Blueprint('home', __name__)

@home_bp.route('/', methods=['GET'])
def redirect_home():
    return redirect(url_for('home.index'))

@home_bp.route('/home', methods=['GET'])
def index():
    # To filter using tags
    logged_in = 'user_id' in session
    tag_ids = request.args.get('tags')

    if tag_ids:
        # If tag button is clicked
        # Convert the comma-separated string into a list of tag IDs
        tag_ids = tag_ids.split(',')

        # Filter posts by the selected tags
        posts = []
        all_posts = storage.all(Post).values()
        for post in all_posts:
            post_tag_ids = [str(tag.id) for tag in post.tags]  # Get the tag IDs for each post
            if any(tag_id in post_tag_ids for tag_id in tag_ids):
                posts.append(post)
    else:
        # No tags selected, show all posts
        posts = storage.all(Post).values()
    follower_map = {}
    tags = storage.all(Tag).values()  # Fetch all tags for display
    for post in posts:
        post.contents = sorted(post.contents, key=lambda c: c.paragraph)
        followers = storage.all(Follower).values()
        # Filter followers for the current post's author
        follower_map[post.author.id] = [follower.follower_id for follower in followers if follower.followed_id == post.author.id]

    return render_template('index.html', posts=posts, tags=tags, follower_map=follower_map, logged_in=logged_in)
