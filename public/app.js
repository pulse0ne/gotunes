/**
 * Created by tsned on 3/21/17.
 */
(function (angular) {
    'use strict';
    const app = angular.module('gotunes', [
        'ngRoute',
        'ngAnimate',
        'ngMaterial',
        'ngWebsocket',
        'ngIdle',
        'ngFitText',
        'gotunes.track.slider',
        'ngLongPress',
        'ngLocalStorage',
        'QuickList'
    ]);

    /**
     * Setup routing
     */
    app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
        $routeProvider.when('/', {
            templateUrl: 'views/nav/nav.html',
            controller: 'gotunes.controller.nav'
        });
        $routeProvider.when('/idle', {
            templateUrl: 'views/idle/idle.html',
            controller: 'gotunes.controller.idle'
        });
        $routeProvider.otherwise({
            redirectTo: '/'
        });

        $locationProvider.html5Mode(true);
    }]);

    /**
     * Constants
     */
    app.constant('SVGs', {
        play: 'M 0 0 L 0 32 L 32 16 z',
        pause: 'M 2 0 L 2 32 L 12 32 L 12 0 L 2 0 M 20 0 L 20 32 L 30 32 L 30 0 L 20 0'
    });

    app.constant('themes', [
        { id: 0, name: 'Default', url: 'assets/css/theme-default.css' },
        { id: 1, name: 'Odin', url: 'assets/css/theme-odin.css' },
        { id: 2, name: 'Aether', url: 'assets/css/theme-aether.css' }
    ]);

    app.factory('gotunes.default.settings', ['themes', function (themes) {
        return {
            theme: themes[0],
            disabledNav: [],
            idleDelay: 20,
            idleEnabled: true,
            dynamicTitle: true
        };
    }]);

    app.directive('settingsEntry', function () {
        return {
            restrict: 'E',
            replace: true,
            scope: { settingsLabel: '@' },
            transclude: true,
            template:
                '<div layout="row" layout-align="center center" class="setting-container">' +
                '    <span class="settings-label">{{settingsLabel}}:</span>' +
                '</div>',
            link: function (scope, element, attrs, ctl, transclude) {
                transclude(function (item) { element.append(item) });
            }
        };
    });

    /**
     * Ws handling service
     */
    app.service('gotunes.service', [
        '$websocket',
        '$rootScope',
        function ($websocket, $rootScope) {
            const self = this;
            const MessageType = (window.enums || {}).MessageType || {};

            const ws = $websocket.$new({
                url: 'ws://' + window.location.host + '/websocket',
                reconnect: true,
                reconnectInterval: 1000,
                maxReconnectInterval: 10000,
                doubleInterval: true,
                immediate: false
            });

            self.sendCommand = function (cmdType, data) {
                ws.$send({type: MessageType.COMMAND, payload: {command: cmdType, data: data}});
            };

            ws.$on('$message', function (message) {
                message = JSON.parse(message);
                $rootScope.$broadcast('ws.message', message);
            });

            ws.$on('$open', function () {
                $rootScope.$broadcast('ws.open');
            });

            ws.$on('$close', function (code) {
                $rootScope.$broadcast('ws.close', code);
            });

            ws.$on('$error', function () {
                $rootScope.$broadcast('ws.error');
            });

            self.connect = ws.$open;

            self.UUID = (function () {
                    let self = {};
                    let t = [];
                    let rand = new Uint32Array(4);

                    for (let i = 0; i < 256; i++) {
                        t[i] = (i < 16 ? '0' : '') + i.toString(16);
                    }
                    let crypto = window.crypto || window.msCrypto; // for IE

                    self.generate = function () {
                        crypto.getRandomValues(rand);
                        let w = rand[0];
                        let x = rand[1];
                        let y = rand[2];
                        let z = rand[3];

                        return t[w & 0xff] + t[w >> 8 & 0xff] + t[w >> 16 & 0xff] + t[w >> 24 & 0xff] + '-' +
                            t[x & 0xff] + t[x >> 8 & 0xff] + '-' +
                            t[x >> 16 & 0x0f | 0x40] + t[x >> 24 & 0xff] + '-' +
                            t[y & 0x3f | 0x80] + t[y >> 8 & 0xff] + '-' +
                            t[y >> 16 & 0xff] + t[y >> 24 & 0xff] + t[z & 0xff] + t[z >> 8 & 0xff] + t[z >> 16 & 0xff] + t[z >> 24 & 0xff];
                    };

                    return self;
                })();
        }
    ]);

    app.controller('gotunes.controller.main', [
        '$scope',
        '$rootScope',
        '$timeout',
        '$location',
        '$document',
        '$window',
        '$mdDialog',
        '$localStorage',
        'ngIdle',
        'gotunes.service',
        'gotunes.default.settings',
        'themes',
        'SVGs',
        function ($scope, $rootScope, $timeout, $location, $document, $window, $mdDialog, $localStorage, ngIdle, gotunes, defaults, themes, svg) {
            const enums = window.enums || {};
            const Command = enums.Command || {};
            const PlayState = $scope.ps = enums.PlayState || {};
            const ViewType = $scope.vt = enums.ViewType || {};
            const MessageType = enums.MessageType || {};
            const unicodePlay = '\u25B6 ';
            let lastVolume = 50;
            let domSlider;

            $scope.themes = themes;
            $scope.settings = $localStorage.load('gotunes.settings') || defaults;
            if ($scope.settings.theme.id === 0) $scope.settings.theme = themes[0]; // hack to populate settings dropdown
            $scope.disableAllTracks = $scope.settings.disabledNav.includes(ViewType.ALL_TRACKS);
            $scope.disableAllAlbums = $scope.settings.disabledNav.includes(ViewType.ALL_ALBUMS);
            $scope.disableAllArtists = $scope.settings.disabledNav.includes(ViewType.ALL_ARTISTS);
            $scope.wsConnected = false;
            $scope.nowplaying = {
                playstate: PlayState.PAUSED,
                time: { current: 0, total: 0 },
                track: null
            };
            $scope.view = { type: null, data: [] };
            $scope.svgPath = svg.play;

            $window.oncontextmenu = function (event) {
                event.preventDefault();
                event.stopPropagation();
                return false;
            };

            $scope.$watch('nowplaying.volume', function (newVal) {
                if (newVal !== lastVolume) {
                    domSlider.noUiSlider.set(Math.max(Math.min(newVal, 100), 0));
                }
            });

            $scope.$watch('nowplaying.playstate', function (newVal) {
                $scope.svgPath = newVal !== PlayState.PLAYING ? svg.play : svg.pause;
                if (newVal === PlayState.PLAYING && $document[0].title.substring(0, unicodePlay.length) !== unicodePlay) {
                    $document[0].title = unicodePlay + $document[0].title;
                } else if ($document[0].title.substring(0, unicodePlay.length) === unicodePlay) {
                    $document[0].title = $document[0].title.substring(unicodePlay.length);
                }
            });

            $scope.$watch('nowplaying.track', function (nv, ov) {
                let track = $scope.nowplaying.track;
                if ($scope.settings.dynamicTitle && nv && nv !== ov && (track.artist || track.title)) {
                    let artist = track.artist || 'Unknown';
                    let song = track.title || track.filename || 'Unknown';
                    let ps = $document[0].title.substring(0, 1) === unicodePlay ? unicodePlay : '';
                    $document[0].title = ps + '"' + song + '" by ' + artist;
                    $rootScope.showStatusBar = $location.path() !== '/idle';
                    $rootScope.statusBarText = '"' + song + '" by ' + artist;
                } else {
                    $document[0].title = 'gotunes';
                    $rootScope.showStatusBar = false;
                }
            }, true);

            $scope.$watch('settings', function (nv) {
                $localStorage.save('gotunes.settings', nv);
                $scope.disableAllTracks = nv.disabledNav.includes(ViewType.ALL_TRACKS);
                $scope.disableAllAlbums = nv.disabledNav.includes(ViewType.ALL_ALBUMS);
                $scope.disableAllArtists = nv.disabledNav.includes(ViewType.ALL_ARTISTS);
            }, true);

            $scope.$on('ngIdle', function () {
                if (!$scope.settings.idleEnabled) return;
                if ($location.path() !== '/idle' && $scope.nowplaying.playstate === PlayState.PLAYING) {
                    $rootScope.showStatusBar = false;
                    $location.path('/idle');
                }
            });

            $scope.$on('ws.open', function () {
                gotunes.sendCommand(Command.REQUEST_VIEW, ViewType.QUEUE);
            });

            $scope.$on('ws.message', function (evt, msg) {
                switch(msg.type) {
                    case MessageType.VIEW_UPDATE:
                        $scope.view.type = msg.payload.type;
                        $scope.view.data = msg.payload.data;
                        break;
                    case MessageType.NOW_PLAYING:
                        $scope.nowplaying = msg.payload;
                        break;
                    default:
                        break;
                }
                $scope.$apply();
            });

            $scope.onSeek = function (percent) {
                let pos = $scope.nowplaying.time.total * (percent / 100);
                gotunes.sendCommand(Command.SEEK_TO, pos);
            };

            $scope.togglePlaystate = function () {
                let s = ($scope.nowplaying || {}).playstate === PlayState.PLAYING ? PlayState.PAUSED : PlayState.PLAYING;
                gotunes.sendCommand(Command.SET_PLAYSTATE, s);
            };

            $scope.playNext = angular.bind(this, gotunes.sendCommand, Command.PLAY_NEXT);
            $scope.playPrevious = angular.bind(this, gotunes.sendCommand, Command.PLAY_PREV);

            $scope.playInPosition = function (pos) {
                gotunes.sendCommand(Command.PLAY_QUEUE_FROM_POSITION, pos);
            };

            $scope.loadPlaylist = function (name) {
                gotunes.sendCommand(Command.LOAD_PLAYLIST, name);
            };

            $scope.openQueueMenu = function (track) {
                console.log('TODO');
            };

            $scope.savePlaylist = function () {
                $mdDialog.show(
                    $mdDialog.prompt()
                        .textContent('Enter playlist name:')
                        .placeholder('playlist1')
                        .ariaLabel('playlist name')
                        .cancel('Cancel')
                        .ok('Save')
                ).then(function (name) {
                    if (!name || name === '') {
                        name = gotunes.UUID.generate();
                    }
                    gotunes.sendCommand(Command.SAVE_PLAYLIST, name);
                });
            };

            $scope.openSettings = function () {
                $scope.view.type = ViewType.SETTINGS;
            };

            // TODO: slider doesn't work very well...need to fix it
            domSlider = $document[0].getElementById('slider');
            noUiSlider.create(domSlider, {
                start: lastVolume,
                connect: [true, false],
                range: { min: 0, max: 100 }
            });

            domSlider.noUiSlider.on('start', function (val) { lastVolume = val });
            domSlider.noUiSlider.on('end', function (val) { gotunes.sendCommand(Command.SET_VOLUME, Math.floor(val)) });

            let ignored = angular.element('#sticky-footer');

            gotunes.connect();
            ngIdle.setIdle($scope.settings.idleDelay);
            ngIdle.ignore(ignored);
            ngIdle.start();

            $scope.test = function (e) { console.log('test1 ' + e) };
            $scope.test2 = function (e) { console.log('test2 ' + e) };
        }
    ]);

    app.filter('trackTime', function () {
        return function (val) {
            let dur = moment.duration(val * 1000);
            let time = moment(dur.asMilliseconds()).format('mm:ss');
            return time === 'Invalid date' ? '00:00' : time;
        }
    });
})(window.angular);
