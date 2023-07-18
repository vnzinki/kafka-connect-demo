import time
from datetime import datetime
from random import randint

import schedule
import sqlalchemy as db
from sqlalchemy import Column, Integer
from sqlalchemy.orm import declarative_base, sessionmaker

engine = db.create_engine("mysql://user:password@127.0.0.1:3306/database")
connection = engine.connect()
session = sessionmaker(autocommit=False, autoflush=False, bind=engine)()


Base = declarative_base()


class UserModel(Base):
    __tablename__ = "user"

    id = Column(Integer, primary_key=True, index=True)
    sensitive_data = Column(Integer, nullable=False)
    time = Column(Integer)


class SessionModel(Base):
    __tablename__ = "session"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer)
    time = Column(Integer)


Base.metadata.create_all(bind=engine)


def userProducer():
    record = UserModel()
    record.time = int(time.time())
    record.sensitive_data = randint(1, 99999)
    session.add(record)
    session.commit()
    print("Added User record: {}".format(datetime.fromtimestamp(record.time)))


def sessionProducer():
    record = SessionModel()
    record.time = int(time.time())
    record.user_id = randint(1, 100)
    session.add(record)
    session.commit()
    print("Added Session record: {}".format(datetime.fromtimestamp(record.time)))


if __name__ == "__main__":
    schedule.every(5).seconds.do(userProducer)
    schedule.every(1).second.do(sessionProducer)

    while True:
        schedule.run_pending()
        time.sleep(1)
