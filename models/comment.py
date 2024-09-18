#!/usr/bin/python3
"""Defines the Comment class"""
from sqlalchemy import Column, Text, ForeignKey, String
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class Comment(BaseModel, Base):
    """Represents a comment on a post"""
    __tablename__ = 'comments'

    content = Column(Text, nullable=False)
    author_id = Column(String(60), ForeignKey('users.id'), nullable=False)
    post_id = Column(String(60), ForeignKey('posts.id'), nullable=False)

    # Relationships
    author = relationship("User", back_populates="comments")
    post = relationship("Post", back_populates="comments")
