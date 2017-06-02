(function (exports) {
    'use strict';

    function name (val) { return Object.keys(this).find(function (k) { return this[k] === val }, this) }

    function enumify (arr, start) {
        let i = (start === undefined) ? 0 : start;
        let obj = arr.reduce(function (a, c) { a[c] = i++; return a; }, {});
        obj.name = name.bind(obj);
        Object.freeze(obj);
        return obj;
    }

    exports.PlayState = enumify([
        "STOPPED",
        "PLAYING",
        "PAUSED"
    ], 1);

    exports.MessageType = enumify([
        "NOW_PLAYING",
        "VIEW_UPDATE",
        "COMMAND"
    ], 1);

    exports.Command = enumify([
        "SET_PLAYSTATE",
        "SEEK_TO",
        "PLAY_NEXT",
        "PLAY_PREV",
        "PLAY_QUEUE_FROM_POSITION",
        "SET_VOLUME",
        "SET_SHUFFLE",
        "SET_REPEAT_MODE",
        "REQUEST_VIEW",
        "ADD_TO_QUEUE",
        "ADD_TO_QUEUE_AND_PLAY",
        "REMOVE_FROM_QUEUE",
        "CLEAR_QUEUE",
        "SAVE_PLAYLIST",
        "DELETE_PLAYLIST",
        "RENAME_PLAYLIST",
        "LOAD_PLAYLIST"
    ], 1);

    exports.ViewType = enumify([
        "ALL_ARTISTS",
        "ARTIST_DETAIL",
        "ALL_ALBUMS",
        "ALBUM_DETAIL",
        "ALL_TRACKS",
        "PLAYLIST",
        "PLAYLIST_DETAIL",
        "QUEUE",
        "SETTINGS" // local only
    ], 1);

    exports.RepeatMode = enumify([
        "OFF",
        "ALL",
        "ONE"
    ], 1);
})(typeof exports === 'undefined' ? window.enums = {} : exports);