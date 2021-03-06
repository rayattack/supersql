from supersql.errors import (
    ArgumentError,
    MissingArgumentError,
    MissingCommandError,
    DatabaseError
)
from supersql.connection import connection

from .database import Database
from .table import Table

_pg = "postgres"
_ox = "oracle"
_my = "mysql"
_ms = "mssql"

SUPPORTED_VENDORS = {
    "postgres": _pg,
    "postgresql": _pg,
    "oracle": _ox,
    "oracledb": _ox,
    "mariadb": _my,
    "mysql": _my,
    "mssql": _ms,
    "sqlserver": _ms,
    "sqlite": "sqlite"
}



# The underscore denotes presence of space before chars for cognitive reasons
_AND = " AND"
DELETE = "DELETE"
SELECT = "SELECT"
_FROM_ = " FROM "
WHERE = "WHERE"


class Query(object):
    """Query objects are the pipes through which ALL communication to
    the database is made.

    A query is the work horse for supersql as it is in SQL. Queries
    can be connected to the database or isolated and used for only
    generating vendor specific SQL strings.
    """

    def __init__(self, vendor, user=None, password=None, host=None, silent=True, pool=False):
        """Query constructor
        Sets up a query engine for use by saving initialization
        parameters internally for use in connecting to the backing
        database engine as required to execute queries.

        ..params:

        engine {str | required}:
            The database engine i.e. postgres, oracle,
            mysql, mssql etc.

        user {str | optional}:
            The user to connect as
        
        password {str | optional}:
            The password for the provided user
        
        host {str | optional}:
            Where is the server seated, what ip, port and what is the
            name of the database?
        
        silent {bool | optional}:
            Flag to specify if there should be syntax check i.e. should
            supersql check for an error in your prepared query before even
            sending it out to the database engine?
            Defaults to `True` i.e. do not check for errors
        
        pool {bool | optional}:
            Should the connection be made using connection pooling or not
        """
        if vendor not in SUPPORTED_VENDORS:
            raise NotImplementedError(f"{vendor} is not a supersql supported engine")
        self._vendor = vendor
        self._user = user
        self._password = password
        self._host = host
        self._silent = silent
        self._pool = pool
        self._sql = []

        self._pristine = True
        self._disparity = 0  # how many tables is this query for?
        self._from = []
        self._callstack = []

        self._tablenames = set()
        self._orphans = set()
        self._alias = None

    def _clone(self):
        """
        # ! Why copy of Query object is returned
        # A copy of query is returned so internal query variables do not
        # leak out to other query objects.
        #
        # i.e. it is possible to declare a global query object with connection
        # configuration and reuse without fear of internal state corruption
        """
        return type(self)(
            vendor=self._vendor,
            user=self._user,
            password=self._password,
            host=self._host,
            silent=self._silent
        )

    def _conditionator(self, condition):
        if isinstance(condition, str):
            return f" {condition}"
        else:
            try:
                return f" {condition.print()}"
            except AttributeError:
                msg = "Condition clause `WHERE` can only process strings or column comparison operations"
                raise ArgumentError(msg)

    def database(self, name):
        """
        Get the database provided by name for inspection or operation
        """
        return Database()

    def execute(self, *args, **kwargs):
        """
        Flushes the SQL command to the server for execution

        ..raises:
        {ConnectionError} On failed connection attempts to the database engine

        {StatementError} When the SQL Statement is malformed i.e. before being
            sent to the database engine and only if `silent = False`
        
        {CommandError} When the database server could not execute the command
            sent from preparation of your query. Wraps the message of the
            database server internally for easy debugging.
        """
        return connection(self)

    def get_tablename(self, table):
        """
        Irrespective of the type of naming convention used i.e.
        3 part "schema.table.columnname" or 4 part "db.schema.table.columnname"
        this method expects the second part from the right to always be 
        table name
        """
        if isinstance(table, Table):
            return Table.__tn__()
        elif isinstance(table, str):
            processed = table.split(".")
            parts = len(processed)
            return None if parts < 2 else processed[-2]
        else:
            column = table
            return column._imeta.__tn__()

    def print(self):
        """
        Prints the current SQL statement as it exists on the query object

        ..return {str}  String representation of the SQL command to be sent
            to the database server if `execute` method is called.
        """
        return "".join(self._sql)

    def run(self, *args, **kwargs):
        """
        This works the same as execute but returns iterable results
        that can be traversed via a cursor as opposed to fetching all results
        into memory.
        
        For errors and parameter options see :func: `~supersql.core.query.Query.execute`
        """
        return self.execute(*args, **kwargs)

    def was_called(self, command):
        return command in self._callstack

    def warn(self, command):
        if self._callstack[-1] == command:
            msg = f'Invalid Query Chaining: repeated {command} more than once'
            raise SQLError(msg)

    def AND(self, condition):
        self._sql.append(_AND)
        sql_snippet = self._conditionator(condition)
        self._sql.append(sql_snippet)
        return self

    def AS(self, alias):
        """
        # ? Revisit this docstring to make the op for this method clearer and more understandable
        If select was last called then use this as the alias for select
        """
        self._alias = alias
        return self

    def CREATE(self, table, safe=True):
        """
        Create a database bound entity table/database/extension/view etc
        currently supports only table creation.

        ..params:

        table {str, Table, Database, View, Extension | required}

        safe {boolean | optional}
            Set to false to set as unsafe as raise an error if table already
            exists and a create statement is issued otherwise 'IF NOT EXISTS' is
            added to the create table statement.
        """
        this = self._clone()
        this._sql.append(table.ddl(this._vendor, safe))
        return this

    def DELETE(self, table):
        return self.DELETE_FROM(table)

    def DELETE_FROM(self, table):
        this = self._clone()
        this._callstack.append(DELETE)
        if isinstance(table, str):
            tablename = table
        elif isinstance(table, Table):
            tablename = table.__tn__()

        this._tablenames.add(tablename)
        sqlstatement = f'DELETE FROM {tablename}'
        this._sql.append(sqlstatement)
        return this

    def FETCH(self, *args):...

    def FROM(self, *args, **kwargs):
        """
        Adds the left and right spaced _FROM_ constant to the pseudo call stack first
        then evaluates the argument(s) passed in for type
        to resolve the tablename and alias i.e. _from_ `abc as alias, bdd as alias`
        or _from_ `abc` before appending that to the _sql array
        """
        self._callstack.append(_FROM_)

        num_of_args = len(args)
        num_of_tables = len(self._tablenames)
        # msg = f"Found different tables in SELECT statement but only 1 command received in FROM"
        msg = f"tables:{num_of_tables}, args:{num_of_args}"

        if num_of_args != num_of_tables and num_of_tables > 0:
            raise MissingArgumentError(msg)

        for source in args:
            has_alias = None
            if isinstance(source, str):
                _a = None
                _q = source
            elif isinstance(source, Query):
                _a = source._alias
                _q = f"({source.print()})"
            elif isinstance(source, Table):
                td = source
                _a = source._alias
                _q = source.__tn__()

            _from = f"{_q} AS {_a}" if _a and num_of_tables > 1 else f"{_q}"
            self._from.append(_from)

        _sql = ", ".join(
            [f"{table}" for table in self._from]
        ) if len(self._from) > 1 else "".join(table for table in self._from)
        self._sql.extend([_FROM_, _sql])

        return self

    def GROUP_BY(self, *args):...

    def INSERT(self, *args):...

    def JOIN(self, *args):...

    def LIMIT(self, *args):...

    def ORDER_BY(self, *args):...

    def SELECT(self, *args):
        this = self._clone()
        this._callstack.append(SELECT)

        num_of_args = len(args)
        if num_of_args == 0:
            separator = "*"
        elif num_of_args == 1:
            arg = args[0]
            if isinstance(arg, str):
                _table = this.get_tablename(arg)
                this._tablenames.add(
                    _table
                ) if _table else this._orphans.add(arg)
                separator = arg
            elif isinstance(arg, Table):
                separator = "*"
                this._tablenames.add(arg.__tn__())
            else:
                separator = arg._name
                this._tablenames.add(arg._meta.__tn__())
        else:
            cols = []
            unique_tablenames = set([this.get_tablename(table) for table in args])
            is_heterogeneous = len(unique_tablenames) > 1

            for member in args:
                if isinstance(member, str):
                    _table = this.get_tablename(member)
                    this._tablenames.add(_table) if _table else None

                    cols.append(member)
                elif isinstance(member, Table):
                    # activate alias on tables if is_heterogeneous i.e.
                    # more than one table seen in select
                    if is_heterogeneous:
                        member.AS(member._alias if member._alias else member.__tn__())

                    this._tablenames.add(member.__tn__())
                    cols.extend(member.columns())
                else:
                    # Column object i.e. String, Int, Varchar etc.
                    alias = member._imeta._alias
                    tablename = member._imeta.__tn__()
                    member._imeta.AS(alias if alias else tablename)

                    this._tablenames.add(tablename)
                    if is_heterogeneous:
                        cols.append(f"{member._imeta._alias}.{member._name}")
                    else:
                        cols.append(f"{member._name}")

            separator = ", ".join(
                [col for col in cols]
            )

        _select_statement = f"SELECT {separator}"

        # if this._pristine:
        #     this._sql.append(_select_statement)
        #     this._pristine = False
        this._sql.append(_select_statement)
        return this

    def UNION(self, *args):...

    def UPDATE(self, *args):...

    def UPSERT(self, *args):...

    def WHERE(self, condition):
        """
        Check if _from was called and call it only if not a
        DELETE statement in play.

        Add WHERE constant to the pseudo call stack
        
        Process the `where ` condition by checking if condition is a
            - string: add as is with spaces as necessary
            - supersql datatype comparator: calls print() to get the SQL string repr.
        """
        if _FROM_ not in self._callstack:
            if DELETE not in self._callstack:
                tablenames = self._tablenames
                self = self.FROM(*tablenames)

        self._callstack.append(WHERE)
        self._sql.append(f" {WHERE}")
        sql_snippet = self._conditionator(condition)
        self._sql.append(sql_snippet)

        return self

    def WITHOUT(self, *args):...

    def WITH_RECURSIVE(self, *args):...
