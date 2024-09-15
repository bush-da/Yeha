from os import getenv
from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session, sessionmaker
from models.base_model import Base
from models.user import User
from models.post import Post
from models.comment import Comment
from models.category import Category
from models.content import Content
from models.like import Like
from models.tag import Tag
from models.report import Report
from models.flagged import FlaggedContent
from models.follower import Follower
from models.post_tag import PostTag

class DBStorage:
    __engine = None
    __session = None

    def __init__(self):
        """Initialize the database storage."""
        self.__engine = create_engine(
            "mysql+mysqldb://{}:{}@{}/{}".
            format(getenv("YEHA_MYSQL_USER"),
                   getenv("YEHA_MYSQL_PWD"),
                   getenv("YEHA_MYSQL_HOST"),
                   getenv("YEHA_MYSQL_DB")),
            pool_pre_ping=True
        )
        if getenv("YEHA_ENV") == "test":
            Base.metadata.drop_all(self.__engine)

    def all(self, cls=None):
        """Query all objects of a specific class, or all objects if no class is specified."""
        if cls is None:
            # Query all models
            objs = []
            for model in [User, Post, Comment, Category, Content, Like, Tag, Report, FlaggedContent, Follower, PostTag]:
                objs.extend(self.__session.query(model).all())
        else:
            if isinstance(cls, str):
                cls = self.get_model(cls)
            objs = self.__session.query(cls).all()
        return {"{}.{}".format(type(o).__name__, o.id): o for o in objs}

    def new(self, obj):
        """Add a new object to the session."""
        self.__session.add(obj)

    def save(self):
        """Commit the current session."""
        self.__session.commit()

    def delete(self, obj=None):
        """Delete an object from the session."""
        if obj is not None:
            self.__session.delete(obj)

    def reload(self):
        """Reload the database session."""
        Base.metadata.create_all(self.__engine)
        session_factory = sessionmaker(bind=self.__engine, expire_on_commit=False)
        Session = scoped_session(session_factory)
        self.__session = Session()

    def close(self):
        """Close the current session."""
        self.__session.remove()

    def get_model(self, name):
        """Get the model class from its name."""
        model_map = {
            "User": User,
            "Post": Post,
            "Comment": Comment,
            "Category": Category,
            "Content": Content,
            "Like": Like,
            "Tag": Tag,
            "Report": Report,
            "FlaggedContent": FlaggedContent,
            "Follower": Follower,
            "PostTag": PostTag
        }
        return model_map.get(name)
