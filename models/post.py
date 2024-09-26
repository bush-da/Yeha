#!/usr/bin/python3
"""Defines the Post class"""
from sqlalchemy import Column, String, Text, ForeignKey
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base
from models.like import Like

from models.post_tag import post_tag

class Post(BaseModel, Base):
    """Represents a blog post"""
    __tablename__ = 'posts'

    title = Column(String(255), nullable=False)
    author_id = Column(String(60), ForeignKey('users.id'), nullable=False)
    category_id = Column(String(60), ForeignKey('categories.id'), nullable=True)

    # Relationships
    category = relationship('Category', back_populates='posts')
    contents = relationship("Content", back_populates="post", cascade="all, delete-orphan")
    author = relationship("User", back_populates="posts")
    comments = relationship("Comment", back_populates="post", cascade="all, delete-orphan")
    likes = relationship("Like", back_populates="post", cascade="all, delete-orphan")
    tags = relationship("Tag", secondary=post_tag, back_populates="posts")

    def count_comment(self):
        """Return the number of comment for the post."""
        return len(self.comments)

    def count_likes(self):
        """Return the number of likes for the post."""
        return len(self.likes)

    def is_liked(self, user, post):
        from models import storage
        likes = storage.all(Like).values()
        for like in likes:
            if like.post_id == post:
                if like.author_id == user:
                    return True
        return False

    def add_tag(self, tag):
        """Add a tag to the post."""
        if tag not in self.tags:
            self.tags.append(tag)
            models.storage.save()

    def remove_tag(self, tag):
        """Remove a tag from the post."""
        if tag in self.tags:
            self.tags.remove(tag)
            models.storage.save()
