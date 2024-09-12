#!/usr/bin/python3
"""Defines the BaseModel class for Yeha"""
import models
from uuid import uuid4
from datetime import datetime
from sqlalchemy import Column, String, DateTime
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class BaseModel:
    """A base class for all models in the Yeha project.

    Attributes:
        id (sqlalchemy String): The unique identifier for the model instance.
        created_at (sqlalchemy DateTime): The timestamp when the instance was created.
        updated_at (sqlalchemy DateTime): The timestamp when the instance was last updated.
    """

    id = Column(String(60), primary_key=True, nullable=False)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)

    def __init__(self, *args, **kwargs):
        """Initialize a new instance of BaseModel.

        Args:
            *args: Variable length argument list.
            **kwargs: Arbitrary keyword arguments.
        """

        if not kwargs.get('id'):
            self.id = str(uuid4())
        if not kwargs.get('created_at'):
            self.created_at = self.updated_at = datatime.utcnow()
        else:
            self.created_at = datetime.strptime(kwargs['updated_at'], "%Y-%m-%dT%H:%M:%S.%f")
            self.updated_at = datetime.strptime(kwargs['updated_at'], "%Y-%m-%dT%H:%M:%S.%f")
        for key, value in kwargs.items():
            setattr(self, key, value)

    def save(self):
        """Update the udpated_at time and save the instance"""
        self.updated_at = datetime.utcnow()
        models.storage.new(self)
        models.storage.save()

    def to_dict(self):
        """Convert the instance into a dictionary format"""

        my_dict = self.__dict__.copy()
        my_dict['__class__'] = self.__class__.__name__
        my_dict['created_at'] = self.created_at.isoformat()
        my_dict['updated_at'] = self.updated_at.isoformat()

    def delete(self):
        """Delete the current instance from the storage"""
        models.storage.delete(self)

    def __str__(self):
        """Return string representation of the BaseModel instance"""
        d = self.__dict__.copy()
        return f'[{self.__class__.__name__}] ({self.id}) {d}'
