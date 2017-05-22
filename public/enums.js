(function (exports) {

    function name (val) { return Object.keys(this).find(function (k) { return this[k] === val }, this) }

    function _enumify (obj) {
        obj.name = name.bind(obj);
        Object.freeze(obj);
        return obj;
    }

    exports.PlayState = _enumify({
        STOPPED: 1,
        PLAYING: 2,
        PAUSED: 3
    });

    exports.MessageType = _enumify({
        NOW_PLAYING: 1,
        VIEW_UPDATE: 2,
        COMMAND: 3
    });

    exports.Command = _enumify({
        SET_PLAYSTATE: 1,
        SEEK_TO: 2,
        PLAY_NEXT: 3,
        PLAY_PREV: 4,
        PLAY_FROM_CONTEXT: 5,
        SET_VOLUME: 6,
        SET_CONTEXT: 7,
        REQUEST_VIEW: 8,
        NEW_PLAYLIST: 9,
        SAVE_PLAYLIST: 10,
        ADD_TO_PLAYLIST: 11
    });

    exports.ContextType = _enumify({
        ALL_ARTISTS: 1,
        ARTIST_DETAIL: 2,
        ALL_ALBUMS: 3,
        ALBUM_DETAIL: 4,
        ALL_TRACKS: 5,
        PLAYLIST: 6,
        PLAYLIST_DETAIL: 7,
    });

    exports.RepeatMode = _enumify({
        OFF: 1,
        ALL: 2,
        ONE: 3
    });
})(typeof exports === 'undefined' ? window.enums = {} : exports);