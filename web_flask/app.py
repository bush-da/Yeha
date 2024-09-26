from flask import Flask
from models import storage
from web_flask.routes.auth_routes import auth_bp
from web_flask.routes.post_routes import post_bp
from web_flask.routes.comment_routes import comment_bp
from web_flask.routes.like_routes import like_bp
from web_flask.routes.about import about_bp
from web_flask.routes.home import home_bp
from web_flask.routes.profile_routes import profile_bp
from web_flask.routes.follow_routes import follow_bp

def create_app():
    """Create and configure the Flask application."""
    app = Flask(__name__)
    # Configuration
    app.config.from_object('web_flask.config.Config')

    # Register Blueprints
    app.register_blueprint(home_bp)
    app.register_blueprint(about_bp)
    app.register_blueprint(auth_bp)
    app.register_blueprint(post_bp)
    app.register_blueprint(comment_bp)
    app.register_blueprint(like_bp)
    app.register_blueprint(profile_bp)
    app.register_blueprint(follow_bp)

    @app.teardown_appcontext
    def teardown(exception):
        """Close the database session."""
        storage.close()

    return app

if __name__ == "__main__":
    app = create_app()
    app.run(host='0.0.0.0', port=5000, debug=True)
