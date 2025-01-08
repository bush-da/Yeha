#!/usr/bin/python3
"""Defines the Follower class for user following system"""
from sqlalchemy import Column, String, ForeignKey
from models.base_model import BaseModel, Base
from sqlalchemy.orm import relationship

class Follower(Base):
    """Represents a user following another user"""
    __tablename__ = 'followers'

    follower_id = Column(
        String(60),
        ForeignKey('users.id', ondelete='CASCADE'),
        primary_key=True
    )
    followed_id = Column(
        String(60),
        ForeignKey('users.id', ondelete='CASCADE'),
        primary_key=True
    )
