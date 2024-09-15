#!/usr/bin/python3
"""Defines the Category class"""
from sqlalchemy import Column, String
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class Category(BaseModel, Base):
    """Represents a category for blog posts"""
    __tablename__ = 'categories'

    name = Column(String(45), nullable=False, unique=True)

    # Relationships
    posts = relationship("Post", back_populates="category", cascade="all, delete-orphan")
