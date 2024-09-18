#!/usr/bin/python3
"""Defines the Tag class"""
from sqlalchemy import Column, String
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base
from models.post_tag import post_tag

class Tag(BaseModel, Base):
    """Represents a tag for posts (e.g., 'Python', 'Web Development')"""
    __tablename__ = 'tags'

    name = Column(String(45), nullable=False, unique=True)

    # Relationships
    posts = relationship("Post", secondary=post_tag, back_populates="tags")
