#!/usr/bin/python3
"""Defines the Post class"""
from sqlalchemy import Column, String, Text, ForeignKey
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base
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

    def add_like(self, user):
        """Add a like to the post by a user."""
        from models.like import Like
        existing_like = next((like for like in self.likes if like.author_id == user.id), None)
        if not existing_like:
            like = Like(author_id=user.id, post_id=self.id)
            models.storage.new(like)
            models.storage.save()

    def remove_like(self, user):
        """Remove a like from the post by a user."""
        from models.like import Like
        like = models.storage.get(Like, (user.id, self.id))
        if like:
            models.storage.delete(like)
            models.storage.save()

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
