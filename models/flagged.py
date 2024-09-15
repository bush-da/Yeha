#!/usr/bin/python3
import models
from models.base_model import Base, BaseModel
from models.user import User
from models.comment import Comment
from models.post import Post
from sqlalchemy import Boolean

class FlaggedContent(BaseModel, Base):
    __tablename__ = 'flagged_content'

    user_id = Column(String(60), ForeignKey('users.id'), nullable=False)
    post_id = Column(String(60), ForeignKey('posts.id'), nullable=True)  # Either a post or a comment is flagged
    comment_id = Column(String(60), ForeignKey('comments.id'), nullable=True)
    reason = Column(Text, nullable=False)  # Reason for flagging
    reviewed = Column(Boolean, default=False)  # Flag to see if admin reviewed it
    action_taken = Column(String(64), nullable=True)  # Admin's action ('deleted', 'warned', etc.)

    user = relationship('User', backref='flagged_reports')
    post = relationship('Post', backref='flags', cascade="all, delete-orphan")
    comment = relationship('Comment', backref='flags', cascade="all, delete-orphan")
