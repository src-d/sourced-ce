from superset import conf, db
from superset.models import core as models
from superset.connectors.connector_registry import ConnectorRegistry


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

    return dbobj


def add_gitbase_tables():
    dbobj = get_or_create_gitbase_db()
    TBL = ConnectorRegistry.sources['table']
    for table in dbobj.all_table_names_in_schema('gitbase'):
        # table_name should match the one in the datasource for fetch_metadata to work
        if db.session.query(TBL).filter_by(table_name=table).first():
            continue
        if db.session.query(TBL).filter_by(table_name='gitbase.' + table).first():
            continue

        # create table with original name and fetch columns
        tbl = TBL(table_name=table)
        tbl.database = dbobj
        db.session.add(tbl)
        db.session.commit()
        tbl.fetch_metadata()

        # rename with prefix and set source
        tbl.table_name = 'gitbase.' + table
        tbl.sql = 'select * from ' + table
        db.session.add(dbobj)
        db.session.commit()


get_or_create_gitbase_db()
add_gitbase_tables()
