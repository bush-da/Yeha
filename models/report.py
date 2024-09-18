#!/usr/bin/python3
"""Defines the Report class for flagging posts"""
from sqlalchemy import Column, Text, ForeignKey, String
from sqlalchemy.orm import relationship
from models.base_model import BaseModel, Base

class Report(BaseModel, Base):
    """Represents a report for flagging inappropriate posts or comments"""
    __tablename__ = 'reports'

    user_id = Column(String(60), ForeignKey('users.id'), nullable=False)
    post_id = Column(String(60), ForeignKey('posts.id'), nullable=True)
    comment_id = Column(String(60), ForeignKey('comments.id'), nullable=True)
    reason = Column(Text, nullable=False)

    # Relationships
    user = relationship("User")
    post = relationship("Post")
    comment = relationship("Comment")
