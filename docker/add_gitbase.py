from superset import conf, db
from superset.models import core as models


def get_or_create_gitbase_db():
    database_name = 'gitbase'

    dbobj = db.session.query(models.Database).filter_by(
        database_name=database_name).first()
    if not dbobj:
        dbobj = models.Database(
            database_name=database_name,
            expose_in_sqllab=True,
            allow_run_async=True,
            allow_dml=True)
    dbobj.set_sqlalchemy_uri(conf.get('GITBASE_DATABASE_URI'))
    db.session.add(dbobj)
    db.session.commit()


get_or_create_gitbase_db()
