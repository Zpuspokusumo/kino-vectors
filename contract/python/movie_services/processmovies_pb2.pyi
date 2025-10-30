from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class MovieInfo(_message.Message):
    __slots__ = ("title", "director", "year", "genre", "actors", "summary")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DIRECTOR_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    GENRE_FIELD_NUMBER: _ClassVar[int]
    ACTORS_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    title: str
    director: str
    year: str
    genre: _containers.RepeatedScalarFieldContainer[str]
    actors: _containers.RepeatedScalarFieldContainer[str]
    summary: str
    def __init__(self, title: _Optional[str] = ..., director: _Optional[str] = ..., year: _Optional[str] = ..., genre: _Optional[_Iterable[str]] = ..., actors: _Optional[_Iterable[str]] = ..., summary: _Optional[str] = ...) -> None: ...

class MovieInfos(_message.Message):
    __slots__ = ("Movies",)
    MOVIES_FIELD_NUMBER: _ClassVar[int]
    Movies: _containers.RepeatedCompositeFieldContainer[MovieInfo]
    def __init__(self, Movies: _Optional[_Iterable[_Union[MovieInfo, _Mapping]]] = ...) -> None: ...

class ProcessMovieResponse(_message.Message):
    __slots__ = ("status", "message", "items_processed", "unprocessed_items")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    ITEMS_PROCESSED_FIELD_NUMBER: _ClassVar[int]
    UNPROCESSED_ITEMS_FIELD_NUMBER: _ClassVar[int]
    status: int
    message: str
    items_processed: int
    unprocessed_items: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, status: _Optional[int] = ..., message: _Optional[str] = ..., items_processed: _Optional[int] = ..., unprocessed_items: _Optional[_Iterable[str]] = ...) -> None: ...

class RecommendMoviesRequest(_message.Message):
    __slots__ = ("text_query", "genres", "year_gte", "year_lte")
    TEXT_QUERY_FIELD_NUMBER: _ClassVar[int]
    GENRES_FIELD_NUMBER: _ClassVar[int]
    YEAR_GTE_FIELD_NUMBER: _ClassVar[int]
    YEAR_LTE_FIELD_NUMBER: _ClassVar[int]
    text_query: str
    genres: _containers.RepeatedScalarFieldContainer[str]
    year_gte: str
    year_lte: str
    def __init__(self, text_query: _Optional[str] = ..., genres: _Optional[_Iterable[str]] = ..., year_gte: _Optional[str] = ..., year_lte: _Optional[str] = ...) -> None: ...

class RecommendMoviesResponse(_message.Message):
    __slots__ = ("status", "quantity", "Movies")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    QUANTITY_FIELD_NUMBER: _ClassVar[int]
    MOVIES_FIELD_NUMBER: _ClassVar[int]
    status: int
    quantity: int
    Movies: _containers.RepeatedCompositeFieldContainer[MovieInfo]
    def __init__(self, status: _Optional[int] = ..., quantity: _Optional[int] = ..., Movies: _Optional[_Iterable[_Union[MovieInfo, _Mapping]]] = ...) -> None: ...
